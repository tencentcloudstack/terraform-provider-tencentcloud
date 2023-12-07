Provides a resource to create a tsf config_template

Example Usage

```hcl
resource "tencentcloud_tsf_config_template" "config_template" {
  config_template_name = "terraform-template-name"
  config_template_type = "Ribbon"
  config_template_value = <<-EOT
    ribbon.ReadTimeout: 5000
    ribbon.ConnectTimeout: 2000
    ribbon.MaxAutoRetries: 0
    ribbon.MaxAutoRetriesNextServer: 1
    ribbon.OkToRetryOnAllOperations: true
  EOT
  config_template_desc = "terraform-test"
}
```