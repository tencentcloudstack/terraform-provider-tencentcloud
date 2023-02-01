---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_namespace"
sidebar_current: "docs-tencentcloud-resource-tsf_namespace"
description: |-
  Provides a resource to create a tsf namespace
---

# tencentcloud_tsf_namespace

Provides a resource to create a tsf namespace

## Example Usage

```hcl
resource "tencentcloud_tsf_namespace" "namespace" {
  namespace_name = "namespace-name"
  # cluster_id = "cls-xxxx"
  namespace_desc = "namespace desc"
  # namespace_resource_type = ""
  namespace_type = "DEF"
  # namespace_id = ""
  is_ha_enable = "0"
  # program_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `namespace_name` - (Required, String) namespace name.
* `cluster_id` - (Optional, String) cluster ID.
* `is_ha_enable` - (Optional, String) whether to enable high availability.
* `namespace_desc` - (Optional, String) namespace description.
* `namespace_id` - (Optional, String) Namespace ID.
* `namespace_resource_type` - (Optional, String) namespace resource type (default is DEF).
* `namespace_type` - (Optional, String) Whether it is a global namespace (the default is DEF, which means a common namespace; GLOBAL means a global namespace).
* `program_id_list` - (Optional, Set: [`String`]) Program id list.
* `program_id` - (Optional, String) ID of the dataset to be bound.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - creation time.
* `delete_flag` - Delete ID.
* `is_default` - default namespace.
* `kube_inject_enable` - KubeInjectEnable value.
* `namespace_code` - Namespace encoding.
* `namespace_status` - namespace status.
* `update_time` - update time.


## Import

tsf namespace can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_namespace.namespace namespace_id
```

