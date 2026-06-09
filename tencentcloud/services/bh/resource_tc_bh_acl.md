Provides a BH (Bastion Host) acl.

Example Usage

```hcl
resource "tencentcloud_bh_acl" "example" {
  name                    = "tf-example"
  allow_disk_redirect     = true
  allow_any_account       = false
  allow_clip_file_up      = true
  allow_clip_file_down    = false
  allow_clip_text_up      = false
  allow_clip_text_down    = true
  allow_file_up           = false
  allow_file_down         = true
  allow_disk_file_up      = true
  allow_disk_file_down    = false
  allow_shell_file_up     = false
  allow_shell_file_down   = false
  allow_file_del          = false
  allow_access_credential = true
  allow_keyboard_logger   = false
}
```

Import

BH acl can be imported using the id, e.g.

```
terraform import tencentcloud_bh_acl.example 1374
```
