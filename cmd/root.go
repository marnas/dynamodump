/*
Copyright © 2019 AltoStack <info@altostack.io>

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
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	dynamoTableAccountID string
	dynamoTableName      string
	dynamoBatchSize      int64
	dynamoAppendRestore  bool
	dynamoTableRegion    string
	forceRestore         bool
	roleAssumed          string
	s3BucketAccountID    string
	s3BucketName         string
	s3BucketFolderName   string
	s3BucketRegion       string
	s3DateSuffix         bool
	waitTime             int64
)

var rootCmd = &cobra.Command{
	Use:   "dynamodump",
	Short: "AWS DynamoDB Backup and Restores",
	Long: `
		Dynamodump allows for easier and cheaper actions and restores of DynamoDB Tables.
		Backups are compatible with the AWS DataPipeline functionality.
				
		It is also capable of restoring a actions from s3 to a given table both from this
		tool or from a actions generated using the AWS DataPipeline functionality.
		to quickly create a Cobra application.
  `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Defining persistent flags, which
	// will be global for the application.
	rootCmd.PersistentFlags().StringVarP(&dynamoTableName, "dynamo-table-name", "t", "", "Name of the Dynamo table to actions. Environment variable: DYN_DYNAMO_TABLE_NAME (required)")
	rootCmd.PersistentFlags().Int64VarP(&dynamoBatchSize, "dynamo-table-batch-size", "s", 1000, "Max number of records to read from the Dynamo table at once. Environment variable: DYN_DYNAMO_TABLE_BATCH_SIZE")
	rootCmd.PersistentFlags().StringVarP(&dynamoTableAccountID, "dynamo-table-account-id", "x", "", "AccountID that will be used to access the dynamoDB")
	rootCmd.PersistentFlags().StringVarP(&dynamoTableRegion, "dynamo-table-region", "o", "", "AWS region of the Dynamo table. Environment variable: DYN_DYNAMO_TABLE_REGION (required)")
	rootCmd.PersistentFlags().Int64VarP(&waitTime, "dynamo-table-batch-wait-time", "w", 100, "Number of milliseconds to wait between batches. If a ProvisionedThroughputExceededException is encountered, the script will wait twice that amount of time before retrying. Environment variable: DYN_WAIT_TIME")
	rootCmd.PersistentFlags().StringVarP(&roleAssumed, "assume-role", "g", "OrganizationAccountAccessRole", "Role that will be used to access the s3 Bucket")
	rootCmd.PersistentFlags().StringVarP(&s3BucketAccountID, "s3-bucket-account-id", "e", "", "AccountID that will be used to access the s3 Bucket")
	rootCmd.PersistentFlags().StringVarP(&s3BucketName, "s3-bucket-name", "b", "", "Name of the S3 bucket where to put the actions. Environment variable: DYN_S3_BUCKET_NAME (required)")
	rootCmd.PersistentFlags().StringVarP(&s3BucketFolderName, "s3-bucket-folder-name", "f", "", "Path inside the S3 bucket where to put actions. Environment variable: DYN_S3_BUCKET_FOLDER_NAME (required)")
	rootCmd.PersistentFlags().StringVarP(&s3BucketRegion, "s3-bucket-region", "d", "", "AWS region of the s3 Bucket. Environment variable: DYN_S3_BUCKET_REGION (required)")

	rootCmd.MarkPersistentFlagRequired("dynamo-table-name")
	rootCmd.MarkPersistentFlagRequired("dynamo-table-region")
	rootCmd.MarkPersistentFlagRequired("s3-bucket-name")
	rootCmd.MarkPersistentFlagRequired("s3-bucket-region")
	rootCmd.MarkPersistentFlagRequired("s3-bucket-folder-name")
}
