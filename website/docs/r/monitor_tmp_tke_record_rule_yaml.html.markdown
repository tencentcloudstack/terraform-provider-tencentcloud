---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_record_rule_yaml"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_record_rule_yaml"
description: |-
  Provides a resource to create a tke tmpRecordRule
---

# tencentcloud_monitor_tmp_tke_record_rule_yaml

Provides a resource to create a tke tmpRecordRule

## Example Usage

```hcl
# create monitor
variable "cluster_type" {
  default = "tke"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}

# create record rule
resource "tencentcloud_monitor_tmp_tke_record_rule_yaml" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  content     = <<-EOT
        apiVersion: monitoring.coreos.com/v1
        kind: PrometheusRule
        metadata:
          name: example-record
        spec:
          groups:
            - name: kube-apiserver.rules
              rules:
                - expr: sum(metrics_test)
                  labels:
                    verb: read
                  record: 'apiserver_request:burnrate1d'
    EOT
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Contents of record rules in yaml format.
* `instance_id` - (Required, String) Instance Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - An ID identify the cluster, like cls-xxxxxx.
* `name` - Name of the instance.
* `template_id` - Used for the argument, if the configuration comes to the template, the template id.
* `update_time` - Last modified time of record rule.


