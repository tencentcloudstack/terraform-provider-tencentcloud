Provides a resource to create a ses black_list

~> **NOTE:** Used to remove email addresses from blacklists.

Example Usage

```hcl
resource "tencentcloud_ses_black_list_delete" "black_list" {
  email_address = "terraform-tf@gmail.com"
}
```