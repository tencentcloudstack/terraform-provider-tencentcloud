resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = tencentcloud_gaap_proxy.foo.id
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  name      = "ci-test-gaap-sr"
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_policies "foo" {
  id = tencentcloud_gaap_security_policy.foo.id
}

data tencentcloud_gaap_security_rules "ruleId" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  rule_id   = tencentcloud_gaap_security_rule.foo.id
}