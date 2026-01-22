Provides a resource to create a WeData code file

Example Usage

```hcl
resource "tencentcloud_wedata_code_folder" "example" {
  project_id         = "2983848457986924544"
  folder_name        = "tf_example"
  parent_folder_path = "/"
}

resource "tencentcloud_wedata_code_file" "example" {
  project_id         = "2983848457986924544"
  code_file_name     = "tf_example_code_file"
  parent_folder_path = tencentcloud_wedata_code_folder.example.path
  code_file_content  = "Hello Terraform"
}
```

Import

WeData code file can be imported using the projectId#codeFileId, e.g.

```
terraform import tencentcloud_wedata_code_file.example 1470547050521227264#2bfa8813-344f-4858-a2cc-7a07bd10ac1d
```
