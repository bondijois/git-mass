package cmd

import (
	"context"
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

var ctx = context.Background()

func getGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func verifyAuth(username, token string) bool {
	client := getGithubClient(token)
	user, _, err := client.Users.Get(ctx, username)
	if err != nil {
		log.Fatalln("Invalid credentials. Please update username and token")
		return false
	}
	fmt.Printf("Successfully authenticated as user: %v\n", user.GetLogin())
	return true
}

func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func fetchAllRepos(token, organization string) []*github.Repository {
	client := getGithubClient(token)
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}
	var allRepos []*github.Repository
	i := 1
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, organization, opt)
		if err != nil {
			log.Fatalln(err)
		}
		allRepos = append(allRepos, repos...)
		fmt.Printf("Fetching repos from page %d\n", i)
		i++
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos
}

func cloneRepo(username string, token string, path string, repo *github.Repository, errors chan error) {
	_, err := git.PlainClone(
		path, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: username,
				Password: token,
			},
			URL: repo.GetCloneURL(),
		},
	)
	if err != nil {
		errors <- err
	} else {
		errors <- nil
	}
}

func pullRepo(username string, token string, path string, errors chan error) {
	r, openErr := git.PlainOpen(path)
	if openErr != nil {
		errors <- openErr
	}
	w, workErr := r.Worktree()
	if workErr != nil {
		errors <- workErr
	}
	pullErr := w.Pull(
		&git.PullOptions{
			Auth: &http.BasicAuth{
				Username: username,
				Password: token,
			},
		},
	)
	if pullErr != nil {
		errors <- pullErr
	} else {
		errors <- nil
	}
}
