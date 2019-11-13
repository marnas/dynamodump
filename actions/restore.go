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
	"fmt"
	"log"
	"time"

	"github.com/AltoStack/dynamodump/core"
)

func TableRestore(tableName string, batchSize int64, waitPeriod time.Duration, bucket, prefix string, appendToTable bool) {
	proc := core.NewAwsHelper("")

	// Check if the table exists and has data in it. If so, abort
	itemsCount, err := proc.CheckTableEmpty(tableName)
	if err != nil {
		log.Fatalf("[ERROR] Unable to retrieve the target table informations: %s\nAborting...\n", err)
	}
	switch {
	case itemsCount > 0 && !appendToTable:
		log.Fatalf("[ERROR] The target table is not empty. Aborting...\n")
	case itemsCount == -1:
		log.Fatalf("[ERROR] The target table does not exists. Aborting...\n")
	case itemsCount < -1:
		log.Fatalf("[ERROR] The target table is not in ACTIVE state, so not writable. Aborting...\n")
	}

	// Check if a file "_SUCCESS" is present in the directory
	if exists, err := proc.ExistsInS3(bucket, fmt.Sprintf("%s/_SUCCESS", prefix)); !exists {
		switch {
		case err != nil:
			log.Fatalf("[ERROR] Unable to retrieve the _SUCCESS flag information: %s\nAborting...\n", err)
		case !exists:
			log.Fatalf("[ERROR] Unable to find a _SUCCESS flag in the provided folder. Are you sure the actions was successful?\nAborting...\n")
		}
	}

	// Pull the manifest from s3 and load it to memory
	err = proc.LoadManifestFromS3(bucket, fmt.Sprintf("%s/manifest", prefix))
	if err != nil {
		log.Fatalf("[ERROR] Unable to load the manifest flag information: %s\nAborting...\n", err)
	}

	// For each file in the manifest pull the file, decode each line and add them to a batch and push them into the table (batch size, then wait and continue)
	err = proc.S3ToDynamo(tableName, batchSize, waitPeriod)
	if err != nil {
		log.Fatalf("[ERROR] Unable to import the full s3 actions to Dynamo: %s\nAborting...\n", err)
	}
}