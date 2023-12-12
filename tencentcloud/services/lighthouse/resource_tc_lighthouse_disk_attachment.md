Provides a resource to create a lighthouse disk_attachment

Example Usage

```hcl
resource "tencentcloud_lighthouse_disk_attachment" "disk_attachment" {
  disk_id = "lhdisk-xxxxxx"
  instance_id = "lhins-xxxxxx"
}
```

Import

lighthouse disk_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_disk_attachment.disk_attachment disk_attachment_id
```