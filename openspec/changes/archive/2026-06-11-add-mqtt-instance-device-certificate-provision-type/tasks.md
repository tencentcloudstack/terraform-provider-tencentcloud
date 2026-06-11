## 1. Schema and CRUD Modification

- [x] 1.1 Modify `device_certificate_provision_type` schema field in `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go` from `Computed: true` to `Optional: true, Computed: true`
- [x] 1.2 Add `request.DeviceCertificateProvisionType` assignment in the Update function (`ResourceTencentCloudMqttInstanceUpdate`) using the `d.GetOk` pattern, consistent with other string fields like `x509_mode`

## 2. Unit Tests

- [x] 2.1 Add unit test cases in `tencentcloud/services/mqtt/resource_tc_mqtt_instance_test.go` to cover setting and updating `device_certificate_provision_type` using gomonkey mock approach

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/mqtt/resource_tc_mqtt_instance.md` to reflect that `device_certificate_provision_type` is now a configurable parameter with example usage
