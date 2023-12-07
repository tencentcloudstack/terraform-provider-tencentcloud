Provide a resource to create a DnsPod Domain instance.

Example Usage

```hcl
resource "tencentcloud_dnspod_domain_instance" "foo" {
  domain = "hello.com"
  remark = "this is demo"
}
```

Import

DnsPod Domain instance can be imported, e.g.

```
$ terraform import tencentcloud_dnspod_domain_instance.foo domain
```