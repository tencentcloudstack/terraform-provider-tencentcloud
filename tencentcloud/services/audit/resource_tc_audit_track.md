Provides a resource to create a audit track

Example Usage

```hcl
resource "tencentcloud_audit_track" "example" {
  action_type = "Read"
  event_names = [
    "*",
  ]
  name                  = "terraform_track"
  resource_type         = "*"
  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "db90b92c-91d2-46b0-94ac-debbbb21dc4e"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cls"
  }
}
```

Specify storage user

```hcl
resource "tencentcloud_audit_track" "example" {
  action_type = "Read"
  event_names = [
    "*",
  ]
  name                  = "terraform_track"
  resource_type         = "*"
  status                = 1
  track_for_all_members = 0

  storage {
    storage_name       = "db90b92c-91d2-46b0-94ac-debbbb21dc4e"
    storage_prefix     = "cloudaudit"
    storage_region     = "ap-guangzhou"
    storage_type       = "cos"
    storage_account_id = "100037717137"
    storage_app_id     = "1309116520"
  }
}
```

Import

audit track can be imported using the id, e.g.
```
$ terraform import tencentcloud_audit_track.example 24283
```