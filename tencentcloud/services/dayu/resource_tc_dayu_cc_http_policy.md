Use this resource to create a dayu CC self-define http policy

Example Usage

```hcl
resource "tencentcloud_dayu_cc_http_policy" "test_bgpip" {
  resource_type = "bgpip"
  resource_id   = "bgpip-00000294"
  name          = "policy_match"
  smode         = "matching"
  action        = "drop"
  switch        = true
  rule_list {
    skey     = "host"
    operator = "include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_net" {
  resource_type = "net"
  resource_id   = "net-0000007e"
  name          = "policy_match"
  smode         = "matching"
  action        = "drop"
  switch        = true
  rule_list {
    skey     = "cgi"
    operator = "equal"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_bgpmultip" {
  resource_type = "bgp-multip"
  resource_id   = "bgp-0000008o"
  name          = "policy_match"
  smode         = "matching"
  action        = "alg"
  switch        = true
  ip            = "111.230.178.25"

  rule_list {
    skey     = "referer"
    operator = "not_include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_bgp" {
  resource_type = "bgp"
  resource_id   = "bgp-000006mq"
  name          = "policy_match"
  smode         = "matching"
  action        = "alg"
  switch        = true

  rule_list {
    skey     = "ua"
    operator = "not_include"
    value    = "123"
  }
}
```