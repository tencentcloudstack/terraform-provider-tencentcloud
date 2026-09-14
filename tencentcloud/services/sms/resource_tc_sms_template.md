Provides a resource to create a sms template

Example Usage

Create a sms template

```hcl
resource "tencentcloud_sms_template" "template" {
  template_name = "tf_example_sms_template"
  template_content = "example for sms template"
  international = 0 # Mainland China SMS
  sms_type = 0 # regular SMS
  remark = "terraform example"
}

```
The `status_code` field is a computed output attribute that shows the template's review status:
- `0`: approved and effective
- `1`: pending review
- `2`: approved pending activation
- `-1`: review failed or rejected