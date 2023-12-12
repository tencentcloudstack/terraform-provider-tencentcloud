Use this data source to query detailed information of ckafka group_offsets

Example Usage

```hcl
data "tencentcloud_ckafka_group_offsets" "group_offsets" {
  instance_id = "ckafka-xxxxxx"
  group = "xxxxxx"
}
```