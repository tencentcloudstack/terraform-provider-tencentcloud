## ADDED Requirements

### Requirement: Data source SHALL return all basic information fields
The `tencentcloud_clb_instances` data source SHALL return all basic information fields from the DescribeLoadBalancers API response.

#### Scenario: Return forward field
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `forward` field indicating whether it's a load balancer (1) or traditional load balancer (0)

#### Scenario: Return domain fields
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `domain` field for legacy domain (gradually deprecated)
- **THEN** the data source MUST return `load_balancer_domain` field for the current domain

### Requirement: Data source SHALL return all network configuration fields
The data source SHALL return comprehensive network configuration information including IPv6, ISP, and advanced networking features.

#### Scenario: Return IPv6 configuration
- **WHEN** querying CLB instances with IPv6 enabled
- **THEN** the data source MUST return `address_ipv6` field with the IPv6 address
- **THEN** the data source MUST return `ipv6_mode` field with value "IPv6Nat64" or "IPv6FullChain"
- **THEN** the data source MUST return `mix_ip_target` field indicating if mixed IPv4/IPv6 target binding is enabled

#### Scenario: Return anycast configuration
- **WHEN** querying anycast CLB instances
- **THEN** the data source MUST return `anycast_zone` field with the publishing region
- **THEN** empty string MUST be returned for non-anycast load balancers

#### Scenario: Return network egress
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `egress` field indicating network egress configuration

#### Scenario: Return local BGP status
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `local_bgp` field indicating if IP type is local BGP

### Requirement: Data source SHALL return all billing and lifecycle fields
The data source SHALL return complete billing information including charge type, prepaid attributes, and expiration time.

#### Scenario: Return charge type
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `charge_type` field with value "PREPAID" or "POSTPAID_BY_HOUR"

#### Scenario: Return prepaid attributes for prepaid instances
- **WHEN** querying prepaid CLB instances
- **THEN** the data source MUST return `prepaid_period` field with the purchased period in months
- **THEN** the data source MUST return `prepaid_renew_flag` field with auto-renewal setting
- **THEN** the data source MUST return `prepaid_cur_instance_deadline` field with instance deadline time

#### Scenario: Handle null prepaid attributes for postpaid instances
- **WHEN** querying postpaid CLB instances
- **THEN** the prepaid_* fields MUST NOT be set or return null values

#### Scenario: Return expiration time
- **WHEN** querying prepaid CLB instances
- **THEN** the data source MUST return `expire_time` field in format "YYYY-MM-DD HH:mm:ss"
- **THEN** null MUST be returned for postpaid instances

### Requirement: Data source SHALL return all logging configuration fields
The data source SHALL return CLS (Cloud Log Service) configuration for both access logs and health check logs.

#### Scenario: Return access log configuration
- **WHEN** querying CLB instances with access logging enabled
- **THEN** the data source MUST return `log_set_id` field with the logset ID
- **THEN** the data source MUST return `log_topic_id` field with the log topic ID

#### Scenario: Return health check log configuration
- **WHEN** querying CLB instances with health check logging enabled
- **THEN** the data source MUST return `health_log_set_id` field with the health check logset ID
- **THEN** the data source MUST return `health_log_topic_id` field with the health check log topic ID

#### Scenario: Handle null log configuration
- **WHEN** querying CLB instances without logging enabled
- **THEN** all log-related fields MUST return null or not be set

### Requirement: Data source SHALL return all security and isolation fields
The data source SHALL return security-related information including DDoS protection, isolation status, and blocking status.

#### Scenario: Return high-defense and SNAT configuration
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `open_bgp` field indicating high-defense LB (1) or regular (0)
- **THEN** the data source MUST return `snat` field indicating if SNAT is enabled
- **THEN** the data source MUST return `snat_pro` field indicating if SnatPro is enabled
- **THEN** the data source MUST return `snat_ips` field as JSON string containing SnatIp list when SnatPro is enabled

#### Scenario: Return isolation status
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `isolation` field with value 0 (not isolated) or 1 (isolated)
- **THEN** the data source MUST return `isolated_time` field in format "YYYY-MM-DD HH:mm:ss" for isolated instances

#### Scenario: Return blocking status
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `is_block` field indicating if VIP is blocked
- **THEN** the data source MUST return `is_block_time` field in format "YYYY-MM-DD HH:mm:ss" with block/unblock time

#### Scenario: Return DDoS protection capability
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `is_ddos` field indicating if DDoS protection package can be bound

### Requirement: Data source SHALL return performance and capacity fields
The data source SHALL return performance specification, instance type, and backend target count information.

#### Scenario: Return SLA type
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `sla_type` field with performance capacity spec
- **THEN** valid values MUST include: "clb.c1.small", "clb.c2.medium", "clb.c3.small", "clb.c3.medium", "clb.c4.small", "clb.c4.medium", "clb.c4.large", "clb.c4.xlarge", or empty string for non-performance instances

