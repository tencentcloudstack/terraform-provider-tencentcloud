Provides a resource to create a chdfs file_system

Example Usage

```hcl
resource "tencentcloud_chdfs_file_system" "file_system" {
  capacity_quota           = 10995116277760
  description              = "file system for terraform test"
  enable_ranger            = true
  file_system_name         = "terraform-test"
  posix_acl                = false
  ranger_service_addresses = [
    "127.0.0.1:80",
    "127.0.0.1:8000",
  ]
  super_users              = [
    "terraform",
    "iac",
  ]
}
```

Import

chdfs file_system can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_file_system.file_system file_system_id
```