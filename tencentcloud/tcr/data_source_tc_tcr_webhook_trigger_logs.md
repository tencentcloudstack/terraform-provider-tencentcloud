Use this data source to query detailed information of tencentcloud_tcr_webhook_trigger_logs

Example Usage

```hcl
data "tencentcloud_tcr_webhook_trigger_logs" "my_logs" {
  registry_id = local.tcr_id
  namespace = var.tcr_namespace
  trigger_id = var.trigger_id
    tags = {
    "createdBy" = "terraform"
  }
}
```