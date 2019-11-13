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
package actions

import (
	"log"
	"os"
	"time"

	"github.com/AltoStack/dynamodump/core"
)

// Table manages the consumer from a given DynamoDB table and a producer
// to a given s3 bucket
func TableBackup(tableName string, batchSize int64, waitPeriod time.Duration, bucket, prefix string, addDate, crossRegions bool, dynamoRegion string, s3Region string) {
	if addDate {
		t := time.Now().UTC()
		prefix += "/" + t.Format("2006-01-02-15-04-05")
	}

	if crossRegions {
		if dynamoRegion == "" || s3Region == "" {
			log.Fatal("Error. Missing fields dynamoRegion or s3Region with cross regions flag enabled")
			os.Exit(-1)
		}
	} else if dynamoRegion != "" || s3Region != "" {
		log.Fatal("Error. Cross regions disabled, please either enable it or remove s3-bucket-region and dynamo-table-region flags")
		os.Exit(-1)
	}

	proc := core.NewAwsHelper(dynamoRegion)
	dest := core.NewAwsHelper(s3Region)

	go proc.ChannelToS3(bucket, prefix, 10*1024*1024, dest)

	err := proc.TableToChannel(tableName, batchSize, waitPeriod)
	if err != nil {
		log.Fatal(err.Error())
	}

	proc.Wg.Wait()
}
