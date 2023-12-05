Provides a resource to restore object

Example Usage

```hcl
resource "tencentcloud_cos_object_restore_operation" "object_restore" {
    bucket = "keep-test-1308919341"
    key = "test-restore.txt"
    tier = "Expedited"
    days = 2
}
```