/*
Copyright © 2020 AltoStack <info@altostack.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package core

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// AwsHelper supports a set of helpers around DynamoDB and s3
type AwsHelper struct {
	AwsSession client.ConfigProvider
	DynamoSvc  dynamodbiface.DynamoDBAPI
	Wg         sync.WaitGroup
	DataPipe   chan map[string]*dynamodb.AttributeValue
	ManifestS3 S3Manifest
	RoleCreds  credentials.Credentials
}

// Creates a new AwsHelper, initializing an AWS session and a few
// objects like a channel or a DynamoDB client
func NewSession(region, accountID, accountRole string) *AwsHelper {
	awsSess, err := session.NewSessionWithOptions(session.Options{
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(region),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		log.Fatal(err)
	}

	dataPipe := make(chan map[string]*dynamodb.AttributeValue)

	var dynamoSvc dynamodbiface.DynamoDBAPI
	var creds credentials.Credentials

	if accountID != "" {
		arn := "arn:aws:iam::" + accountID + ":role/" + accountRole
		creds = *stscreds.NewCredentials(awsSess, arn)
		dynamoSvc = dynamodb.New(awsSess, &aws.Config{Credentials: &creds})
	} else {
		dynamoSvc = dynamodb.New(awsSess)
	}
	return &AwsHelper{AwsSession: awsSess, DataPipe: dataPipe, DynamoSvc: dynamoSvc, RoleCreds: creds}
}
