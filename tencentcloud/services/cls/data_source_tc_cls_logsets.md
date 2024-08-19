Use this data source to query detailed information of cls logsets

Example Usage

Query all cls logsets

```hcl
data "tencentcloud_cls_logsets" "logsets" {}
```

Query by filters

```hcl
# Query by `logsetName`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "logsetName"
    values = ["cls_service_logging"]
  }
}

# Query by `logsetId`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "logsetId"
    values = ["50d499a8-c4c0-4442-aa04-e8aa8a02437d"]
  }
}

# Query by `tagKey`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "tagKey"
    values = ["createdBy"]
  }
}

# Query by `tag:tagKey`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "tag:createdBy"
    values = ["terraform"]
  }
}
```