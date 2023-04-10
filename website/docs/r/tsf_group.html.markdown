---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_group"
sidebar_current: "docs-tencentcloud-resource-tsf_group"
description: |-
  Provides a resource to create a tsf group
---

# tencentcloud_tsf_group

Provides a resource to create a tsf group

## Example Usage

```hcl
resource "tencentcloud_tsf_group" "group" {
  application_id = "application-xxx"
  namespace_id   = "namespace-aemrxxx"
  group_name     = "terraform-test"
  cluster_id     = "cluster-vwgjxxxx"
  group_desc     = "terraform desc"
  alias          = "terraform test"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String) The application ID to which the group belongs.
* `cluster_id` - (Required, String) Cluster ID.
* `group_name` - (Required, String) Group name field, length 1~60, beginning with a letter or underscore, can contain alphanumeric underscore.
* `namespace_id` - (Required, String) ID of the namespace to which the group belongs.
* `alias` - (Optional, String) Deployment Group Notes.
* `group_desc` - (Optional, String) Group description.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `group_resource_type` - Deployment Group Resource Type.


## Import

tsf group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_group.group group-axxx
```

