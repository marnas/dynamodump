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
package cmd

import (
	"time"

	"github.com/AltoStack/dynamodump/actions"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().BoolVarP(&dynamoAppendRestore, "dynamo-append-restore", "z", false, "Appends the rows to a non-empty table when restoring instead of aborting. Environment variable: DYN_DYNAMO_RESTORE_APPEND")
	restoreCmd.Flags().BoolVarP(&forceRestore, "force-restore", "p", false, "Force restore even if the _SUCCESS file is absent")
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a DynamoDB Table from S3",
	Run: func(cmd *cobra.Command, args []string) {
		actions.TableRestore(dynamoTableName, dynamoBatchSize, time.Duration(waitTime)*time.Millisecond, s3BucketName, s3BucketFolderName, dynamoAppendRestore, forceRestore, dynamoTableAccountID, dynamoTableRegion, roleAssumed, s3BucketAccountID, s3BucketRegion)
	},
}
