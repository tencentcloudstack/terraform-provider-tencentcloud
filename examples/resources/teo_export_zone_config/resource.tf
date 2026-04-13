# terraform-provider-tencentcloud Example for teo_export_zone_config

resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-xxxxxxxxxxxxx"

  # Optional: specify which configuration types to export
  # If not specified, all configuration types will be exported
  export_types = ["L7AccelerationConfig"]
}

output "exported_config" {
  description = "The exported zone configuration"
  value       = tencentcloud_teo_export_zone_config.example.content
}

output "exported_config_json" {
  description = "The exported zone configuration in JSON format"
  value       = jsondecode(tencentcloud_teo_export_zone_config.example.content)
}
