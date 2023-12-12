Provides a resource to create a cam tag_role

Example Usage

```hcl
resource "tencentcloud_cam_tag_role_attachment" "tag_role" {
  tags {
    key = "test1"
    value = "test1"
  }
  role_id = "test-cam-tag"
}
```

Import

cam tag_role can be imported using the id, e.g.

```
terraform import tencentcloud_cam_tag_role_attachment.tag_role tag_role_id
```