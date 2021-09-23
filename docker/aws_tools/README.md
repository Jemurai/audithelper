# AWS Tools Container

This directory aims to build a docker container with all the standard AWS
tools I use for AWS assessments. Here is a rundown of which tools are included.

1. [ScoutSuite](https://github.com/nccgroup/ScoutSuite): Excellent overall
   configuration evaluation tool which is the basis of the [JASP](https://jasp.cloud)
   ranger worker.
1. [Cloudsplaining](https://github.com/salesforce/cloudsplaining/): Analyzes
   IAM policies, users, groups, and roles to identify service and resource
   wildcards and other risks based on IAM configuration.
1. [Prowler](https://github.com/toniblyx/prowler): Another good tool for
   checking for configuration issues in an AWS account. This one is shell based
   but is quite thorough and has coverage that some of the other tools do not.
1. [aws-shell](https://github.com/awslabs/aws-shell): The best AWS CLi which
   provides autocomplete, history, and built-in man page browser as well as
   example commands.
1. [Cartography](https://github.com/lyft/cartography): Generates a graph
   database to allow for some advanced authorization scenarios. NB: this is not
   yet in a working state with the current container configuration.

## How to use it

The easiest way to get up and running is to use [aws-vault](https://github.com/99designs/aws-vault)
and [docker-compose](https://docs.docker.com/compose/). The image can also be
run directly, but using compose you can eliminate the need to enumerate the
necessary environment variables on the command-line.

```sh
## Build the containers (Only necessary to do this once or when changed)
docker-compose build

## Run the tools via aws-vault
aws-vault exec SuperAwesomeAuditProfile -- docker-compose run --rm tools
```

This will drop you at a command-line in the container with your AWS credentials
available, all tools on the PATH, and a volume mapping from `/work` in the
container to a `work` subdirectory where you launched the container.

Run all the tools you need to run and exit the shell when finished. Docker will
automatically shutdown and destroy the container. When you are completely
finished, tear down everything by running `docker-compose down`.
