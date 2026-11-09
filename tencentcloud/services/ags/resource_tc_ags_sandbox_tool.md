Provides a resource to create an AGS sandbox tool.

-> **Note:** `tool_name`, `tool_type`, `default_timeout`, `client_token`, `role_arn`, `storage_mounts`, `log_configuration` and `persistent` are create-only: the update API does not accept them, so changing any of them destroys the existing sandbox tool and creates a new one.

-> **Note:** `default_timeout` and `client_token` are only accepted by the create API. The describe API reports the timeout in seconds (`DefaultTimeoutSeconds`) instead of the original duration string, so the values configured here are kept as they are and drift on them cannot be detected.

-> **Note:** Removing `description`, `tags`, `custom_configuration`, `computer_configuration` or `persistent` from the configuration keeps the value that currently exists on the cloud side, because the update API does not transmit empty values. To clear those attributes, delete and recreate the sandbox tool.

-> **Note:** Destroying this resource deletes the sandbox tool itself only. It does not terminate sandbox instances that were created from the tool and it does not release any storage (COS bucket, CFS file system or volume) referenced by `storage_mounts`.

Example Usage

```hcl

resource "tencentcloud_ags_sandbox_tool" "custom" {
  client_token = null
  custom_configuration = {
    args                = ["-port", "8080"]
    command             = ["/bin/go-httpbin"]
    dns_config          = null
    env                 = null
    image               = "ccr.ccs.tencentyun.com/test/go-httpbin:v2.16.1"
    image_registry_type = "personal"
    ports = [
      {
        name     = "http"
        port     = 8080
        protocol = "TCP"
      },
    ]
    probe = {
      failure_threshold = 60
      http_get = {
        path   = "/status/200"
        port   = 8080
        scheme = "HTTP"
      }
      probe_period_ms   = 1000
      probe_timeout_ms  = 1000
      ready_timeout_ms  = 30000
      success_threshold = 1
    }
    resources = {
      cpu     = "1"
      memory  = "1Gi"
      storage = null
    }
    security_context = null
  }
  default_timeout   = null
  description       = "Guoning standard web sandbox: intranet TCR go-httpbin v2.16.1, 1C1Gi, HTTP /status/200:8080"
  log_configuration = null
  network_configuration = {
    network_mode = "PUBLIC"
    vpc_config   = null
  }
  persistent     = false
  role_arn       = "qcs::cam::uin/xxxxxx:roleName/cos-tcr-full-ags"
  storage_mounts = null
  tool_name = var.tool_name
  tool_type = "custom"
}
```

Import

AGS sandbox tool can be imported using the sandbox tool id, e.g.

```
$ terraform import tencentcloud_ags_sandbox_tool.example tool-xxxxxxxx
```
