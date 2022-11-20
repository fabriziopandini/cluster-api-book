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
	. "github.com/onsi/gomega"
	"net/url"
	"os"
	"path/filepath"
	"testing"
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
				fatalError: "error parsing url: parse \"$$$%%%???\": invalid URL escape \"%%%\"",
			},
		},
		{
			name: "https url",
			path: "/root/test.md",
			url:  "https://www.google.com",
			wantUrl: link{
				rawLink: "https://www.google.com",
				URL:     mustParseUrl("https://www.google.com"),
			},
		},
		{
			name: "file url",
			path: "/root/test.md",
			url:  "another-page.md",
			wantUrl: link{
				rawLink:    "another-page.md",
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
				fatalError: "error parsing url: parse \"$$$%%%???\": invalid URL escape \"%%%\"",
			},
		},
		{
			name: "https url",
			path: "/root/hugo/content/en/test.md",
			url:  "https://www.google.com",
			wantUrl: link{
				rawLink: "https://www.google.com",
				URL:     mustParseUrl("https://www.google.com"),
			},
		},
		{
			name: "relative path url",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "another-page.md",
			wantUrl: link{
				rawLink: "another-page.md",
				URL:     mustParseUrl("/root/hugo/content/en/folder/another-page.md"),
			},
		},
		{
			name: "relative path url with anchor on the current page",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "test.md#anchor",
			wantUrl: link{
				rawLink: "test.md#anchor",
				URL:     mustParseUrl("/root/hugo/content/en/folder/test.md#anchor"),
			},
		},
		{
			name: "url with anchor on the current page (without path)",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "#anchor",
			wantUrl: link{
				rawLink: "#anchor",
				URL:     mustParseUrl("/root/hugo/content/en/folder/test.md#anchor"),
			},
		},
		{
			name: "relative path url with anchor to another page",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "another.md#anchor",
			wantUrl: link{
				rawLink: "another.md#anchor",
				URL:     mustParseUrl("/root/hugo/content/en/folder/another.md#anchor"),
			},
		},
		{
			name: "absolute path url linking a file in content/language root",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"/another-page.md\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"/another-page.md\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/another-page.md"),
			},
		},
		{
			name: "absolute path url linking a file in a folder nested into content/language root",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"/folder/another-page.md\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"/folder/another-page.md\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/folder/another-page.md"),
			},
		},
		{
			name: "ref with relative url linking a file in content/language root",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"../another-page.md\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"another-page.md\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/another-page.md"),
			},
		},
		{
			name: "ref with relative path linking a file in a folder nested into content/language root",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"another-page.md\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"another-page.md\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/folder/another-page.md"),
			},
		},
		{
			name: "ref with absolute path linking a file in content/language root",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"/folder/another-page.md\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"/folder/another-page.md\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/folder/another-page.md"),
			},
		},
		{
			name: "ref with absolute path linking a file in a folder nested into content/language root",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"/folder/another-page.md\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"/folder/another-page.md\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/folder/another-page.md"),
			},
		},
		{
			name: "ref linking to an anchor on the current page (with path)",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"test.md#anchor\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"test.md#anchor\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/folder/test.md#anchor"),
			},
		},
		{
			name: "ref linking to an anchor on the current page (without path)",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"#anchor\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"#anchor\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/folder/test.md#anchor"),
			},
		},
		{
			name: "ref linking to an anchor on another page",
			path: "/root/hugo/content/en/folder/test.md",
			url:  "{{< ref \"another.md#anchor\" >}}",
			wantUrl: link{
				rawLink: "{{< ref \"another.md#anchor\" >}}",
				URL:     mustParseUrl("/root/hugo/content/en/folder/another.md#anchor"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			page := newPage(tt.path)
			page.addLink(tt.url)
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

	anotherp := newPage(filepath.Join(contentDir, "en/another.md"))
	anotherp.anchors = []string{"anchor"}

	tests := []struct {
		name      string
		page      func() page
		wantLinks []link
	}{
		// pages inside the hugo website
		{
			name: "page with a valid ref linking to the current page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"test.md\" >}}")
				return p
			},
			wantLinks: []link{
				{
					rawLink: "{{< ref \"test.md\" >}}",
					URL:     mustParseUrl(filepath.Join(contentDir, "en/test.md")),
				},
			},
		},
		{
			name: "page with a valid ref linking to the another page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"another.md\" >}}")
				return p
			},
			wantLinks: []link{
				{
					rawLink: "{{< ref \"another.md\" >}}",
					URL:     mustParseUrl(filepath.Join(contentDir, "en/another.md")),
				},
			},
		},
		{
			name: "page with a valid ref linking to an anchor on the current page (without path)",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"#anchor\" >}}")
				p.anchors = []string{"anchor"}
				return p
			},
			wantLinks: []link{
				{
					rawLink: "{{< ref \"#anchor\" >}}",
					URL:     mustParseUrl(filepath.Join(contentDir, "en/test.md#anchor")),
				},
			},
		},
		{
			name: "page with a valid ref linking to an anchor on the current page (without path)",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"#anchor\" >}}")
				p.anchors = []string{"anchor"}
				return p
			},
			wantLinks: []link{
				{
					rawLink: "{{< ref \"#anchor\" >}}",
					URL:     mustParseUrl(filepath.Join(contentDir, "en/test.md#anchor")),
				},
			},
		},
		{
			name: "page with a valid ref linking to an anchor on the another page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"another.md#anchor\" >}}")
				return p
			},
			wantLinks: []link{
				{
					rawLink: "{{< ref \"another.md#anchor\" >}}",
					URL:     mustParseUrl(filepath.Join(contentDir, "en/another.md#anchor")),
				},
			},
		},
		{
			name: "page with an invalid ref",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"invalid.md\" >}}")
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "{{< ref \"invalid.md\" >}}",
					URL:        mustParseUrl(filepath.Join(contentDir, "en/invalid.md")),
					fatalError: fmt.Sprintf("%s does not exist", filepath.Join(contentDir, "en/invalid.md")),
				},
			},
		},
		{
			name: "page with a valid ref linking to an invalid anchor on the current page (without path)",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"#invalid\" >}}")
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "{{< ref \"#invalid\" >}}",
					URL:        mustParseUrl(filepath.Join(contentDir, "en/test.md#invalid")),
					fatalError: "#invalid does exists in the target page",
				},
			},
		},
		{
			name: "page with a valid ref linking to an invalid anchor on the another page",
			page: func() page {
				p := newPage(filepath.Join(contentDir, "en/test.md"))
				p.addLink("{{< ref \"another.md#invalid\" >}}")
				return p
			},
			wantLinks: []link{
				{
					rawLink:    "{{< ref \"another.md#invalid\" >}}",
					URL:        mustParseUrl(filepath.Join(contentDir, "en/another.md#invalid")),
					fatalError: "#invalid does exists in the target page",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			pagesLock.Lock()
			p := tt.page()
			pages = []*page{&p, &anotherp}
			pagesByPath = map[string]*page{p.path: &p, anotherp.path: &anotherp}
			pagesLock.Unlock()

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
