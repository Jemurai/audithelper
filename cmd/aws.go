//Package cmd

package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Get basic audit info from AWS",
	Long: `Get information about users and 
	policies from AWS to support an audit.`,
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := session.NewSession(&aws.Config{})
		if err != nil {
			log.Debug("Got error creating session:")
			log.Error(err.Error())
		}
		log.Debug(viper.GetString("AWS_REGION"))
		checkUsers(sess)
	},
}

func checkUsers(sess *session.Session) {
	svc := iam.New(sess)
	result, err := svc.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(100),
	})
	if err != nil {
		log.Error("Error", err)
		return
	}
	for i, user := range result.Users {
		if user == nil {
			continue
		}
		visitIAMUser(i, user)
		keys, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
			UserName: user.UserName,
		})
		if err != nil {
			log.Error("Error", err)
			return
		}
		visitAWSKeys(keys.AccessKeyMetadata)
	}
}

func visitAWSKeys(keys []*iam.AccessKeyMetadata) {
	for _, key := range keys {
		if key == nil {
			continue
		}
		m := map[string]interface{}{
			"user":      *key.UserName,
			"status":    *key.Status,
			"createdAt": *key.CreateDate,
			"id":        *key.AccessKeyId,
		}
		b, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			log.Error("error:", err)
		}
		fmt.Print(string(b))
	}
}

func visitIAMUser(i int, user *iam.User) {
	t := time.Time{}
	if user.PasswordLastUsed != nil {
		log.Debugf("Time is not zero")
		t = *user.PasswordLastUsed
	}

	m := map[string]interface{}{
		"user":             *user.UserName,
		"createdAt":        *user.CreateDate,
		"passwordLastUsed": t,
	}
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		log.Error("error:", err)
	}
	fmt.Print(string(b))
	//	fmt.Printf("%d user %s created %v and pass last used %v\n", i, *user.UserName, user.CreateDate, user.PasswordLastUsed)
}

func init() {
	rootCmd.AddCommand(awsCmd)
}
