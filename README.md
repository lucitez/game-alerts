# TODOS

- Add env vars and secrets to the container
- Create a database to check who has subscriptions and don't send the same email twice
- Build support for sending emails to multiple coaches
- Write tests

### Eventually
- Use terraform to provision cloud function, db, etc

# Development

Spin up the database container using docker compose:
```shell
$ docker compose up
```

Seed the database:
```shell
$ psql -h localhost -U postgres -d game_alerts -f seed.sql
```

Connect to the database:
```shell
$ psql -h localhost -U postgres -w pwd -d game_alerts
```

Start Cloud Function:
`FUNCTION_TARGET=SendGameAlerts LOCAL_ONLY=true go run ./cmd/main.go`

Curl the cloud function:
```shell
curl --location 'localhost:8080' \
--header 'Content-Type: application/json' \
--header 'ce-id: 123451234512345' \
--header 'ce-specversion: 1.0' \
--header 'ce-time: 2020-01-02T12:34:56.789Z' \
--header 'ce-type: google.cloud.pubsub.topic.v1.messagePublished' \
--header 'ce-source: //pubsub.googleapis.com/projects/game-alerts-431604/topics/game-alerts' \
--data '{}'
```

# Helpful Docs

- [Functions Framework Go](https://github.com/GoogleCloudPlatform/functions-framework-go)