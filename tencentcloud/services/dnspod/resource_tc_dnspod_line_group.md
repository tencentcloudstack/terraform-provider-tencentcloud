Provides a resource to create a DNSPod line group.

Example Usage

```hcl
resource "tencentcloud_dnspod_line_group" "example" {
  domain = "example.com"
  name   = "telecom_group"
  lines  = ["电信", "移动"]
}
```

Import

DNSPod line group can be imported using the id (format: `{domain}#{line_group_id}`), e.g.

```
$ terraform import tencentcloud_dnspod_line_group.example example.com#123
```
