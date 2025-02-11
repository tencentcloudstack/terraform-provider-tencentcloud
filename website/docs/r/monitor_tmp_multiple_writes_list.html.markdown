---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_multiple_writes_list"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_multiple_writes_list"
description: |-
  Provides a resource to create a monitor multiple writes list
---

# tencentcloud_monitor_tmp_multiple_writes_list

Provides a resource to create a monitor multiple writes list

~> **NOTE:** For the same instance of prometheus, resource `tencentcloud_monitor_tmp_multiple_writes` and resource `tencentcloud_monitor_tmp_multiple_writes_list` cannot be used simultaneously. Resource `tencentcloud_monitor_tmp_multiple_writes` will been deprecated in version v1.81.166, Please use resource `tencentcloud_monitor_tmp_multiple_writes_list` instead.

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_multiple_writes_list" "example" {
  instance_id = "prom-gzg3f1em"

  remote_writes {
    url = "http://172.16.0.111:9090/api/v1/prom/write"
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
    headers {
      key   = "Key"
      value = "Value"
    }
  }

  remote_writes {
    url                = "http://172.16.0.111:8080/api/v1/prom/write"
    url_relabel_config = "# 添加 label\n#- target_label: key\n#  replacement: value\n# 丢弃指标\n#- source_labels: [__name__]\n#  regex: kubelet_.+;\n#  action: drop"
    headers {
      key   = "Key"
      value = "Value"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `remote_writes` - (Optional, Set) Data multiple write configuration.

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

monitor multiple writes list can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_multiple_writes_list.example prom-gzg3f1em
```

