---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_xml_config"
sidebar_current: "docs-tencentcloud-resource-clickhouse_xml_config"
description: |-
  Provides a resource to create a clickhouse xml_config
---

# tencentcloud_clickhouse_xml_config

Provides a resource to create a clickhouse xml_config

## Example Usage

```hcl
resource "tencentcloud_clickhouse_xml_config" "xml_config" {
  instance_id = "cdwch-datuhk3z"
  modify_conf_context {
    file_name      = "metrika.xml"
    new_conf_value = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHlhbmRleD4KICAgIDx6b29rZWVwZXItc2VydmVycz4KICAgIDwvem9va2VlcGVyLXNlcnZlcnM+CjwveWFuZGV4Pgo="
    file_path      = "/etc/clickhouse-server"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Cluster ID.
* `modify_conf_context` - (Required, List) Configuration file modification information.

The `modify_conf_context` object supports the following:

* `file_name` - (Required, String) Configuration file name.
* `new_conf_value` - (Required, String) New content of configuration file, base64 encoded.
* `file_path` - (Optional, String) Path to save configuration file.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clickhouse xml_config can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_xml_config.xml_config cdwch-datuhk3z#metrika.xml
```

