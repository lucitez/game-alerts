# TODOS

- Add sql support to production
- Look into scheduling and running using purely gha
- Write tests

### Eventually
- Use terraform to provision cloud function, db, etc

# Development

Spin up the supabase db container

```shell
$ supabase start
```

Migrate and seed local db
```shell
$ supabase db reset
```

Connect to local db
```shell
$ supabase status
$ psql DB_URL
```

Start Cloud Function:
`go run ./cmd/main.go`

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
- [Supabase](https://supabase.com/docs/guides/database/overview)