Provides a resource to create a TEO security client attester.

Example Usage

```hcl
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-3fkff38fyw8s"

  client_attesters {
    name              = "tf-example"
    attester_source   = "TC-RCE"
    attester_duration = "300s"

    tc_rce_option {
      channel = "12399223"
      region  = "ap-beijing"
    }
  }
}
```

Example Usage - TC-EO-CAPTCHA with usage limit and mobile captcha mode

```hcl
resource "tencentcloud_teo_security_client_attester" "eo_example" {
  zone_id = "zone-3fkff38fyw8s"

  client_attesters {
    name             = "tf-eo-example"
    attester_source   = "TC-EO-CAPTCHA"
    attester_duration = "300s"
    usage_limit       = "on"
    max_usage_count   = 100

    tc_eo_captcha_option {
      captcha_mode        = "Invisible"
      mobile_captcha_mode = "Visible"
    }
  }
}
```

Import

TEO security client attester can be imported using the zoneId#clientAttesterId, e.g.

```
terraform import tencentcloud_teo_security_client_attester.example zone-3fkff38fyw8s#attest-0000361666
```
