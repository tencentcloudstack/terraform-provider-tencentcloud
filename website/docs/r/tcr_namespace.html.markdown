---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_namespace"
sidebar_current: "docs-tencentcloud-resource-tcr_namespace"
description: |-
  Use this resource to create tcr namespace.
---

# tencentcloud_tcr_namespace

Use this resource to create tcr namespace.

## Example Usage

### Create a tcr namespace instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "example"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }

  tags = {
    "env"       = "production"
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the TCR instance.
* `name` - (Required, String, ForceNew) Name of the TCR namespace. Valid length is [2~30]. It can only contain lowercase letters, numbers and separators (`.`, `_`, `-`), and cannot start, end or continue with separators.
* `cve_whitelist_items` - (Optional, List) Vulnerability Whitelist.
* `is_auto_scan` - (Optional, Bool) Scanning level, `True` is automatic, `False` is manual. Default is `false`.
* `is_prevent_vul` - (Optional, Bool) Blocking switch, `True` is open, `False` is closed. Default is `false`.
* `is_public` - (Optional, Bool) Indicate that the namespace is public or not. Default is `false`.
* `severity` - (Optional, String) Block vulnerability level, currently only supports `low`, `medium`, `high`.
* `tags` - (Optional, Map) Tag key-value pairs for the TCR namespace.

The `cve_whitelist_items` object supports the following:

* `cve_id` - (Optional, String) Vulnerability Whitelist ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNamespace](https://cloud.tencent.com/document/api/1141/42729) - This API is used to create a namespace in an Enterprise Edition instance.
* Read: [DescribeNamespaces](https://cloud.tencent.com/document/api/1141/42765) - This API is used to query the namespace list or the information of the specified namespace.
* Update: [ModifyNamespace](https://cloud.tencent.com/document/api/1141/42727) - This API is used to update the information of a namespace. Currently, only the namespace access level can be modified.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteNamespace](https://cloud.tencent.com/document/api/1141/42728) - This API is used to delete a namespace.
* Delete: [DescribeNamespaces](https://cloud.tencent.com/document/api/1141/42765) - This API is used to query the namespace list or the information of the specified namespace.

## Import

tcr namespace can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_namespace.example tcr_instance_id#namespace_name
```