#### Scenario: Return exclusive instance type
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `exclusive` field with value 1 (exclusive) or 0 (non-exclusive)

#### Scenario: Return target count
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `target_count` field with the number of bound backend services

### Requirement: Data source SHALL return cluster and deployment fields
The data source SHALL return cluster information, NFV status, and availability zone configuration.

#### Scenario: Return multiple cluster IDs
- **WHEN** querying CLB instances associated with multiple clusters
- **THEN** the data source MUST return `cluster_ids` field as string array containing all cluster IDs
- **THEN** existing `cluster_id` field MUST continue to return the first cluster ID for backward compatibility

#### Scenario: Return cluster tag
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `cluster_tag` field with the layer-7 exclusive tag

#### Scenario: Return NFV information
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `nfv_info` field with value "l7nfv" for layer-7 NFV or empty string

#### Scenario: Return backup availability zones
- **WHEN** querying CLB instances with multi-AZ deployment
- **THEN** the data source MUST return `backup_zone_set` field as list of maps containing ZoneInfo objects
- **THEN** each ZoneInfo MUST include: zone_id, zone, zone_name, zone_region, local_zone fields

#### Scenario: Return availability zone affinity configuration
- **WHEN** querying CLB instances with AZ affinity enabled
- **THEN** the data source MUST return `available_zone_affinity_info` field as JSON string
- **THEN** JSON MUST contain: Enable (bool), ExitRatio (uint64), ReentryRatio (uint64)

### Requirement: Data source SHALL return advanced configuration fields
The data source SHALL return advanced settings including configuration ID, traffic passing, attribute flags, and exclusive cluster details.

#### Scenario: Return configuration ID
- **WHEN** querying CLB instances with personalized configuration
- **THEN** the data source MUST return `config_id` field with the configuration ID

#### Scenario: Return traffic passing setting
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `load_balancer_pass_to_target` field indicating if traffic from LB to backend is allowed

#### Scenario: Return attribute flags
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `attribute_flags` field as string array
- **THEN** possible values MUST include: "DeleteProtect", "UserInVisible", "BlockStatus", "NoLBNat", "BanStatus", "ShiftupFlag", "Stop", "NoVpcGw", "SgInTgw", "SharedLimitFlag", "WafFlag", "IsDomainCLB", "IPv6Snat", "HideDomain", "JumboFrame", "NoLBNatL4IPdc", "VpcGwL3Service", "Ipv62Flag", "Ipv62ExclusiveFlag", "BgpPro", "ToaClean"

#### Scenario: Return exclusive cluster details
- **WHEN** querying exclusive CLB instances
- **THEN** the data source MUST return `exclusive_cluster` field as JSON string containing ExclusiveCluster object
- **THEN** JSON MUST include L4Clusters, L7Clusters, ClassicalCluster information

#### Scenario: Return extra information
- **WHEN** querying CLB instances
- **THEN** the data source MUST return `extra_info` field as JSON string for internal use
- **THEN** general users MAY ignore this field

#### Scenario: Return endpoint association
- **WHEN** querying CLB instances associated with endpoints
- **THEN** the data source MUST return `associate_endpoint` field with the endpoint ID

### Requirement: Data source SHALL handle null values gracefully
All optional fields MUST be checked for null before accessing to prevent runtime panics.

#### Scenario: Safe null pointer access
- **WHEN** API returns null for optional fields
- **THEN** the data source MUST check for nil before dereferencing pointers
- **THEN** null fields MUST NOT be set in the Terraform state or set with appropriate zero values

#### Scenario: Nested object null handling
- **WHEN** API returns null for nested objects like PrepaidAttributes or ExclusiveCluster
- **THEN** the data source MUST check parent object for nil before accessing child fields
- **THEN** no fields from null objects MUST be set in the state

### Requirement: Data source MUST maintain backward compatibility
All new fields MUST be Computed attributes and existing fields MUST NOT be modified or removed.

#### Scenario: Preserve existing field behavior
- **WHEN** adding new fields to the schema
- **THEN** all existing fields MUST retain their current names, types, and descriptions
- **THEN** existing user configurations MUST continue to work without modification

#### Scenario: New fields are Computed
- **WHEN** defining new schema fields
- **THEN** all new fields MUST be marked as Computed: true
- **THEN** no new Required or Optional input fields MUST be added

### Requirement: Documentation MUST describe all fields comprehensively
The data source documentation MUST include clear descriptions for all returned fields.

#### Scenario: Document field descriptions
- **WHEN** adding new fields to the schema
- **THEN** each field MUST have a Description attribute explaining its purpose
- **THEN** descriptions MUST include valid value ranges or examples where applicable
- **THEN** descriptions MUST note when fields may return null

#### Scenario: Update generated documentation
- **WHEN** schema changes are complete
- **THEN** `make doc` MUST be run to regenerate `website/docs/d/clb_instances.html.markdown`
- **THEN** generated documentation MUST include all new fields in the Attributes Reference section
