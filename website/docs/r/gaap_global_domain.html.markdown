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

* Create: [CreateGlobalDomain](https://cloud.tencent.com/document/api/608/79567) - 创建域名
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeGlobalDomains](https://cloud.tencent.com/document/api/608/79562) - 查询域名列表
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [DescribeGlobalDomains](https://cloud.tencent.com/document/api/608/79562) - 查询域名列表
* Update: [DisableGlobalDomain](https://cloud.tencent.com/document/api/608/79561) - 暂停域名解析
* Update: [EnableGlobalDomain](https://cloud.tencent.com/document/api/608/79560) - 开启域名解析
* Update: [ModifyGlobalDomainAttribute](https://cloud.tencent.com/document/api/608/79559) - 修改域名属性
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Delete: [DeleteGlobalDomain](https://cloud.tencent.com/document/api/608/79565) - 删除域名
* Delete: [DescribeGlobalDomains](https://cloud.tencent.com/document/api/608/79562) - 查询域名列表
* Delete: [DisableGlobalDomain](https://cloud.tencent.com/document/api/608/79561) - 暂停域名解析

## Import

gaap global_domain can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_global_domain.global_domain ${projectId}#${domainId}
```

