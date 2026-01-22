Provides a resource to create a CLB customized config which type is `SERVER` or `LOCATION`.

Example Usage

If config_type is SERVER

```hcl
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "SERVER"
}

output "configId" {
  value = tencentcloud_clb_customized_config_v2.example.config_id
}
```

If config_type is LOCATION

```hcl
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "LOCATION"
}

output "configId" {
  value = tencentcloud_clb_customized_config_v2.example.config_id
}
```

Import

CLB customized V2 config can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_customized_config_v2.example pz-diowqstq#SERVER

Or

$ terraform import tencentcloud_clb_customized_config_v2.example pz-4r10y4b2#LOCATION
```