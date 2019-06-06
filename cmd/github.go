// Package cmd

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Get information from GitHub",
	Long:  `Pull information out of GitHub as part of a user audit.`,
	Run: func(cmd *cobra.Command, args []string) {
		public := viper.GetString("public")
		org := viper.GetString("github-org")
		user := viper.GetString("github-user")
		githubToken := viper.GetString("github-token")
		log.Debug("github called with", org, " and ", user, " and ", public)

		client := getGithubClient(githubToken)
		if org != "" {
			getOrgUsers(org, client)
			getOrgRepos(org, client)
			getOrgUserRepos(org, client)
		}
		if user != "" {
			getUserRepos(user, client)
		}
	},
}

func getUserRepos(user string, client *github.Client) {
	log.Debug("Getting", user, " repositories:")
	opt := &github.RepositoryListOptions{}
	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, user, opt)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(repos); i++ {
		visitRepo(*repos[i])
	}
}

func getOrgRepos(org string, client *github.Client) {
	log.Debug("Getting", org, " repositories:")
	//	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	opt := &github.RepositoryListByOrgOptions{}
	ctx := context.Background()
	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(repos); i++ {
		visitRepo(*repos[i])
	}
}

func getOrgUserRepos(org string, client *github.Client) {
	log.Debug("Getting", org, " user repositories:")
	opt := &github.ListMembersOptions{}
	ctx := context.Background()
	users, _, err := client.Organizations.ListMembers(ctx, org, opt)
	if err != nil {
		fmt.Println(err)
	}
	log.Debug("Handling", len(users), " users from ", org)
	for j := 0; j < len(users); j++ {
		log.Debug(*users[j].Login)
		getUserRepos(*users[j].Login, client)
	}
}

func getOrgUsers(org string, client *github.Client) {
	log.Debug("Getting", org, " users:")
	opt := &github.ListMembersOptions{}
	ctx := context.Background()
	users, _, err := client.Organizations.ListMembers(ctx, org, opt)
	if err != nil {
		fmt.Println(err)
	}
	log.Debug("Handling", len(users), " users from ", org)
	for j := 0; j < len(users); j++ {
		visitGithubUser(*users[j])
	}
}

func visitGithubUser(user github.User) {
	m := map[string]interface{}{
		"login": *user.Login,
		"url":   *user.URL,
	}
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		log.Error("error:", err)
	}
	fmt.Print(string(b))
	//	fmt.Println(*repo.Name, "\t\t\t", *repo.Private, "\t", *repo.UpdatedAt, "\t", *repo.CloneURL)
}

func visitRepo(repo github.Repository) {
	m := map[string]interface{}{"name": *repo.Name,
		"url":     *repo.CloneURL,
		"owner":   *repo.Owner.Login,
		"update":  *repo.UpdatedAt,
		"create":  *repo.CreatedAt,
		"private": *repo.Private,
	}
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		log.Error("error:", err)
	}
	fmt.Print(string(b))
	//	fmt.Println(*repo.Name, "\t\t\t", *repo.Private, "\t", *repo.UpdatedAt, "\t", *repo.CloneURL)
}

func init() {
	rootCmd.AddCommand(githubCmd)

	githubCmd.Flags().String("github-org", "", "The organization to audit.")
	githubCmd.Flags().String("github-user", "", "The user to audit.")
	githubCmd.Flags().String("github-token", "", "Github token")
	githubCmd.Flags().String("github-public", "", "Should the output only include public repos?")

	viper.BindPFlag("github-org", githubCmd.Flags().Lookup("github-org"))
	viper.BindPFlag("github-user", githubCmd.Flags().Lookup("github-user"))
	viper.BindPFlag("github-token", githubCmd.Flags().Lookup("github-token"))
	viper.BindPFlag("github-public", githubCmd.Flags().Lookup("github-public"))

}

func getGithubClient(token string) *github.Client {
	if token == "" {
		log.Info("Warning: empty token so searching public.")
		githubClient := github.NewClient(nil)
		return githubClient
	}
	// If the token is defined, get an OAuth client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauth2Client := oauth2.NewClient(context.Background(), ts)
	githubClient := github.NewClient(oauth2Client)
	return githubClient
}
