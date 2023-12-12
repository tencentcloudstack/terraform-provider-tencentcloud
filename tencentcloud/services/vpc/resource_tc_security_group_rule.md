Provides a resource to create security group rule.

~> **NOTE:** This resource will be offline and no longer supported, beacause single security rule is hardly ordered. Please use 'tencentcloud_security_group_lite_rule' instead.

Example Usage

Source is CIDR ip

```hcl
resource "tencentcloud_security_group" "sglab_1" {
  name        = "mysg_1"
  description = "favourite sg_1"
  project_id  = 0
}

resource "tencentcloud_security_group_rule" "sglab_1" {
  security_group_id = tencentcloud_security_group.sglab_1.id
  type              = "ingress"
  cidr_ip           = "10.0.0.0/16"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  description       = "favourite sg rule_1"
}
```

Source is a security group id

```hcl
resource "tencentcloud_security_group" "sglab_2" {
  name        = "mysg_2"
  description = "favourite sg_2"
  project_id  = 0
}

resource "tencentcloud_security_group" "sglab_3" {
  name        = "mysg_3"
  description = "favourite sg_3"
  project_id  = 0
}

resource "tencentcloud_security_group_rule" "sglab_2" {
  security_group_id = tencentcloud_security_group.sglab_2.id
  type              = "ingress"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  source_sgid       = tencentcloud_security_group.sglab_3.id
  description       = "favourite sg rule_2"
}
```