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
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

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

	// lineNumber where the link has been found.
	lineNumber int

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
	pages = append(pages, &p)
	if pagesByPath == nil {
		pagesByPath = map[string]*page{}
	}
	pagesByPath[p.path] = &p
}

func (p *page) addLink(l string, lineNumber int) {
	u, err := url.Parse(l)
	if err != nil {
		p.links = append(p.links, link{rawLink: l, lineNumber: lineNumber, fatalError: fmt.Sprintf("error parsing url: %v", err)})
		return
	}

	// if it is a file url (no scheme is considered file url)
	if u.Scheme == "" {
		// Error if file url is used in pages outside the hugo website.
		// TODO: think about pages outside hugo content/language folder, should we support file url? how this behaves in github?
		if !p.isHugoPage {
			p.links = append(p.links, link{rawLink: l, lineNumber: lineNumber, fatalError: "scheme is required on links outside the hugo website"})
			return
		}

		// Parse the link extracting the key parts.
		path, fragment, language, err := parseLink(l)
		if err != nil {
			p.links = append(p.links, link{rawLink: l, lineNumber: lineNumber, fatalError: err.Error()})
			return
		}
		if path == "" {
			// if path is empty the link is a fragment pointing to an anchor on the current page (e.g. #anchor).
			// NOTE: drop .md from the page name so it aligns to how links behaves in hugo
			// TODO: think about pages outside hugo content/language folder, should we support file url? how this behaves in github?
			path = strings.TrimSuffix(filepath.Base(p.hugoPath), ".md")
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

		// If the target page is a directory, add _index.md
		isDir, err := isDirectory(rawURL)
		if err != nil {
			p.links = append(p.links, link{rawLink: l, lineNumber: lineNumber, fatalError: fmt.Sprintf("error checking if path is a directory: %v", err)})
			return
		}
		if isDir {
			rawURL = filepath.Join(rawURL, "_index.md")
		}

		// if it is not a dirctory, then it is an .md file
		// TODO: what about html files
		if !isDir {
			rawURL += ".md"
		}

		if fragment != "" {
			rawURL += fragment
		}
		URL, err := url.Parse(rawURL)
		if err != nil {
			p.links = append(p.links, link{rawLink: l, lineNumber: lineNumber, fatalError: fmt.Sprintf("error parsing url: %v", err)})
			return
		}
		p.links = append(p.links, link{URL: URL, rawLink: l, lineNumber: lineNumber})
		return
	}

	// otherwise it is an http/https url, use as it is.
	// TODO: link title, e.g. [Duck Duck Go](https://duckduckgo.com "The best search engine for privacy")
	p.links = append(p.links, link{rawLink: l, lineNumber: 1, URL: u})
}

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return fileInfo.IsDir(), err
}

func (p *page) logPath() string {
	if p.isHugoPage {
		return fmt.Sprintf("<site>/content/%s%s", p.hugoLanguage, p.hugoPath)
	}
	return fmt.Sprintf("<root>/%s", strings.TrimPrefix(p.path, *root))
}

// This pattern applies to the addr part of [text](addr) and searches for {{< tag "value" >}}, captures both tag and value values.
// ^ and $ are used to avoid more tags on
var refRx = regexp.MustCompile(`^\s*\{\{<\s*([\S\#]+)\s+\"([^\s=]+)\"\s*>\}\}\s*$`)

func parseLink(rawLink string) (path, fragment, language string, err error) {
	// if the rawLink is a ref/refLink shortcode.
	// NOTE: this makes .md files easier to write/read; it is also aligned with common practice in use for the K8s website.
	refs := refRx.FindAllStringSubmatch(rawLink, -1)
	if len(refs) == 1 && (refs[0][1] == "ref" || refs[0][1] == "refLink") {
		return "", "", "", errors.Errorf("ref/refLink shortcodes must not be used, use %q instead", refs[0][2])
	}

	// Otherwise it is a plain markdown link.
	path, fragment = splitPathAndFragment(rawLink)

	// In hugo the folder name must be used when referring to "_index.md"
	if filepath.Base(path) == "_index.md" {
		return "", "", "", errors.Errorf("links must not end with _index.md, use \"%s/\" instead", filepath.Dir(path))
	}

	// In hugo the file name must not have the .md extension.
	if filepath.Ext(path) == ".md" {
		return "", "", "", errors.Errorf("links must not have .md extension, use \"%s%s\" instead", strings.TrimSuffix(path, ".md"), fragment)
	}

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
	if err := filepath.Walk(*root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				addPage(newPageWithFatalError(path, fmt.Sprintf("Error walking path %s: %v", path, err)))
				return nil
			}

			if filepath.Ext(path) == ".md" {
				addPage(readMarkdownPage(path))
			}
			return nil
		}); err != nil {
		return errors.Errorf("Error walking path %s: %v", *root, err)
	}
	return nil
}

// readMarkdownPage reads a markdown page
func readMarkdownPage(path string) page {
	p := newPage(path)

	// Gets the page content.
	content, err := os.ReadFile(path)
	if err != nil {
		p.fatalError = fmt.Sprintf("Error reading content: %v", err)
		return p
	}

	// Gets the list of anchors in the page.
	p.anchors = readMarkdownAnchors(string(content))

	// Gets the list of links in the page.
	for i, line := range strings.Split(string(content), "\n") {
		links := readMarkdownLineLinks(line)
		for _, l := range links {
			p.addLink(l, i+1)
		}
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
var referencelRx = regexp.MustCompile(`^\s+\[[^\]]+\]\:\s+(.+)$`)

func readMarkdownLineLinks(line string) (links []string) {
	mv := lRx.FindAllStringSubmatch(line, -1)
	for _, m := range mv {
		links = append(links, m[1])
	}
	mv = referencelRx.FindAllStringSubmatch(line, -1)
	for _, m := range mv {
		links = append(links, m[1])
	}
	return
}

// linkcheckAll all pages.
func linkcheckAll() error {
	for i := range pages {
		p := pages[i]

		// Perform page link check, which can take some time depending by the number of urls.
		linkcheckPage(p.path)

		// When page validation is completed, update page.
		// TODO: check if we need this because linkcheckPage changes the page in place...
		pages[i] = p
	}
	return nil
}

func linkcheckPage(path string) {
	p, ok := pagesByPath[path]
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
				l.fatalError = fmt.Sprintf("the link resolves to %s which does not exist", strings.TrimPrefix(l.URL.Path, *root))
				p.links[i] = l
				continue
			}

			targetp, ok := pagesByPath[l.URL.Path]
			if !ok {
				// TODO: this should never happen (if we protect from link outside root). Might be we should panic here...
				l.fatalError = fmt.Sprintf("the link resolves to %s which has not been processed by linkcheck", l.URL.Path)
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
					l.fatalError = fmt.Sprintf("%s%s does exists in %s", anchorSeparator, l.URL.Fragment, targetp.logPath())
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

	sort.SliceStable(pages, func(i, j int) bool { return pages[i].path < pages[j].path })

	fmt.Println()

	a := 0
	l := 0
	for i := range pages {
		p := pages[i]
		a += len(p.anchors)
		l += len(p.links)

		s := ""
		prints := false
		s += fmt.Sprintf("PAGE: %s\n", p.logPath())
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
					t += fmt.Sprintf(" - ERROR: line %d, %s: %s\n", l.lineNumber, l.rawLink, l.fatalError)
					break
				default:
					if *verbose {
						t += fmt.Sprintf(" - OK: line %d, %s\n", l.lineNumber, l.rawLink)
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
	fmt.Printf("Total page processed: %d links: %d anchors: %d \n", len(pages), l, a)

}
