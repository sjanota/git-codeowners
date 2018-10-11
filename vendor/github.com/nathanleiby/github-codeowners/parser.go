package parser

import (
	"bufio"
	"io"
	"strings"
)

// Mapping is a Path to []Owner pair
type Mapping struct {
	Path   string
	Owners []Owner
}

type Owner string

func Parse(r io.Reader) ([]Mapping, error) {
	owners := []Mapping{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		owner, ok := parseLine(line)
		if !ok {
			continue
		}
		owners = append(owners, owner)
	}

	if scanner.Err() != nil {
		return []Mapping{}, scanner.Err()
	}

	return owners, nil
}

func parseLine(line string) (Mapping, bool) {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "#") {
		return Mapping{}, false
	}

	fields := strings.Fields(trimmed)
	if !(len(fields) >= 2) {
		return Mapping{}, false
	}

	o := Mapping{}
	for idx, field := range fields {
		if idx == 0 {
			o.Path = field
			continue
		}

		if strings.HasPrefix(field, "#") {
			// mid line comment
			break
		}

		u := Owner(field)
		if u.IsValid() {
			o.Owners = append(o.Owners, u)
		}
	}

	return o, true
}

func (u Owner) IsGithubOwner() bool {
	return strings.HasPrefix(string(u), "@") && !strings.Contains(string(u), "/")
}

func (u Owner) IsGithubTeam() bool {
	return strings.HasPrefix(string(u), "@") && strings.Contains(string(u), "/")
}

func (u Owner) IsEmailAddress() bool {
	return !strings.HasPrefix(string(u), "@") && strings.Contains(string(u), "@")
}

func (u Owner) IsValid() bool {
	return u.IsGithubOwner() || u.IsGithubTeam() || u.IsEmailAddress()
}
