Provides a resource to create a BH asset sync flag config

Example Usage

```hcl
resource "tencentcloud_bh_asset_sync_flag_config" "example" {
  auto_sync = true
}
```

Import

BH asset sync flag config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_bh_asset_sync_flag_config.example yci5a1o76a5HzCaqJM2bQA==
```
