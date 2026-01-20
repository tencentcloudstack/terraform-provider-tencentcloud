Provides a resource to create a Waf object

~> **NOTE:** If you need to change field `instance_id`, you need to keep `status` at `0`; If you need to change field `proxy(ip_headers)`, you need to keep `status` at `1`.

Example Usage

Bind current account resources

```hcl
resource "tencentcloud_waf_object" "example" {
  object_id   = "lb-9h5x9lze"
  instance_id = "waf_2kxtlbky11b2v4fe"
  status      = 1
  proxy       = 3
  ip_headers = [
    "my-header1",
    "my-header2",
    "my-header3",
  ]
}
```

Bind other member account resources

```hcl
resource "tencentcloud_waf_object" "example" {
  object_id     = "lb-0ljh3xew"
  instance_id   = "waf_2kxtlbky11b2v4fe"
  member_app_id = 1306832456
  member_uin    = "100987654164"
  status        = 1
  proxy         = 1
}
```

Import

Waf object can be imported using the id, e.g.

```
terraform import tencentcloud_waf_object.example lb-9h5x9lze
```
