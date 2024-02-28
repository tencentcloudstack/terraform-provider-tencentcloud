Use this data source to query detailed information of privatedns private_zone_list

Example Usage

Get All PrivateZones:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {}
```

Get PrivateZone By ZoneId:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "ZoneId"
    values = ["zone-6xg5xgky1"]
  }
}
```

Get PrivateZone By Domain:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "Domain"
    values = ["domain.com"]
  }
}
```

Get PrivateZone By Vpc:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "Vpc"
    values = ["vpc-axrsmmrv"]
  }
}
```
