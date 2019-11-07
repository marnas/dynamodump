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
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringVarP(&dynamoTableName, "dynamo-table-name", "t", "", "Name of the Dynamo table to actions. Environment variable: DYN_DYNAMO_TABLE_NAME")
	restoreCmd.Flags().Int64VarP(&dynamoBatchSize, "dynamo-table-batch-size", "s", 1000, "Max number of records to read from the Dynamo table at once. Environment variable: DYN_DYNAMO_TABLE_BATCH_SIZE")
	restoreCmd.Flags().BoolVarP(&dynamoAppendRestore, "dynamo-append-restore", "a", false, "Appends the rows to a non-empty table when restoring instead of aborting. Environment variable: DYN_DYNAMO_RESTORE_APPEND")
	restoreCmd.Flags().Int64VarP(&waitTime, "dynamo-table-batch-wait-time", "w", 100, "Number of milliseconds to wait between batches. If a ProvisionedThroughputExceededException is encountered, "+
		"the script will wait twice that amount of time before retrying. Environment variable: DYN_WAIT_TIME")
	restoreCmd.Flags().StringVarP(&s3BucketName, "s3-bucket-name", "b", "", "Name of the S3 bucket where to put the actions. Environment variable: DYN_S3_BUCKET_NAME")
	restoreCmd.Flags().StringVarP(&s3BucketFolderName, "s3-bucket-folder-name", "f", "", "Path inside the S3 bucket where to put actions. Environment variable: DYN_S3_BUCKET_FOLDER_NAME")
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a DynamoDB Table from S3",
	Long:  `WIP`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.TableRestore(dynamoTableName, dynamoBatchSize, time.Duration(waitTime)*time.Millisecond, s3BucketName, s3BucketFolderName, dynamoAppendRestore)
	},
}
