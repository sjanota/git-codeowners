# git-codeowners

## Installation
```
go get github.com/sjanota/git-codeowners
```

## Usage
git-codeowners is simple tool that helps you find out who to ask for review. When you run it in your repository when 
on branch different then master it will find out what changes you made, consult it with CODEOWNERS file and tell 
you which rules from CODEOWNERS file applies to which file. Example:

```
$ git codeowners
=> /proj1/ [@sjanota @octocat]
	* proj1/resources/config.yaml.tpl
	* proj1/src/main.go
=> /resources/proj2 [@octocat]
	* resources/proj2/resources/config.yaml.tpl
	* resources/proj2/src/main.go
```

You may also provide list of people who have already approved your PR. Rules that shows them as owners will be filtered 
out:

```
$ git codeowners @sjanota
=> /resources/proj2 [@octocat]
	* resources/proj2/resources/config.yaml.tpl
	* resources/proj2/src/main.go
```
