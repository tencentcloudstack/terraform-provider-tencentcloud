Provides a resource to create a ses email_address

Example Usage

```hcl
resource "tencentcloud_ses_email_address" "email_address" {
  email_address     = "aaa@iac-tf.cloud"
  email_sender_name = "aaa"
}

```
Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_email_address.email_address aaa@iac-tf.cloud
```