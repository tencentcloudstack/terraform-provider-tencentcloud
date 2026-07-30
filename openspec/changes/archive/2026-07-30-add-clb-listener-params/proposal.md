## Why

The CLB Listener resource (`tencentcloud_clb_listener`) is missing several parameters that are already supported by the TencentCloud CLB API, including listener-level connection limits (`MaxConn`, `MaxCps`), proxy protocol support (`ProxyProtocol`), and data compression mode (`DataCompressMode`). These parameters are particularly important for performance capacity-type CLB instances and TCP_SSL/QUIC listeners. Adding them enables Terraform users to fully configure these listener features without resorting to manual API calls.

## What Changes

- Add `MaxConn` parameter (TypeInt, Optional/Computed) - listener-level maximum concurrent connections. Supported by CreateListener, DescribeListeners, and ModifyListener APIs.
- Add `MaxCps` parameter (TypeInt, Optional/Computed) - listener-level maximum new connections per second. Supported by CreateListener, DescribeListeners, and ModifyListener APIs.
- Add `ProxyProtocol` parameter (TypeBool, Optional) - enable proxy protocol for TCP_SSL and QUIC listeners. Supported by CreateListener and ModifyListener APIs. Note: this field is not returned by DescribeListeners API, so it is write-only (not set in Read).
- Add `DataCompressMode` parameter (TypeString, Optional/Computed) - data compression mode. Valid values: `transparent`, `compatibility`. Supported by CreateListener, DescribeListeners, and ModifyListener APIs.

## Capabilities

### New Capabilities
- `clb-listener-params`: Add MaxConn, MaxCps, ProxyProtocol, and DataCompressMode parameters to the tencentcloud_clb_listener resource.

### Modified Capabilities
<!-- No existing capabilities are being modified at the spec level -->

## Impact

- **Affected code**: `tencentcloud/services/clb/resource_tc_clb_listener.go` - Schema, Create, Read, Update functions
- **Affected tests**: `tencentcloud/services/clb/resource_tc_clb_listener_test.go` - Test cases for new parameters
- **Affected docs**: `tencentcloud/services/clb/resource_tc_clb_listener.md` - Documentation for new parameters
- **Dependencies**: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317` SDK (no version change needed, parameters already exist in vendor)