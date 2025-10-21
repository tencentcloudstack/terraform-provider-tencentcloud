---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_deploy_yarn_operation"
sidebar_current: "docs-tencentcloud-resource-emr_deploy_yarn_operation"
description: |-
  Provides a resource to deploy a emr yarn
---

# tencentcloud_emr_deploy_yarn_operation

Provides a resource to deploy a emr yarn

## Example Usage

```hcl
resource "tencentcloud_emr_deploy_yarn_operation" "emr_yarn" {
  instance_id = "emr-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) EMR Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



