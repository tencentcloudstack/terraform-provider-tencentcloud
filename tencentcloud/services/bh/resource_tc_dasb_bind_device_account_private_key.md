Provides a resource to create a dasb bind_device_account_private_key

Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_account_private_key" "example" {
  device_account_id    = 16
  private_key          = "MIICXAIBAAKBgQCqGKukO1De7zhZj6+H0qtjTkVxwTCpvKe4eCZ0FPqri0cb2JZfXJ/DgYSF6vUpwmJG8wVQZKjeGcjDOL5UlsuusFncCzWBQ7RKNUSesmQRMSGkVb1/3j+skZ6UtW+5u09lHNsj6tQ51s1SPrCBkedbNf0Tp0GbMJDyR4e9T04ZZwIDAQABAoGAFijko56+qGyN8M0RVyaRAXz++xTqHBLh"
  private_key_password = "TerraformPassword"
}
```