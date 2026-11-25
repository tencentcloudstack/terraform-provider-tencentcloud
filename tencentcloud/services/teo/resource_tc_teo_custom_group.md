Provides a resource to create a TEO (EdgeOne) custom group. A custom group is a reusable collection of items (TLS fingerprints, UA values or Referer values) that can be referenced by security rules. The group type is immutable after creation; the group name and the item list can be updated in place.

~> **NOTE:** The cloud API keeps the previous value when `group_name` or `group_items` is omitted on update. Since both attributes are required, the configured list fully replaces the cloud-side items (`GroupItemOperation=Replace`). Emptying the item list is rejected by the provider because the behavior of `Replace` with an empty list is not confirmed by the TEO API; remove unneeded items individually or recreate the group instead.

~> **NOTE:** Deleting the resource only removes the custom group itself. Whether a group that is still referenced (`ref_count > 0`) can be deleted is decided by the cloud API.

Example Usage

```hcl
resource "tencentcloud_teo_custom_group" "example" {
  zone_id    = "zone-27h0vbm5w1e"
  group_name = "tf-example-group"
  group_type = "TLSFingerprint"
  group_items = [
    {
      content   = "cd08b7ecb0e6f5a72c6f5a8ea17e2e5a"
      item_type = "JA3"
    },
    {
      content   = "t13d1714h2_8daaf6b1a14185f4"
      item_type = "JA4"
    },
  ]
}
```

Import

TEO custom group can be imported using the composite id in `<zone_id>#<group_id>` format, e.g.

````
terraform import tencentcloud_teo_custom_group.example zone-27h0vbm5w1e#cg-2qbzqqti0e9z
````
