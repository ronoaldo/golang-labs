# Dataflow Pipeline Demo

Simple demo for a branching pipeline in Go, that is also
used as a templated pipeline with Dataflow Flex Templates.

## Pipeline

The pipeline uses some DoFns and the I/Os for Bigquery and TextIO.
Some helpful commands to run the pipeline locally are available
in the provided `Makefile`.

**TIP**: render the pipeline using `-runner=dot -dot_file=pipeline.dot`
to debug incompatible types of PCollections while developing.

## Flex Template

Sample steps extracted from https://cloud.google.com/dataflow/docs/guides/templates/using-flex-templates:

0. Setup

```bash
export PROJECT=
export BUCKET=
export REGISTRY=df-templates
```

1. Build the pipeline binary:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pipeline
```

2. Create a Bucket

```bash
gcloud storage buckets create gs://$BUCKET
```

3. Create an artifact repo

```bash
gcloud artifacts repositories create $REGISTRY \
 --repository-format=docker \
 --location=us-central1
```

4. Configure auth

```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

5. Build the image

```bash
gcloud builds submit \
    --tag us-central1-docker.pkg.dev/$PROJECT/$REGISTRY/dataflow/gcp-releases:latest .

gcloud dataflow flex-template build gs://$BUCKET/template/gcp-releases.json \
    --image "us-central1-docker.pkg.dev/$PROJECT/$REGISTRY/dataflow/gcp-releases:latest" \
    --sdk-language "GO" \
    --metadata-file "metadata.json"
```

6. Submit your job

```bash
gcloud dataflow flex-template run "gcp-releases-`date +%Y%m%d-%H%M%S`" \
    --template-file-gcs-location "gs://$BUCKET/template/gcp-releases.json" \
    --parameters output="gs://$BUCKET/gcp-releases/" \
    --additional-user-labels "owner=ronoaldo" \
    --region "us-central1"
```