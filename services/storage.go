package services

import (
	"context"
	// "log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetStorageClient() (*minio.Client, string) {
	// ctx := context.Background()
	endpoint := "localhost:9090"
	accessKeyID := "Oc0Dy1XcUCjTkzKTyx9N"
	secretAccessKey := "h0ocsyImxYFS86jjRCBZrzjzEMFatg57fia0uInW"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic("Failed to connect to storage")
	}
	_, err =  minioClient.HealthCheck(3)
	if err != nil {
		panic("Failed HealthCheck")
	}

	bucketName := "mymusic"
	// location := "us-east-1"

	// err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	// if err != nil {
	// 	// Check to see if we already own this bucket (which happens if you run this twice)
	// 	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		log.Printf("We already own %s\n", bucketName)
	// 	} else {
	// 		log.Fatalln(err)
	// 	}
	// } else {
	// 	log.Printf("Successfully created %s\n", bucketName)
	// }

	return minioClient, bucketName
}

func UploadFileStorage(file *multipart.FileHeader, buffer multipart.File) (minio.UploadInfo, error) {
	client, bucket := GetStorageClient()

	objectName := file.Filename
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	ctx := context.Background()
	info, err := client.PutObject(
		ctx,
		bucket,
		objectName,
		fileBuffer,
		fileSize,
		minio.PutObjectOptions{ContentType: contentType},
	)

	return info, err
}
