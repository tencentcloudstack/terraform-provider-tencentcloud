# Proposal: Add Organization IP Whitelist Config Resource

## What

Add a new Terraform resource `tencentcloud_organization_ip_whitelist_config` to manage IP whitelist for organization CIC (Cloud Identity Center) zones.

## Why

Currently, the provider does not support managing IP whitelist configurations for organization zones. Users need to manage IP whitelists through other means (console or raw API calls). This feature provides Terraform-native management of organization IP whitelists.

## Background

The Organization service provides IP whitelist functionality for CIC zones. The IP whitelist restricts access to organization resources based on source IP addresses.

### API Details

**GetIPWhitelist API:**
- Input: `ZoneId` (String)
- Output: `IpWhitelist` (Array of String)

**UpdateIPWhitelist API:**
- Input: `ZoneId` (String), `IpWhitelist.N` (Array of String, max 100 entries)
- Output: `Success` (Boolean)

### Schema Design

| Field | Type | Required | ForceNew | Description |
|-------|------|----------|----------|-------------|
| `zone_id` | String | Yes | Yes | Zone ID (unique identifier) |
| `ip_whitelist` | List(String) | Yes | No | IP whitelist entries |

## Scope

- Create resource implementation
- Create markdown documentation
- Create unit test code
- Register resource in provider

## Out of Scope

- Data source implementation (not required)
- Import functionality (not supported by API)
