Use this data source to query detailed information of ckafka group_info

Example Usage

```hcl
data "tencentcloud_ckafka_group_info" "group_info" {
  instance_id = "ckafka-xxxxxx"
  group_list = ["xxxxxx"]
}
```