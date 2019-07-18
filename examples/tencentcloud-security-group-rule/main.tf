resource "tencentcloud_security_group" "default" {
  name        = "${var.security_group_name}"
  description = "New security group"
}

resource "tencentcloud_security_group" "default2" {
  name        = "${var.security_group_name}"
  description = "Anthor security group"
}

resource "tencentcloud_security_group_rule" "http-in" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}

resource "tencentcloud_security_group_rule" "ssh-in" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "22"
  policy            = "accept"
}

resource "tencentcloud_security_group_rule" "egress-drop" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "egress"
  cidr_ip           = "10.2.3.0/24"
  ip_protocol       = "udp"
  port_range        = "3000-4000"
  policy            = "drop"
}

resource "tencentcloud_security_group_rule" "sourcesgid-in" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  source_sgid       = "${tencentcloud_security_group.default2.id}"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}
