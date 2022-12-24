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
	cmd := exec.Command(`git for-each-ref --no-merged=master --count=50 --sort=-committerdate refs/heads/ --format="%(refname:short)"`)
	output := strings.Builder{}
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		if !strings.Contains(err.Error(), " no such file or directory") {
			log.Fatal(err)
		}
	}

	BranchCache = strings.Split(output.String(), "\n")

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
		"git-switch",
		complete.Command{
			Args: RecentBranch{},
		},
	)

	// AddFlags adds the completion flags to the program flags,
	// in case of using non-default flag set, it is possible to pass
	// it as an argument.
	// it is possible to set custom flags name
	// so when one will type 'self -h', he will see '-complete' to install the
	// completion and -uncomplete to uninstall it.

	// if the completion did not do anything, we can run our program logic here.
	PopulateBranches()
	if cmp.Run() {
		fmt.Println("taco", "chicken")
		os.Exit(1)
	}
}
