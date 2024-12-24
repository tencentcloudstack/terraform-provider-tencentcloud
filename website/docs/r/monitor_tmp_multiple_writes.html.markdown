---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_multiple_writes"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_multiple_writes"
description: |-
  Provides a resource to create a monitor monitor_tmp_multiple_writes
---

# tencentcloud_monitor_tmp_multiple_writes

Provides a resource to create a monitor monitor_tmp_multiple_writes

~> **NOTE:** When using `<<EOT`, please pay attention to spaces, line breaks, indentation, etc.

~> **NOTE:** When importing, the unique id is separated by the first `#`.

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_multiple_writes" "monitor_tmp_multiple_writes" {
  instance_id = "prom-l9cl1ptk"

  remote_writes {
    label          = null
    max_block_size = null
    url            = "http://172.16.0.111:9090/api/v1/prom/write"
    url_relabel_config = trimspace(<<-EOT
            # 添加 label
            # - target_label: key
            #  replacement: value
            # 丢弃指标
            #- source_labels: [__name__]
            #  regex: kubelet_.+;
            #  action: drop
        EOT
    )
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `remote_writes` - (Optional, List) Data multiple write configuration.

The `basic_auth` object of `remote_writes` supports the following:

* `password` - (Optional, String) Password.
* `user_name` - (Optional, String) User name.

The `headers` object of `remote_writes` supports the following:

* `key` - (Required, String) HTTP header key.
* `value` - (Optional, String) HTTP header value.

The `remote_writes` object supports the following:

* `url` - (Required, String) Data multiple write url.
* `basic_auth` - (Optional, List) Authentication information.
* `headers` - (Optional, List) HTTP additional headers.
* `label` - (Optional, String) Label.
* `max_block_size` - (Optional, String) Maximum block.
* `url_relabel_config` - (Optional, String) RelabelConfig.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor monitor_tmp_multiple_writes can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes prom-l9cl1ptk#http://172.16.0.111:9090/api/v1/prom/write
```

