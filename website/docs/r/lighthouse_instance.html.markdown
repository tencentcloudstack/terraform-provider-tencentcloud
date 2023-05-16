---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_instance"
sidebar_current: "docs-tencentcloud-resource-lighthouse_instance"
description: |-
  Provides a resource to create a lighthouse instance.
---

# tencentcloud_lighthouse_instance

Provides a resource to create a lighthouse instance.

## Example Usage

```hcl
resource "tencentcloud_lighthouse_instance" "lighthouse" {
  bundle_id    = "bundle2022_gen_01"
  blueprint_id = "lhbp-f1lkcd41"

  period     = 1
  renew_flag = "NOTIFY_AND_AUTO_RENEW"

  instance_name = "hello world"
  zone          = "ap-guangzhou-3"

  containers {
    container_image = "ccr.ccs.tencentyun.com/qcloud/nginx"
    container_name  = "nginx"
    envs {
      key   = "key"
      value = "value"
    }
    envs {
      key   = "key2"
      value = "value2"
    }
    publish_ports {
      host_port      = 80
      container_port = 80
      ip             = "127.0.0.1"
      protocol       = "tcp"
    }
    publish_ports {
      host_port      = 36000
      container_port = 36000
      ip             = "127.0.0.1"
      protocol       = "tcp"
    }
    volumes {
      container_path = "/data"
      host_path      = "/tmp"
    }
    volumes {
      container_path = "/var"
      host_path      = "/tmp"
    }
    command = "ls -l"
  }

  containers {
    container_image = "ccr.ccs.tencentyun.com/qcloud/resty"
    container_name  = "resty"
    envs {
      key   = "key2"
      value = "value2"
    }
    publish_ports {
      host_port      = 80
      container_port = 80
      ip             = "127.0.0.1"
      protocol       = "udp"
    }

    volumes {
      container_path = "/var"
      host_path      = "/tmp"
    }
    command = "echo \"hello\""
  }
}
```

## Argument Reference

The following arguments are supported:

* `blueprint_id` - (Required, String) ID of the Lighthouse image.
* `bundle_id` - (Required, String) ID of the Lighthouse package.
* `instance_name` - (Required, String) The display name of the Lighthouse instance.
* `period` - (Required, Int) Subscription period in months. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.
* `renew_flag` - (Required, String) Auto-Renewal flag. Valid values: NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically; NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically. You need to manually renew DISABLE_NOTIFY_AND_AUTO_RENEW: neither notify upon expiration nor renew automatically. Default value: NOTIFY_AND_MANUAL_RENEW.
* `client_token` - (Optional, String) A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.
* `containers` - (Optional, List) Configuration of the containers to create.
* `dry_run` - (Optional, Bool) Whether the request is a dry run only.true: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available. If the dry run fails, the corresponding error code will be returned.If the dry run succeeds, the RequestId will be returned.false (default value): send a normal request and create instance(s) if all the requirements are met.
* `is_update_bundle_id_auto_voucher` - (Optional, Bool) Whether the voucher is deducted automatically when update bundle id. Value range: `true`: indicates automatic deduction of vouchers, `false`: does not automatically deduct vouchers. Default value: `false`.
* `login_configuration` - (Optional, List) Login password of the instance. It is only available for Windows instances. If it is not specified, it means that the user choose to set the login password after the instance creation.
* `permit_default_key_pair_login` - (Optional, String) Whether to allow login using the default key pair. `YES`: allow login; `NO`: disable login. Default: `YES`.
* `zone` - (Optional, String) List of availability zones. A random AZ is selected by default.

The `containers` object supports the following:

* `command` - (Optional, String) The command to run.
* `container_image` - (Optional, String) Container image address.
* `container_name` - (Optional, String) Container name.
* `envs` - (Optional, List) List of environment variables.
* `publish_ports` - (Optional, List) List of mappings of container ports and host ports.
* `volumes` - (Optional, List) List of container mount volumes.

The `envs` object supports the following:

* `key` - (Required, String) Environment variable key.
* `value` - (Required, String) Environment variable value.

The `login_configuration` object supports the following:

* `auto_generate_password` - (Required, String) whether auto generate password. if false, need set password.
* `password` - (Optional, String) Login password.

The `publish_ports` object supports the following:

* `container_port` - (Required, Int) Container port.
* `host_port` - (Required, Int) Host port.
* `ip` - (Optional, String) External IP. It defaults to 0.0.0.0.
* `protocol` - (Optional, String) The protocol defaults to tcp. Valid values: tcp, udp and sctp.

The `volumes` object supports the following:

* `container_path` - (Required, String) Container path.
* `host_path` - (Required, String) Host path.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



