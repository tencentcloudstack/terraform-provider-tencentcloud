Use this data source to query detailed information of TEO available origin acl family

Example Usage

Query available origin acl family by zone id

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

Query available origin acl family by filters

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"
  filters {
    name   = "OriginACLFamily"
    values = ["gaz"]
  }
}
```
