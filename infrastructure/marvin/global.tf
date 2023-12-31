# data "azurerm_resource_group" "global" {
#   name = "global"
# }

# data "azurerm_dns_zone" "stumpy_fr" {
#   name = "stumpy.fr"
# }

data "cloudflare_zone" "stumpy_fr" {
  name = "stumpy.fr"
}
