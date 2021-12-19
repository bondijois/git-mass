package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var orgCmd = &cobra.Command{
	Use:   "orgs",
	Short: "List organizations",
	Long:  `Lists all public organizations the user is part of`,
	Run: func(cmd *cobra.Command, args []string) {
		username := config.GetString("github.username")
		token := config.GetString("github.token")
		if verifyAuth(username, token) == true {
			client := getGithubClient(token)
			org, _, err := client.Organizations.List(ctx, username, nil)
			fmt.Printf("%v public organizations found\n", len(org))
			if err != nil {
				log.Fatalln(err)
			}
			for _, name := range org {
				fmt.Println("* ", name.GetLogin())
			}
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
}
