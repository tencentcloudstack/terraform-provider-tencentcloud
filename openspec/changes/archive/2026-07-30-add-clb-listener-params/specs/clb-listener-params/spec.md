## ADDED Requirements

### Requirement: CLB listener supports MaxConn parameter
The `tencentcloud_clb_listener` resource SHALL support a `max_conn` parameter of type `TypeInt`, Optional and Computed. The parameter SHALL be passed to the CreateListener and ModifyListener APIs as `request.MaxConn`, and SHALL be read from the DescribeListeners API response (`response.Listeners.MaxConn`).

#### Scenario: Create listener with MaxConn
- **WHEN** user specifies `max_conn = 1000` in the listener configuration
- **THEN** the CreateListener API is called with `MaxConn = 1000`
- **AND** the value is persisted in Terraform state

#### Scenario: Read listener with MaxConn
- **WHEN** the DescribeListeners API returns a listener with `MaxConn = 1000`
- **THEN** the `max_conn` attribute in Terraform state is set to `1000`

#### Scenario: Update listener MaxConn
- **WHEN** user changes `max_conn` from `1000` to `500`
- **THEN** the ModifyListener API is called with `MaxConn = 500`

### Requirement: CLB listener supports MaxCps parameter
The `tencentcloud_clb_listener` resource SHALL support a `max_cps` parameter of type `TypeInt`, Optional and Computed. The parameter SHALL be passed to the CreateListener and ModifyListener APIs as `request.MaxCps`, and SHALL be read from the DescribeListeners API response (`response.Listeners.MaxCps`).

#### Scenario: Create listener with MaxCps
- **WHEN** user specifies `max_cps = 100` in the listener configuration
- **THEN** the CreateListener API is called with `MaxCps = 100`
- **AND** the value is persisted in Terraform state

#### Scenario: Read listener with MaxCps
- **WHEN** the DescribeListeners API returns a listener with `MaxCps = 100`
- **THEN** the `max_cps` attribute in Terraform state is set to `100`

#### Scenario: Update listener MaxCps
- **WHEN** user changes `max_cps` from `100` to `50`
- **THEN** the ModifyListener API is called with `MaxCps = 50`

### Requirement: CLB listener supports ProxyProtocol parameter
The `tencentcloud_clb_listener` resource SHALL support a `proxy_protocol` parameter of type `TypeBool`, Optional (not Computed). The parameter SHALL be passed to the CreateListener and ModifyListener APIs as `request.ProxyProtocol`. Since the DescribeListeners API does not return this field, the Terraform Read function SHALL NOT set `proxy_protocol` in state.

#### Scenario: Create TCP_SSL listener with ProxyProtocol
- **WHEN** user specifies `proxy_protocol = true` for a TCP_SSL listener
- **THEN** the CreateListener API is called with `ProxyProtocol = true`
- **AND** the value is persisted in Terraform state as provided

#### Scenario: Update listener ProxyProtocol
- **WHEN** user changes `proxy_protocol` from `false` to `true`
- **THEN** the ModifyListener API is called with `ProxyProtocol = true`

### Requirement: CLB listener supports DataCompressMode parameter
The `tencentcloud_clb_listener` resource SHALL support a `data_compress_mode` parameter of type `TypeString`, Optional and Computed. The parameter SHALL be passed to the CreateListener and ModifyListener APIs as `request.DataCompressMode`, and SHALL be read from the DescribeListeners API response (`response.Listeners.DataCompressMode`).

#### Scenario: Create listener with DataCompressMode
- **WHEN** user specifies `data_compress_mode = "transparent"` in the listener configuration
- **THEN** the CreateListener API is called with `DataCompressMode = "transparent"`
- **AND** the value is persisted in Terraform state

#### Scenario: Read listener with DataCompressMode
- **WHEN** the DescribeListeners API returns a listener with `DataCompressMode = "compatibility"`
- **THEN** the `data_compress_mode` attribute in Terraform state is set to `"compatibility"`

#### Scenario: Update listener DataCompressMode
- **WHEN** user changes `data_compress_mode` from `"transparent"` to `"compatibility"`
- **THEN** the ModifyListener API is called with `DataCompressMode = "compatibility"`