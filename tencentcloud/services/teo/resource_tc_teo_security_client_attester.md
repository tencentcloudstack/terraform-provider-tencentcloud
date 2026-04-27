Provides a resource to create a TEO (EdgeOne) security client attester for managing client attestation options.

Example Usage

TC-RCE Authentication

```hcl
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  client_attesters {
    name              = "test-rce"
    attester_source   = "TC-RCE"
    attester_duration = "60s"

    tc_rce_option {
      channel = "channel-1"
      region  = "ap-beijing"
    }
  }
}
```

TC-CAPTCHA Authentication

```hcl
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  client_attesters {
    name              = "test-captcha"
    attester_source   = "TC-CAPTCHA"
    attester_duration = "120s"

    tc_captcha_option {
      captcha_app_id = "199999999"
      app_secret_key = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    }
  }
}
```

TC-EO-CAPTCHA Authentication

```hcl
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  client_attesters {
    name              = "test-eo-captcha"
    attester_source   = "TC-EO-CAPTCHA"
    attester_duration = "300s"

    tc_eo_captcha_option {
      captcha_mode = "Invisible"
    }
  }
}
```

Multiple Attesters

```hcl
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  client_attesters {
    name              = "test-rce"
    attester_source   = "TC-RCE"
    attester_duration = "60s"

    tc_rce_option {
      channel = "channel-1"
      region  = "ap-beijing"
    }
  }

  client_attesters {
    name              = "test-eo-captcha"
    attester_source   = "TC-EO-CAPTCHA"
    attester_duration = "120s"

    tc_eo_captcha_option {
      captcha_mode = "Adaptive"
    }
  }
}
```

Import

teo security_client_attester can be imported using the id, e.g.

```
terraform import tencentcloud_teo_security_client_attester.example zone_id#att-001,att-002
```
