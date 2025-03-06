Provides a resource to create a CAM tag role

Example Usage

Create by role_id

```hcl
resource "tencentcloud_cam_tag_role_attachment" "example" {
  role_id = "4611686018441060141"

  tags {
    key   = "tagKey"
    value = "tagValue"
  }
}
```

Create by role_name

```hcl
resource "tencentcloud_cam_tag_role_attachment" "example" {
  role_name = "tf-example"

  tags {
    key   = "tagKey"
    value = "tagValue"
  }
}
```

Import

CAM tag role can be imported using the id, e.g.

```
# Please use role_name#role_id
terraform import tencentcloud_cam_tag_role_attachment.example tf-example#4611686018441060141
```