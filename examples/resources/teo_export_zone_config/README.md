# TEO Export Zone Config Example

This example demonstrates how to export EdgeOne (TEO) zone configuration using the `tencentcloud_teo_export_zone_config` resource.

## Overview

The `tencentcloud_teo_export_zone_config` resource allows you to export zone configurations from EdgeOne. This is useful for:
- Backup and version control of your zone configurations
- Migrating configurations between zones or environments
- Auditing and reviewing configuration changes
- Implementing infrastructure as code practices

## Features Demonstrated

1. **Basic Export**: Export all configuration types for a zone
2. **Selective Export**: Export specific configuration types (e.g., L7AccelerationConfig)
3. **Accessing Exported Content**: Retrieve and use the exported configuration
4. **JSON Parsing**: Decode the exported JSON configuration for further processing

## Usage

### Prerequisites

- Configure your Tencent Cloud credentials using environment variables:
  ```bash
  export TENCENTCLOUD_SECRET_ID="your-secret-id"
  export TENCENTCLOUD_SECRET_KEY="your-secret-key"
  ```

- You need a valid TEO zone ID. You can find it in the TEO console or using the `tencentcloud_teo_zones` data source.

### Run the Example

```bash
terraform init
terraform plan
terraform apply
```

## Example Configuration

### Basic Export (All Configuration Types)

```hcl
resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-xxxxxxxxxxxxx"
}
```

### Selective Export (Specific Configuration Types)

```hcl
resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-xxxxxxxxxxxxx"
  export_types = ["L7AccelerationConfig"]
}
```

### Access Exported Content

```hcl
output "exported_config" {
  description = "The exported zone configuration"
  value       = tencentcloud_teo_export_zone_config.example.content
}

output "exported_config_json" {
  description = "The exported zone configuration in JSON format"
  value       = jsondecode(tencentcloud_teo_export_zone_config.example.content)
}
```

## Configuration Types

The following configuration types can be exported:

- `L7AccelerationConfig`: L7 acceleration configuration (includes "Site Acceleration - Global Acceleration Configuration" and "Site Acceleration - Rule Engine")

Note: Additional configuration types may be added in the future. When exporting all types, please note the size of the exported file. It is recommended to specify the configuration types to be exported to control the size of the request response packet load.

## Resource ID Format

The resource ID is constructed as: `zoneId#exportTypes`

Example:
- `zone-abc123#` (export all types)
- `zone-abc123#L7AccelerationConfig` (export only L7AccelerationConfig)

## Notes

- This resource is a virtual resource that exports configuration on demand. It does not persist any state on the cloud side.
- The `zone_id` parameter is immutable (ForceNew), changing it will create a new resource.
- The `export_types` parameter is also immutable (ForceNew), changing it will create a new resource.
- The exported content is returned as a JSON string. You can use `jsondecode()` function to parse it in Terraform.

## Import

You can import an existing export configuration using:

```bash
terraform import tencentcloud_teo_export_zone_config.example zone-id#export-type
```

Example:
```bash
terraform import tencentcloud_teo_export_zone_config.example zone-abc123#L7AccelerationConfig
```

## API Reference

- Resource: [tencentcloud_teo_export_zone_config](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs/resources/teo_export_zone_config)
- API: [ExportZoneConfig](https://cloud.tencent.com/document/product/1552/xxxxx)
