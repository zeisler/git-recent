package main

import (
	"fmt"
	"github.com/posener/complete"
	"log"
	"os"
	"os/exec"
	"strings"
)

var BranchCache []string

type RecentBranch struct {
}

func DefaultBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output := strings.Builder{}
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	names := strings.Split(strings.TrimSuffix(output.String(), "\n"), "/")

	if len(names) <= 1 {
		log.Fatal(fmt.Sprintf("unknown default branch: '%v'", names))
	}

	return names[len(names)-1]
}

func PopulateBranches() {
	cmd := exec.Command("git", "for-each-ref", fmt.Sprintf("--no-merged=%v", DefaultBranch()), "--count=20", "--sort=-committerdate", "refs/heads/", `--format="%(refname:short)"`)
	output := strings.Builder{}
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
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
		var branch string
		if len(os.Args) > 1 {
			branch = os.Args[1]
		} else if len(BranchCache) > 0 {
			branch = BranchCache[0]
		}
		if branch != "" {
			cmd := exec.Command("git", "checkout", branch)
			b, err := cmd.CombinedOutput()
			fmt.Println(string(b))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
