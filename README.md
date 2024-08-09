# TODOS

- Containerize
- Provision the deployment
- Deploy via GH actions
- Create db and check if i've already sent this alert
- Write tests
- Get league id from payload instead of env

# Deployment

```
gcloud functions deploy SendGameAlert \
--gen2 \
--region=us-west1 \
--runtime=go122 \
--entry-point=SendGameAlert \
--trigger-topic=game-alerts
```
