---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_jwks_authenticator"
sidebar_current: "docs-tencentcloud-resource-mqtt_jwks_authenticator"
description: |-
  Provides a resource to create a MQTT jwks authenticator
---

# tencentcloud_mqtt_jwks_authenticator

Provides a resource to create a MQTT jwks authenticator

## Example Usage

```hcl
resource "tencentcloud_mqtt_jwks_authenticator" "example" {
  instance_id      = "mqtt-zxjwkr98"
  from             = "username"
  endpoint         = "https://example.com"
  refresh_interval = 60
  remark           = "Remark."
}
```

### Or

```hcl
resource "tencentcloud_mqtt_jwks_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  from        = "password"
  text        = "your text content"
  remark      = "Remark."
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `endpoint` - (Optional, String) JWKS endpoint.
* `from` - (Optional, String) Pass the key of JWT when connecting the device; Username - passed using the username field; Password - Pass using password field.
* `refresh_interval` - (Optional, Int) JWKS refresh interval. unit: s.
* `remark` - (Optional, String) Remark.
* `text` - (Optional, String) JWKS text.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateJWKSAuthenticator](https://www.tencentcloud.com/search/?wd=CreateJWKSAuthenticator&type=document&lang=en)
* Read: [DescribeAuthenticator](https://www.tencentcloud.com/search/?wd=DescribeAuthenticator&type=document&lang=en)
* Update: [ModifyJWKSAuthenticator](https://www.tencentcloud.com/search/?wd=ModifyJWKSAuthenticator&type=document&lang=en)
* Delete: [DeleteAuthenticator](https://www.tencentcloud.com/search/?wd=DeleteAuthenticator&type=document&lang=en)

## Import

MQTT jwks authenticator can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_jwks_authenticator.example mqtt-zxjwkr98
```

