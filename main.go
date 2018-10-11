package main

import (
	"os"
	"log"
	"fmt"
	"github.com/nathanleiby/github-codeowners"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
	"strings"
)

type match struct {
	mapping *parser.Mapping
	files []string
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get cwd: %v", err)
	}

	coMappings := loadCodeOwnersMappings(wd)
	changedFiles := getChangedFiles(wd)

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

func getChangedFiles(wd string) []string {
	repo, err := git.PlainOpen(wd)
	if err != nil {
		log.Fatalf("Error while opening repository: %v", err)
	}

	masterTree, err := getBranchTree(repo, "master")
	if err != nil {
		log.Fatalf("Error while reading master tree: %v", err)
	}

	head, err := repo.Head()
	if err != nil {
		log.Fatalf("Error while reading HEAD: %v", err)
	}

	headTree, err := getReferenceTree(repo, head)
	if err != nil {
		log.Fatalf("Error while reading HEAD tree: %v", err)
	}

	patch ,err := masterTree.Patch(headTree)
	if err != nil {
		log.Fatalf("Error while calculating patch from master to HEAD: %v", err)
	}

	changed := make(map[string]bool)
	for _, patch := range patch.FilePatches() {
		from, to := patch.Files()
		if to != nil && to.Path() != "" {
			changed[to.Path()] = true
		}
		if from != nil && from.Path() != "" {
			changed[from.Path()] = true
		}
	}

	result := make([]string, 0)
	for f := range changed {
		result = append(result, f)
	}

	return result
}

func getBranchTree(repo *git.Repository, branchName string) (*object.Tree, error) {
	branch, err := repo.Branch(branchName)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Reference(branch.Merge, false)
	if err != nil {
		return nil, err
	}

	return getReferenceTree(repo, ref)
}

func getReferenceTree(repo *git.Repository, ref *plumbing.Reference) (*object.Tree, error) {
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	return commit.Tree()
}
