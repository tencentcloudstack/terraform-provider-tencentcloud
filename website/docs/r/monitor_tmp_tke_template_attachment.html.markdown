---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_template_attachment"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_template_attachment"
description: |-
  Provides a resource to create a tmp tke template attachment
---

# tencentcloud_monitor_tmp_tke_template_attachment

Provides a resource to create a tmp tke template attachment

## Example Usage

```hcl
# create tke
variable "default_instance_type" {
  default = "SA1.MEDIUM2"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    password = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_second
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    cam_role_name = "CVM_QcsRole"
    password      = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}

# create monitor
variable "zone" {
  default = "ap-guangzhou"
}

variable "cluster_type" {
  default = "tke"
}

resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = local.first_vpc_id
  subnet_id           = local.first_subnet_id
  data_retention_time = 30
  zone                = var.availability_zone_second
  tags = {
    "createdBy" = "terraform"
  }
}

# tmp tke bind
resource "tencentcloud_monitor_tmp_tke_cluster_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id

  agents {
    region          = var.zone
    cluster_type    = var.cluster_type
    cluster_id      = tencentcloud_kubernetes_cluster.example.id
    enable_external = false
  }
}

# create monitor template
resource "tencentcloud_monitor_tmp_tke_template" "foo" {
  template {
    name     = "tf-template"
    level    = "cluster"
    describe = "template"
    service_monitors {
      name   = "tf-ServiceMonitor"
      config = <<-EOT
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: example-service-monitor
  namespace: monitoring
  labels:
    k8s-app: example-service
spec:
  selector:
    matchLabels:
      k8s-app: example-service
  namespaceSelector:
    matchNames:
      - default
  endpoints:
  - port: http-metrics
    interval: 30s
    path: /metrics
    scheme: http
    bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    tlsConfig:
      insecureSkipVerify: true
      EOT
    }

    pod_monitors {
      name   = "tf-PodMonitors"
      config = <<-EOT
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: example-pod-monitor
  namespace: monitoring
  labels:
    k8s-app: example-pod
spec:
  selector:
    matchLabels:
      k8s-app: example-pod
  namespaceSelector:
    matchNames:
      - default
  podMetricsEndpoints:
  - port: http-metrics
    interval: 30s
    path: /metrics
    scheme: http
    bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    tlsConfig:
      insecureSkipVerify: true
EOT
    }

    pod_monitors {
      name   = "tf-RawJobs"
      config = <<-EOT
scrape_configs:
  - job_name: 'example-job'
    scrape_interval: 30s
    static_configs:
      - targets: ['example-service.default.svc.cluster.local:8080']
    metrics_path: /metrics
    scheme: http
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    tls_config:
      insecure_skip_verify: true
EOT
    }
  }
}

resource "tencentcloud_monitor_tmp_tke_template_attachment" "temp_attachment" {
  template_id = tencentcloud_monitor_tmp_tke_template.foo.id

  targets {
    cluster_type = var.cluster_type
    cluster_id   = tencentcloud_kubernetes_cluster.example.id
    region       = var.zone
    instance_id  = tencentcloud_monitor_tmp_instance.foo.id
  }

  depends_on = [tencentcloud_monitor_tmp_tke_cluster_agent.foo]
}
```

## Argument Reference

The following arguments are supported:

* `targets` - (Required, List, ForceNew) Sync target details.
* `template_id` - (Required, String, ForceNew) The ID of the template, which is used for the outgoing reference.

The `targets` object supports the following:

* `instance_id` - (Required, String) instance id.
* `region` - (Required, String) target area.
* `cluster_id` - (Optional, String) ID of the cluster.
* `cluster_name` - (Optional, String) Name the cluster.
* `cluster_type` - (Optional, String) Cluster type.
* `instance_name` - (Optional, String) Name of the prometheus instance.
* `sync_time` - (Optional, String) Last sync template time.
* `version` - (Optional, String) Template version currently in use.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



