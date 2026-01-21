# MQTT Message Enrichment Rule Example

This example demonstrates how to use the `tencentcloud_mqtt_message_enrichment_rule` resource to create and manage MQTT message enrichment rules.

## Key Features

- **Message Enrichment**: Add metadata, timestamps, and custom properties to MQTT messages
- **Rule Prioritization**: Control rule execution order with priority settings
- **Flexible Matching**: Use JSON-based conditions to match specific clients, topics, and users
- **Custom Actions**: Define response topics and user properties for enriched messages

## Important Notes

### Base64 Encoding

The MQTT API requires `condition` and `actions` fields to be Base64-encoded JSON. Terraform provides helper functions to make this easier:

```hcl
# Use base64encode(jsonencode()) to convert JSON to Base64
condition = base64encode(jsonencode({
  clientId = "sensor-*"
  topic    = "sensors/+/temperature"
}))
```

### Condition Structure

The `condition` field supports the following JSON structure:

```json
{
  "clientId": "client-pattern",  // Client ID pattern (supports wildcards)
  "username": "username-pattern", // Username pattern
  "topic": "topic/pattern/+"     // Topic pattern (supports MQTT wildcards)
}
```

### Actions Structure

The `actions` field defines what happens to matched messages:

```json
{
  "messageExpiryInterval": 3600,           // Message expiry in seconds
  "responseTopic": "enriched/${clientid}", // Target topic (supports variables)
  "correlationData": "${traceid}",         // Correlation data (supports variables)
  "userProperty": [                        // Custom properties to add
    {
      "key": "enriched-time",
      "value": "${timestamp}"
    }
  ]
}
```

### Variables

The following variables are available in actions:

- `${clientid}`: Original client ID
- `${traceid}`: Trace ID for correlation
- `${timestamp}`: Current timestamp

### Priority

- Lower numbers = higher priority
- High priority rules can override low priority rules
- UserProperty fields are merged across rules

### Status Values

- `0`: Undefined
- `1`: Activated (default)
- `2`: Deactivated

## Usage

1. Update the instance name in the data source
2. Customize the rule conditions and actions for your use case
3. Run `terraform plan` to review changes
4. Run `terraform apply` to create the rules

## Best Practices

1. **Use descriptive rule names** that indicate their purpose
2. **Set appropriate priorities** to ensure correct rule execution order
3. **Test conditions carefully** to avoid unintended message matching
4. **Monitor rule performance** and adjust priorities as needed
5. **Use meaningful user properties** to add valuable metadata