---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_service_operation"
sidebar_current: "docs-tencentcloud-resource-ssm_service_operation"
description: |-
  Provides a resource to create a SSM (Secrets Manager) service operation to activate (open) the SSM service in PostPay (pay-as-you-go) mode via the CreateService API. Because CreateService is asynchronous, this operation polls GetServiceStatus until the service is enabled.
---

# tencentcloud_ssm_service_operation

Provides a resource to create a SSM (Secrets Manager) service operation to activate (open) the SSM service in PostPay (pay-as-you-go) mode via the CreateService API. Because CreateService is asynchronous, this operation polls GetServiceStatus until the service is enabled.

## Example Usage

```hcl
resource "tencentcloud_ssm_service_operation" "example" {
  pay_mode = "PostPay"
}
```

## Argument Reference

The following arguments are supported:

* `pay_mode` - (Required, String, ForceNew) Billing mode for opening SSM. Only `PostPay` (pay-as-you-go) is supported; the value is case-sensitive.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateService](https://www.tencentcloud.com/search/?wd=CreateService&type=document&lang=en)
  - [GetServiceStatus](https://www.tencentcloud.com/document/api/1078/38648)

