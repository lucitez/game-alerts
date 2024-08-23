# TODOS

- Containerize
- Provision the deployment
- Create db and check if i've already sent this alert
- Write tests
- Add a database to store coaches emails, league id, etc.
- Build support for sending emails to different coaches.

```
curl --location 'localhost:8080' \
--header 'Content-Type: application/json' \
--header 'ce-id: 123451234512345' \
--header 'ce-specversion: 1.0' \
--header 'ce-time: 2020-01-02T12:34:56.789Z' \
--header 'ce-type: google.cloud.pubsub.topic.v1.messagePublished' \
--header 'ce-source: //pubsub.googleapis.com/projects/game-alerts-431604/topics/game-alerts' \
--data '{}'
```
