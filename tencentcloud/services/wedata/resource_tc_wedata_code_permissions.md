Provides a resource to create a WeData code permissions

Example Usage

```hcl
resource "tencentcloud_wedata_code_permissions" "example" {
  project_id = "3108707295180644352"
  authorize_permission_objects {
    resource {
      resource_type        = ""
      resource_id          = ""
      resource_id_for_path = ""
      resource_cfs_path    = ""
    }

    authorize_subjects {
      subject_type   = ""
      subject_values = []
      privileges     = []
    }
  }
}
```

Import

WeData code permissions can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_code_permissions.example wedata_code_permissions_id
```
