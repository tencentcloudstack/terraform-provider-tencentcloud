---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_microservice"
sidebar_current: "docs-tencentcloud-resource-tsf_microservice"
description: |-
  Provides a resource to create a tsf microservice
---

# tencentcloud_tsf_microservice

Provides a resource to create a tsf microservice

## Example Usage

```hcl
resource "tencentcloud_tsf_microservice" "microservice" {
  namespace_id      = "namespace-vjlkzkgy"
  microservice_name = "test-microservice"
  microservice_desc = "desc-microservice"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `microservice_name` - (Required, String) Microservice name.
* `namespace_id` - (Required, String) Namespace ID.
* `microservice_desc` - (Optional, String) Microservice description information.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateMicroservice](https://www.tencentcloud.com/search/?wd=CreateMicroservice&type=document&lang=en)
* Create: [DescribeMicroservices](https://www.tencentcloud.com/search/?wd=DescribeMicroservices&type=document&lang=en)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Read: [DescribeMicroservices](https://www.tencentcloud.com/search/?wd=DescribeMicroservices&type=document&lang=en)
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
* Update: [ModifyMicroservice](https://www.tencentcloud.com/search/?wd=ModifyMicroservice&type=document&lang=en)
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Delete: [DeleteMicroservice](https://www.tencentcloud.com/search/?wd=DeleteMicroservice&type=document&lang=en)

## Import

tsf microservice can be imported using the namespaceId#microserviceId, e.g.

```
terraform import tencentcloud_tsf_microservice.microservice namespace-vjlkzkgy#ms-vjeb43lw
```

