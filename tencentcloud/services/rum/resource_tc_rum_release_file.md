Provides a resource to create a rum release file

Example Usage

```hcl
resource "tencentcloud_rum_release_file" "example" {
  project_id      = 123
  version         = "1.0"
  file_key        = "120000-last-1632921299138-index.js.map"
  file_name       = "index.js.map"
  file_hash       = "b148c43fd81d845ba7cc6907928ce430"
  release_file_id = 1
}
```

Import

rum release_file can be imported using the id, e.g.

```
terraform import tencentcloud_rum_release_file.example projectId#releaseFileId
```
