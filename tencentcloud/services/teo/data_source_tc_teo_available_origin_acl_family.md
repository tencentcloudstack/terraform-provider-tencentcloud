Use this data source to query available origin ACL families of TEO origin protection

Example Usage

Query all available origin ACL families by zone id

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

Query available origin ACL families by filter

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"

  filters {
    name   = "OriginACLFamily"
    values = ["gaz", "mlc"]
  }
}
```