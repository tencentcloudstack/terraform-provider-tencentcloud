Provides an action to activate (open) the SSM (Secrets Manager) service (PostPay / pay-as-you-go) via the CreateService API. Because CreateService is asynchronous, the action polls GetServiceStatus until the service is enabled. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_ssm_service" "example" {
  config {
    pay_mode = "PostPay"
  }
}
```
