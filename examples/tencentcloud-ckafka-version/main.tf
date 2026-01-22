# Example usage of tencentcloud_ckafka_version data source

# Query CKafka instances first
data "tencentcloud_ckafka_instances" "instances" {
  instance_name = "your-ckafka-instance-name"
}

# Query version information for the first instance
data "tencentcloud_ckafka_version" "version" {
  instance_id = data.tencentcloud_ckafka_instances.instances.instance_list[0].instance_id
}

# Output current version information
output "current_kafka_version" {
  description = "Current Kafka version"
  value       = data.tencentcloud_ckafka_version.version.kafka_version
}

output "current_broker_version" {
  description = "Current broker version"
  value       = data.tencentcloud_ckafka_version.version.cur_broker_version
}

output "available_broker_versions" {
  description = "Available broker versions for upgrade"
  value       = data.tencentcloud_ckafka_version.version.latest_broker_versions
}

# Example: Save version information to file
data "tencentcloud_ckafka_version" "version_with_output" {
  instance_id        = data.tencentcloud_ckafka_instances.instances.instance_list[0].instance_id
  result_output_file = "ckafka_version_info.json"
}