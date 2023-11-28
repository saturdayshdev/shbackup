package main

import (
	"log"
	"os"

	"github.com/robfig/cron/v3"
	backup "github.com/saturdayshdev/shbackup/internal/backup"
	docker "github.com/saturdayshdev/shbackup/internal/docker"
)

func main() {
	storage, err := backup.CreateStorageClient(backup.StorageClientConfig{
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

	docker, err := docker.CreateClient()
	if err != nil {
		panic(err)
	}

	c := cron.New()
	c.AddFunc(os.Getenv("BACKUP_CRON"), func() {
		containers, err := docker.GetContainers()
		if err != nil {
			log.Println(err)
			return
		}

		for _, container := range containers {
			labels := container.Labels
			if labels["shbackup.enabled"] != "true" {
				continue
			}

			config, err := backup.GetBackupConfig(&container)
			if err != nil {
				log.Println(err)
				continue
			}

			err = backup.BackupDatabase(docker, storage, config)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	})
	c.Start()

	select {}
}
