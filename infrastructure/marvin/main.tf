data "terraform_remote_state" "google_workspace" {
  backend = "azurerm"

  config = {
    resource_group_name  = "tfstate"
    storage_account_name = "tfstatefreier"
    container_name       = "tfstate"
    key                  = "terraform_google_workspace.tfstate"
  }
}
