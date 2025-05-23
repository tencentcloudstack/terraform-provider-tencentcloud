Use this data source to query detailed information of Private Dns end points

Example Usage

Query all private dns end points

```hcl
data "tencentcloud_private_dns_end_points" "example" {}
```

Query all private dns end points by filters

```hcl
data "tencentcloud_private_dns_end_points" "example" {
  filters {
    name   = "EndPointName"
    values = ["tf-example"]
  }

  filters {
    name   = "EndPointId"
    values = ["eid-72dc11b8f3"]
  }

  filters {
    name   = "EndPointServiceId"
    values = ["vpcsvc-61wcwmar"]
  }

  filters {
    name   = "EndPointVip"
    values = [
      "172.10.10.1"
    ]
  }
}
```
