---
subcategory: "TencentCloud ServiceMesh(TCM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcm_access_log_config"
sidebar_current: "docs-tencentcloud-resource-tcm_access_log_config"
description: |-
  Provides a resource to create a tcm access_log_config
---

# tencentcloud_tcm_access_log_config

Provides a resource to create a tcm access_log_config

## Example Usage

```hcl
resource "tencentcloud_tcm_access_log_config" "access_log_config" {
  address       = "10.0.0.1"
  enable        = true
  enable_server = true
  enable_stdout = true
  encoding      = "JSON"
  format        = "{\n\t\"authority\": \"%REQ(:AUTHORITY)%\",\n\t\"bytes_received\": \"%BYTES_RECEIVED%\",\n\t\"bytes_sent\": \"%BYTES_SENT%\",\n\t\"downstream_local_address\": \"%DOWNSTREAM_LOCAL_ADDRESS%\",\n\t\"downstream_remote_address\": \"%DOWNSTREAM_REMOTE_ADDRESS%\",\n\t\"duration\": \"%DURATION%\",\n\t\"istio_policy_status\": \"%DYNAMIC_METADATA(istio.mixer:status)%\",\n\t\"method\": \"%REQ(:METHOD)%\",\n\t\"path\": \"%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%\",\n\t\"protocol\": \"%PROTOCOL%\",\n\t\"request_id\": \"%REQ(X-REQUEST-ID)%\",\n\t\"requested_server_name\": \"%REQUESTED_SERVER_NAME%\",\n\t\"response_code\": \"%RESPONSE_CODE%\",\n\t\"response_flags\": \"%RESPONSE_FLAGS%\",\n\t\"route_name\": \"%ROUTE_NAME%\",\n\t\"start_time\": \"%START_TIME%\",\n\t\"upstream_cluster\": \"%UPSTREAM_CLUSTER%\",\n\t\"upstream_host\": \"%UPSTREAM_HOST%\",\n\t\"upstream_local_address\": \"%UPSTREAM_LOCAL_ADDRESS%\",\n\t\"upstream_service_time\": \"%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%\",\n\t\"upstream_transport_failure_reason\": \"%UPSTREAM_TRANSPORT_FAILURE_REASON%\",\n\t\"user_agent\": \"%REQ(USER-AGENT)%\",\n\t\"x_forwarded_for\": \"%REQ(X-FORWARDED-FOR)%\"\n}\n"
  mesh_name     = "mesh-rofjmxxx"
  template      = "istio"

  cls {
    enable = false
    # log_set = "SCF_logset_NLCsDxxx"
    # topic   = "SCF_logtopic_rPWZpxxx"
  }

  selected_range {
    all = true
  }
}

resource "tencentcloud_tcm_access_log_config" "delete_log_config" {
  enable_server = false
  enable_stdout = false
  mesh_name     = "mesh-rofjmux7"

  cls {
    enable = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `mesh_name` - (Required, String) Mesh ID.
* `address` - (Optional, String) Third party grpc server address.
* `cls` - (Optional, List) CLS config.
* `enable_server` - (Optional, Bool) Whether enable third party grpc server.
* `enable_stdout` - (Optional, Bool) Whether enable stdout.
* `enable` - (Optional, Bool) Whether enable log.
* `encoding` - (Optional, String) Log encoding, TEXT or JSON.
* `format` - (Optional, String) Log format.
* `selected_range` - (Optional, List) Selected range.
* `template` - (Optional, String) Log template, istio/trace/custome.

The `cls` object supports the following:

* `enable` - (Required, Bool) Whether enable CLS.
* `log_set` - (Optional, String) Log set of CLS.
* `topic` - (Optional, String) Log topic of CLS.

The `items` object supports the following:

* `gateways` - (Optional, Set) Ingress gateway list.
* `namespace` - (Optional, String) Namespace.

The `selected_range` object supports the following:

* `all` - (Optional, Bool) Select all if true, default false.
* `items` - (Optional, List) Items.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcm access_log_config can be imported using the mesh_id(mesh_name), e.g.
```
$ terraform import tencentcloud_tcm_access_log_config.access_log_config mesh-rofjmxxx
```

