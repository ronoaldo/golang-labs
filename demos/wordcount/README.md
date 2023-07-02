# wordcount

Exemplo de construção e execução de uma pipeline Apache Beam
com a linguagem Go.

Ref: https://github.com/apache/beam/blob/master/sdks/go/examples/wordcount/wordcount.go


## Executando no Google Cloud

```bash

export PROJECT=
export BUCKET=

./wordcount \
    --output gs://$BUCKET/wordcount-go/output/wordcount.txt \
    --runner DataflowRunner \
    --async \
    --project=$PROJECT \
    --region=us-central1 \
    --temp_location=gs://$BUCKET/wordcount-go/temp/ \
    --staging_location=gs://$BUCKET/wordcount-go/staging/
```