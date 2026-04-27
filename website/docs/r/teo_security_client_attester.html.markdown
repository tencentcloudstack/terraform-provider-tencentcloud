---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_client_attester"
sidebar_current: "docs-tencentcloud-resource-teo_security_client_attester"
description: |-
  Provides a resource to create a TEO (EdgeOne) security client attester for managing client attestation options.
---

# tencentcloud_teo_security_client_attester

Provides a resource to create a TEO (EdgeOne) security client attester for managing client attestation options.

## Example Usage

### TC-RCE Authentication

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

### TC-CAPTCHA Authentication

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

### TC-EO-CAPTCHA Authentication

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

### Multiple Attesters

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

## Argument Reference

The following arguments are supported:

* `client_attesters` - (Required, List) Client attestation option list.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `client_attesters` object supports the following:

* `attester_source` - (Required, String) Authentication method. Values: TC-RCE, TC-CAPTCHA, TC-EO-CAPTCHA.
* `name` - (Required, String) Attestation option name.
* `attester_duration` - (Optional, String) Authentication validity duration. Default 60s. Supported units: s (60-43200), m (1-720), h (1-12).
* `tc_captcha_option` - (Optional, List) TC-CAPTCHA authentication configuration, required when attester_source is TC-CAPTCHA.
* `tc_eo_captcha_option` - (Optional, List) TC-EO-CAPTCHA authentication configuration, required when attester_source is TC-EO-CAPTCHA.
* `tc_rce_option` - (Optional, List) TC-RCE authentication configuration, required when attester_source is TC-RCE.

The `tc_captcha_option` object of `client_attesters` supports the following:

* `app_secret_key` - (Required, String) AppSecretKey information.
* `captcha_app_id` - (Required, String) CaptchaAppId information.

The `tc_eo_captcha_option` object of `client_attesters` supports the following:

* `captcha_mode` - (Required, String) EdgeOne human-machine verification mode. Values: Invisible, Adaptive.

The `tc_rce_option` object of `client_attesters` supports the following:

* `channel` - (Required, String) Channel information.
* `region` - (Required, String) RCE Channel region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `client_attester_ids` - List of client attestation option IDs.


## Import

teo security_client_attester can be imported using the id, e.g.

```
terraform import tencentcloud_teo_security_client_attester.example zone_id#att-001,att-002
```

