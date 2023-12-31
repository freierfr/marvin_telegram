module "marvin" {
  source = "../../marvin"

  env=var.env
  location=var.location
  OPENAI_API_KEY = var.OPENAI_API_KEY
  TELEGRAM_TOKEN = var.TELEGRAM_TOKEN
  REDIS_HOST = var.REDIS_HOST
  REDIS_LOGIN = var.REDIS_LOGIN
  REDIS_PASSWORD = var.REDIS_PASSWORD
  CLOUDFLARE_API_TOKEN = var.CLOUDFLARE_API_TOKEN
}
