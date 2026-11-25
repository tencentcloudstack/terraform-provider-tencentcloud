Provides a resource to start an AGS sandbox instance.

Example Usage

```hcl
resource "tencentcloud_ags_sandbox_instance" "example" {
  tool_name = "tf-example-tool"
  timeout   = "30m"
  auth_mode = "DEFAULT"

  mount_options = [
    {
      name       = "workspace"
      mount_path = "/workspace"
      read_only  = false
    },
  ]

  custom_configuration = {
    image               = "ccr.ccs.tencentyun.com/tf-example/go-httpbin:v2.16.1"
    image_registry_type = "personal"
    command             = ["/bin/go-httpbin"]
    args                = ["-port", "8080"]
    ports = [
      {
        name     = "http"
        port     = 8080
        protocol = "TCP"
      },
    ]
    resources = {
      cpu    = "1"
      memory = "1Gi"
    }
    env = [
      {
        name  = "EXAMPLE_KEY"
        value = "example-value"
      },
    ]
  }

  metadata = [
    {
      name  = "run_id"
      value = "tf-example"
    },
  ]
}
```

Import

AGS sandbox instance can be imported using the id, e.g.

```
terraform import tencentcloud_ags_sandbox_instance.example instance-xxxxxxxx
```
