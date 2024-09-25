Provides a resource to restore object

~> **NOTE:** The current resource does not support cdc.

Example Usage

```hcl
resource "tencentcloud_cos_object_restore_operation" "object_restore" {
    bucket = "keep-test-1308919341"
    key = "test-restore.txt"
    tier = "Expedited"
    days = 2
}
```