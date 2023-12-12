Provides a resource to create a CAM policy.

Example Usage

```hcl
resource "tencentcloud_cam_policy" "foo" {
  name        = "tf_cam_policy"
  document    = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": [
        "name/sts:AssumeRole"
      ],
      "effect": "allow",
      "resource": [
        "*"
      ]
    }
  ]
}
EOF
  description = "tf_test"
}
```

Import

CAM policy can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_policy.foo 26655801
```