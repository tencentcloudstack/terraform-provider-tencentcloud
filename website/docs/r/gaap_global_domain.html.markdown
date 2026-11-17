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

* Create:
  - [CreateGlobalDomain](https://www.tencentcloud.com/search/?wd=CreateGlobalDomain&type=document&lang=en)
  - [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Read:
  - [DescribeGlobalDomains](https://www.tencentcloud.com/search/?wd=DescribeGlobalDomains&type=document&lang=en)
  - [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
* Update:
  - [DescribeGlobalDomains](https://www.tencentcloud.com/search/?wd=DescribeGlobalDomains&type=document&lang=en)
  - [DisableGlobalDomain](https://www.tencentcloud.com/search/?wd=DisableGlobalDomain&type=document&lang=en)
  - [EnableGlobalDomain](https://www.tencentcloud.com/search/?wd=EnableGlobalDomain&type=document&lang=en)
  - [ModifyGlobalDomainAttribute](https://www.tencentcloud.com/search/?wd=ModifyGlobalDomainAttribute&type=document&lang=en)
  - [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Delete:
  - [DeleteGlobalDomain](https://www.tencentcloud.com/search/?wd=DeleteGlobalDomain&type=document&lang=en)
  - [DescribeGlobalDomains](https://www.tencentcloud.com/search/?wd=DescribeGlobalDomains&type=document&lang=en)
  - [DisableGlobalDomain](https://www.tencentcloud.com/search/?wd=DisableGlobalDomain&type=document&lang=en)

## Import

gaap global_domain can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_global_domain.global_domain ${projectId}#${domainId}
```

