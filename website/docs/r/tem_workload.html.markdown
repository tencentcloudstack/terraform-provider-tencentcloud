---
subcategory: "TencentCloud Elastic Microservice(TEM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_workload"
sidebar_current: "docs-tencentcloud-resource-tem_workload"
description: |-
  Provides a resource to create a tem workload
---

# tencentcloud_tem_workload

Provides a resource to create a tem workload

## Example Usage

```hcl
resource "tencentcloud_tem_workload" "workload" {
  application_id = "app-3j29aa2p"
  environment_id = "en-853mggjm"
  deploy_version = "hello-world"
  deploy_mode    = "IMAGE"
  img_repo       = "tem_demo/tem_demo"
  init_pod_num   = 1
  cpu_spec       = 1
  memory_spec    = 1
  liveness {
    type                  = "HttpGet"
    protocol              = "HTTP"
    path                  = "/"
    port                  = 8201
    initial_delay_seconds = 0
    timeout_seconds       = 1
    period_seconds        = 10

  }
  readiness {
    type                  = "HttpGet"
    protocol              = "HTTP"
    path                  = "/"
    port                  = 8201
    initial_delay_seconds = 0
    timeout_seconds       = 1
    period_seconds        = 10

  }
  startup_probe {
    type                  = "HttpGet"
    protocol              = "HTTP"
    path                  = "/"
    port                  = 8201
    initial_delay_seconds = 0
    timeout_seconds       = 1
    period_seconds        = 10

  }
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String, ForceNew) application ID.
* `cpu_spec` - (Required, Float64) cpu.
* `deploy_mode` - (Required, String) deploy mode, support IMAGE.
* `deploy_version` - (Required, String) deploy version.
* `environment_id` - (Required, String, ForceNew) environment ID.
* `img_repo` - (Required, String) repository name.
* `init_pod_num` - (Required, Int) initial pod number.
* `memory_spec` - (Required, Float64) mem.
* `deploy_strategy_conf` - (Optional, List) deploy strategy.
* `env_conf` - (Optional, List) .
* `liveness` - (Optional, List) liveness config.
* `post_start` - (Optional, String) mem.
* `pre_stop` - (Optional, String) mem.
* `readiness` - (Optional, List) .
* `repo_server` - (Optional, String) repo server addr when deploy by image.
* `repo_type` - (Optional, Int) repo type when deploy: 0: tcr personal; 1: tcr enterprise; 2: public repository; 3: tem host tcr; 4: demo repo.
* `security_group_ids` - (Optional, Set: [`String`]) security groups.
* `startup_probe` - (Optional, List) .
* `storage_confs` - (Optional, List) storage configuration.
* `storage_mount_confs` - (Optional, List) storage mount configuration.
* `tcr_instance_id` - (Optional, String) tcr instance id when deploy by image.

The `deploy_strategy_conf` object supports the following:

* `deploy_strategy_type` - (Required, Int) strategy type, 0 means auto, 1 means manual, 2 means manual with beta batch.
* `total_batch_count` - (Required, Int) total batch number.
* `batch_interval` - (Optional, Int) interval between batches.
* `beta_batch_num` - (Optional, Int) beta batch number.
* `force` - (Optional, Bool) force update.
* `min_available` - (Optional, Int) minimal available instances duration deployment.

The `env_conf` object supports the following:

* `key` - (Required, String) env key.
* `value` - (Required, String) env value.
* `config` - (Optional, String) referenced config name when type=referenced.
* `secret` - (Optional, String) referenced secret name when type=referenced.
* `type` - (Optional, String) env type, support default, referenced.

The `liveness` object supports the following:

* `type` - (Required, String) check type, support HttpGet, TcpSocket and Exec.
* `exec` - (Optional, String) script.
* `initial_delay_seconds` - (Optional, Int) initial delay seconds for liveness check.
* `path` - (Optional, String) path.
* `period_seconds` - (Optional, Int) period seconds for liveness check.
* `port` - (Optional, Int) liveness check port.
* `protocol` - (Optional, String) protocol.
* `timeout_seconds` - (Optional, Int) timeout seconds for liveness check.

The `readiness` object supports the following:

* `type` - (Required, String) check type, support HttpGet, TcpSocket and Exec.
* `exec` - (Optional, String) script.
* `initial_delay_seconds` - (Optional, Int) initial delay seconds for readiness check.
* `path` - (Optional, String) path.
* `period_seconds` - (Optional, Int) period seconds for readiness check.
* `port` - (Optional, Int) readiness check port.
* `protocol` - (Optional, String) protocol.
* `timeout_seconds` - (Optional, Int) timeout seconds for readiness check.

The `startup_probe` object supports the following:

* `type` - (Required, String) check type, support HttpGet, TcpSocket and Exec.
* `exec` - (Optional, String) script.
* `initial_delay_seconds` - (Optional, Int) initial delay seconds for startup check.
* `path` - (Optional, String) path.
* `period_seconds` - (Optional, Int) period seconds for startup check.
* `port` - (Optional, Int) startup check port.
* `protocol` - (Optional, String) protocol.
* `timeout_seconds` - (Optional, Int) timeout seconds for startup check.

The `storage_confs` object supports the following:

* `storage_vol_ip` - (Required, String) volume ip.
* `storage_vol_name` - (Required, String) volume name.
* `storage_vol_path` - (Required, String) volume path.

The `storage_mount_confs` object supports the following:

* `mount_path` - (Required, String) mount path.
* `volume_name` - (Required, String) volume name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tem workload can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_workload.workload envirnomentId#applicationId
```

