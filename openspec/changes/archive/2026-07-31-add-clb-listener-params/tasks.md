## 1. Schema Definition

- [x] 1.1 Add `max_conn` (TypeInt, Optional, Computed) field to the resource schema in `resource_tc_clb_listener.go`
- [x] 1.2 Add `max_cps` (TypeInt, Optional, Computed) field to the resource schema in `resource_tc_clb_listener.go`
- [x] 1.3 Add `proxy_protocol` (TypeBool, Optional) field to the resource schema in `resource_tc_clb_listener.go`
- [x] 1.4 Add `data_compress_mode` (TypeString, Optional, Computed) field to the resource schema in `resource_tc_clb_listener.go`

## 2. Create Implementation

- [x] 2.1 Add `MaxConn`, `MaxCps`, `ProxyProtocol`, `DataCompressMode` to `resourceTencentCloudClbListenerCreate` function, following existing patterns with `d.GetOkExists`

## 3. Read Implementation

- [x] 3.1 Add `MaxConn`, `MaxCps`, `DataCompressMode` to `resourceTencentCloudClbListenerRead` function with nil-check before `d.Set`
- [x] 3.2 Add `ProxyProtocol` to `resourceTencentCloudClbListenerRead` function: set to `true` if `AttrFlags` contains "ProxyProtocol", otherwise `false`

## 4. Update Implementation

- [x] 4.1 Add `MaxConn`, `MaxCps`, `ProxyProtocol`, `DataCompressMode` to `resourceTencentCloudClbListenerUpdate` function with `d.HasChange` checks

## 5. Documentation

- [x] 5.1 Update `resource_tc_clb_listener.md` with Example Usage for the new parameters
- [x] 5.2 Run `make doc` to generate `website/docs/` markdown documentation (delegated to tfpacer-finalize skill)

## 6. Testing

- [x] 6.1 Add test cases for new parameters in `resource_tc_clb_listener_test.go`