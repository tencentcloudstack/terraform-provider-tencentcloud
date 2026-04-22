Use this data source to query detailed information of TEO default certificates

Example Usage

```hcl
data "tencentcloud_teo_default_certificate" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

Query with filters

```hcl
data "tencentcloud_teo_default_certificate" "example" {
  filters {
    name = "zone-id"
    values = [
      "zone-2qtuhspy7cr6"
    ]
  }
}
```
