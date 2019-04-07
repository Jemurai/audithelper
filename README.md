# GAA = Go Away Auditor

Basically, make it quick and easy to get information to support audits.

## Running

You can use gaa to audit github, aws or google apps.  To do so, you
need to set up access.

## Common Features

1. json diff?
1. ??

## GitHub

### Access
To get a GitHub OAuth token, use these instructions:
https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line

### Command

`gaa github --github-org Jemurai`

### What you get?

You get a list of repositories with metadata for any user associated with your organization.  The _use case_ is that you want to ensure that the repos your team has, and that are public, are as intended.

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
