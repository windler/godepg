package matcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type GoPackagesMatcher struct {
	text string
}

var (
	_          Matcher = &FilterMatcher{}
	goPackages *[]string
)

func NewGoPackagesMatcher(text string) *GoPackagesMatcher {
	return &GoPackagesMatcher{
		text: text,
	}
}

func (f *GoPackagesMatcher) Matches() bool {
	for _, m := range getGoPackages() {
		if m == f.text || strings.HasPrefix(f.text, m+"/") {
			return true
		}
	}
	return false
}

type GithubGolangContent struct {
	Type string
	Name string
}

func getGoPackages() []string {
	if goPackages == nil {
		goPackages = &[]string{}

		response, err := http.Get("https://api.github.com/repos/golang/go/contents/src")
		if err != nil {
			log.Fatal(err.Error())
			return *goPackages
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)

		var content []GithubGolangContent
		err = json.Unmarshal(body, &content)
		for _, c := range content {
			if c.Type == "dir" {
				*goPackages = append(*goPackages, c.Name)
			}
		}
	}

	return *goPackages

}
