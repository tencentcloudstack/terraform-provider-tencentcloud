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
resource "tencentcloud_monitor_tmp_tke_record_rule_yaml" "foo" {
  instance_id = ""
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


