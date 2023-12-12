Provides a CBS storage attachment resource.

Example Usage

```hcl
resource "tencentcloud_cbs_storage_attachment" "attachment" {
  storage_id  = "disk-kdt0sq6m"
  instance_id = "ins-jqlegd42"
}
```

Import

CBS storage attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage_attachment.attachment disk-41s6jwy4
```