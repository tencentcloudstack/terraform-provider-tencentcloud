Provides a resource to create a CLB customized config which type is `CLB`.

Example Usage

Create clb customized config without CLB instance

```hcl
resource "tencentcloud_clb_customized_config" "example" {
  config_name    = "tf-example"
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
}
```

Create clb customized config with CLB instances

```hcl
resource "tencentcloud_clb_customized_config" "example" {
  config_name    = "tf-example"
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  load_balancer_ids = [
    "lb-l6cp6jt4",
    "lb-muk4zzxi",
  ]
}
```

Import

CLB customized config can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_customized_config.example pz-diowqstq
```