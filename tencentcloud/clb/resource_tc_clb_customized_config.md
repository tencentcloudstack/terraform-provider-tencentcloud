Provides a resource to create a CLB customized config.

Example Usage

```hcl
resource "tencentcloud_clb_customized_config" "foo" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "helloWorld"
  load_balancer_ids = [
    "${tencentcloud_clb_instance.internal_clb.id}",
    "${tencentcloud_clb_instance.internal_clb2.id}",
  ]
}
```
Import

CLB customized config can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_customized_config.foo pz-diowqstq
```