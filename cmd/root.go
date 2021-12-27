package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var config = viper.New()

var rootCmd = &cobra.Command{
	Use:   "git-mass",
	Short: "Download all repositories of an organization",
	Long:  `A simple tool to download all repositories(public & private) of an organization`,
	Run: func(cmd *cobra.Command, args []string) {
		username := config.GetString("github.username")
		token := config.GetString("github.token")
		if username != "" && token != "" {
			fmt.Println("username:", username)
			fmt.Println("password:", token)
			log.Fatalln("Must use sub command")
		} else {
			log.Fatalln("Credentials not found. Use the 'config' subcommand to set it")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	config.AddConfigPath(home)
	config.SetConfigName("/.git-mass")
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err == nil {
		return
		// log.Println("Using config file:", config.ConfigFileUsed())
	}
}
