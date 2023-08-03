---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_recording_rule"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_recording_rule"
description: |-
  Provides a resource to create a monitor tmp recordingRule
---

# tencentcloud_monitor_tmp_recording_rule

Provides a resource to create a monitor tmp recordingRule

## Example Usage

```hcl
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

resource "tencentcloud_monitor_tmp_recording_rule" "recordingRule" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  rule_state  = 2
  name        = "tf-recording-rule"
  group       = <<EOF
---
name: example-test
rules:
  - record: job:http_inprogress_requests:sum
    expr: sum by (job) (http_inprogress_requests)
EOF

}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String) Recording rule group.
* `instance_id` - (Required, String) Instance id.
* `name` - (Required, String) Recording rule name.
* `rule_state` - (Optional, Int) Rule state.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor recordingRule can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_recording_rule.recordingRule instanceId#recordingRule_id
```

