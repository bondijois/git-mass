package cmd

import (
	"fmt"
	"github.com/google/go-github/v41/github"
	"github.com/spf13/cobra"
	"log"
	"sync"
	"time"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone repositories",
	Long:  `clones all repositories of an organization`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		organization := config.GetString("organization")
		if organization == "" {
			log.Fatalln("organization is not set")
		} else {
			fmt.Println("Organization:", organization)
		}
		username := config.GetString("github.username")
		token := config.GetString("github.token")
		if verifyAuth(username, token) == true {
			allRepos := fetchAllRepos(token, organization)
			fmt.Println("Total Repositories: ", len(allRepos))
			createDir(organization)
			var wg sync.WaitGroup
			wg.Add(len(allRepos))
			cloneErrors, pullErrors := make(chan error, len(allRepos)), make(chan error, len(allRepos))
			for i := range allRepos {
				go func(repo *github.Repository) {
					defer wg.Done()
					path := "./" + organization + "/" + repo.GetName()
					go cloneRepo(username, token, path, repo, cloneErrors)
					cloneError := <-cloneErrors
					if cloneError != nil {
						go pullRepo(username, token, path, pullErrors)
						pullError := <-pullErrors
						if pullError != nil {
							fmt.Printf("- %s : repository %s\n", repo.GetCloneURL(), pullError)
						} else {
							fmt.Printf("+ %s : pulled latest commit successfully\n", repo.GetCloneURL())
						}
					} else {
						fmt.Printf("+ %s : cloned successfully\n", repo.GetCloneURL())
					}
				}(allRepos[i])
			}
			wg.Wait()
			close(cloneErrors)
			close(pullErrors)
			elapsed := time.Since(start)
			fmt.Printf("Time Taken: %s\n", elapsed)
		}
	},
}

func init() {
	cloneCmd.PersistentFlags().StringP("org", "o", config.GetString("organization"), "Organization")
	config.BindPFlag("organization", cloneCmd.PersistentFlags().Lookup("org"))
	rootCmd.AddCommand(cloneCmd)
}
