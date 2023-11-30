Provides a resource to create a rum offline_log_config_attachment

Example Usage

```hcl
resource "tencentcloud_rum_offline_log_config_attachment" "offline_log_config_attachment" {
  project_key = "ZEYrYfvaYQ30jRdmPx"
  unique_id = "100027012454"
}

```
Import

rum offline_log_config_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_offline_log_config_attachment.offline_log_config_attachment ZEYrYfvaYQ30jRdmPx#100027012454
```