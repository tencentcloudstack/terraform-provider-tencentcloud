Use this data source to query the events list supported by the audit.

Example Usage
```hcl
data "tencentcloud_audit_events" "events" {
  start_time  = "1727433841"
  end_time    = "1727437441"
  max_results = 50

  lookup_attributes {
    attribute_key   = "ResourceType"
    attribute_value = "cvm"
  }

  lookup_attributes {
    attribute_key   = "OnlyRecordNotSeen"
    attribute_value = "0"
  }

  lookup_attributes {
    attribute_key   = "EventPlatform"
    attribute_value = "0"
  }

  is_return_location = 1
}
```