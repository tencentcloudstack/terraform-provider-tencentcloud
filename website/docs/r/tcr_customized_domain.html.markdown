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
  name          = "tf-example-tcr"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_customized_domain" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  domain_name    = "www.test.com"
  certificate_id = "your_cert_id"
  tags = {
    "createdBy" = "terraform"
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



## Import

tcr customized_domain can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_customized_domain.customized_domain customized_domain_id
```

