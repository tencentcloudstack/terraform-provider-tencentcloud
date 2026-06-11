Use this data source to query detailed information of MongoDB (mongodb) DB instance node property

Example Usage

```hcl
data "tencentcloud_mongodb_db_instance_node_property" "example" {
  instance_id = "cmgo-5aqo4yf7"
}
```

Example Usage with filters

```hcl
data "tencentcloud_mongodb_db_instance_node_property" "example" {
  instance_id  = "cmgo-5aqo4yf7"
  roles        = ["PRIMARY", "SECONDARY"]
  only_hidden  = false
  priority     = 1
  votes        = 1
  tags {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```
