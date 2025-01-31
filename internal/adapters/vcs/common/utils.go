package common

import "strings"

func ParseRepoString(repo string) (owner string, repoName string) {
	parts := strings.Split(repo, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", repo
}
