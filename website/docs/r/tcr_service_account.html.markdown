---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_service_account"
sidebar_current: "docs-tencentcloud-resource-tcr_service_account"
description: |-
  Provides a resource to create a tcr service account.
---

# tencentcloud_tcr_service_account

Provides a resource to create a tcr service account.

## Example Usage

### Create custom account with specified duration days

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr-instance"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_test_tcr_namespace"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "tf_example_cve_id"
  }
}

resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example_account"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  duration    = 10
  disable     = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

### With specified expiration time

```hcl
resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example_account"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  expires_at  = 1676897989000 //time stamp
  disable     = false
  tags = {
    "createdBy" = "terraform"
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
* `tags` - (Optional, Map) Tag description list.

The `permissions` object supports the following:

* `actions` - (Required, Set) Actions, currently only support: tcr:PushRepository, tcr:PullRepository. Note: This field may return null, indicating that no valid value can be obtained.
* `resource` - (Required, String) resource path, currently only supports Namespace. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `password` - Password of the service account.


## Import

tcr service_account can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_service_account.service_account registry_id#account_name
```

