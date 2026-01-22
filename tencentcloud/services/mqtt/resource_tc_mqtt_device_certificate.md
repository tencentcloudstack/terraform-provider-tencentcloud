Provides a resource to create a MQTT device certificate

Example Usage

```hcl
resource "tencentcloud_mqtt_device_certificate" "example" {
  instance_id        = "mqtt-zxjwkr98"
  device_certificate = <<-EOF
-----BEGIN CERTIFICATE-----
MIIDgzCCAmugAwIBAgIQbWhvyXL8dmDtyID8f0kLlTANBgkqhkiG9w0BAQsFADA/
MQswCQYDVQQGEwJDTjEbMBkGA1UEChMSTXkgQ0EgT3JnYW5pemF0aW9uMRMwEQYD
VQQDEwpNeSBDQSBSb290MB4XDTI1MDUwOTA4MDMzOVoXDTI2MDUwOTA4MDMzOVow
WjELMAkGA1UEBhMCQ04xHDAaBgNVBAoTE0NsaWVudCBPcmdhbml6YXRpb24xLTAr
BgNVBAMTJGQyNGE1MDYxLTNjYjktNGVkZi04MGJhLTBmODNkY2IyNDM2MTCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALZXOJ47sZemb1/wDoBWE7Kgy2EO
UsbsUu/YYhvHWR/ePGdbXsBx0C4fdD38aqZkSRK7R/YsgarzFRF8ozKMRvyvdga4
YReXQaOhop3HL8oVZm/NW8x7GyruD4D1CP6/odtFWtG2JWf1UH/L3YUieR3D9X7S
LFZXBICdN69qnP05hUIiRiQ7yRfhs6sWdCH8YPTu6LXintWGHAg9RCw/8ewuwh/P
g4WGej4ycQcwBQ85zNMF0zXmkNWE4BdJvO/+2TgN0S6rXkRH0sBrghQeURLzmzv8
5HUIj740wEZpC37SLeZbYnp2RpbDAooTOBkyrLpJ5d0bV4441GpjwGQAuAUCAwEA
AaNgMF4wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEF
BQcDAjAMBgNVHRMBAf8EAjAAMB8GA1UdIwQYMBaAFM+RdbWwsz5TpRVX8ipQqUAF
urZVMA0GCSqGSIb3DQEBCwUAA4IBAQBOLOJKnmOtiiSlk6a4cNAQROWwDxQeWlZz
4NuPGPpjx3OHQZTi9PGeeJJtL6VyPTBdrETjfriTU+vzsYEpYs303B04hcCpHMgc
SMS14V8iSuRnXPXpSrX2/a3B6KNTeXd5662k1FCwZG/bGVvE/Q1sAu6Ls/1Q1XfY
stvJQTb4MEKa64d1e+58yTp2UhmxyfWTFy7LqguIGZgTd8Oz8ISJjBg0ca+Co/gN
uD7+CB4HqiiiN3o3meihJePo68foyvwnGntrx0KKlas8NJxCkWmM/HHpwjxz7eJZ
ulX9ykqE3WqMkWMcVTzx/wAhvixKckQD3+bZzBvOqerMpkRMpGOC
-----END CERTIFICATE-----
EOF
  ca_sn              = "1"
  client_id          = "d24a5061-3cb9-4edf-80ba-0f83dcb24361"
  format             = "PEM"
  status             = "ACTIVE"
}
```

Import

MQTT device certificate can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_device_certificate.example mqtt-zxjwkr98#6d686fc972fc7660edc880fc7f490b95
```
