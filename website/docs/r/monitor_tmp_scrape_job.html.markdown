---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_scrape_job"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_scrape_job"
description: |-
  Provides a resource to create a monitor tmpScrapeJob
---

# tencentcloud_monitor_tmp_scrape_job

Provides a resource to create a monitor tmpScrapeJob

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

resource "tencentcloud_monitor_tmp_cvm_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  name        = "tf-agent"
}

resource "tencentcloud_monitor_tmp_scrape_job" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  agent_id    = tencentcloud_monitor_tmp_cvm_agent.foo.agent_id
  config      = <<-EOT
job_name: demo-config
honor_timestamps: true
metrics_path: /metrics
scheme: https
EOT
}
```

## Argument Reference

The following arguments are supported:

* `agent_id` - (Required, String) Agent id.
* `instance_id` - (Required, String) Instance id.
* `config` - (Optional, String) Job content.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmpScrapeJob can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_scrape_job.tmpScrapeJob tmpScrapeJob_id
```

