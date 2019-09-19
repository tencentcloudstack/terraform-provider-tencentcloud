resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_lite_rule" "lite-rule" {
  security_group_id = "${tencentcloud_security_group.foo.id}"

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
    "DROP#8.8.8.8#80,90#UDP",
    "ACCEPT#0.0.0.0/0#80-90#TCP",
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
    "ACCEPT#10.0.0.0/8#ALL#ICMP",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]
}

resource "tencentcloud_security_group_lite_rule" "lite-rule-ingress" {
  security_group_id = "${tencentcloud_security_group.foo.id}"

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
  ]
}

resource "tencentcloud_security_group_lite_rule" "lite-rule-egress" {
  security_group_id = "${tencentcloud_security_group.foo.id}"

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
  ]
}