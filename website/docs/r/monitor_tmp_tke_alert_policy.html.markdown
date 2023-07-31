---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_alert_policy"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_alert_policy"
description: |-
  Provides a resource to create a tke tmpAlertPolicy
---

# tencentcloud_monitor_tmp_tke_alert_policy

Provides a resource to create a tke tmpAlertPolicy

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
resource "tencentcloud_monitor_tmp_tke_alert_policy" "basic" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  alert_rule {
    name = "alert_rule-test"
    rules {
      name     = "rules-test"
      rule     = "(count(kube_node_status_allocatable_cpu_cores) by (cluster) -1)   / count(kube_node_status_allocatable_cpu_cores) by (cluster)"
      template = "The CPU requested by the Pod in the cluster {{ $labels.cluster }} is overloaded, and the current CPU application ratio is {{ $value | humanizePercentage }}"
      for      = "5m"
      labels {
        name  = "severity"
        value = "warning"
      }
    }
    notification {
      type    = "amp"
      enabled = true
      alert_manager {
        url = "xxx"
      }
    }
  }

  depends_on = [tencentcloud_monitor_tmp_tke_cluster_agent.foo]
}
```

## Argument Reference

The following arguments are supported:

* `alert_rule` - (Required, List) Alarm notification channels.
* `instance_id` - (Required, String) Instance Id.

The `alert_manager` object supports the following:

* `url` - (Required, String) Alertmanager url.
* `cluster_id` - (Optional, String) The ID of the cluster where the alertmanager is deployed. Note: This field may return null, indicating that a valid value could not be retrieved.
* `cluster_type` - (Optional, String) Alertmanager is deployed in the cluster type. Note: This field may return null, indicating that a valid value could not be retrieved.

The `alert_rule` object supports the following:

* `name` - (Required, String) Policy name.
* `rules` - (Required, List) A list of rules.
* `cluster_id` - (Optional, String) If the alarm policy is derived from the CRD resource definition of the user cluster, the ClusterId is the cluster ID to which it belongs.
* `id` - (Optional, String) Alarm policy ID. Note: This field may return null, indicating that a valid value could not be retrieved.
* `notification` - (Optional, List) Alarm channels, which may be returned using null in the template.
* `template_id` - (Optional, String) If the alarm is sent from a template, the TemplateId is the template id.
* `updated_at` - (Optional, String) Last modified time.

The `annotations` object supports the following:

* `name` - (Required, String) Name of map.
* `value` - (Required, String) Value of map.

The `labels` object supports the following:

* `name` - (Required, String) Name of map.
* `value` - (Required, String) Value of map.

The `notification` object supports the following:

* `enabled` - (Required, Bool) Whether it is enabled.
* `type` - (Required, String) The channel type, which defaults to amp, supports the following `amp`, `webhook`, `alertmanager`.
* `alert_manager` - (Optional, List) If Type is alertmanager, the field is required. Note: This field may return null, indicating that a valid value could not be retrieved..
* `notify_way` - (Optional, Set) Alarm notification method. At present, there are SMS, EMAIL, CALL, WECHAT methods.
* `phone_arrive_notice` - (Optional, Bool) Telephone alerts reach notifications.
* `phone_circle_interval` - (Optional, Int) Effective end timeTelephone alarm wheel interval. Units: Seconds.
* `phone_circle_times` - (Optional, Int) PhoneCircleTimes.
* `phone_inner_interval` - (Optional, Int) Telephone alarm wheel intervals. Units: Seconds.
* `phone_notify_order` - (Optional, Set) Telephone alarm sequence.
* `receiver_groups` - (Optional, Set) Alert Receiving Group (User Group).
* `repeat_interval` - (Optional, String) Convergence time.
* `time_range_end` - (Optional, String) Effective end time.
* `time_range_start` - (Optional, String) The time from which it takes effect.
* `web_hook` - (Optional, String) If Type is webhook, the field is required. Note: This field may return null, indicating that a valid value could not be retrieved.

The `rules` object supports the following:

* `for` - (Required, String) Time of duration.
* `labels` - (Required, List) Extra labels.
* `name` - (Required, String) Rule name.
* `rule` - (Required, String) Prometheus statement.
* `template` - (Required, String) Alert sending template.
* `annotations` - (Optional, List) Refer to annotations in prometheus rule.
* `describe` - (Optional, String) A description of the rule.
* `rule_state` - (Optional, Int) Alarm rule status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



