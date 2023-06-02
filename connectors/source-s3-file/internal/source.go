// Copyright 2023 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.Source = &s3FileSource{}

func NewExampleSource() cdkgo.Source {
	return &s3FileSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type s3FileSource struct {
	config *exampleConfig
	events chan *cdkgo.Tuple
	number int
}

func (s *s3FileSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*exampleConfig)
	go s.loopProduceEvent()
	return nil
}

func (s *s3FileSource) Name() string {
	return "s3FileSource"
}

func (s *s3FileSource) Destroy() error {
	return nil
}

func (s *s3FileSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *s3FileSource) loopProduceEvent() {
	// Set up an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("AKIA4IXKMXC7HGMO24GN", "5AqIo6Iuq2LF7aEc8pq9ldav6XfCfOAQD1XIVyxX", ""),
	})
	if err != nil {
		log.Println("Failed to create AWS session:", err)
		return
	}

	// Set up an S3 client
	svc := s3.New(sess)

	// Set up the S3 bucket name
	bucketName := "dj-new-test"

	// List existing files in the bucket
	existingFiles, err := listFiles(svc, bucketName)
	if err != nil {
		log.Println("Failed to list existing files:", err)
		return
	}

	for _, file := range existingFiles {
		log.Println("File Name:", file.FileName)
		log.Println("File Content:", string(file.FileContent))

		// Create a CloudEvent for the file
		event := ce.NewEvent()
		event.SetSource("s3FileSource")
		event.SetType("fileUploaded")
		event.SetData(ce.ApplicationJSON, file)

		// Send the event through the events channel
		b, _ := json.Marshal(event)
		success := func() {
			fmt.Println("send event success: " + string(b))
		}
		failed := func(err error) {
			fmt.Println("send event failed: " + string(b) + ", error: " + err.Error())
		}
		s.events <- cdkgo.NewTuple(&event, success, failed)
	}

	// Start a goroutine to check for new file uploads periodically
	go func() {
		for {
			newFiles, err := checkForNewFiles(svc, bucketName)
			if err != nil {
				log.Println("Failed to check for new files:", err)
			}
			for _, file := range newFiles {
				log.Println("New file uploaded:", file.FileName)
				log.Println("File Content:", string(file.FileContent))

				// Create a CloudEvent for the new file
				event := ce.NewEvent()
				event.SetSource("s3FileSource")
				event.SetType("fileUploaded")
				event.SetData(ce.ApplicationJSON, file)

				// Send the event through the events channel
				b, _ := json.Marshal(event)
				success := func() {
					fmt.Println("send event success: " + string(b))
				}
				failed := func(err error) {
					fmt.Println("send event failed: " + string(b) + ", error: " + err.Error())
				}
				s.events <- cdkgo.NewTuple(&event, success, failed)
			}
			time.Sleep(10 * time.Second) // Adjust the interval as needed
		}
	}()
}

// File represents a file in the S3 bucket
type File struct {
	FileName    string
	FileContent []byte
}

// List existing files in the S3 bucket and retrieve their contents as byte slices
func listFiles(svc *s3.S3, bucketName string) ([]File, error) {
	var files []File

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Existing Files:")
	for _, item := range resp.Contents {
		file := File{
			FileName: *item.Key,
		}

		// Download the file content
		objOutput, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    item.Key,
		})
		if err != nil {
			log.Println("Failed to download file:", err)
			continue
		}

		// Read the file content into a byte slice
		fileContent, err := ioutil.ReadAll(objOutput.Body)
		if err != nil {
			log.Println("Failed to read file content:", err)
			continue
		}

		file.FileContent = fileContent
		files = append(files, file)
	}

	return files, nil
}

// Check for new file uploads in the S3 bucket and retrieve their contents as byte slices
func checkForNewFiles(svc *s3.S3, bucketName string) ([]File, error) {
	var files []File

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Checking for New Files...")
	for _, item := range resp.Contents {
		// Process only new files (ignore existing files)
		if item.LastModified.After(time.Now().Add(-10 * time.Second)) { // Adjust the duration as needed
			file := File{
				FileName: *item.Key,
			}

			// Download the file content
			objOutput, err := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    item.Key,
			})
			if err != nil {
				log.Println("Failed to download file:", err)
				continue
			}

			// Read the file content into a byte slice
			fileContent, err := ioutil.ReadAll(objOutput.Body)
			if err != nil {
				log.Println("Failed to read file content:", err)
				continue
			}

			file.FileContent = fileContent
			files = append(files, file)
		}
	}

	return files, nil
}
