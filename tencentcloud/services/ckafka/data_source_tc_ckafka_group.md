Use this data source to query detailed information of ckafka group

Example Usage

```hcl
data "tencentcloud_ckafka_group" "group" {
  instance_id = "ckafka-xxxxxxx"
  search_word = "xxxxxx"
}
```