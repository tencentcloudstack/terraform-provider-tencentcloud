Use this data source to query detailed information of CLB listener

Example Usage

```hcl
data "tencentcloud_clb_listeners" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-mwr6vbtv"
  protocol    = "TCP"
  port        = 80
}
```