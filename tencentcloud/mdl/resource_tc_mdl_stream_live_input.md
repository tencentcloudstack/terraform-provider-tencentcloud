Provides a resource to create a mdl streamlive_input

Example Usage

```hcl
resource "tencentcloud_mdl_stream_live_input" "stream_live_input" {
  name               = "terraform_test"
  type               = "RTP_PUSH"
  security_group_ids = [
    "6405DF9D000007DFB4EC"
  ]
}
```

Import

mdl stream_live_input can be imported using the id, e.g.

```
terraform import tencentcloud_mdl_stream_live_input.stream_live_input id
```