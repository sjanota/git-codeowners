package main

import (
	"os"
	"log"
	"fmt"
	"github.com/nathanleiby/github-codeowners"
	"strings"
	"os/exec"
	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
)

type match struct {
	mapping *parser.Mapping
	files []string
}

func main() {
	root := getRootDir()
	coMappings := loadCodeOwnersMappings(root)
	changedFiles := getChangedFiles()

	matchingMappings := make(map[string]*match)

	approved := os.Args[1:]

	for _, changedFile := range changedFiles {
		mapping := findMatchingMapping(changedFile, coMappings)
		if mapping != nil && !isApproved(mapping, approved) {
			m, ok := matchingMappings[mapping.Path]
			if !ok {
				m = &match{
					mapping: mapping,
					files: make([]string, 0),
				}
				matchingMappings[mapping.Path] = m
			}
			m.files = append(m.files, changedFile)
		}
	}

	for _, mapping := range matchingMappings {
		fmt.Printf("=> %s %v\n", mapping.mapping.Path, mapping.mapping.Owners)
		for _, f := range mapping.files {
			fmt.Printf("\t* %s\n", f)
		}
	}
}

func isApproved(mapping *parser.Mapping, approved []string) bool {
	for _, owner := range mapping.Owners {
		for _, a := range approved {
			if string(owner) == a || string(owner) == "@" + a {
				return true
			}
		}
	}
	return false
}

func findMatchingMapping(path string, mappings []parser.Mapping) *parser.Mapping {
	for i := range mappings {
		mapping := mappings[len(mappings)-1-i]
		pattern := gitignore.ParsePattern(mapping.Path, []string{})
		if result := pattern.Match(strings.Split(path, "/"), false); result == gitignore.Exclude {
			return &mapping
		}
	}
	return nil
}

func loadCodeOwnersMappings(wd string) []parser.Mapping {
	coFilePath := fmt.Sprintf("%s/CODEOWNERS", wd)
	coFile, err := os.Open(coFilePath)
	if os.IsNotExist(err) {
		log.Fatal("There is no CODEOWNERS file")
	}
	if err != nil {
		log.Fatalf("Error while opening CODEOWNERS file: %v", err)
	}

	coMapping, err := parser.Parse(coFile)
	if err != nil {
		log.Fatalf("Error while parsing CODEOWNERS file: %v", err)
	}

	return coMapping
}

func getChangedFiles() []string {
	cmd := exec.Command("git", "diff", "--name-only", "master...HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error while getting diff: %v, %s", err, string(out))
	}

	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

func getRootDir() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error while getting root: %v, %s", err, string(out))
	}

	return strings.TrimSpace(string(out))
}
