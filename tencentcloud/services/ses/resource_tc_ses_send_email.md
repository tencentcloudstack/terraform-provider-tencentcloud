Provides a resource to create a ses send_email

Example Usage

```hcl
resource "tencentcloud_ses_send_email" "send_email" {
  from_email_address = "aaa@iac-tf.cloud"
  destination        = ["1055482519@qq.com"]
  subject            = "test subject"
  reply_to_addresses = "aaa@iac-tf.cloud"

  template {
    template_id   = 99629
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  unsubscribe  = "1"
  trigger_type = 1
}
```