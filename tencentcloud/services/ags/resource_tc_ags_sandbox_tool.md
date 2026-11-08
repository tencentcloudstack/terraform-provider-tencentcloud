Provides a resource to create an AGS sandbox tool.

-> **Note:** `tool_name`, `tool_type`, `default_timeout`, `client_token`, `role_arn`, `storage_mounts`, `log_configuration` and `persistent` are create-only: the update API does not accept them, so changing any of them destroys the existing sandbox tool and creates a new one.

-> **Note:** `default_timeout` and `client_token` are only accepted by the create API. The describe API reports the timeout in seconds (`DefaultTimeoutSeconds`) instead of the original duration string, so the values configured here are kept as they are and drift on them cannot be detected.

-> **Note:** Removing `description`, `tags`, `custom_configuration`, `computer_configuration` or `persistent` from the configuration keeps the value that currently exists on the cloud side, because the update API does not transmit empty values. To clear those attributes, delete and recreate the sandbox tool.

-> **Note:** Destroying this resource deletes the sandbox tool itself only. It does not terminate sandbox instances that were created from the tool and it does not release any storage (COS bucket, CFS file system or volume) referenced by `storage_mounts`.

Example Usage

Create a custom sandbox tool with a public network

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "tf-example-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "PUBLIC"
  }

  description = "sandbox tool created by terraform"

  custom_configuration {
    image               = "ccr.ccs.tencentyun.com/tf-example/agent:latest"
    image_registry_type = "personal"
    command             = ["/bin/sh", "-c"]
    args                = ["python3 -m http.server 8080"]

    resources {
      cpu     = "1"
      memory  = "2Gi"
      storage = "5Gi"
    }

    ports {
      name     = "http"
      port     = 8080
      protocol = "TCP"
    }

    env {
      name  = "LOG_LEVEL"
      value = "info"
    }
  }
}
```

Create a sandbox tool in a VPC with a storage mount and log push

```hcl
resource "tencentcloud_ags_sandbox_tool" "vpc" {
  tool_name  = "tf-example-vpc-tool"
  tool_type  = "browser"
  persistent = true

  network_configuration {
    network_mode = "VPC"

    vpc_config {
      subnet_ids         = ["subnet-xxxxxxxx"]
      security_group_ids = ["sg-xxxxxxxx"]
    }
  }

  storage_mounts {
    name       = "workspace"
    mount_path = "/workspace"
    read_only  = false

    storage_source {
      cos {
        endpoint    = "cos.ap-guangzhou.myqcloud.com"
        bucket_name = "tf-example-bucket"
        bucket_path = "/workspace"
      }
    }
  }

  log_configuration {
    cls_config {
      topic_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    }

    log_sources {
      files = ["/logs/app.log"]
    }
  }

  tags {
    key   = "created_by"
    value = "terraform"
  }
}
```

Import

AGS sandbox tool can be imported using the sandbox tool id, e.g.

```
$ terraform import tencentcloud_ags_sandbox_tool.example tool-xxxxxxxx
```
