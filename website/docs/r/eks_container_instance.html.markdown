---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eks_container_instance"
sidebar_current: "docs-tencentcloud-resource-eks_container_instance"
description: |-
  Provides an elastic kubernetes service container instance.
---

# tencentcloud_eks_container_instance

Provides an elastic kubernetes service container instance.

## Example Usage

```hcl
data "tencentcloud_security_groups" "group" {
}

data "tencentcloud_availability_zones_by_product" "zone" {
  product = "cvm"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/24"
  name       = "tf-test-eksci"
}

resource "tencentcloud_subnet" "sub" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
  cidr_block        = "10.0.0.0/24"
  name              = "sub"
  vpc_id            = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_cbs_storage" "cbs" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
  storage_name      = "cbs1"
  storage_size      = 10
  storage_type      = "CLOUD_PREMIUM"
}

resource "tencentcloud_eks_container_instance" "eci1" {
  name            = "foo"
  vpc_id          = tencentcloud_vpc.vpc.id
  subnet_id       = tencentcloud_subnet.sub.id
  cpu             = 2
  cpu_type        = "intel"
  restart_policy  = "Always"
  memory          = 4
  security_groups = [data.tencentcloud_security_groups.group.security_groups[0].security_group_id]
  cbs_volume {
    name    = "vol1"
    disk_id = tencentcloud_cbs_storage.cbs.id
  }
  container {
    name  = "redis1"
    image = "redis"
    liveness_probe {
      init_delay_seconds = 1
      timeout_seconds    = 3
      period_seconds     = 11
      success_threshold  = 1
      failure_threshold  = 3
      http_get_path      = "/"
      http_get_port      = 443
      http_get_scheme    = "HTTPS"
    }
    readiness_probe {
      init_delay_seconds = 1
      timeout_seconds    = 3
      period_seconds     = 10
      success_threshold  = 1
      failure_threshold  = 3
      tcp_socket_port    = 81
    }
  }
  container {
    name  = "nginx"
    image = "nginx"
  }
  init_container {
    name  = "alpine"
    image = "alpine:latest"
  }
}
```

## Argument Reference

The following arguments are supported:

* `container` - (Required, List) List of container.
* `cpu` - (Required, Float64) The number of CPU cores. Check https://intl.cloud.tencent.com/document/product/457/34057 for specification references.
* `memory` - (Required, Float64) Memory size. Check https://intl.cloud.tencent.com/document/product/457/34057 for specification references.
* `name` - (Required, String) Name of EKS container instance.
* `security_groups` - (Required, List: [`String`]) List of security group id.
* `subnet_id` - (Required, String) Subnet ID of container instance.
* `vpc_id` - (Required, String) VPC ID.
* `auto_create_eip` - (Optional, Bool) Indicates whether to create EIP instead of specify existing EIPs. Conflict with `existed_eip_ids`.
* `cam_role_name` - (Optional, String) CAM role name authorized to access.
* `cbs_volume` - (Optional, List) List of CBS volume.
* `cpu_type` - (Optional, String) Type of cpu, which can set to `intel` or `amd`. It also support backup list like `amd,intel` which indicates using `intel` when `amd` sold out.
* `dns_config_options` - (Optional, Map, ForceNew) Map of DNS config options.
* `dns_names_servers` - (Optional, List: [`String`], ForceNew) IP Addresses of DNS Servers.
* `dns_searches` - (Optional, List: [`String`], ForceNew) List of DNS Search Domain.
* `eip_delete_policy` - (Optional, Bool) Indicates weather the EIP release or not after instance deleted. Conflict with `existed_eip_ids`.
* `eip_max_bandwidth_out` - (Optional, Int) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bits per second). Conflict with `existed_eip_ids`.
* `eip_service_provider` - (Optional, String) EIP service provider. Default is `BGP`, values `CMCC`,`CTCC`,`CUCC` are available for whitelist customer. Conflict with `existed_eip_ids`.
* `existed_eip_ids` - (Optional, List: [`String`]) Existed EIP ID List which used to bind container instance. Conflict with `auto_create_eip` and auto create EIP options.
* `gpu_count` - (Optional, Int) Count of GPU. Check https://intl.cloud.tencent.com/document/product/457/34057 for specification references.
* `gpu_type` - (Optional, String) Type of GPU. Check https://intl.cloud.tencent.com/document/product/457/34057 for specification references.
* `image_registry_credential` - (Optional, List) List of credentials which pull from image registry.
* `init_container` - (Optional, List) List of initialized container.
* `nfs_volume` - (Optional, List) List of NFS volume.
* `restart_policy` - (Optional, String) Container instance restart policy. Available values: `Always`, `Never`, `OnFailure`.

