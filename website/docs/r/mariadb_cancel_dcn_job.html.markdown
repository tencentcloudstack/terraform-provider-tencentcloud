---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_cancel_dcn_job"
sidebar_current: "docs-tencentcloud-resource-mariadb_cancel_dcn_job"
description: |-
  Provides a resource to create a mariadb cancel_dcn_job
---

# tencentcloud_mariadb_cancel_dcn_job

Provides a resource to create a mariadb cancel_dcn_job

## Example Usage

```hcl
resource "tencentcloud_mariadb_cancel_dcn_job" "cancel_dcn_job" {
  instance_id = "tdsql-9vqvls95"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



