---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_global_notification"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_global_notification"
description: |-
  Provides a resource to create a tmp tke global notification
---

# tencentcloud_monitor_tmp_tke_global_notification

Provides a resource to create a tmp tke global notification

## Example Usage

```hcl
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

# create record rule
resource "tencentcloud_monitor_tmp_tke_global_notification" "basic" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  notification {
    enabled = true
    type    = "webhook"
    alert_manager {
      cluster_id   = ""
      cluster_type = ""
      url          = ""
    }
    web_hook              = ""
    repeat_interval       = "5m"
    time_range_start      = "00:00:00"
    time_range_end        = "23:59:59"
    notify_way            = ["SMS", "EMAIL"]
    receiver_groups       = []
    phone_notify_order    = []
    phone_circle_times    = 0
    phone_inner_interval  = 0
    phone_circle_interval = 0
    phone_arrive_notice   = false
  }

  depends_on = [tencentcloud_monitor_tmp_tke_cluster_agent.foo]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance Id.
* `notification` - (Required, List) Alarm notification channels.

The `alert_manager` object supports the following:

* `url` - (Required, String) Alert manager url.
* `cluster_id` - (Optional, String) Cluster id.
* `cluster_type` - (Optional, String) Cluster type.

The `notification` object supports the following:

* `enabled` - (Required, Bool) Alarm notification switch.
* `type` - (Required, String) Alarm notification type, Valid values: `amp`, `webhook`, `alertmanager`.
* `alert_manager` - (Optional, List) Alert manager, if Type is `alertmanager`, this field is required.
* `notify_way` - (Optional, Set) Alarm notification method, Valid values: `SMS`, `EMAIL`, `CALL`, `WECHAT`.
* `phone_arrive_notice` - (Optional, Bool) Phone Alarm Reach Notification, NotifyWay is `CALL`, and this parameter is used.
* `phone_circle_interval` - (Optional, Int) Telephone alarm off-wheel interval, NotifyWay is `CALL`, and this parameter is used.
* `phone_circle_times` - (Optional, Int) Number of phone alerts (user group), NotifyWay is `CALL`, and this parameter is used.
* `phone_inner_interval` - (Optional, Int) Interval between telephone alarm rounds, NotifyWay is `CALL`, and this parameter is used.
* `phone_notify_order` - (Optional, Set) Phone alert sequence, NotifyWay is `CALL`, and this parameter is used.
* `receiver_groups` - (Optional, Set) Alarm receiving group(user group).
* `repeat_interval` - (Optional, String) Convergence time.
* `time_range_end` - (Optional, String) Effective end time.
* `time_range_start` - (Optional, String) Effective start time.
* `web_hook` - (Optional, String) Web hook, if Type is `webhook`, this field is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



