Provides a resource to create a clb target_group_attachments

Example Usage

```hcl
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  load_balancer_id = "lb-phbx2420"
  associations {
		listener_id = "lbl-m2q6sp9m"
		target_group_id = "lbtg-5xunivs0"
		location_id = "loc-jjqr0ric"

  }
}
```

