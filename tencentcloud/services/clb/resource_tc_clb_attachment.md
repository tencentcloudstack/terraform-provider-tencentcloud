Provides a resource to create a CLB attachment.

~> **NOTE:** This resource is designed to manage the entire set of binding relationships associated with a particular CLB (Cloud Load Balancer). As such, it does not allow the simultaneous use of this resource for the same CLB across different contexts or environments.


Example Usage

Bind a Cvm instance
```hcl
resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }
}
```

Bind multiple Cvm instances
```hcl
resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }
  
  targets {
    instance_id = "ins-ekloqpa1"
    port        = 81
    weight      = 10
  }
}
```

Bind backend target is ENI
```hcl
resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"
  
  targets {
    eni_ip      = "example-ip"
    port        = 23
    weight      = 50
  }
}
```
Import

CLB attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_attachment.foo loc-4xxr2cy7#lbl-hh141sn9#lb-7a0t6zqb
```