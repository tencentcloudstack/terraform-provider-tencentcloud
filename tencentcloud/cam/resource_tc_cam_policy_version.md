Provides a resource to create a cam policy_version

Example Usage

```hcl
resource "tencentcloud_cam_policy_version" "policy_version" {
  policy_id = 171173780
  policy_document = jsonencode({
    "version": "2.0",
    "statement": [
      {
        "effect": "allow",
        "action": [
          "sts:AssumeRole"
        ],
        "resource": [
          "*"
        ]
      },
      {
        "effect": "allow",
        "action": [
          "cos:PutObject"
        ],
        "resource": [
          "*"
        ]
      },
      {
        "effect": "deny",
        "action": [
          "aa:*"
        ],
        "resource": [
          "*"
        ]
      },
      {
        "effect": "deny",
        "action": [
          "aa:*"
        ],
        "resource": [
          "*"
        ]
      }
    ]
  })
  set_as_default = "false"
}
```

Import

cam policy_version can be imported using the id, e.g.

```
terraform import tencentcloud_cam_policy_version.policy_version policy_version_id
```