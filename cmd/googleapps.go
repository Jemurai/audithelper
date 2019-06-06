// Package cmd

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

// googleappsCmd represents the googleapps command
var googleappsCmd = &cobra.Command{
	Use:   "googleapps",
	Short: "Audit google apps",
	Long:  `Keep track of google apps users.`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		// If modifying these scopes, delete your previously saved token.json.
		config, err := google.ConfigFromJSON(b, admin.AdminDirectoryUserReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		client := getClient(config)

		srv, err := admin.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve directory Client %v", err)
		}

		r, err := srv.Users.List().Customer("my_customer").MaxResults(100).
			OrderBy("email").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve users in domain: %v", err)
		}

		if len(r.Users) == 0 {
			log.Debugln("No users found.")
		} else {
			log.Debugln("Users:")
			for _, u := range r.Users {
				m := map[string]interface{}{
					"email":     u.PrimaryEmail,
					"name":      u.Name.FullName,
					"lastLogin": u.LastLoginTime,
					"admin":     u.IsAdmin,
					"suspended": u.Suspended,
				}
				b, err := json.MarshalIndent(m, "", " ")
				if err != nil {
					log.Error("error:", err)
				}
				fmt.Print(string(b))

				// fmt.Printf("%s (%s) Admin? %v Suspended? %v Last Login: %v\n", u.PrimaryEmail, u.Name.FullName, u.IsAdmin, u.Suspended, u.LastLoginTime)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(googleappsCmd)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
