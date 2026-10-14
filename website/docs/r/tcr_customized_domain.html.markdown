---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_customized_domain"
sidebar_current: "docs-tencentcloud-resource-tcr_customized_domain"
description: |-
  Provides a resource to create a tcr customized domain
---

# tencentcloud_tcr_customized_domain

Provides a resource to create a tcr customized domain

## Example Usage

### Create a tcr customized domain

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "premium"
  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_tcr_customized_domain" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  domain_name    = "www.demo.com"
  certificate_id = "your_cert_id"
  tags = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) certificate id.
* `domain_name` - (Required, String, ForceNew) custom domain name.
* `registry_id` - (Required, String, ForceNew) instance id.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79579) - This API is used to create a custom domain name.
* Create: [DescribeInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79573) - This API is used to query the list of custom domain names of an instance.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79573) - This API is used to query the list of custom domain names of an instance.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79576) - This API is used to delete a custom domain name.

## Import

tcr customized domain can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_customized_domain.example tcr-fjvvsfdh#www.demo.com
```

