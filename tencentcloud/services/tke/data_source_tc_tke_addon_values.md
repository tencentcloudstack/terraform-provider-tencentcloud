Use this data source to query detailed information of TKE addon values.

Example Usage

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "ip-masq-agent"
}

locals {
  current_values   = jsondecode(data.tencentcloud_tke_addon_values.example.values)
  available_values = jsondecode(data.tencentcloud_tke_addon_values.example.default_values)
}

output "addon_current_values" {
  value = local.current_values
}

output "addon_available_values" {
  value = local.available_values
}
```
