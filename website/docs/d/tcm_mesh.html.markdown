---
subcategory: "TencentCloud ServiceMesh(TCM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcm_mesh"
sidebar_current: "docs-tencentcloud-datasource-tcm_mesh"
description: |-
  Use this data source to query detailed information of tcm mesh
---

# tencentcloud_tcm_mesh

Use this data source to query detailed information of tcm mesh

## Example Usage

```hcl
data "tencentcloud_tcm_mesh" "mesh" {
  mesh_id      = ["mesh-xxxxxx"]
  mesh_name    = ["KEEP_MASH"]
  tags         = ["key"]
  mesh_cluster = ["cls-xxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `mesh_cluster` - (Optional, Set: [`String`]) Mesh name.
* `mesh_id` - (Optional, Set: [`String`]) Mesh instance Id.
* `mesh_name` - (Optional, Set: [`String`]) Display name.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Set: [`String`]) tag key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `mesh_list` - The mesh information is queriedNote: This field may return null, indicating that a valid value is not available.
  * `config` - Mesh configuration.
    * `istio` - Istio configuration.
      * `disable_http_retry` - Disable http retry.
      * `disable_policy_checks` - Disable policy checks.
      * `enable_pilot_http` - Enable HTTP/1.0 support.
      * `outbound_traffic_policy` - Outbound traffic policy.
      * `smart_dns` - SmartDNS configuration.
        * `istio_meta_dns_auto_allocate` - Enable auto allocate address.
        * `istio_meta_dns_capture` - Enable dns proxy.
  * `display_name` - Mesh name.
  * `mesh_id` - Mesh instance Id.
  * `tag_list` - A list of associated tags.
    * `key` - Tag key.
    * `passthrough` - Passthrough to other related product.
    * `value` - Tag value.
  * `type` - Mesh type.  Value range:- `STANDALONE`: Standalone mesh- `HOSTED`: hosted the mesh.
  * `version` - Mesh version.


