---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_global_domain"
sidebar_current: "docs-tencentcloud-resource-gaap_global_domain"
description: |-
  Provides a resource to create a gaap global domain
---

# tencentcloud_gaap_global_domain

Provides a resource to create a gaap global domain

## Example Usage

```hcl
resource "tencentcloud_gaap_global_domain" "global_domain" {
  project_id    = 0
  default_value = "xxxxxx.com"
  alias         = "demo"
  tags = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `default_value` - (Required, String) Domain name default entry.
* `project_id` - (Required, Int) Domain Name Project ID.
* `alias` - (Optional, String) alias.
* `status` - (Optional, String) Global domain statue. Available values: open and close, default is open.
* `tags` - (Optional, Map) Instance tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateGlobalDomain](https://cloud.tencent.com/document/api/608/79567)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeGlobalDomains](https://cloud.tencent.com/document/api/608/79562)
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [DescribeGlobalDomains](https://cloud.tencent.com/document/api/608/79562)
* Update: [DisableGlobalDomain](https://cloud.tencent.com/document/api/608/79561)
* Update: [EnableGlobalDomain](https://cloud.tencent.com/document/api/608/79560)
* Update: [ModifyGlobalDomainAttribute](https://cloud.tencent.com/document/api/608/79559)
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteGlobalDomain](https://cloud.tencent.com/document/api/608/79565)
* Delete: [DescribeGlobalDomains](https://cloud.tencent.com/document/api/608/79562)
* Delete: [DisableGlobalDomain](https://cloud.tencent.com/document/api/608/79561)

## Import

gaap global_domain can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_global_domain.global_domain ${projectId}#${domainId}
```

