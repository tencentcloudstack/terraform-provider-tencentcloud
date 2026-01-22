Use this data source to query detailed information of BH devices

Example Usage

Query all bh devices

```hcl
data "tencentcloud_bh_devices" "example" {}
```

Query bh devices by filters

```hcl
data "tencentcloud_bh_devices" "example" {
  id_set = [
    107,
    108,
    109,
    110,
  ]

  name = "tf-example"

  ap_code_set = [
    "ap-guangzhou",
    "ap-beijing",
    "ap-shanghai",
  ]

  kind_set = [
    1, 2, 3, 4
  ]

  filters {
    name   = "InstanceId"
    values = ["ext-21ae68e02-4570-1"]
  }

  tag_filters {
    tag_key = "tagKey"
    tag_value = [
      "tagValue1",
      "tagValue2",
    ]
  }
}
```
