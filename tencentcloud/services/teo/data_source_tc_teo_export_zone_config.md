Use this data source to export TEO (EdgeOne) site configuration

Example Usage

Export all types of zone configuration

```hcl
data "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

Export specific types of zone configuration

```hcl
data "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-3fkff38fyw8s"
  types    = ["L7AccelerationConfig"]
}
```
