# Marvin - Telegram bot

A Telegram bot that integrates with OpenAI's official ChatGPT and DALL·E APIs to provide answers.
Ready to use with minimal configuration required.

## Features

- [x] Deployed as an Azure Function App
- [x] Image generation via the `/image` command
- [x] Whitelist of users that can use the bot
- [x] basic memory of conversations, using redis as a backend
- [x] Group chat support with optional `keyword` prefix

## Configuration

The bot is configured via environment variables, which need to be set in the Azure Function App configuration (we use Terraform to deploy the infrastructure, see bellow).

| Variable | Description | Default |
| -------- | ----------- | ------- |
| OPENAI_MODEL | OpenAI model to use | gpt-4 |
| BOT_PROMPT | Prompt to use when generating text | You are a helpful assistant and your name is Marvin |
| GROUP_TRIGGER_KEYWORD | Keyword to use in group chats | marvin |
| OPENAI_API_KEY | OpenAI API key | |
| TELEGRAM_TOKEN | Telegram bot token | |
| ALLOWED_TELEGRAM_USER_IDS | Comma-separated list of Telegram user IDs that can use the bot | |
| REDIS_URL | Redis URL | redis-10579.c56.east-us.azure.cloud.redislabs.com:10579 |
| REDIS_LOGIN | Redis login | |
| REDIS_PASSWORD | Redis password | |


## Deployment

You need to have the following tools installed:

- [Terraform](https://www.terraform.io/downloads.html)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)

### Login to Azure

```bash
az login
```

### Deploy infrastructure

```bash
cd infrastructure/env/(dev|prod)
terraform init
terraform plan
terraform apply
```


## Installation

You need to have the following tools installed:

- [Go](https://golang.org/doc/install)
- [protoc](https://grpc.io/docs/protoc-installation/)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)
- [Azure functions core tools](https://learn.microsoft.com/en-us/azure/azure-functions/functions-run-local)

```bash
protoc -I=proto/ --go_out=. proto/*.proto
GOOS=linux GOARCH=amd64 go build handler.go
func azure functionapp publish marvin-(dev|prd)-function-app
```

# Disclaimer

This is a personal project and is not affiliated with OpenAI in any way.
