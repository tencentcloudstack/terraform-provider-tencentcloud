Provides a resource to create a ses email address

Example Usage

Create ses email address

```hcl
resource "tencentcloud_ses_email_address" "example" {
  email_address     = "demo@iac-terraform.cloud"
  email_sender_name = "root"
}
```

Set smtp password

```hcl
resource "tencentcloud_ses_email_address" "example" {
  email_address     = "demo@iac-terraform.cloud"
  email_sender_name = "root"
  smtp_password     = "Password@123"
}
```

Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_email_address.example demo@iac-terraform.cloud
```