Use this data source to query detailed information of TEO environments, including the environment list, the total environment count (`total_count`), and the source version (`source_version`) of each effective config group version.

Example Usage

```hcl
data "tencentcloud_teo_environments" "teo_environments" {
  zone_id = "zone-2qtuhspy7cr6"
}
```