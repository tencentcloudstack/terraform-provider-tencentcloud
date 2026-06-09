Use this data source to query detailed information of Config discovered resources.

Example Usage

Query all discovered resources

```hcl
data "tencentcloud_config_discovered_resources" "example" {}
```

Query by resource ID filter

```hcl
data "tencentcloud_config_discovered_resources" "example" {
  filters {
    name   = "resourceId"
    values = ["ins-pbu2hghz"]
  }
  order_type = "desc"
}
```

Query by resource name and tags

```hcl
data "tencentcloud_config_discovered_resources" "example" {
  filters {
    name   = "resourceName"
    values = ["my-cvm"]
  }

  tags {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```
