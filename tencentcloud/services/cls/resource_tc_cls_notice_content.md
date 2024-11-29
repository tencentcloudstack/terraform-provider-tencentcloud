Provides a resource to create a cls notice content

Example Usage

```hcl
resource "tencentcloud_cls_notice_content" "example" {
  name = "tf-example"
  type = 0
  notice_contents {
    type = "Email"

    trigger_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }

    recovery_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }
  }
}
```

Import

cls notice content can be imported using the id, e.g.

```
terraform import tencentcloud_cls_notice_content.example noticetemplate-b417f32a-bdf9-46c5-933e-28c23cd7a6b7
```
