package utils

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/backpulse/core/models"
	"google.golang.org/api/option"
	"gopkg.in/mgo.v2/bson"
)

// GetGoogleCloudClient : Return google cloud client
func GetGoogleCloudClient(ctx context.Context) (*storage.Client, error) {

	// Use ./google_credentials.json by default
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Println("Local file credentials")
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./google_credentials.json")
		client, err := storage.NewClient(ctx)
		return client, err
	}

	// Use env variables
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))))
	return client, err
}

// UpdateFilename : update filename of a file
func UpdateFilename(fileID bson.ObjectId, filename string) error {
	config := GetConfig()
	bucketName := config.BucketName

	ctx := context.Background()

	client, err := GetGoogleCloudClient(ctx)

	bucket := client.Bucket(bucketName)
	object := bucket.Object(fileID.Hex())

	log.Print(fileID.Hex())
	_, err = object.Update(ctx, storage.ObjectAttrsToUpdate{
		ContentDisposition: "inline; filename=\"" + filename + "\"",
	})
	log.Print(err)
	return err
}

// UploadFile : upload file to google cloud
func UploadFile(file multipart.File, fileName string) (bson.ObjectId, error) {

	config := GetConfig()
	bucketName := config.BucketName

	ctx := context.Background()
	client, err := GetGoogleCloudClient(ctx)

	objectID := bson.NewObjectId()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectID.Hex())

	object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)

	wc := object.NewWriter(ctx)
	_, err = io.Copy(wc, file)
	if err != nil {
		return "", err
	}

	err = wc.Close()
	if err != nil {
		return "", err
	}
	update, err := object.Update(ctx, storage.ObjectAttrsToUpdate{
		ContentDisposition: "inline; filename=\"" + fileName + "\"",
	})
	log.Println("update", update)
	log.Println("error", err)

	return objectID, nil
}

// RemoveGoogleCloudPhotos : remove photos from google cloud to save space
func RemoveGoogleCloudPhotos(photos []models.Photo) error {
	config := GetConfig()
	bucketName := config.BucketName

	ctx := context.Background()

	client, _ := GetGoogleCloudClient(ctx)

	bucket := client.Bucket(bucketName)
	for _, photo := range photos {
		log.Println(photo.ID.Hex())
		object := bucket.Object(photo.ID.Hex())
		err := object.Delete(ctx)
		if err != nil {
			log.Println(err)

			return err
		}
	}
	return nil
}
