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
	"testing"

	. "github.com/onsi/gomega"
)

func Test_newPage(t *testing.T) {
	cancel := setFlags("/root", "hugo", []string{"en"})
	defer cancel()

	tests := []struct {
		name     string
		path     string
		wantPage page
	}{
		{
			name: "page outside the hugo website",
			path: "/root/test.md",
			wantPage: page{
				path:       "/root/test.md",
				isHugoPage: false,
			},
		},
		{
			name: "page in the hugo website",
			path: "/root/hugo/content/en/test.md",
			wantPage: page{
				path:         "/root/hugo/content/en/test.md",
				isHugoPage:   true,
				hugoLanguage: "en",
				hugoPath:     "/test.md",
			},
		},
		{
			name: "page in the hugo website - invalid language",
			path: "/root/hugo/content/it/test.md",
			wantPage: page{
				path:       "/root/hugo/content/it/test.md",
				fatalError: "hugo page /it/test.md does not belong to one of the know languages: en",
				isHugoPage: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			page := newPage(tt.path)
			g.Expect(page).To(Equal(tt.wantPage))
		})
	}
}

func Test_addUrl(t *testing.T) {
	cancel := setFlags("/root", "hugo", []string{"en"})
	defer cancel()

	tests := []struct {
		name    string
		path    string
		url     string
		wantUrl link
	}{
		// pages outside the hugo website
		{
			name: "invalid url",
			path: "/root/test.md",
			url:  "$$$%%%???",
			wantUrl: link{
				rawLink:    "$$$%%%???",
				lineNumber: 1,
				fatalError: "error parsing url: parse \"$$$%%%???\": invalid URL escape \"%%%\"",
			},
		},
		{
			name: "https url",
			path: "/root/test.md",
			url:  "https://www.google.com",
			wantUrl: link{
				rawLink:    "https://www.google.com",
				lineNumber: 1,
				URL:        mustParseUrl("https://www.google.com"),
			},
		},
		{
			name: "file url",
			path: "/root/test.md",
			url:  "another-page.md",
			wantUrl: link{
				rawLink:    "another-page.md",
				lineNumber: 1,
				fatalError: "scheme is required on links outside the hugo website",
			},
		},

		// pages inside the hugo website
		{
			name: "invalid url",
			path: "/root/hugo/content/en/test.md",
			url:  "$$$%%%???",
			wantUrl: link{
				rawLink:    "$$$%%%???",
				lineNumber: 1,
				fatalError: "error parsing url: parse \"$$$%%%???\": invalid URL escape \"%%%\"",
			},
		},
		{
			name: "url with ref/reflink",
			path: "/root/hugo/content/en/test.md",
			url:  "{{< ref \"something\" >}}",
			wantUrl: link{
				rawLink:    "{{< ref \"something\" >}}",
				lineNumber: 1,
				fatalError: "ref/refLink shortcodes must not be used, use \"something\" instead",
			},
		},
		{
			name: "url with _index.md",
			path: "/root/hugo/content/en/test.md",
			url:  "something/_index.md",
			wantUrl: link{
				rawLink:    "something/_index.md",
				lineNumber: 1,
				fatalError: "links must not end with _index.md, use \"something/\" instead",
			},
		},
		{
			name: "https url",
			path: "/root/hugo/content/en/test.md",
			url:  "https://www.google.com",
			wantUrl: link{
				rawLink:    "https://www.google.com",
				lineNumber: 1,
				URL:        mustParseUrl("https://www.google.com"),
			},
		},
		{
			name: "relative path url",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "another-page",
			wantUrl: link{
				rawLink:    "another-page",
				lineNumber: 1,
				URL:        mustParseUrl("/root/hugo/content/en/folder/another-page.md"),
			},
		},
		{
			name: "relative path url with anchor on the current page (with path)",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "test#anchor",
			wantUrl: link{
				rawLink:    "test#anchor",
				lineNumber: 1,
				URL:        mustParseUrl("/root/hugo/content/en/folder/test.md#anchor"),
			},
		},
		{
			name: "relative path url with anchor on the current page (without path)",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "#anchor",
			wantUrl: link{
				rawLink:    "#anchor",
				lineNumber: 1,
				URL:        mustParseUrl("/root/hugo/content/en/folder/test.md#anchor"),
			},
		},
		{
			name: "relative path url with anchor on another page",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "another#anchor",
			wantUrl: link{
				rawLink:    "another#anchor",
				lineNumber: 1,
				URL:        mustParseUrl("/root/hugo/content/en/folder/another.md#anchor"),
			},
		},
		{
			name: "absolute path url linking a file in content/en",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "/another-page",
			wantUrl: link{
				rawLink:    "/another-page",
				lineNumber: 1,
				URL:        mustParseUrl("/root/hugo/content/en/another-page.md"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			page := newPage(tt.path)
			page.addLink(tt.url, 1)
			g.Expect(page.links[0]).To(Equal(tt.wantUrl))
		})
	}
}

func Test_linkcheckPage(t *testing.T) {
	g := NewWithT(t)

	root, err := os.MkdirTemp("", "linkcheck")
	g.Expect(err).ToNot(HaveOccurred())
	defer os.RemoveAll(root)

	hugoFolder := "hugo"
	hugoLanguages := []string{"en"}

	cancel := setFlags(root, hugoFolder, hugoLanguages)
	defer cancel()

	contentDir := filepath.Join(root, hugoFolder, contentFolder)

	touch(g, filepath.Join(contentDir, "en/test.md"))
	touch(g, filepath.Join(contentDir, "en/another.md"))
	touch(g, filepath.Join(contentDir, "en/folder/_index.md"))

	anotherp := newPage(filepath.Join(contentDir, "en/another.md"))
	anotherp.anchors = []string{"anchor"}

	indexp := newPage(filepath.Join(contentDir, "en/folder/_index.md"))
	indexp.anchors = []string{"anchor"}

	tests := []struct {
		name      string
		page      func() page
		wantLinks []link
	}{
		// TODO: pages outside the hugo website

		// pages inside the hugo website
		{
			name: "page with a valid ref linking to the current page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("test", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "test",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/test.md")),
				},
			},
		},
		{
			name: "page with a valid ref linking to the another page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("another", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "another",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/another.md")),
				},
			},
		},
		{
			name: "page with a valid ref linking to an anchor on the current page (with path)",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("test#anchor", 1)
				p.anchors = []string{"anchor"}
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "test#anchor",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/test.md#anchor")),
				},
			},
		},
		{
			name: "page with a valid ref linking to an anchor on the current page (without path)",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("#anchor", 1)
				p.anchors = []string{"anchor"}
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "#anchor",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/test.md#anchor")),
				},
			},
		},
		{
			name: "page with a valid ref linking to an anchor on the another page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("another#anchor", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "another#anchor",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/another.md#anchor")),
				},
			},
		},
		{
			name: "page with an link to _index.md",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("/folder", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "/folder",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/folder/_index.md")),
				},
			},
		},
		{
			name: "page with an link to an anchor in _index.md",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("/folder#anchor", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "/folder#anchor",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/folder/_index.md#anchor")),
				},
			},
		},
		{
			name: "page with an link to an anchor in _index.md",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/folder/test.md"))
				p.addLink("../folder#anchor", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "../folder#anchor",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/folder/_index.md#anchor")),
				},
			},
		},
		{
			name: "page with an invalid ref",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("invalid", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "invalid",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/invalid.md")),
					fatalError: fmt.Sprintf("the link resolves to %s which does not exist", "/hugo/content/en/invalid.md"),
				},
			},
		},
		{
			name: "page with a valid ref linking to an invalid anchor on the current page (without path)",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("#invalid", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "#invalid",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/test.md#invalid")),
					fatalError: "#invalid does exists in <site>/content/en/test.md",
				},
			},
		},
		{
			name: "page with a valid ref linking to an invalid anchor on the another page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("another#invalid", 1)
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "another#invalid",
					lineNumber: 1,
					URL:        mustParseUrl(filepath.Join(contentDir, "en/another.md#invalid")),
					fatalError: "#invalid does exists in <site>/content/en/another.md",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			p := tt.page()
			pages = []*page{&p, &anotherp, &indexp}
			pagesByPath = map[string]*page{p.path: &p, anotherp.path: &anotherp, indexp.path: &indexp}

			linkcheckPage(p.path)

			g.Expect(p.links).To(Equal(tt.wantLinks))
		})
	}
}

func setFlags(rootValue, hugoFolderValue string, hugoLanguagesValue []string) (resetFlags func()) {
	rootBefore := root
	hugoFolderBefore := hugoFolder
	hugoLanguagesBefore := hugoLanguages

	root = &rootValue
	hugoFolder = &hugoFolderValue
	hugoLanguages = &hugoLanguagesValue

	return func() {
		root = rootBefore
		hugoFolder = hugoFolderBefore
		hugoLanguages = hugoLanguagesBefore
	}
}

func mustParseUrl(l string) *url.URL {
	u, err := url.Parse(l)
	if err != nil {
		panic(err.Error())
	}
	return u
}

func touch(g *WithT, path string) {
	g.Expect(os.MkdirAll(filepath.Dir(path), os.ModePerm)).ToNot(HaveOccurred())
	g.Expect(os.WriteFile(path, []byte(""), 0600)).To(Succeed())
}
