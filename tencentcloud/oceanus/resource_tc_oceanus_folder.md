Provides a resource to create a oceanus folder

Example Usage

```hcl
resource "tencentcloud_oceanus_folder" "example" {
  folder_name   = "tf_example"
  parent_id     = "folder-lfqkt11s"
  folder_type   = 0
  work_space_id = "space-125703345ap-shenzhen-fsi"
}
```

Import

oceanus folder can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_folder.example space-125703345ap-shenzhen-fsi#folder-f40fq79g#0
```
