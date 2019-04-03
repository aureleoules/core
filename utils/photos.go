package utils

import (
	"context"
	"io"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// UpdateFilename : update filename of a file
func UpdateFilename(fileID bson.ObjectId, filename string) error {
	config := GetConfig()
	bucketName := config.BucketName

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

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

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

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

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

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
