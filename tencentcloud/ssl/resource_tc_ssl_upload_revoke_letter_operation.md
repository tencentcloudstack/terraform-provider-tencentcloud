Provides a resource to create a ssl upload_revoke_letter

Example Usage

```hcl
resource "tencentcloud_ssl_upload_revoke_letter_operation" "upload_revoke_letter" {
  certificate_id = "8xRYdDlc"
  revoke_letter = filebase64("./c.pdf")
}
```

Import

ssl upload_revoke_letter can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_upload_revoke_letter_operation.upload_revoke_letter upload_revoke_letter_id
```