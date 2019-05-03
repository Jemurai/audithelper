# GAA = Go Away Auditor

Basically, make it quick and easy to get information to support audits.

## Getting GAA

You can build gaa just by cloning, installing dependencies and running `go build`.  If you want to just run from source, you can just clone then run `go run main.go <platform>` which is fine for some folks.

If you want to get a prebuilt release version, you can get it from here for your platform.
https://github.com/Jemurai/gaa/releases

## Running

You can use gaa to audit github, aws or google apps.  To do so, you
need to set up access.  The following sections show how to set up access and run for each different platform.  In principle, it is just `gaa <platform>`.

## GitHub

### Access
To get a GitHub OAuth token, use these instructions:
https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line

Once you have a token, you can put it in a .gaa.yaml file in your home directory.  _Note that the github token should be treated as a secret and handled accordingly._

In other words, your ~/.gaa.yaml file might look like this:

```
github-token: b4a9b....
github-org: Jemurai
```

### Command

`gaa github --github-org Jemurai`

### What you get?

You get a list of repositories with metadata for any user associated with your organization.  The _use case_ is that you want to ensure that the repos your team has, and that are public, are as intended.

The idea would be that you cross check the users with your organizational user list and then make sure the repos have the correct visibility.

## AWS

### Access

We recommend using the excellent [aws-vault](https://github.com/99designs/aws-vault) library from 99 Designs to run any AWS tasks.

Based on a combination of aws-vault and ~/.aws/config profiles, when we run with the AWS command shown below, the process takes all of the information from the environment and we don't need to pass further information.

Generally, we are reading out of the AWS account so you'll want to run with ReadOnly or SecurityAudit privileges.

See this documentation on how to set up STS assume role:  
https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-role.html

### Command

`aws-vault exec jemurai-mkonda -- gaa aws`

### What you get?

What gaa does with AWS is:

1. List users and basic information to be able to see change over time.

## Google Apps

### Access

Based on [these instructions](https://developers.google.com/admin-sdk/directory/v1/quickstart/go) we need to:

1. Create client credentials:  click on enable Directory API, then Download Client Configuration and place that in a file (credentials.json) in the directory you plan to run gaa from.

2. When initially running gaa, a browser window will launch.  Click through the web prompt to allow google to issue you an OAuth2 token.

_Note that the credentials.json file should be treated as a secret and handled accordingly._

## Microsoft O365

_TODO THIS IS NOT EVEN STARTED_

https://github.com/mhoc/msgoraph
https://github.com/Azure/azure-sdk-for-go

