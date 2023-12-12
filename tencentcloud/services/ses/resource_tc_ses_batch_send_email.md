Provides a resource to create a ses batch_send_email

Example Usage

```hcl
resource "tencentcloud_ses_batch_send_email" "batch_send_email" {
  from_email_address = "aaa@iac-tf.cloud"
  receiver_id        = 1063742
  subject            = "terraform test"
  task_type          = 1
  reply_to_addresses = "reply@mail.qcloud.com"
  template {
    template_id   = 99629
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"

  }

  cycle_param {
    begin_time = "2023-09-07 15:10:00"
    interval_time = 1
  }
  timed_param {
    begin_time = "2023-09-07 15:20:00"
  }
  unsubscribe = "0"
  ad_location = 0
}
```