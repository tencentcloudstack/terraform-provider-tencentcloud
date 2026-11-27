Provides a resource to start an AGS sandbox instance.

~> **NOTE:** The `x-ags-` name prefix of `metadata` is reserved by the system. Such entries are added by the cloud side (for example `x-ags-precreate-match`), are not tracked in the state and must not be configured, otherwise a permanent diff is reported.

~> **NOTE:** The timeout of a sandbox instance that is not persistent is enforced by the cloud side: the instance is reclaimed automatically when `timeout` expires. Such an instance is removed from the state during the next refresh, so the next `terraform apply` creates a new sandbox instance instead of updating the existing one.

Example Usage

```hcl

resource "tencentcloud_ags_sandbox_tool" "custom1" {
  custom_configuration = {
    args                = ["/nix/runtime/nexus/deployment/runtime/start.sh", "-H", "0.0.0.0", "-p", "8000", "-v", "3.12"]
    command             = ["/nix/runtime/nix-env/bin/bash"]
    dns_config          = null
    env                 = null
    image               = "ccr.ccs.tencentyun.com/finofliu/go-httpbin:2.19-alpine-with-envd"
    image_registry_type = "personal"
    ports = [
      {
        name     = "http"
        port     = 8000
        protocol = "TCP"
      },
    ]
    probe = {
      failure_threshold = 10
      http_get = {
        path   = "/health"
        port   = 8000
        scheme = "HTTP"
      }
      probe_period_ms   = 3000
      probe_timeout_ms  = 1000
      ready_timeout_ms  = 30000
      success_threshold = 1
    }
    resources = {
      cpu     = "2"
      memory  = "2Gi"
    }
  }
  description       = "image-volume regression baseline"
  network_configuration = {
    network_mode = "PUBLIC"
  }
  persistent = false
  role_arn   = "qcs::cam::uin/3321337994:roleName/cos-tcr-full-ags"
  storage_mounts = [
    {
      mount_path = "/nix"
      name       = "nix-runtime"
      read_only  = true
      storage_source = {
        image = {
          image_registry_type = "enterprise"
          reference           = "tencentags-guangzhou-994-test.tencentcloudcr.com/ags-e2e/nexus-runtime:0.2.2"
          sub_path            = "nix"
        }
      }
    },
    {
      mount_path = "/mnt/same-name-runtime"
      name       = "runtime"
      read_only  = true
      storage_source = {
        image = {
          image_registry_type = "enterprise"
          reference           = "tencentags-guangzhou-994-test.tencentcloudcr.com/ags-e2e/nexus-runtime:0.2.2"
          sub_path            = "nix/runtime/nix-env/bin"
        }
      }
    },
  ]
  tool_name = "tf-test"
  tool_type = "custom"
}

resource "tencentcloud_ags_sandbox_instance" "example1" {
  auth_mode = "DEFAULT"
  custom_configuration = {
    args                = ["/nix/runtime/nexus/deployment/runtime/start.sh", "-H", "0.0.0.0", "-p", "8000", "-v", "3.12"]
    command             = ["/nix/runtime/nix-env/bin/bash"]
    image               = "ccr.ccs.tencentyun.com/finofliu/go-httpbin:2.19-alpine-with-envd"
    image_registry_type = "personal"
    ports = [
      {
        name     = "http"
        port     = 8000
        protocol = "TCP"
      },
    ]
    probe = {
      failure_threshold = 10
      http_get = {
        path   = "/health"
        port   = 8000
        scheme = "HTTP"
      }
      probe_period_ms   = 3000
      probe_timeout_ms  = 1000
      ready_timeout_ms  = 30000
      success_threshold = 1
    }
    resources = {
      cpu    = "2"
      memory = "2Gi"
    }
  }
  mount_options = [
    {
      mount_path = "/nix"
      name       = "nix-runtime"
    },
    {
      mount_path = "/mnt/same-name-runtime"
      name       = "runtime"
    },
  ]
  timeout   = "5m"
  tool_name = tencentcloud_ags_sandbox_tool.custom1.tool_name
}
```

Import

AGS sandbox instance can be imported using the id, e.g.

```
terraform import tencentcloud_ags_sandbox_instance.example instance-xxxxxxxx
```
