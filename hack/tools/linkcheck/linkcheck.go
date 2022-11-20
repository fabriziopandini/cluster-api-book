//gox:build tools
// +xbuild tools

/*
Copyright 2022 The Kubernetes Authors.

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
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"k8s.io/utils/pointer"
)

const (
	contentFolder   = "content"
	anchorSeparator = "#"
)

var (
	root          = pflag.String("root", ".", "root path to walk for linting .m files")
	hugoFolder    = pflag.String("hugo-folder", "", "path to the folder contaning the hugo website")
	hugoLanguages = pflag.StringSlice("hugo-languages", []string{"en"}, "list of languages supported by the hugo website") // TODO: infer from config.toml.
	verbose       = pflag.Bool("verbose", false, "verbose")
)

var (
	// pages being processes.
	pages       []*page
	pagesByPath map[string]*page

	// protect pages from concurrent access.
	pagesLock sync.RWMutex
)

// page define a page validated by linkcheck.
type page struct {
	// path to the file the page has been read from.
	path string

	// fatalError if set, defines an error in reading or processing the page that prevents further processing.
	fatalError string

	// isHugoPage is true when the page is defined inside content/language folder of the hugo website.
	isHugoPage bool

	// hugoLanguage is the language of the page in the hugo website.
	hugoLanguage string

	// hugoPath is path of the page relative to the content/language folder of the hugo website.
	hugoPath string

	// links contains the list of links defined in the page.
	links []link

	// anchors contains the list of anchors (~headers) defined in the page.
	anchors []string
}

// link define a link on a page validated by linkcheck.
type link struct {
	// rawLink is the link as it is defined in the page.
	rawLink string

	// fatalError if set, defines an error in reading or processing the link that prevents further processing.
	fatalError string

	// URL derived from the link.
	// NOTE: for localLinks (link to files) the link path is translated to an absolute path.
	URL *url.URL
}

func newPage(path string) page {
	p := page{path: path}

	// If the page is inside the hugo content dir.
	contentDir := filepath.Join(*root, *hugoFolder, contentFolder)
	if strings.HasPrefix(path, contentDir) {
		p.isHugoPage = true

		// Identify the page language or error out if the page does not belong to one of the know languages.
		switch len(*hugoLanguages) {
		case 0:
			// TODO: handle non localized hugo websites
		default:
			for _, l := range *hugoLanguages {
				languageDir := filepath.Join(contentDir, l)
				if strings.HasPrefix(path, languageDir) {
					p.hugoLanguage = l
					p.hugoPath = strings.TrimPrefix(path, languageDir)
				}
			}
			if p.hugoLanguage == "" {
				p.fatalError = fmt.Sprintf("hugo page %s does not belong to one of the know languages: %s", strings.TrimPrefix(path, contentDir), strings.Join(*hugoLanguages, ", "))
			}
		}
	}
	return p
}

func newPageWithFatalError(path string, error string) page {
	return page{path: path, fatalError: error}
}

func addPage(p page) {
	pagesLock.Lock()
	defer pagesLock.Unlock()
	pages = append(pages, &p)
	if pagesByPath == nil {
		pagesByPath = map[string]*page{}
	}
	pagesByPath[p.path] = &p
}

func (p *page) addLink(l string) {
	u, err := url.Parse(l)
	if err != nil {
		p.links = append(p.links, link{rawLink: l, fatalError: fmt.Sprintf("error parsing url: %v", err)})
		return
	}

	// if it is a file url (no scheme is considered file url)
	if u.Scheme == "" {
		// Error if file url is used in pages outside the hugo website.
		// TODO: think about pages outside hugo content/language folder, should we support file url? how this behaves in github?
		if !p.isHugoPage {
			p.links = append(p.links, link{rawLink: l, fatalError: "scheme is required on links outside the hugo website"})
			return
		}

		// Parse the link extracting the key parts.
		path, fragment, language, err := parseLink(l)
		if err != nil {
			p.links = append(p.links, link{rawLink: l, fatalError: err.Error()})
			return
		}
		if path == "" {
			// if path is empty the link is a fragment pointing to an anchor on the current page (e.g. #anchor).
			path = filepath.Base(p.hugoPath)
		}

		// Compute the path of the target page, transforming relative paths to absolute ones.
		if !filepath.IsAbs(path) {
			// TODO: validate the relative path is contained in the hugo root
			path = filepath.Join(filepath.Dir(p.hugoPath), path)
		}

		// Compute the content dir where the target page will be hosted.
		contentDir := filepath.Join(*root, *hugoFolder, contentFolder)

		// Compute the language of the target page.
		// Use the language detected from the hugoRef or default to the same language of the page where the link is defined.
		if language == "" {
			language = p.hugoLanguage
		}

		// Compute the url pointing to the target page.
		rawURL := filepath.Join(contentDir, language, path)
		if fragment != "" {
			rawURL += fragment
		}
		URL, err := url.Parse(rawURL)
		if err != nil {
			p.links = append(p.links, link{rawLink: l, fatalError: fmt.Sprintf("error parsing ref/refLink path: %v", err)})
			return
		}
		p.links = append(p.links, link{URL: URL, rawLink: l})
		return
	}

	// otherwise it is an http/https url, use as it is.
	// TODO: link title, e.g. [Duck Duck Go](https://duckduckgo.com "The best search engine for privacy")
	p.links = append(p.links, link{rawLink: l, URL: u})
}

// This pattern applies to the addr part of [text](addr) and searches for {{< tag "value" >}}, captures both tag and value values.
// ^ and $ are used to avoid more tags on
var refRx = regexp.MustCompile(`^\s*\{\{<\s*([\S\#]+)\s+\"([^\s=]+)\"\s*>\}\}\s*$`)

func parseLink(rawLink string) (path, fragment, language string, err error) {
	// if the rawLink is a ref/refLink shortcode.
	refs := refRx.FindAllStringSubmatch(rawLink, -1)
	// TODO: ref/refLink path=
	// TODO: ref/refLink path= language=
	if len(refs) == 1 && (refs[0][1] == "ref" || refs[0][1] == "refLink") {
		// Grab the path and the fragment, if any.
		path, fragment := splitPathAndFragment(refs[0][2])
		return path, fragment, "", nil
	}

	// Otherwise it means the link is not a ref/refLink shortcode; we assume it is a plain markdown link.
	path, fragment = splitPathAndFragment(rawLink)
	return path, fragment, "", nil
}

func splitPathAndFragment(addr string) (string, string) {
	path := addr
	fragment := ""
	if strings.Contains(path, anchorSeparator) {
		s := strings.Split(path, anchorSeparator)
		fragment = strings.TrimPrefix(path, s[0])
		path = s[0]
	}
	// TODO: special case for _index.md: index.md can be reference either by its path or by its containing folder without the ending /
	// TODO: link title, e.g. [page](page.md "a tooltip")
	return path, fragment
}

// readAll markdown pages from the root folder.
func readAll() error {
	// Handle a collection of read tasks working in parallel.
	var readmdwg sync.WaitGroup
	if err := filepath.Walk(*root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				addPage(newPageWithFatalError(path, fmt.Sprintf("Error walking path %s: %v", path, err)))
				return nil
			}

			if filepath.Ext(path) == ".md" {
				// For each markdown page to read, create a read task.
				readmdwg.Add(1)
				go func() {
					addPage(readMarkdownPage(path))
					readmdwg.Done()
				}()
			}
			return nil
		}); err != nil {
		return errors.Errorf("Error walking path %s: %v", *root, err)
	}

	// Waits for all the read tasks to complete.
	readmdwg.Wait()
	return nil
}

// readMarkdownPage reads a markdown page
func readMarkdownPage(path string) page {
	p := newPage(path)

	// Gets the page content.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		p.fatalError = fmt.Sprintf("Error reading content: %v", err)
		return p
	}

	// Gets the list of anchors in the page.
	p.anchors = readMarkdownAnchors(string(content))

	// Gets the list of links in the page.
	links := readMarkdownLinks(string(content))
	for _, l := range links {
		p.addLink(l)
	}
	return p
}

// Search for markdown headers.
// (?m) is required to force multiline search due to ^ and $ used to exclude other things on the same line.
var anchorRx = regexp.MustCompile(`(?m)^\s*\#+\s*(.+)$`)

func readMarkdownAnchors(body string) (anchors []string) {
	// TODO: check if we need to do something for repeated anchors
	mv := anchorRx.FindAllStringSubmatch(body, -1)
	for _, m := range mv {
		ref := strings.ToLower(strings.TrimSpace(m[1]))
		ref = strings.ReplaceAll(ref, " ", "-")
		ref = strings.ReplaceAll(ref, "/", "")
		anchors = append(anchors, ref)
	}
	return
}

// Search for links in the format [text](addr), captures addr value.
// [^\!] is required to drop image links ![]()
var lRx = regexp.MustCompile(`[^\!]\[[^\]]+\]\(([^\)]+)\)`)

// Search for reference links in the format [text]: addr, captures addr value.
// (?m) is required to force multiline search due to ^ and $ used to exclude other things on the same line.
var referencelRx = regexp.MustCompile(`(?m)^\[[^\]]+\]\:\s+(.+)$`)

func readMarkdownLinks(body string) (links []string) {
	seen := map[string]struct{}{}
	mv := lRx.FindAllStringSubmatch(body, -1)
	for _, m := range mv {
		ref := m[1]
		if _, ok := seen[ref]; !ok {
			seen[ref] = struct{}{}
			links = append(links, ref)
		}
	}
	mv = referencelRx.FindAllStringSubmatch(body, -1)
	for _, m := range mv {
		ref := m[1]
		if _, ok := seen[ref]; !ok {
			seen[ref] = struct{}{}
			links = append(links, ref)
		}
	}
	return
}

// linkcheckAll all pages.
func linkcheckAll() error {
	// Handle a collection of validate tasks working in parallel.
	var validatewg sync.WaitGroup

	pagesLock.RLock()
	for i, p := range pages {
		validatewg.Add(1)
		go func() {
			// Perform page link check, which can take some time depending by the number of urls.
			linkcheckPage(p.path)

			// When page validation is completed, update page.
			// TODO: check if we need this because linkcheckPage changes the page in place...
			pagesLock.Lock()
			defer pagesLock.Unlock()
			pages[i] = p

			validatewg.Done()
		}()
	}
	pagesLock.RUnlock()

	// Waits for all the validate tasks to complete.
	validatewg.Wait()
	return nil
}

func linkcheckPage(path string) {
	pagesLock.RLock()
	p, ok := pagesByPath[path]
	pagesLock.RUnlock()
	if !ok {
		panic(fmt.Sprintf("linkcheckPage %s which has not been read before", path))
	}

	// If the page already has been marked with a fatal error, link should not be checked.
	if p.fatalError != "" {
		return
	}

	for i, l := range p.links {
		// If the link already has been marked with a fatal error, skip it.
		if l.fatalError != "" {
			continue
		}

		// If it is a file url (no scheme is considered file url)
		if l.URL.Scheme == "" {
			// Check the links targets an existing page.
			if _, err := os.Stat(l.URL.Path); errors.Is(err, os.ErrNotExist) {
				l.fatalError = fmt.Sprintf("%s does not exist", l.URL.Path)
				p.links[i] = l
				continue
			}

			pagesLock.RLock()
			targetp, ok := pagesByPath[l.URL.Path]
			pagesLock.RUnlock()
			if !ok {
				// TODO: this should never happen (if we protect from link outside root). Might be we should panic here...
				l.fatalError = fmt.Sprintf("%s does has not been processed by linkcheck", l.URL.Path)
				p.links[i] = l
				continue
			}

			// If the link targets an anchor, check it exists.
			if l.URL.Fragment != "" {
				found := false
				for _, a := range targetp.anchors {
					if l.URL.Fragment == a {
						found = true
						break
					}
				}
				if !found {
					l.fatalError = fmt.Sprintf("%s%s does exists in the target page", anchorSeparator, l.URL.Fragment)
					p.links[i] = l
					continue
				}
			}
		}
		// TODO: handle http or https; use a map of links to avoid duplicated http calls.
	}
	return
}

func main() {
	pflag.Parse()
	if *root == "." {
		path, err := os.Getwd()
		if err != nil {
			fmt.Printf("ERROR: failed to get current working directory: %v\n", err)
			os.Exit(1)
		}
		root = pointer.String(path)
	}
	if !filepath.IsAbs(*root) {
		path, err := filepath.Abs(*root)
		if err != nil {
			fmt.Printf("ERROR: failed to convert root to an absolute path: %v\n", err)
			os.Exit(1)
		}
		root = pointer.String(path)
	}

	if err := readAll(); err != nil {
		fmt.Printf("ERROR: failed to read pages: %v\n", err)
		os.Exit(1)
	}

	if err := linkcheckAll(); err != nil {
		fmt.Printf("ERROR: failed to check links on pages: %v\n", err)
		os.Exit(1)
	}

	pagesLock.Lock()
	defer pagesLock.Unlock()
	sort.SliceStable(pages, func(i, j int) bool { return pages[i].path < pages[j].path })

	fmt.Println()
	for _, p := range pages {
		s := ""
		prints := false
		s += fmt.Sprintf("PAGE: %s\n", p.path)
		switch {
		case p.fatalError != "":
			prints = true
			s += fmt.Sprintln()
			s += fmt.Sprintf(" - ERROR: %s\n", p.fatalError)
			break
		default:
			t := ""
			errorst := 0
			for _, l := range p.links {
				switch {
				case l.fatalError != "":
					prints = true
					errorst++
					t += fmt.Sprintf(" - ERROR: %s: %s\n", l.rawLink, l.fatalError)
					break
				default:
					if *verbose {
						t += fmt.Sprintf(" - OK: %s\n", l.rawLink)
					}
				}
			}
			switch errorst {
			case 0:
				s += fmt.Sprintf("      %d links, no errors\n\n", len(p.links))
				break
			default:
				s += fmt.Sprintf("      %d links, %d errors\n\n", len(p.links), errorst)
			}
			if t != "" {
				s += fmt.Sprintf("%s\n", t)
			}
		}

		if *verbose || prints {
			fmt.Print(s)
		}
	}
	fmt.Printf("Total page processed: %d\n", len(pages))

}
