Provides a resource to create a ses receiver

Example Usage

```hcl
resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc = "description"

  data {
    email = "abc@abc.com"
  }

  data {
    email = "abcd@abcd.com"
  }
}
```

Create a template with `template_data`
```hcl
resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc = "description"

  data {
    email = "abc@abc.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  data {
    email = "abcd@abcd.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }
}
```
Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_receiver.receiver receiverId
```