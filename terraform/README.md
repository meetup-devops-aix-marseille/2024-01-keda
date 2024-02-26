# Terraform

Login to GCP and set the default project and region.

```bash
gcloud config configurations activate $GCP_CONFIG_NAME
gcloud auth application-default login
gcloud config set project $GCP_PROJECT
gcloud config set compute/region $GCP_REGION
```

## Get gke credentials

```bash
gcloud container clusters get-credentials --project meetup-devops meetup-keda --region europe-west1
```
