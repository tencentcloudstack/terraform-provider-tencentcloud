Provides a resource to create a cfs user_quota

Example Usage

```hcl
resource "tencentcloud_cfs_user_quota" "user_quota" {
  file_system_id = "cfs-4636029bc"
  user_type = "Uid"
  user_id = "2159973417"
  capacity_hard_limit = 10
  file_hard_limit = 10000
}
```

Import

cfs user_quota can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_user_quota.user_quota fileSystemId#userType#userId
```
