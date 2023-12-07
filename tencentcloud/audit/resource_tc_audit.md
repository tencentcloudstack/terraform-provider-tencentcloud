Provides a resource to create an audit.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_audit_track.

Example Usage

```hcl
resource "tencentcloud_audit" "foo" {
  name                 = "audittest"
  cos_bucket           = "test"
  cos_region           = "ap-hongkong"
  log_file_prefix      = "test"
  audit_switch         = true
  read_write_attribute = 3
}
```

Import

Audit can be imported using the id, e.g.

```
$ terraform import tencentcloud_audit.foo audit-test
```