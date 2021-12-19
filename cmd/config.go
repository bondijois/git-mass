package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

type credential struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type credentials struct {
	credential `json:"github"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure credentials",
	Long: `Creates a .cred file in the current directory with your github credentials. 
Make sure you delete this file once done, or set token expiry on github`,
	Run: func(cmd *cobra.Command, args []string) {
		storedUname := config.GetString("github.username")
		storedToken := config.GetString("github.token")

		username := config.GetString("username")
		token := config.GetString("token")
		if username != "" && token != "" {
			fmt.Println("Verifying credentials")
		} else if storedUname == "" || storedToken == "" {
			log.Fatalln("Credentials not found. Use -u and -t flags to set username and token")
		} else {
			fmt.Printf("Stored credentials found at: %s\n", config.ConfigFileUsed())
		}
		if username != "" && token != "" {
			data := credentials{
				credential{
					Username: username,
					Token:    token,
				},
			}
			if verifyAuth(username, token) == true {
				fmt.Println("Creating/Updating credential file")
				credFile, e := os.Create(".creds")
				if e != nil {
					log.Fatalln(e)
				}
				credErr := credFile.Close()
				if credErr != nil {
					return
				}
				file, _ := json.MarshalIndent(data, "", " ")
				_ = ioutil.WriteFile(".creds", file, 0644)
			}
		} else {
			if verifyAuth(storedUname, storedToken) == true {
				return
			}
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringP("username", "u", config.GetString("username"), "Username")
	configCmd.PersistentFlags().StringP("token", "t", config.GetString("token"), "Token")
	config.BindPFlag("username", configCmd.PersistentFlags().Lookup("username"))
	config.BindPFlag("token", configCmd.PersistentFlags().Lookup("token"))
	rootCmd.AddCommand(configCmd)
}
