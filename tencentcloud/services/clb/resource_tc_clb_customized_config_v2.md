Provides a resource to create a CLB customized V2 config.

Example Usage

Create clb customized V2 config without CLB instance

```hcl
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "CLB"
}

output "configId" {
  value = tencentcloud_clb_customized_config_v2.example.config_id
}
```

Import

CLB customized V2 config can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_customized_config_v2.example pz-diowqstq
```