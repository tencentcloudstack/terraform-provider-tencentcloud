Use this data source to query the list of EdgeOne (TEO) custom groups.

Example Usage

Query all custom groups of a site

```hcl
data "tencentcloud_teo_custom_groups" "example" {
  zone_id = "zone-27h0vbm5w1e"
}
```

Query custom groups filtered by group type

```hcl
data "tencentcloud_teo_custom_groups" "example" {
  zone_id = "zone-27h0vbm5w1e"

  group_types = ["TLSFingerprint", "UA"]
}
```
