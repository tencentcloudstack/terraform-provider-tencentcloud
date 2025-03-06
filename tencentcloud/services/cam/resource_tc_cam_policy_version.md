Provides a resource to create a CAM policy version

Example Usage

```hcl
resource "tencentcloud_cam_policy_version" "example" {
  policy_id      = 171173780
  set_as_default = "false"
  policy_document = jsonencode({
    "version" : "3.0",
    "statement" : [
      {
        "effect" : "allow",
        "action" : [
          "sts:AssumeRole"
        ],
        "resource" : [
          "*"
        ]
      },
      {
        "effect" : "allow",
        "action" : [
          "cos:PutObject"
        ],
        "resource" : [
          "*"
        ]
      },
      {
        "effect" : "deny",
        "action" : [
          "aa:*"
        ],
        "resource" : [
          "*"
        ]
      },
      {
        "effect" : "deny",
        "action" : [
          "aa:*"
        ],
        "resource" : [
          "*"
        ]
      }
    ]
  })
}
```

Import

CAM policy version can be imported using the id, e.g.

```
terraform import tencentcloud_cam_policy_version.example 234290251#3
```