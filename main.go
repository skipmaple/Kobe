package main

import (
	"github.com/skipmaple/kobe/gitcode"
	"github.com/skipmaple/kobe/github"
	"os"
	"strconv"
	"strings"
)

func main() {
	host := os.Args[1]
	if host == "gitcode" {
		privateToken := os.Args[2]
		if len(privateToken) == 0 {
			return
		}

		projectId, _ := strconv.Atoi(os.Args[3])

		issueId := 1
		if len(os.Args[4]) > 0 {
			issueId, _ = strconv.Atoi(os.Args[4])
		}

		city := os.Args[5]
		if len(city) == 0 {
			city = ""
		}

		gitcode.GetUp(privateToken, projectId, issueId, city)

	} else if host == "github" {
		privateToken := os.Args[2]
		if len(privateToken) == 0 {
			return
		}

		fullPath := os.Args[3]
		arr := strings.Split(fullPath, "/")
		owner, repo := arr[0], arr[1]

		city := os.Args[4]
		if len(city) == 0 {
			city = ""
		}

		github.GetUp(privateToken, owner, repo, city)
	}

}
