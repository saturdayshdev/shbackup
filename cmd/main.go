package main

import (
	"log"
	"os"

	"github.com/robfig/cron/v3"
	lib "github.com/saturdayshdev/shbackup/internal"
)

func main() {
	storage, err := lib.CreateStorageClient(lib.StorageClientConfig{
		BucketName:     os.Getenv("BUCKET_NAME"),
		BucketLocation: os.Getenv("BUCKET_REGION"),
		BucketClass:    os.Getenv("BUCKET_CLASS"),
		ProjectID:      os.Getenv("PROJECT_ID"),
		PrivateKeyID:   os.Getenv("PRIVATE_KEY_ID"),
		PrivateKey:     os.Getenv("PRIVATE_KEY"),
		ClientID:       os.Getenv("CLIENT_ID"),
		ClientEmail:    os.Getenv("CLIENT_EMAIL"),
	})
	if err != nil {
		panic(err)
	}

	docker, err := lib.CreateDockerClient()
	if err != nil {
		panic(err)
	}

	c := cron.New()
	c.AddFunc(os.Getenv("BACKUP_CRON"), func() {
		containers, err := docker.GetContainers()
		if err != nil {
			log.Println(err)
		}

		for _, container := range containers {
			labels := container.Labels
			if labels["shbackup.enabled"] != "true" {
				continue
			}

			config, err := lib.GetBackupConfig(labels, &container)
			if err != nil {
				log.Println(err)
			}

			err = lib.BackupDatabase(docker, storage, config)
			if err != nil {
				log.Println(err)
			}
		}
	})
	c.Start()
	select {}
}
