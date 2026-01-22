Use this data source to query detailed information of teo zoneAvailablePlans

Example Usage

```hcl
data "tencentcloud_teo_zones" "teo_zones" {
  filters {
    name = "zone-id"
    values = ["zone-39quuimqg8r6"]
  }

  filters {
    name = "tag-key"
    values = ["createdBy"]
  }

  filters {
    name = "tag-value"
    values = ["terraform"]
  }
}
```