# TDMQ API Research Results

## Research Summary

This document summarizes the API research for optimizing RabbitMQ VIP instance update logic.

## Create API: CreateRabbitMQVipInstance

The Create API supports the following parameters:

### Core Configuration Parameters (Supported)
- `ZoneIds` - Availability zones (Required)
- `VpcId` - VPC ID (Required)
- `SubnetId` - Subnet ID (Required)
- `ClusterName` - Cluster name (Required)
- `NodeSpec` - Node specification (Optional, default: rabbit-vip-basic-1)
- `NodeNum` - Node count (Optional, default: 1 for single zone, 3 for multi-zone)
- `StorageSize` - Single node storage (Optional, default: 200G)
- `EnablePublicAccess` - Enable public network access (Optional, default: false)
- `Bandwidth` - Public network bandwidth in Mbps (Optional)
- `ClusterVersion` - Cluster version (Optional, default: 3.8.30)

### Other Parameters (Supported)
- `EnableCreateDefaultHaMirrorQueue` - Enable default mirror queue
- `AutoRenewFlag` - Auto renew flag
- `TimeSpan` - Purchase duration
- `PayMode` - Payment mode (0: postpaid, 1: prepaid)
- `ResourceTags` - Resource tags
- `EnableDeletionProtection` - Enable deletion protection

## Modify API: ModifyRabbitMQVipInstance

The Modify API supports the following parameters:

### Supported Parameters
- `InstanceId` - Instance ID (Required)
- `ClusterName` - Cluster name (Optional)
- `Remark` - Remark (Optional)
- `EnableDeletionProtection` - Enable deletion protection (Optional)
- `RemoveAllTags` - Remove all tags (Optional, default: false)
- `Tags` - Tags (Optional, full replacement)
- `EnableRiskWarning` - Enable risk warning (Optional)

### NOT Supported Parameters
The following core configuration parameters are **NOT** supported by the Modify API:
- ❌ `ZoneIds` - Cannot modify availability zones
- ❌ `VpcId` - Cannot modify VPC ID
- ❌ `SubnetId` - Cannot modify subnet ID
- ❌ `NodeSpec` - Cannot modify node specification
- ❌ `NodeNum` - Cannot modify node count
- ❌ `StorageSize` - Cannot modify storage size
- ❌ `Bandwidth` - Cannot modify public network bandwidth
- ❌ `EnablePublicAccess` - Cannot toggle public network access
- ❌ `ClusterVersion` - Cannot modify cluster version
- ❌ `PayMode` - Cannot modify payment mode
- ❌ `TimeSpan` - Cannot modify purchase duration
- ❌ `AutoRenewFlag` - Cannot modify auto renew flag
- ❌ `EnableCreateDefaultHaMirrorQueue` - Cannot modify mirror queue setting

## Alternative APIs

### ModifyCluster
This API is for Pulsar clusters, not RabbitMQ VIP instances. It supports:
- `ClusterId` - Pulsar cluster ID
- `ClusterName` - Cluster name
- `Remark` - Remark
- `PublicAccessEnabled` - Enable public access (can only be true)

This API is not applicable to RabbitMQ VIP instances.

### ModifyPublicNetworkSecurityPolicy
This API modifies public network security policies, not core instance parameters.

### ModifyRabbitMQPermission / ModifyRabbitMQUser / ModifyRabbitMQVirtualHost
These APIs are for RabbitMQ resources (permissions, users, virtual hosts), not instance configuration.

## Conclusion

Based on the API research:

1. **ModifyRabbitMQVipInstance API is very limited** - Only supports metadata updates (cluster name, tags, remark, deletion protection, risk warning)
2. **No alternative API for core parameter updates** - No other TDMQ API supports modifying node_spec, node_num, storage_size, band_width, or enable_public_access
3. **ForceNew is required for core parameters** - The following parameters must use ForceNew mechanism:
   - `zone_ids` (immutable - requires recreation)
   - `vpc_id` (immutable - requires recreation)
   - `subnet_id` (immutable - requires recreation)
   - `cluster_version` (immutable - requires recreation)
   - `node_spec` (requires ForceNew - no API support)
   - `node_num` (requires ForceNew - no API support)
   - `storage_size` (requires ForceNew - no API support)
   - `band_width` (requires ForceNew - no API support)
   - `enable_public_access` (requires ForceNew - no API support)

## Implementation Strategy

Given the API limitations, the implementation strategy is:

1. **Add ForceNew to immutable parameters** - Add ForceNew: true to parameters that require recreation
2. **Keep existing update logic** - Keep support for cluster_name and resource_tags updates
3. **Improve error handling** - Provide clear error messages when ForceNew parameters change
4. **Update documentation** - Clearly document which parameters require recreation
