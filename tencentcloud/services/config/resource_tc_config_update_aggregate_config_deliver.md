Provides a resource to manage the aggregate delivery settings (投递设置) of a Tencent Cloud Config account group (账号组).

Example Usage

```hcl
resource "tencentcloud_config_update_aggregate_config_deliver" "example" {
  account_group_id     = "ca-ag-xxxxxx"
  status               = 1
  deliver_name         = "tf-example-aggregate-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_uin          = 0
  deliver_content_type = 3
}
```

Import

Config aggregate deliver can be imported using the account_group_id, e.g.

```shell
terraform import tencentcloud_config_update_aggregate_config_deliver.example ca-ag-xxxxxx
```