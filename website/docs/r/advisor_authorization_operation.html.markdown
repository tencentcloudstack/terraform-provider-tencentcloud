---
subcategory: "Advisor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_advisor_authorization_operation"
sidebar_current: "docs-tencentcloud-resource-advisor_authorization_operation"
description: |-
  Provides a resource to create a Advisor authorization operation
---

# tencentcloud_advisor_authorization_operation

Provides a resource to create a Advisor authorization operation

## Example Usage

```hcl
resource "tencentcloud_advisor_authorization_operation" "example" {}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAdvisorAuthorization](https://cloud.tencent.com/document/api/1264/127787) - Enable Tencent Cloud Smart Advisor authorization. This will synchronously enable report interpretation and cloud architecture collaboration permissions.

