# dynamodump: DynamoDB Backups and Restores

## Table of Contents

  * [Background](#background)
  * [What is it?](#what-is-it)
  * [Why creating this tool?](#why-creating-this-tool)
  * [How to use it?](#how-to-use-it)
  * [Contributing to the project](#contributing-to-the-project)
  
## Background

This is a fork of [dynamodbdump](https://github.com/VEVO/dynamodbdump) that aims to complete some of 
it's [TODOs](https://github.com/VEVO/dynamodbdump/blob/master/TODO.md).

## What is it?

This tool performs a backup of a given DynamoDB table and pushes it to a given folder in s3
in a format compatible with the AWS DataPipeline functionality.

It is also capable of restoring a backup from s3 to a given table both from
this tool or from a backup generated using the AWS DataPipeline functionality.

## Why create this tool?

Using the AWS DataPipelines to backup DynamoDB tables spawns EMR clusters which
can take some time, and for small tables it will cost you 20min of EMR runs for
just a few seconds of backup time, which makes no sense.

This tool can be run in a command-line, in a docker container and ending up on a
Kubernetes cron job very easily, allowing you to leverage your existing
architecture without additional costs.

## How to use it?

🏗 WIP 🚧
