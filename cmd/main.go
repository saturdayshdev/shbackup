package main

import (
	"log"
	"os"
	"sync"

	"github.com/robfig/cron/v3"
	backup "github.com/saturdayshdev/shbackup/internal/backup"
	docker "github.com/saturdayshdev/shbackup/internal/docker"
)

func main() {
	storageClient, err := backup.CreateStorageClient(backup.StorageClientConfig{
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

	dockerClient, err := docker.CreateClient()
	if err != nil {
		panic(err)
	}

	c := cron.New()
	c.AddFunc(os.Getenv("BACKUP_CRON"), func() {
		containers, err := dockerClient.GetContainers()
		if err != nil {
			panic(err)
		}

		var wg sync.WaitGroup
		for _, container := range containers {
			wg.Add(1)
			go func(container docker.Container) {
				defer wg.Done()

				config, err := backup.GetBackupConfig(&container)
				if err != nil {
					log.Println(err)
					return
				}

				err = backup.BackupDatabase(dockerClient, storageClient, config)
				if err != nil {
					log.Println(err)
				}
			}(container)
		}
	})
	c.Start()

	select {}
}
