---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_backup_stream"
sidebar_current: "docs-tencentcloud-resource-css_backup_stream"
description: |-
  Provides a resource to create a css backup_stream
---

# tencentcloud_css_backup_stream

Provides a resource to create a css backup_stream

~> **NOTE:** This resource is only valid when the push stream. When the push stream ends, it will be deleted.

## Example Usage

```hcl
resource "tencentcloud_css_backup_stream" "backup_stream" {
  push_domain_name  = "177154.push.tlivecloud.com"
  app_name          = "live"
  stream_name       = "1308919341_test"
  upstream_sequence = "2209501773993286139"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String, ForceNew) App name.
* `push_domain_name` - (Required, String, ForceNew) Push domain.
* `stream_name` - (Required, String, ForceNew) Stream id.
* `upstream_sequence` - (Required, String) Sequence.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css backup_stream can be imported using the id, e.g.

```
terraform import tencentcloud_css_backup_stream.backup_stream pushDomainName#appName#streamName
```

