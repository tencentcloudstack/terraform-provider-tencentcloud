Use this data source to query detailed information of Private Dns private zone list

Example Usage

Query All private zones:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {}
```

Query private zones by ZoneId:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "ZoneId"
    values = ["zone-6xg5xgky1"]
  }
}
```

Query private zones by Domain:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "Domain"
    values = ["domain.com"]
  }
}
```

Query private zones by Vpc:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "Vpc"
    values = ["vpc-axrsmmrv"]
  }
}
```
