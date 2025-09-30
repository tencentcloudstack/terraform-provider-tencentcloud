Provides a resource to create a WeData code file

Example Usage

```hcl
resource "tencentcloud_wedata_code_file" "example" {
  project_id         = ""
  code_file_name     = ""
  parent_folder_path = ""
  code_file_config {
    params = ""
    notebook_session_info {
      notebook_session_id   = ""
      notebook_session_name = ""
    }
  }

  code_file_content = ""
}
```

Import

WeData code file can be imported using the projectId#codeFileId, e.g.

```
terraform import tencentcloud_wedata_code_file.example 1470547050521227264#2bfa8813-344f-4858-a2cc-7a07bd10ac1d
```