The `cbs_volume` object supports the following:

* `disk_id` - (Required, String) ID of CBS.
* `name` - (Required, String) Name of CBS volume.

The `container` object supports the following:

* `image` - (Required, String) Image of Container.
* `name` - (Required, String) Name of Container.
* `args` - (Optional, List) Container launch argument list.
* `commands` - (Optional, List) Container launch command list.
* `cpu` - (Optional, Float64) Number of cpu core of container.
* `env_vars` - (Optional, Map) Map of environment variables of container OS.
* `liveness_probe` - (Optional, List) Configuration block of LivenessProbe.
* `memory` - (Optional, Float64) Memory size of container.
* `readiness_probe` - (Optional, List) Configuration block of ReadinessProbe.
* `volume_mount` - (Optional, List) List of volume mount informations.
* `working_dir` - (Optional, String) Container working directory.

The `image_registry_credential` object supports the following:

* `name` - (Optional, String) Name of credential.
* `password` - (Optional, String) Password.
* `server` - (Optional, String) Address of image registry.
* `username` - (Optional, String) Username.

The `init_container` object supports the following:

* `image` - (Required, String) Image of Container.
* `name` - (Required, String) Name of Container.
* `args` - (Optional, List) Container launch argument list.
* `commands` - (Optional, List) Container launch command list.
* `cpu` - (Optional, Float64) Number of cpu core of container.
* `env_vars` - (Optional, Map) Map of environment variables of container OS.
* `memory` - (Optional, Float64) Memory size of container.
* `volume_mount` - (Optional, List) List of volume mount informations.
* `working_dir` - (Optional, String) Container working directory.

The `liveness_probe` object supports the following:

* `exec_commands` - (Optional, List) List of execution commands.
* `failure_threshold` - (Optional, Int) Minimum consecutive failures for the probe to be considered failed after having succeeded.Default: `3`. Minimum value is `1`.
* `http_get_path` - (Optional, String) HttpGet detection path.
* `http_get_port` - (Optional, Int) HttpGet detection port.
* `http_get_scheme` - (Optional, String) HttpGet detection scheme. Available values: `HTTP`, `HTTPS`.
* `init_delay_seconds` - (Optional, Int) Number of seconds after the container has started before probes are initiated.
* `period_seconds` - (Optional, Int) How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is `1`.
* `success_threshold` - (Optional, Int) Minimum consecutive successes for the probe to be considered successful after having failed. Default: `1`. Must be 1 for liveness. Minimum value is `1`.
* `tcp_socket_port` - (Optional, Int) TCP Socket detection port.
* `timeout_seconds` - (Optional, Int) Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is `1`.

The `nfs_volume` object supports the following:

* `name` - (Required, String) Name of NFS volume.
* `path` - (Required, String) NFS volume path.
* `server` - (Required, String) NFS server address.
* `read_only` - (Optional, Bool) Indicates whether the volume is read only. Default is `false`.

The `readiness_probe` object supports the following:

* `exec_commands` - (Optional, List) List of execution commands.
* `failure_threshold` - (Optional, Int) Minimum consecutive failures for the probe to be considered failed after having succeeded.Default: `3`. Minimum value is `1`.
* `http_get_path` - (Optional, String) HttpGet detection path.
* `http_get_port` - (Optional, Int) HttpGet detection port.
* `http_get_scheme` - (Optional, String) HttpGet detection scheme. Available values: `HTTP`, `HTTPS`.
* `init_delay_seconds` - (Optional, Int) Number of seconds after the container has started before probes are initiated.
* `period_seconds` - (Optional, Int) How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is `1`.
* `success_threshold` - (Optional, Int) Minimum consecutive successes for the probe to be considered successful after having failed. Default: `1`. Must be 1 for liveness. Minimum value is `1`.
* `tcp_socket_port` - (Optional, Int) TCP Socket detection port.
* `timeout_seconds` - (Optional, Int) Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is `1`.

The `volume_mount` object supports the following:

* `name` - (Required, String) Volume name.
* `path` - (Required, String) Volume mount path.
* `mount_propagation` - (Optional, String) Volume mount propagation.
* `read_only` - (Optional, Bool) Whether the volume is read-only.
* `sub_path_expr` - (Optional, String) Volume mount sub-path expression.
* `sub_path` - (Optional, String) Volume mount sub-path.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `auto_create_eip_id` - ID of EIP which create automatically.
* `created_time` - Container instance creation time.
* `eip_address` - EIP address.
* `private_ip` - Private IP address.
* `status` - Container instance status.


## Import

```
terraform import tencentcloud_eks_container_instance.foo container-instance-id
```

