# dbdc-db-custom-node-network-mode-eniip Specification

## Purpose
TBD - created by syncing change add-dbdc-db-custom-node-network-mode-eniip. Update Purpose after sync.
## Requirements
### Requirement: network_mode computed field on dbdc db custom node
The `tencentcloud_dbdc_db_custom_node` resource SHALL expose a `network_mode`
parameter (TypeString, Computed) that is refreshed from the `NetworkMode`
field of the `DBCustomNode` struct returned by the `DescribeDBCustomNodes` API.
`network_mode` SHALL NOT be user-settable (no `Optional`), because no
Create/Update API accepts it.

#### Scenario: Read populates network_mode
- **WHEN** the provider reads a `tencentcloud_dbdc_db_custom_node` and the
  `DescribeDBCustomNodes` response's `DBCustomNode.NetworkMode` is not nil
- **THEN** the provider SHALL set `network_mode` in state to the value of
  `DBCustomNode.NetworkMode`

#### Scenario: Read skips network_mode when API omits it
- **WHEN** the provider reads a `tencentcloud_dbdc_db_custom_node` and
  `DBCustomNode.NetworkMode` is nil
- **THEN** the provider SHALL NOT call `d.Set("network_mode", ...)`, leaving
  the field untouched in state

#### Scenario: network_mode is not settable by the user
- **WHEN** a user includes `network_mode` in the `tencentcloud_dbdc_db_custom_node`
  configuration block
- **THEN** the provider SHALL reject the configuration, because `network_mode`
  is a Computed-only field

### Requirement: eni_ip computed field on dbdc db custom node
The `tencentcloud_dbdc_db_custom_node` resource SHALL expose an `eni_ip`
parameter (TypeString, Computed) that is refreshed from the `EniIP` field of
the `DBCustomNode` struct returned by the `DescribeDBCustomNodes` API. `eni_ip`
SHALL NOT be user-settable (no `Optional`), because no Create/Update API
accepts it. `eni_ip` holds the node access IP address when the
`NetworkModeCrossTenantENI` network mode is selected.

#### Scenario: Read populates eni_ip
- **WHEN** the provider reads a `tencentcloud_dbdc_db_custom_node` and the
  `DescribeDBCustomNodes` response's `DBCustomNode.EniIP` is not nil
- **THEN** the provider SHALL set `eni_ip` in state to the value of
  `DBCustomNode.EniIP`

#### Scenario: Read skips eni_ip when API omits it
- **WHEN** the provider reads a `tencentcloud_dbdc_db_custom_node` and
  `DBCustomNode.EniIP` is nil (e.g. when the node uses the
  `NetworkModePrivateLink` mode)
- **THEN** the provider SHALL NOT call `d.Set("eni_ip", ...)`, leaving the
  field untouched in state

#### Scenario: eni_ip is not settable by the user
- **WHEN** a user includes `eni_ip` in the `tencentcloud_dbdc_db_custom_node`
  configuration block
- **THEN** the provider SHALL reject the configuration, because `eni_ip` is a
  Computed-only field
