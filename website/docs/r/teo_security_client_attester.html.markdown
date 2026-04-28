---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_client_attester"
sidebar_current: "docs-tencentcloud-resource-teo_security_client_attester"
description: |-
  Provides a resource to create a TEO security client attester.
---

# tencentcloud_teo_security_client_attester

Provides a resource to create a TEO security client attester.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `client_attesters` - (Required, List) Client attester configuration. Only one attester is allowed per request.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `client_attesters` object supports the following:

* `attester_source` - (Required, String) Authentication method. Valid values: `TC-RCE` (Tencent Cloud RCE), `TC-CAPTCHA` (Tencent CAPTCHA), `TC-EO-CAPTCHA` (EdgeOne CAPTCHA).
* `name` - (Required, String) Attester name.
* `attester_duration` - (Optional, String) Authentication validity duration. Default `60s`. Supported units: `s` (60-43200), `m` (1-720), `h` (1-12). e.g. `300s`, `5m`, `1h`.
* `tc_captcha_option` - (Optional, List) TC-CAPTCHA authentication configuration. Required when `attester_source` is `TC-CAPTCHA`.
* `tc_eo_captcha_option` - (Optional, List) TC-EO-CAPTCHA authentication configuration. Required when `attester_source` is `TC-EO-CAPTCHA`.
* `tc_rce_option` - (Optional, List) TC-RCE authentication configuration. Required when `attester_source` is `TC-RCE`.

The `tc_captcha_option` object of `client_attesters` supports the following:

* `app_secret_key` - (Required, String) AppSecretKey.
* `captcha_app_id` - (Required, String) CaptchaAppId.

The `tc_eo_captcha_option` object of `client_attesters` supports the following:

* `captcha_mode` - (Required, String) EdgeOne CAPTCHA mode. Valid values: `Invisible`, `Adaptive`.

The `tc_rce_option` object of `client_attesters` supports the following:

* `channel` - (Required, String) TC-RCE channel ID.
* `region` - (Optional, String) TC-RCE channel region. Valid values: `ap-beijing`, `ap-jakarta`, `ap-singapore`, `eu-frankfurt`, `na-siliconvalley`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO security client attester can be imported using the zoneId#clientAttesterId, e.g.

```
terraform import tencentcloud_teo_security_client_attester.example zone-3fkff38fyw8s#attest-0000361666
```

