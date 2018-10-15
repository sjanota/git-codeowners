# git-codeowners

## Installation
```
go get github.com/sjanota/git-codeowners
```

## Usage
git-codeowners is simple tool that helps you find out who to ask for review. When you run it in your repository when 
on branch different then master it will find out what changes you made, consult it with CODEOWNERS file and tell 
you which rules from CODEOWNERS file applies to which files. Example:

```
$ git codeowners
=> /proj1/ [@sjanota @octocat]
	* installation/resources/installer-config-cluster.yaml.tpl
	* installation/resources/installer-config-local.yaml.tpl
=> /resources/proj2 [@octocat]
	* resources/core/charts/configurations-generator/templates/deployment.yaml
	* resources/core/charts/configurations-generator/values.yaml
```

You may also provide list of people who have already approved your PR. Rules that shows them as owners will be filtered 
out:

```
$ git codeowners @sjanota
=> /resources/proj2 [@octocat]
	* resources/core/charts/configurations-generator/templates/deployment.yaml
	* resources/core/charts/configurations-generator/values.yaml
```
