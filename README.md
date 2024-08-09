# TODOS

- Containerize
- Provision the deployment
- Deploy via GH actions
- Create db and check if i've already sent this alert
- Write tests
- Add a database to store coaches emails, league id, etc. 
- Build support for sending emails to different coaches. 

# Deployment

```
gcloud functions deploy SendGameAlert \
--gen2 \
--region=us-west1 \
--runtime=go122 \
--entry-point=SendGameAlert \
--trigger-topic=game-alerts
```
