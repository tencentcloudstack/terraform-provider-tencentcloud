# Design: Organization IP Whitelist Config Resource

## Overview

This document describes the design and implementation details for the `tencentcloud_organization_ip_whitelist_config` resource.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Terraform User                           │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                 terraform-provider-tencentcloud                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │    resource_tc_organization_ip_whitelist_config.go      │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐      │   │
│  │  │ Create  │  │  Read   │  │ Update  │  │ Delete  │      │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘      │   │
│  └───────┼────────────┼────────────┼────────────┼────────────┘   │
│          │            │            │            │               │
└──────────┼────────────┼────────────┼────────────┼───────────────┘
           │            │            │            │
           ▼            ▼            ▼            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    OrganizationService                          │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  DescribeOrganizationIPWhitelistConfigById(ctx, id)  │    │
│  │  UpdateOrganizationIPWhitelistConfig(ctx, id, list)   │    │
│  └──────────────────────────┬──────────────────────────────┘    │
└────────────────────────────┼───────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│              TencentCloud Organization API                       │
│  ┌──────────────────┐    ┌──────────────────┐                   │
│  │ GetIPWhitelist   │    │ UpdateIPWhitelist│                   │
│  │   (ZoneId)       │    │ (ZoneId, IPs)    │                   │
│  └──────────────────┘    └──────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
```

## Resource ID Design

**Resource ID**: `zone_id` (String)

The `zone_id` is the unique identifier for the organization zone and will be used as the Terraform resource ID.

## Schema Design

```go
Schema: map[string]*schema.Schema{
    "zone_id": {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "Zone ID.",
    },
    "ip_whitelist": {
        Type:        schema.TypeList,
        Elem:        &schema.Schema{Type: schema.TypeString},
        Required:    true,
        Description: "IP whitelist entries.",
    },
}
```

## API Mapping

| Terraform Field | UpdateIPWhitelist Parameter | GetIPWhitelist Response |
|-----------------|------------------------------|-------------------------|
| `zone_id` | `ZoneId` | `ZoneId` (input) |
| `ip_whitelist` | `IpWhitelist.N` | `IpWhitelist` (array) |

## Implementation Details

### Create Function

1. Read `zone_id` and `ip_whitelist` from resource data
2. Build `UpdateIPWhitelistRequest` with `ZoneId` and `IpWhitelist`
3. Call API via service layer
4. Set resource ID to `zone_id`
5. Call Read function to sync state

### Read Function

1. Extract `zone_id` from resource ID
2. Call `GetIPWhitelist` API with `ZoneId`
3. Set `ip_whitelist` in Terraform state
4. Handle not found case (set ID to empty)

### Update Function

1. Read `zone_id` (unchanged due to ForceNew) and `ip_whitelist`
2. Build `UpdateIPWhitelistRequest`
3. Call API to update whitelist
4. Call Read function to sync state

### Delete Function

1. Extract `zone_id` from resource ID
2. Call `UpdateIPWhitelist` API with empty `IpWhitelist`
3. Return success

## Service Layer Methods

```go
// DescribeOrganizationIPWhitelistConfigById retrieves the IP whitelist for a zone
func (s *OrganizationService) DescribeOrganizationIPWhitelistConfigById(ctx context.Context, zoneId string) ([]string, error)

// UpdateOrganizationIPWhitelistConfig updates the IP whitelist for a zone
func (s *OrganizationService) UpdateOrganizationIPWhitelistConfig(ctx context.Context, zoneId string, ipWhitelist []string) error
```

## Error Handling

- API errors should be retried using `resource.Retry` with appropriate retryable error codes
- Not found errors should be handled gracefully by clearing the resource ID
- Validation errors should be reported before making API calls

## Import

Import is not supported because the GetIPWhitelist API does not support listing all zones. Users must know the `zone_id` to import.
