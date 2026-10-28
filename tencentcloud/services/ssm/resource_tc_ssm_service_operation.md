Provides a resource to create a SSM (Secrets Manager) service operation to activate (open) the SSM service in PostPay (pay-as-you-go) mode via the CreateService API. Because CreateService is asynchronous, this operation polls GetServiceStatus until the service is enabled.

Example Usage

```hcl
resource "tencentcloud_ssm_service_operation" "example" {
  pay_mode = "PostPay"
}
```