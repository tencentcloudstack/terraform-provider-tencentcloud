Provides a resource to start an AGS sandbox instance.

-> **Note:** `tool_id`, `tool_name`, `source_instance_id`, `snapshot_id`, `client_token`, `mount_options`, `custom_configuration`, `auth_mode`, `storage_mounts` and `start_mode` are create-only: the update API does not accept them, so changing any of them destroys the existing sandbox instance and starts a new one.

-> **Note:** `timeout`, `source_instance_id`, `snapshot_id`, `client_token`, `storage_mounts` and `start_mode` are write-only inputs. The describe API reports the timeout in seconds (`timeout_seconds`) instead of the original duration string and does not report the other fields at all, so the configured values are kept in state and drift on them cannot be detected.

-> **Note:** Destroying this resource only stops the sandbox instance itself. It does not touch the source instance, the sandbox tool or any storage (COS bucket, CFS file system or volume) referenced by `storage_mounts`. Whether the data of a stopped instance is retained afterwards is decided by the cloud side.

-> **Note:** The start and stop requests are asynchronous: the create path polls the instance status until it is `RUNNING`, and the destroy path polls it until it is `STOPPED` or no longer visible.

Example Usage

```hcl
resource "tencentcloud_ags_sandbox_instance" "example" {
  tool_id = "tool-xxxxxxxx"
  timeout = "30m"

  metadata = [
    {
      name  = "run_id"
      value = "42"
    },
  ]

  custom_configuration = {
    image               = "ccr.ccs.tencentyun.com/test/go-httpbin:v2.16.1"
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
      cpu     = "1"
      memory  = "1Gi"
      storage = null
    }
  }

  auth_mode = "TOKEN"
}
```

Import

AGS sandbox instance can be imported using the sandbox instance id, e.g.

```
$ terraform import tencentcloud_ags_sandbox_instance.example instance-xxxxxxxx
```
