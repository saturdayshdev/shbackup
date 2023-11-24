package lib

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type StorageClient struct {
	Client    *storage.Client
	Bucket    *storage.BucketHandle
	ProjectID string
}

type StorageClientConfig struct {
	BucketName     string
	BucketLocation string
	BucketClass    string
	ProjectID      string
	PrivateKeyID   string
	PrivateKey     string
	ClientID       string
	ClientEmail    string
}

func (c *StorageClientConfig) CreateCredentialsJSON() ([]byte, error) {
	credentials := struct {
		Type                    string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
	}{
		Type:                    "service_account",
		ProjectID:               c.ProjectID,
		PrivateKeyID:            c.PrivateKeyID,
		PrivateKey:              c.PrivateKey,
		ClientEmail:             c.ClientEmail,
		ClientID:                c.ClientID,
		AuthURI:                 "https://accounts.google.com/o/oauth2/auth",
		TokenURI:                "https://oauth2.googleapis.com/token",
		AuthProviderX509CertURL: "https://www.googleapis.com/oauth2/v1/certs",
		ClientX509CertURL:       "https://www.googleapis.com/robot/v1/metadata/x509/" + c.ClientEmail,
	}

	json, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (c *StorageClient) CreateBucket(name string, location string, class string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	bucket := c.Client.Bucket(name)

	_, err := bucket.Attrs(ctx)
	if err == nil {
		return bucket, nil
	}

	options := &storage.BucketAttrs{Location: location, StorageClass: class}
	err = bucket.Create(ctx, c.ProjectID, options)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (c *StorageClient) ObjectExists(name string) bool {
	_, err := c.Bucket.Object(name).Attrs(context.Background())
	return err == nil
}

func (c *StorageClient) UploadFile(name string, path string, errCh chan error) {
	if c.ObjectExists(name) {
		errCh <- errors.New("object already exists")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		errCh <- err
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		errCh <- err
		return
	}

	buffer := make([]byte, fileInfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		errCh <- err
		return
	}

	go func() {
		writer := c.Bucket.Object(name).NewWriter(context.Background())

		_, err := writer.Write(buffer)
		if err != nil {
			errCh <- err
			return
		}

		err = writer.Close()
		if err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}()
}

func CreateStorageClient(config StorageClientConfig) (*StorageClient, error) {
	credentials, err := config.CreateCredentialsJSON()
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, err
	}

	storageClient := &StorageClient{client, nil, config.ProjectID}
	bucket, err := storageClient.CreateBucket(config.BucketName, config.BucketLocation, config.BucketClass)
	if err != nil {
		return nil, err
	}

	storageClient.Bucket = bucket
	return storageClient, nil
}
