Provide a resource to create a DnsPod record.

~> **NOTE:** Versions before v1.81.43 (including v1.81.43) do not support modifying remark or modifying remark has bug.

Example Usage

```hcl
resource "tencentcloud_dnspod_record" "demo" {
  domain="mikatong.com"
  record_type="A"
  record_line="默认"
  value="1.2.3.9"
  sub_domain="demo"
}
```

Import

DnsPod Domain record can be imported using the Domain#RecordId, e.g.

```
$ terraform import tencentcloud_dnspod_record.demo arunma.com#1194109872
```