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
	"time"

	"github.com/AltoStack/dynamodump/actions"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVarP(&dynamoTableName, "dynamo-table-name", "t", "", "Name of the Dynamo table to actions. Environment variable: DYN_DYNAMO_TABLE_NAME (required)")
	backupCmd.Flags().Int64VarP(&dynamoBatchSize, "dynamo-table-batch-size", "s", 1000, "Max number of records to read from the Dynamo table at once. Environment variable: DYN_DYNAMO_TABLE_BATCH_SIZE")
	backupCmd.Flags().StringVarP(&dynamoTableRegion, "dynamo-table-region", "o", "", "AWS region of the Dynamo table. Environment variable: DYN_DYNAMO_TABLE_REGION (required)")
	backupCmd.Flags().Int64VarP(&waitTime, "dynamo-table-batch-wait-time", "w", 100, "Number of milliseconds to wait between batches. If a ProvisionedThroughputExceededException is encountered, "+
		"the script will wait twice that amount of time before retrying. Environment variable: DYN_WAIT_TIME")
	backupCmd.Flags().StringVarP(&s3BucketName, "s3-bucket-name", "b", "", "Name of the S3 bucket where to put the actions. Environment variable: DYN_S3_BUCKET_NAME (required)")
	backupCmd.Flags().StringVarP(&s3BucketRegion, "s3-bucket-region", "d", "", "AWS region of the s3 Bucket. Environment variable: DYN_S3_BUCKET_REGION (required)")
	backupCmd.Flags().StringVarP(&s3BucketFolderName, "s3-bucket-folder-name", "f", "", "Path inside the S3 bucket where to put actions. Environment variable: DYN_S3_BUCKET_FOLDER_NAME (required)")
	backupCmd.Flags().BoolVarP(&s3DateSuffix, "s3-bucket-folder-name-suffix", "p", false, "Adds an autogenerated suffix folder named using the UTC date in the format YYYY-mm-dd-HH24-MI-SS to the provided S3 folder. Environment variable: DYN_S3_BUCKET_NAME_SUFFIX")

	backupCmd.MarkFlagRequired("dynamo-table-name")
	backupCmd.MarkFlagRequired("dynamo-table-region")
	backupCmd.MarkFlagRequired("s3-bucket-name")
	backupCmd.MarkFlagRequired("s3-bucket-region")
	backupCmd.MarkFlagRequired("s3-bucket-folder-name")
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a DynamoDB Table to S3",
	Run: func(cmd *cobra.Command, args []string) {
		actions.TableBackup(dynamoTableName, dynamoBatchSize, time.Duration(waitTime)*time.Millisecond, s3BucketName, s3BucketFolderName, s3DateSuffix, dynamoTableRegion, s3BucketRegion)
	},
}
