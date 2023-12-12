Provides a resource to create a waf clb instance

~> **NOTE:** Region only supports `ap-guangzhou` and `ap-seoul`.

Example Usage

Create a basic waf premium clb instance

```hcl
resource "tencentcloud_waf_clb_instance" "example" {
  goods_category = "premium_clb"
  instance_name  = "tf-example-clb-waf"
}
```

Create a complete waf ultimate_clb instance

```hcl
resource "tencentcloud_waf_clb_instance" "example" {
  goods_category  = "ultimate_clb"
  instance_name   = "tf-example-clb-waf"
  time_span       = 1
  time_unit       = "m"
  auto_renew_flag = 1
  elastic_mode    = 1
  bot_management  = 1
  api_security    = 1
}
```

Set waf ultimate_clb instance qps limit

```hcl
resource "tencentcloud_waf_clb_instance" "example" {
  goods_category  = "ultimate_clb"
  instance_name   = "tf-example-clb-waf"
  time_span       = 1
  time_unit       = "m"
  auto_renew_flag = 1
  elastic_mode    = 1
  qps_limit       = 200000
  bot_management  = 1
  api_security    = 1
}
```