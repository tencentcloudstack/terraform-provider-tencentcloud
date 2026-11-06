Provides a resource to create and manage an AGS sandbox tool.

Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "my-sandbox-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "PUBLIC"
  }

  description     = "An example sandbox tool"
  default_timeout = "10m"

  tags {
    key   = "env"
    value = "test"
  }

  custom_configuration {
    image               = "nginx:latest"
    image_registry_type = "personal"
    command             = ["/bin/bash"]
    args                = ["-c", "echo hello"]

    env {
      name  = "MY_VAR"
      value = "my_value"
    }

    resources {
      cpu     = "2"
      memory  = "4Gi"
      storage = "5Gi"
    }
  }
}
```

## Import

```hcl
terraform import tencentcloud_ags_sandbox_tool.example <tool-id>
```