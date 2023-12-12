Provides a resource to create security group rule. This resource is similar with tencentcloud_security_group_lite_rule, rules can be ordered and configure descriptions.

~> **NOTE:** This resource must exclusive in one security group, do not declare additional rule resources of this security group elsewhere.

Example Usage

```hcl
resource "tencentcloud_security_group" "base" {
  name        = "test-set-sg"
  description = "Testing Rule Set Security"
}

resource "tencentcloud_security_group" "relative" {
  name        = "for-relative"
  description = "Used for attach security policy"
}

resource "tencentcloud_address_template" "foo" {
  name      = "test-set-aTemp"
  addresses = ["10.0.0.1", "10.0.1.0/24", "10.0.0.1-10.0.0.100"]
}

resource "tencentcloud_address_template_group" "foo" {
  name         = "test-set-atg"
  template_ids = [tencentcloud_address_template.foo.id]
}

resource "tencentcloud_security_group_rule_set" "base" {
  security_group_id = tencentcloud_security_group.base.id

  ingress {
    action      = "ACCEPT"
    cidr_block  = "10.0.0.0/22"
    protocol    = "TCP"
    port        = "80-90"
    description = "A:Allow Ips and 80-90"
  }

  ingress {
    action      = "ACCEPT"
    cidr_block  = "10.0.2.1"
    protocol    = "UDP"
    port        = "8080"
    description = "B:Allow UDP 8080"
  }

  ingress {
    action      = "ACCEPT"
    cidr_block  = "10.0.2.1"
    protocol    = "UDP"
    port        = "8080"
    description = "C:Allow UDP 8080"
  }

  ingress {
    action      = "ACCEPT"
    cidr_block  = "172.18.1.2"
    protocol    = "ALL"
    port        = "ALL"
    description = "D:Allow ALL"
  }

  ingress {
    action             = "DROP"
    protocol           = "TCP"
    port               = "80"
    source_security_id = tencentcloud_security_group.relative.id
    description        = "E:Block relative"
  }

  egress {
    action      = "DROP"
    cidr_block  = "10.0.0.0/16"
    protocol    = "ICMP"
    description = "A:Block ping3"
  }

  egress {
    action              = "DROP"
    address_template_id = tencentcloud_address_template.foo.id
    description         = "B:Allow template"
  }

  egress {
    action                 = "DROP"
    address_template_group = tencentcloud_address_template_group.foo.id
    description            = "C:DROP template group"
  }
}
```

Import

Resource tencentcloud_security_group_rule_set can be imported by passing security grou id:

```
terraform import tencentcloud_security_group_rule_set.sglab_1 sg-xxxxxxxx
```