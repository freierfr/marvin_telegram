variable "env" {
    type = string
}

variable "location" {
  type=string
}

variable "OPENAI_API_KEY" {
  type = string
  sensitive = true
}

variable "TELEGRAM_TOKEN" {
  type = string
  sensitive = true
}

variable "REDIS_HOST" {
  type = string
  sensitive = true
}

variable "REDIS_LOGIN" {
  type = string
}

variable "REDIS_PASSWORD" {
  type = string
  sensitive = true
}

variable "CLOUDFLARE_API_TOKEN" {
  type = string
  sensitive = true
}

