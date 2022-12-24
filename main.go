// Package self
// a program that complete itself
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/posener/complete"
)

var BranchCache []string

type RecentBranch struct {
}

//
//func Match(arg string) bool {
//	for _, branch := range BranchCache {
//		if branch == arg {
//			return true
//		}
//	}
//
//	return false
//}

func PopulateBranches() {
	cmd := exec.Command("git", "for-each-ref", "--no-merged=main", "--count=50", "--sort=-committerdate", "refs/heads/", `--format="%(refname:short)"`)
	output := strings.Builder{}
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		//if !strings.Contains(err.Error(), " no such file or directory") {
		log.Fatal(err)
		//}
	}

	BranchCache = strings.Split(strings.Replace(output.String(), `"`, "", -1), "\n")

	if len(BranchCache) == 0 {
		BranchCache = append(BranchCache, "taco", "chicken")
	}
}

func FilterBranches(arg string) []string {
	var filteredBranches []string
	for _, branch := range BranchCache {
		if branch == arg {
			return []string{branch}
		}
		if strings.Contains(branch, arg) {
			filteredBranches = append(filteredBranches, branch)
		}
	}
	return filteredBranches
}

func (r RecentBranch) Predict(args complete.Args) []string {
	PopulateBranches()
	if args.Last == "" {
		return BranchCache
	}

	return FilterBranches(args.Last)
}

func main() {
	cmp := complete.New(
		"git-recent",
		complete.Command{
			Args: RecentBranch{},
		},
	)

	PopulateBranches()
	if cmp.Run() {
		fmt.Println(strings.Join(BranchCache, " "))
		os.Exit(1)
	} else {
		cmd := exec.Command("git", "checkout", os.Args[1])
		b, err := cmd.CombinedOutput()
		fmt.Println(string(b))
		if err != nil {
			log.Fatal(err)
		}
	}
}
