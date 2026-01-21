# Example usage of tencentcloud_mqtt_message_enrichment_rule resource

# Query MQTT instances first
data "tencentcloud_mqtt_instances" "instances" {
  instance_name = "your-mqtt-instance-name"
}

# Create a message enrichment rule for temperature sensors
resource "tencentcloud_mqtt_message_enrichment_rule" "temperature_enrichment" {
  instance_id = data.tencentcloud_mqtt_instances.instances.instance_list[0].instance_id
  rule_name   = "temperature-sensor-enrichment"
  
  # Rule matching condition (Base64 encoded JSON)
  condition = base64encode(jsonencode({
    clientId = "sensor-*"
    topic    = "sensors/+/temperature"
    username = "iot-device"
  }))
  
  # Rule actions (Base64 encoded JSON)
  actions = base64encode(jsonencode({
    messageExpiryInterval = 3600
    responseTopic        = "enriched/temperature/${clientid}"
    correlationData      = "${traceid}"
    userProperty = [
      {
        key   = "enriched-time"
        value = "${timestamp}"
      },
      {
        key   = "data-source"
        value = "rule-engine"
      },
      {
        key   = "sensor-type"
        value = "temperature"
      }
    ]
  }))
  
  priority = 10
  status   = 1  # 1 = activated
  remark   = "Temperature sensor data enrichment with timestamp and metadata"
}

# Create a message enrichment rule for humidity sensors
resource "tencentcloud_mqtt_message_enrichment_rule" "humidity_enrichment" {
  instance_id = data.tencentcloud_mqtt_instances.instances.instance_list[0].instance_id
  rule_name   = "humidity-sensor-enrichment"
  
  condition = base64encode(jsonencode({
    clientId = "sensor-*"
    topic    = "sensors/+/humidity"
  }))
  
  actions = base64encode(jsonencode({
    messageExpiryInterval = 1800
    responseTopic        = "enriched/humidity/${clientid}"
    correlationData      = "${traceid}"
    userProperty = [
      {
        key   = "enriched-time"
        value = "${timestamp}"
      },
      {
        key   = "data-source"
        value = "rule-engine"
      },
      {
        key   = "sensor-type"
        value = "humidity"
      },
      {
        key   = "unit"
        value = "percentage"
      }
    ]
  }))
  
  priority = 20
  status   = 1
  remark   = "Humidity sensor data enrichment with unit information"
}

# Create a high-priority rule for emergency alerts
resource "tencentcloud_mqtt_message_enrichment_rule" "emergency_alert" {
  instance_id = data.tencentcloud_mqtt_instances.instances.instance_list[0].instance_id
  rule_name   = "emergency-alert-enrichment"
  
  condition = base64encode(jsonencode({
    topic = "alerts/emergency/+"
  }))
  
  actions = base64encode(jsonencode({
    messageExpiryInterval = 60
    responseTopic        = "alerts/enriched/emergency/${clientid}"
    correlationData      = "${traceid}"
    userProperty = [
      {
        key   = "alert-time"
        value = "${timestamp}"
      },
      {
        key   = "priority"
        value = "HIGH"
      },
      {
        key   = "processed-by"
        value = "emergency-rule-engine"
      }
    ]
  }))
  
  priority = 1  # Highest priority
  status   = 1
  remark   = "High-priority emergency alert enrichment"
}

# Output rule information
output "temperature_rule_id" {
  description = "Temperature enrichment rule ID"
  value       = tencentcloud_mqtt_message_enrichment_rule.temperature_enrichment.rule_id
}

output "humidity_rule_id" {
  description = "Humidity enrichment rule ID"
  value       = tencentcloud_mqtt_message_enrichment_rule.humidity_enrichment.rule_id
}

output "emergency_rule_id" {
  description = "Emergency alert enrichment rule ID"
  value       = tencentcloud_mqtt_message_enrichment_rule.emergency_alert.rule_id
}