## shBackup

Service for automatic backups of databases to Google Cloud.

### Usage

If you want a container to be backed-up by shBackup, add the following labels:

```yaml
labels:
  - shbackup.enabled=true
  - shbackup.name=PROJECT_NAME
  - shbackup.strategy=STRATEGY_NAME
  - shbackup.user=USERNAME
  - shbackup.password=PASSWORD
```

Currenty, shBackup supports two strategies: `postgres` and `mysql`. More may come in the future.

You can also run shBackup in a container:

```yaml
services:
  app:
    container_name: shbackup
    image: ghcr.io/saturdayshdev/shbackup:latest
    environment:
      - DOCKER_API_VERSION=1.42
      - BACKUP_CRON=${BACKUP_CRON}
      - BUCKET_NAME=${BUCKET_NAME}
      - BUCKET_REGION=${BUCKET_REGION}
      - BUCKET_CLASS=${BUCKET_CLASS}
      - PROJECT_ID=${PROJECT_ID}
      - PRIVATE_KEY_ID=${PRIVATE_KEY_ID}
      - PRIVATE_KEY=${PRIVATE_KEY}
      - CLIENT_ID=${CLIENT_ID}
      - CLIENT_EMAIL=${CLIENT_EMAIL}
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
```

### Environment

You can see a list of all environment variables below.

| Variable       | Description                                   |
| -------------- | --------------------------------------------- |
| BACKUP_CRON    | Cron expression used by the service.          |
| BUCKET_NAME    | Name of the Cloud Storage bucket.             |
| BUCKET_REGION  | Region of the Cloud Storage bucket.           |
| BUCKET_CLASS   | Class of the Cloud Storage bucket.            |
| PROJECT_ID     | ID of the Google Cloud project.               |
| PRIVATE_KEY_ID | ID of the private key of the Service Account. |
| PRIVATE_KEY    | Private key of the Service Account.           |
| CLIENT_ID      | Client ID of the Service Account.             |
| CLIENT_EMAIL   | Client Email of the Service Account.          |

### Team

- Oskar Wójcikiewicz <oskar@saturdaysheroes.dev>
