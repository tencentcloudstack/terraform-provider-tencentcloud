Provides a resource to create a Config remediation setting.

Example Usage

```hcl
resource "tencentcloud_config_remediation" "example" {
  rule_id                  = "cr-Gxt8pzxgCVZJ0C95H1HO"
  remediation_type         = "SCF"
  remediation_template_id  = "qcs::scf:ap-guangzhou:uin/100000005287:namespace/test/functions/my-remediation-func"
  invoke_type              = "MANUAL_EXECUTION"
  source_type              = "CUSTOM"
}

```

Import

Config remediation can be imported using the id, e.g.

```
terraform import tencentcloud_config_remediation.example crr-lKj43O4nbSD78XYlvGS9
```
