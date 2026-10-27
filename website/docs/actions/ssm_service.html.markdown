---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_service"
sidebar_current: "docs-tencentcloud-action-ssm_service"
description: |-
  Provides an action to activate (open) the SSM (Secrets Manager) service (PostPay / pay-as-you-go) via the CreateService API. Because CreateService is asynchronous, the action polls GetServiceStatus until the service is enabled. This is a one-time operation; no cloud-side state is persisted after the action completes.
---

# tencentcloud_ssm_service

Provides an action to activate (open) the SSM (Secrets Manager) service (PostPay / pay-as-you-go) via the CreateService API. Because CreateService is asynchronous, the action polls GetServiceStatus until the service is enabled. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

## Example Usage

```hcl
action "tencentcloud_ssm_service" "example" {
  config {
    pay_mode     = "PostPay"
    client_token = "ssm-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `client_token` - (Required, String) Client-generated idempotent create token. 1-64 visible ASCII characters (no spaces, Chinese, or control characters); recommended `ssm-<UUID>`.
* `pay_mode` - (Required, String) Billing mode for opening SSM. Only `PostPay` (pay-as-you-go) is supported; the value is case-sensitive.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke:
  - [CreateService](https://www.tencentcloud.com/search/?wd=CreateService&type=document&lang=en)
  - [GetServiceStatus](https://www.tencentcloud.com/document/api/1078/38648)

