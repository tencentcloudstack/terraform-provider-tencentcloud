---
subcategory: "Monitor"
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
resource "tencentcloud_monitor_tmp_scrape_job" "tmpScrapeJob" {
  instance_id = "prom-dko9d0nu"
  agent_id    = "agent-6a7g40k2"
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

