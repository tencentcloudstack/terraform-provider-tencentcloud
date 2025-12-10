Provides a resource to create a Controlcenter batch apply account baselines

Example Usage

```hcl
resource "tencentcloud_batch_apply_account_baselines" "example" {
  member_uin_list = [
    10037652245,
    10037652240,
  ]

  baseline_config_items {
    identifier    = "TCC-AF_SHARE_IMAGE"
    configuration = "{\"Images\":[{\"Region\":\"ap-guangzhou\",\"ImageId\":\"img-mcdsiqrx\",\"ImageName\":\"demo1\"}, {\"Region\":\"ap-guangzhou\",\"ImageId\":\"img-esxgkots\",\"ImageName\":\"demo2\"}]}"
  }
}
```
