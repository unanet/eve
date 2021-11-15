# eve

![pipeline status](https://github.com/unanet/eve/badges/master/pipeline.svg)

## Environment Variables

```bash
EVE_DB_PASSWORD={{eve_password}}
EVE_DB_USERNAME=eve
EVE_DB_HOST={{db_url}}
EVE_DB_NAME=eve

EVE_ADMIN_TOKEN={{admin_token}}

EVE_IDENTITY_CLIENT_ID={{id}}
EVE_IDENTITY_CLIENT_SECRET={{client_secret}}
EVE_IDENTITY_CONN_URL=https://{{domain}}/auth/realms/{{realm}}
EVE_IDENTITY_REDIRECT_URL=https://{{domain}}/oidc/callback

EVE_ARTIFACTORY_BASE_URL=https://{{domain}}jfrog.io/{{company_name}}/api
EVE_ARTIFACTORY_API_KEY={{artifactor_API_key}}
API_Q_URL==https://sqs.us-{{aws_region}}.amazonaws.com/{{account_id}}/{{sqs_queue_name}}.fifo
S3_BUCKET={{bucket_name}}
AWS_REGION=us-east-2
```

The following additional variables can be used for local dev 
```bash
LOCAL_DEV=true # Used to disable cron and deployment queue from starting
```