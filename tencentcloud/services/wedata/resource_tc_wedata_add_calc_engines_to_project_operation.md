Provides a resource to create a WeData add calc engines to project operation

Example Usage

```hcl
resource "tencentcloud_wedata_add_calc_engines_to_project_operation" "example" {
  project_id = "20241107221758402"
  dlc_info {
    compute_resources = [
      "dlc_linau6d4bu8bd5u52ffu52a8"
    ]
    region           = "ap-guangzhou"
    default_database = "default_db"
  }
}
```
