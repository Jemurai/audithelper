// Package cmd

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/reports/v1"
	drive "google.golang.org/api/drive/v3"
)

// googledriveCmd represents the googledrive command
var googledriveCmd = &cobra.Command{
	Use:   "googledrive",
	Short: "Audit google drive",
	Long:  `Keep track of google file shares.`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile("drivecredentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		// If modifying these scopes, delete your previously saved token.json.
		config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, admin.AdminReportsAuditReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		client := getDriveClient(config)

		// checkGoogleDriveFiles(client)
		checkGoogleDriveAuditTrail(client)
	},
}

func init() {
	rootCmd.AddCommand(googledriveCmd)
}

func checkGoogleDriveFiles(client *http.Client) {
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	r, err := srv.Files.List().PageSize(10).
		//Q("visibility = 'shared_externally'").
		Fields("nextPageToken, files(id, name, shared)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("File: %s (%s) Shared: %v Vis: \n", i.Name, i.Id, i.Shared)

			// Useful for debugging response payloads
			// jsn, _ := json.MarshalIndent(i, "", " ")
			// fmt.Println(string(jsn))

			//			driveFile, err := srv.Files.Get(i.Id).Do()
			//			if err != nil {
			//				fmt.Println("Failed to get file details")
			//			}
			//			fmt.Printf("File details: %s\n", driveFile.Description)
		}
	}
}

func checkGoogleDriveAuditTrail(client *http.Client) {
	srv, err := admin.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve reports Client %v", err)
	}

	// showLogins(srv)
	showShares(srv)
}

func showShares(srv *admin.Service) {
	// https://developers.google.com/admin-sdk/reports/v1/appendix/activity/drive#change_user_access
	// https://developers.google.com/admin-sdk/reports/v1/reference
	//.Filters("visibility_change='external'")
	//.Filters("owner='SPIO',visibility='shared_externally'")
	r, err := srv.Activities.List("all", "drive").EventName("change_user_access").MaxResults(1000).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve drive events to domain. %v", err)
	}

	if len(r.Items) == 0 {
		fmt.Println("No events found.")
	} else {
		fmt.Println("Events:")
		for _, a := range r.Items {
			t, err := time.Parse(time.RFC3339Nano, a.Id.Time)
			if err != nil {
				fmt.Println("Unable to parse login time.")
				// Set time to zero.
				t = time.Time{}
			}
			// Useful for debugging response payloads
			// jsn, err := json.MarshalIndent(a, "", " ")
			// fmt.Println(string(jsn))
			// fmt.Printf("%s: %s %s %s %s %s %v\n", t.Format(time.RFC822), a.Actor.Email, a.Events[0].Type, a.Events[0].Name, a.Events[0].Parameters[9].Value, a.Events[0].Parameters[10].Value, a.IpAddress)
			fmt.Printf("%s: %s File: \"%s\" %s With: %s %v\n", t.Format(time.RFC822), a.Actor.Email, a.Events[0].Parameters[9].Value, a.Events[0].Parameters[10].Value, a.Events[0].Parameters[3].Value, a.IpAddress)
		}
	}
}

func showLogins(srv *admin.Service) {
	r, err := srv.Activities.List("all", "login").MaxResults(25).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve logins to domain. %v", err)
	}

	if len(r.Items) == 0 {
		fmt.Println("No logins found.")
	} else {
		fmt.Println("Logins:")
		for _, a := range r.Items {
			t, err := time.Parse(time.RFC3339Nano, a.Id.Time)
			if err != nil {
				fmt.Println("Unable to parse login time.")
				// Set time to zero.
				t = time.Time{}
			}
			fmt.Printf("%s: %s %s\n", t.Format(time.RFC822), a.Actor.Email, a.Events[0].Name)
		}
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getDriveClient(config *oauth2.Config) *http.Client {
	// The file drivetoken.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "drivetoken.json"
	tok, err := tokenDriveFromFile(tokFile)
	if err != nil {
		tok = getDriveTokenFromWeb(config)
		saveDriveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getDriveTokenFromWeb(config *oauth2.Config) *oauth2.Token {
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
func tokenDriveFromFile(file string) (*oauth2.Token, error) {
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
func saveDriveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
