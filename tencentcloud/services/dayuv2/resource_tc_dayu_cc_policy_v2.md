Use this resource to create a dayu CC policy

Example Usage

```hcl
resource "tencentcloud_dayu_cc_policy_v2" "demo" {
  resource_id="bgpip-000004xf"
  business="bgpip"
  thresholds {
    domain="12.com"
    threshold=0
  }
  cc_geo_ip_policys {
    action="drop"
    region_type="china"
    domain="12.com"
    protocol="http"
  }

  cc_black_white_ips {
    protocol="http"
    domain="12.com"
    black_white_ip="1.2.3.4"
    type="black"
  }
  cc_precision_policys{
    policy_action="drop"
    domain="1.com"
    protocol="http"
    ip="162.62.163.34"
    policys {
      field_name="cgi"
      field_type="value"
      value="12123.com"
      value_operator="equal"
    }
  }
  cc_precision_req_limits {
    domain="11.com"
    protocol="http"
    level="loose"
    policys {
      action="alg"
      execute_duration=2
      mode="equal"
      period=5
      request_num=12
      uri="15.com"
    }
  }
}
```