---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_service_account"
sidebar_current: "docs-tencentcloud-resource-tcr_service_account"
description: |-
  Provides a resource to create a TCR service account.
---

# tencentcloud_tcr_service_account

Provides a resource to create a TCR service account.

## Example Usage

### Create custom account with specified duration days

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "standard"
  delete_bucket = true
  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf-example"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
}

resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf-example"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  duration    = 10
  disable     = false
  password    = "Password123"
  tags = {
    createdBy = "Terraform"
  }
}
```

### With specified expiration time

```hcl
resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf-example"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  expires_at  = 1676897989000 //time stamp
  disable     = false
  tags = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Service account name.
* `permissions` - (Required, List) strategy list.
* `registry_id` - (Required, String) instance id.
* `description` - (Optional, String) Service account description.
* `disable` - (Optional, Bool) whether to disable Service accounts.
* `duration` - (Optional, Int) expiration date (unit: day), calculated from the current time, priority is higher than ExpiresAt Service account description.
* `expires_at` - (Optional, Int) Service account expiration time (time stamp, unit: milliseconds).
* `password` - (Optional, String) Password of the service account.
* `tags` - (Optional, Map) Tag description list.

The `permissions` object supports the following:

* `actions` - (Required, Set) Actions, currently support: `tcr:PushRepository`, `tcr:PullRepository`, `tcr:CreateRepository`, `tcr:CreateHelmChart`, `tcr:DescribeHelmCharts`. Note: This field may return null, indicating that no valid value can be obtained.
* `resource` - (Required, String) resource path, currently only supports Namespace. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateServiceAccount](https://cloud.tencent.com/document/api/1141/94608) - This API is used to create a service level account.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Create: [ModifyServiceAccountPassword](https://cloud.tencent.com/document/api/1141/103857) - This API is used to update the password for a service level account.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeServiceAccounts](https://cloud.tencent.com/document/api/1141/94606) - This API is used to query service accounts.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Update: [ModifyServiceAccount](https://cloud.tencent.com/document/api/1141/94605) - This API is used to update a service account.
* Update: [ModifyServiceAccountPassword](https://cloud.tencent.com/document/api/1141/103857) - This API is used to update the password for a service level account.
* Delete: [DeleteServiceAccount](https://cloud.tencent.com/document/api/1141/94607) - This API is used to delete a service account.

## Import

TCR service account can be imported using the registryId#accountName, e.g.

```
terraform import tencentcloud_tcr_service_account.example tcr-ixgt2l0z#tf-example
```

