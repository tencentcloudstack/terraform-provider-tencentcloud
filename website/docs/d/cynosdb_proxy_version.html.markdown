---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_proxy_version"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_proxy_version"
description: |-
  Use this data source to query detailed information of cynosdb proxy_version
---

# tencentcloud_cynosdb_proxy_version

Use this data source to query detailed information of cynosdb proxy_version

## Example Usage

```hcl
data "tencentcloud_cynosdb_proxy_version" "proxy_version" {
  cluster_id     = "cynosdbmysql-bws8h88b"
  proxy_group_id = "cynosdbmysql-proxy-l6zf9t30"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `proxy_group_id` - (Optional, String) Database Agent Group ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `current_proxy_version` - Current proxy version number note: This field may return null, indicating that a valid value cannot be obtained.
* `support_proxy_versions` - Supported Database Agent Version Collection Note: This field may return null, indicating that a valid value cannot be obtained.


