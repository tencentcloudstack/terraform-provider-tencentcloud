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

* Create: [CreateInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79579)
* Create: [DescribeInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79573)
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Read: [DescribeInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79573)
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319)
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Delete: [DeleteInstanceCustomizedDomain](https://cloud.tencent.com/document/api/1141/79576)

## Import

tcr customized domain can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_customized_domain.example tcr-fjvvsfdh#www.demo.com
```

