/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upload called")
		upload(strings.Join(args, " "))
	},
}

// Probably not needed with aws sdk
func getCredentials() (string, string) {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	return awsAccessKeyID, awsSecretAccessKey
}

func uploadFile(fileName string) {
	region := "us-west-2"
	bucket := "dnemes-bundle-uploader-test"
	// key := "myFile"
	// var timeout time.Duration

	if os.Getenv("AWS_REGION") != region {
		if err := os.Setenv("AWS_REGION", region); err != nil {
			log.Fatalf("Could not set AWS region")
		}
	}

	// if err := os.Setenv("AWS_ACCESS_KEY_ID", ""); err != nil {
	// 	log.Fatalf("Could not set AWS_ACCESS_KEY_ID")
	// }

	// if err := os.Setenv("AWS_SECRET_ACCESS_KEY", ""); err != nil {
	// 	log.Fatalf("Could not set AWS_SECRET_ACCESS_KEY")
	// }

	s, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatalf("Could not create AWS session for region %v: %v", region, err)
	}

	// Open the file for use
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open file %v: %v", fileName, err)
	}
	defer f.Close()

	// Get file name and size and read the file content into a buffer
	fileInfo, _ := f.Stat()
	var size int64 = fileInfo.Size()
	var name string = fileInfo.Name()
	buffer := make([]byte, size)
	f.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	if _, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(name),
		ACL:                  aws.String("bucket-owner-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	}); err != nil {
		log.Fatalf("Failed to upload file %v: %v", fileName, err)
	}
}

func upload(fileName string) {
	// awsAccessKeyID, awsSecretAccessKey := getCredentials()
	uploadFile(fileName)

}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
