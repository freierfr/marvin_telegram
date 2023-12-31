terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.21"
      source  = "cloudflare/cloudflare"
    }
  }
  required_version = ">= 1.1.0"
}



data "terraform_remote_state" "google_workspace" {
  backend = "azurerm"

  config = {
    resource_group_name  = "tfstate"
    storage_account_name = "tfstatefreier"
    container_name       = "tfstate"
    key                  = "terraform_google_workspace.tfstate"
  }
}
