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
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/segmentio/ksuid"
)

// S3ManifestEntry represents an entry in the actions manifest stored in the s3 folder of the actions
type S3ManifestEntry struct {
	URL       string `json:"url"`
	Mandatory bool   `json:"mandatory"`
}

// S3Manifest represents the actions manifest stored in the s3 folder of the actions
type S3Manifest struct {
	Name    string            `json:"name"`
	Version int               `json:"version"`
	Entries []S3ManifestEntry `json:"entries"`
}

// genNewFileName returns a UUID used by the data pipelines
func genNewFileName() string {
	uuID := hex.EncodeToString(ksuid.New().Payload())
	return fmt.Sprintf("%s-%s-%s-%s-%s", uuID[:8], uuID[8:12], uuID[12:16], uuID[16:20], uuID[20:])
}

// LoadManifestFromS3 downloads the given manifest file and load it in the
// ManifestS3 attribute of the struct
func (h *AwsHelper) LoadManifestFromS3(bucketName, manifestPath string) error {
	doc, err := h.GetFromS3(bucketName, manifestPath)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == s3.ErrCodeNoSuchKey {
			log.Fatalf("[ERROR] Unable to find a manifest flag in the provided folder. Are you sure the actions was successful?\nAborting...\n")
		}
		log.Fatalf("[ERROR] Unable to retrieve the manifest flag information: %s\nAborting...\n", err)
	}
	defer (*doc).Close()
	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, *doc); err != nil {
		return err
	}

	return json.Unmarshal(buff.Bytes(), &h.ManifestS3)
}

// GetFromS3 download a file from s3 to memory (as the files are small by
// default - just a few Mb).
func (h *AwsHelper) GetFromS3(bucketName, s3Path string) (*io.ReadCloser, error) {
	svc := h.CreateServiceClientValue()
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Path),
	}

	results, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}
	return &results.Body, nil
}

// ExistsInS3 checks that a given path in s3 exists as a file
func (h *AwsHelper) ExistsInS3(bucketName, s3Path string) (bool, error) {
	svc := h.CreateServiceClientValue()
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Path),
	}

	_, err := svc.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UploadToS3 writes the content of a bytes array to the given s3 path
func (h *AwsHelper) UploadToS3(bucketName, s3Key string, data []byte) {
	svc := h.CreateServiceClientValue()
	uploader := s3manager.NewUploaderWithClient(svc)

	upParams := &s3manager.UploadInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(s3Key),
		Body:                 bytes.NewReader(data),
		StorageClass:         aws.String("STANDARD_IA"),
		ServerSideEncryption: aws.String("AES256"),
		// Expire: ...,
		// Tagging:
	}
	// Set file name and content before upload
	log.Printf("Writing file: s3://%s/%s\n", *upParams.Bucket, *upParams.Key)
	_, err := uploader.Upload(upParams)
	if err != nil {
		log.Fatalf("[ERROR] while uploading to s3: %s", err)
	}
}

// ReaderToChannel reads the data from a actions line by line, serializes it and
// sends it to the struct's channel
func (h *AwsHelper) ReaderToChannel(dataReader *io.ReadCloser) error {
	defer (*dataReader).Close()
	scanner := bufio.NewScanner(*dataReader)
	for scanner.Scan() {
		res := map[string]*dynamodb.AttributeValue{}
		data := scanner.Bytes()
		json.Unmarshal(data[:], &res)
		h.DataPipe <- res
	}
	return scanner.Err()
}

// S3ToDynamo pulls the s3 files from AwsHelper.ManifestS3 and import them
// inside the given table using the given batch size (and wait period between
// each batch)
func (h *AwsHelper) S3ToDynamo(tableName string, batchSize int64, waitPeriod time.Duration, destination *AwsHelper) error {
	var err error
	go h.ChannelToTable(tableName, batchSize, waitPeriod, destination)
	h.Wg.Add(1)
	for _, entry := range h.ManifestS3.Entries {
		u, _ := url.Parse(entry.URL)
		if u.Scheme == "s3" {
			data, err := h.GetFromS3(u.Host, u.Path)
			if err != nil {
				return err
			}
			if err = h.ReaderToChannel(data); err != nil {
				break
			}
		}

	}
	close(h.DataPipe)
	h.Wg.Wait()
	return err
}

// DumpBuffer dumps the content of the given buffer to a new randomly generated
// file name in the given s3 path in the given bucket and resets the said buffer
func (h *AwsHelper) DumpBuffer(bucketName, s3Folder string, buff *bytes.Buffer) {
	filePath := fmt.Sprintf("%s/%s", s3Folder, genNewFileName())
	h.UploadToS3(bucketName, filePath, buff.Bytes())
	h.ManifestS3.Entries = append(h.ManifestS3.Entries, S3ManifestEntry{URL: fmt.Sprintf("s3://%s/%s", bucketName, filePath), Mandatory: true})
	buff.Reset()
}

// ChannelToS3 reads from the given channel and sends the data the given bucket
// in files of s3BufferSize max size
func (h *AwsHelper) ChannelToS3(bucketName, s3Folder string, s3BufferSize int, destination *AwsHelper) {
	defer h.Wg.Done()
	// buff is the buffer where the data will be stored while before being sent to s3
	var buff bytes.Buffer
	destination.ManifestS3 = S3Manifest{Version: 3, Name: "DynamoDB-export"}

	for elem := range h.DataPipe {
		data, err := MarshalDynamoAttributeMap(elem)
		if err != nil {
			log.Fatalf("[ERROR] while converting to json: %v\nError: %s\n", elem, err)
		}

		// before overflowing the buffer, dump to s3 and empty it
		if buff.Len()+len(data) >= s3BufferSize && buff.Len() > 0 {
			destination.DumpBuffer(bucketName, s3Folder, &buff)
		}
		// add the data to the buffer
		buff.Write(data)
		buff.WriteString("\n")
	}

	// Upload the rest of the buffer
	destination.DumpBuffer(bucketName, s3Folder, &buff)
	// Signal the success of the actions
	destination.UploadToS3(bucketName, fmt.Sprintf("%s/_SUCCESS", s3Folder), []byte{})
	// Wrap up the manifest of the actions files
	manifestData, err := json.Marshal(destination.ManifestS3)
	if err != nil {
		log.Fatalf("[ERROR] while doing a marshal on the manifest: %v\nError: %s\n", destination.ManifestS3, err)
	}
	destination.UploadToS3(bucketName, fmt.Sprintf("%s/manifest", s3Folder), manifestData)
}

// Check if credentials has been initialised and return a Service Client Value
func (h *AwsHelper) CreateServiceClientValue() s3iface.S3API {
	if h.RoleCreds == (credentials.Credentials{}) {
		return s3.New(h.AwsSession)
	} else {
		return s3.New(h.AwsSession, &aws.Config{Credentials: &h.RoleCreds})
	}
}
