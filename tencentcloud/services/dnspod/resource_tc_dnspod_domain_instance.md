Provide a resource to create a DnsPod Domain instance.

Example Usage

```hcl
resource "tencentcloud_dnspod_domain_instance" "foo" {
  domain = "hello.com"
  remark = "this is demo"
}

# Access computed fields
output "domain_status" {
  value = tencentcloud_dnspod_domain_instance.foo.status
}

output "record_count" {
  value = tencentcloud_dnspod_domain_instance.foo.record_count
}

output "domain_grade" {
  value = tencentcloud_dnspod_domain_instance.foo.grade
}
```

Import

DnsPod Domain instance can be imported, e.g.

```
$ terraform import tencentcloud_dnspod_domain_instance.foo domain
```