Provides a resource to create a TEO inference service.

Example Usage

```hcl
resource "tencentcloud_teo_inference_service" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  name        = "tf-example-inference-service"
  listen_port = 8500
  description = "tf example inference service"

  containers {
    image_type = "TCR"

    tcr_repository_config {
      tcr_type    = "Enterprise"
      image       = "example.tencentcloudcr.com/ns/repo:v1"
      registry_id = "tcr-xxxxxxxx"
      region_name = "ap-guangzhou"
    }

    startup_command = "python app.py"

    environment_variables {
      key   = "MODEL_NAME"
      value = "example-model"
    }
  }

  resource_config {
    scaling_mode      = "Auto"
    hardware_spec_id  = "hs-xxxxxxxx"
    concurrency       = 1

    auto_scaling_config {
      min_instance_count = 1

      scaling_policies {
        policy_name = "tf-example-policy"
        policy_type = "ScheduledScaling"

        scheduled_scaling_policy {
          scheduled_actions {
            cron_expression   = "0 0 * * *"
            min_instance_count = 2
          }

          effective_range {
            effective_type = "LongTerm"
          }

          time_zone = "Asia/Shanghai"
        }
      }
    }
  }

  affinity_config {
    switch        = "On"
    affinity_mode = "SessionId"

    session_id_affinity_config {
      source      = "Header"
      header_name = "EO-Infer-Session-Id"
    }
  }

  request_paths = ["/v1/predict"]
}
```

Import

TEO inference service can be imported using the zoneId#serviceId, e.g.

```
terraform import tencentcloud_teo_inference_service.example zone-3fkff38fyw8s#is-xxxxxxxx
```