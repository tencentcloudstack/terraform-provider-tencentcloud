Provides a resource to create a waf saas instance

~> **NOTE:** Region only supports `ap-guangzhou` and `ap-seoul`.

Example Usage

Create a basic waf premium saas instance

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category = "premium_saas"
  instance_name  = "tf-example-saas-waf"
}
```

Create a complete waf ultimate_saas instance

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category  = "ultimate_saas"
  instance_name   = "tf-example-saas-waf"
  time_span       = 1
  time_unit       = "m"
  auto_renew_flag = 1
  elastic_mode    = 1
  real_region     = "gz"
  bot_management  = 1
  api_security    = 1
}
```

Set waf ultimate_saas instance qps limit

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category  = "ultimate_saas"
  instance_name   = "tf-example-saas-waf"
  time_span       = 1
  time_unit       = "m"
  auto_renew_flag = 1
  elastic_mode    = 1
  real_region     = "gz"
  qps_limit       = 200000
  bot_management  = 1
  api_security    = 1
}
```