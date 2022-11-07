//go:build tools
// +build tools

/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"k8s.io/utils/pointer"
)

var (
	root          = pflag.String("root", ".", "root path to walk for linting .m files")
	hugoFolder    = pflag.String("hugo-folder", "", "path to the folder contaning the hugo website")
	hugoLanguages = pflag.StringSlice("hugo-languages", []string{"en"}, "path to the folder contaning the hugo website")
	verbose       = pflag.Bool("verbose", false, "verbose")
	// TODO: add debug
)

type page struct {
	path string

	// TODO: check if we really need all of this + related flags/processing; it depends by if there are specific behaviors for page inside/outside hugo-folder
	isHugoPage   bool
	hugoLanguage string
	hugoPath     string
}

type relocatedURL struct {
	rawUrl string
	*url.URL
}

func newPage(path string) page {
	p := page{path: path}

	if strings.HasPrefix(path, filepath.Join(*root, *hugoFolder, "content")) {
		p.isHugoPage = true

		if len(*hugoLanguages) > 0 {
			for _, l := range *hugoLanguages {
				if strings.HasPrefix(path, filepath.Join(*root, *hugoFolder, "content", l)) {
					p.hugoLanguage = l
					p.hugoPath = strings.TrimPrefix(path, filepath.Join(*root, *hugoFolder, "content", l))
				}
			}
			if p.hugoLanguage == "" {
				addProblem("Hugo page %s does not belong to one of the know languages: %s", strings.TrimPrefix(path, filepath.Join(*root, *hugoFolder, "content")), strings.Join(*hugoLanguages, ", "))
			}
		}
	}
	return p
}

var (
	readmdwg sync.WaitGroup      // outstanding reads
	readmdq  = make(chan string) // Files to read

	// TODO: move problems at page level
	problemsmu sync.Mutex
	problems   []string
)

func readAll() {
	log.Printf("Reading files from %s", *root)
	err := filepath.Walk(*root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				addProblem("Error walking path %s: %v", path, err)
			}

			if filepath.Ext(path) == ".md" {
				readmdwg.Add(1)
				go func() {
					readmdq <- path
				}()
			}
			return nil
		})

	if err != nil {
		addProblem("Error walking path %s: %v", *root, err)
	}

	readmdwg.Wait()
	close(readmdq)
}

func readMDLoop() {
	for path := range readmdq {
		if err := doReadMD(path); err != nil {
			addProblem("Error reading %s: %v", path, err)
		}
	}
}

func doReadMD(path string) error {
	defer readmdwg.Done()

	p := newPage(path)

	if *verbose {
		log.Printf("Reading %s %s %s", p.path, p.isHugoPage, p.hugoPath)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// TODO: get anchors

	links := localLinks(string(content))
	if *verbose {
		for _, l := range links {
			u, err := url.Parse(l)
			if err != nil {
				addProblem("Error parsing url %q on page %q: %v", l, path, err)
				continue
			}

			ru := relocatedURL{URL: u, rawUrl: l}

			if u.Scheme == "" {
				if !p.isHugoPage {
					// TODO: Check this is a problem when pages are shown in GitHub
					addProblem("url %q on %s desn't have scheme", l, path)
					continue
				}

				// TODO: enforce usage of rel/relref

				if u.Path == "" {
					ru.Path = p.path
				} else {
					ru.Path = filepath.Join(filepath.Dir(p.path), u.Path)

					// TODO: validate the relocated path is contained in the hugo root
				}

			}

			log.Printf("- %s (%s)", ru, ru.rawUrl)
		}
	}

	return nil
}

var aRx = regexp.MustCompile(`\[[^\]]*\]\(([^\)]+)\)`)

func localLinks(body string) (links []string) {
	seen := map[string]struct{}{}
	mv := aRx.FindAllStringSubmatch(body, -1)
	for _, m := range mv {
		ref := m[1]
		if strings.HasPrefix(ref, "/src/") {
			continue
		}
		if _, ok := seen[ref]; !ok {
			seen[ref] = struct{}{}
			links = append(links, m[1])
		}
	}
	return
}

func addProblem(format string, a ...any) {
	problemsmu.Lock()
	defer problemsmu.Unlock()

	msg := fmt.Sprintf(format, a...)
	if *verbose {
		log.Printf(" PROBLEM! %s", msg)
	}
	problems = append(problems, msg)
}

func main() {
	pflag.Parse()
	if *root == "." {
		path, err := os.Getwd()
		if err != nil {
			addProblem("Failed to get current working directory: %v", err)
			return
		}
		root = pointer.String(path)
	}
	if !filepath.IsAbs(*root) {
		path, err := filepath.Abs(*root)
		if err != nil {
			addProblem("Failed to convert root to an absolute path: %v", err)
			return
		}
		root = pointer.String(path)
	}

	go readMDLoop()
	readAll()

	// TODO: Report roblems
	if len(problems) > 0 {
		log.Printf("%d problem detected:", len(problems))
		for _, p := range problems {
			log.Printf(" - %s", p)
		}
		os.Exit(1)
	}
}
