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