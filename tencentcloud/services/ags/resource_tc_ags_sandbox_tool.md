Provides a resource to create a AGS sandbox tool

Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "tf-example"
  tool_type = "custom"
  network_configuration {
    network_mode = "PUBLIC"
  }
  description     = "example sandbox tool"
  default_timeout = "10m"
  tags = {
    "env" = "test"
  }
}

resource "tencentcloud_ags_sandbox_tool" "vpc" {
  tool_name = "tf-vpc-example"
  tool_type = "custom"
  network_configuration {
    network_mode = "VPC"
    vpc_config {
      subnet_ids         = ["subnet-xxxxxxxx"]
      security_group_ids = ["sg-xxxxxxxx"]
    }
  }
  description     = "VPC sandbox tool example"
  default_timeout = "30m"
  custom_configuration {
    image = "ccr.ccs.tencentyun.com/namespace/image:tag"
    resources {
      cpu    = "2"
      memory = "4Gi"
    }
  }
}
```

Import

AGS sandbox tool can be imported using the toolId, e.g.

```
terraform import tencentcloud_ags_sandbox_tool.example sbt-xxxxxxxx
```