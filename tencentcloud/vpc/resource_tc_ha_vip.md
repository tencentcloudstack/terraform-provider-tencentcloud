Provides a resource to create a HA VIP.

Example Usage

```hcl
resource "tencentcloud_ha_vip" "foo" {
  name      = "terraform_test"
  vpc_id    = "vpc-gzea3dd7"
  subnet_id = "subnet-4d4m4cd4s"
  vip       = "10.0.4.16"
}
```

Import

HA VIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_ha_vip.foo havip-kjqwe4ba
```