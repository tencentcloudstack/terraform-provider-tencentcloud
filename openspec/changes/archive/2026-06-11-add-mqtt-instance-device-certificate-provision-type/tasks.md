## 1. Schema and CRUD Modification

- [x] 1.1 Modify `device_certificate_provision_type` schema field in `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go` from `Computed: true` to `Optional: true, Computed: true`
- [x] 1.2 Add `request.DeviceCertificateProvisionType` assignment in the Update function (`ResourceTencentCloudMqttInstanceUpdate`) using the `d.GetOk` pattern, consistent with other string fields like `x509_mode`
- [x] 1.3 Add `message_rate` schema field as `Optional: true, Computed: true` (TypeInt)
- [x] 1.4 Add `use_default_server_cert` schema field as `Optional: true, Computed: true` (TypeBool)
- [x] 1.5 Add `request.MessageRate` assignment in the Update function using `d.GetOkExists` pattern
- [x] 1.6 Add `request.UseDefaultServerCert` assignment in the Update function using `d.GetOkExists` pattern
- [x] 1.7 Add `message_rate` and `use_default_server_cert` to the post-create modify logic in the Create function
- [x] 1.8 Add reading of `MessageRate` and `UseDefaultServerCert` from `DescribeInstance` response in the Read function

## 2. Unit Tests

- [x] 2.1 Add unit test cases in `tencentcloud/services/mqtt/resource_tc_mqtt_instance_test.go` to cover setting and updating `device_certificate_provision_type` using gomonkey mock approach
- [x] 2.2 Add unit test cases for `message_rate` (schema, read, read nil, update)
- [x] 2.3 Add unit test cases for `use_default_server_cert` (schema, read, read nil, update)

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/mqtt/resource_tc_mqtt_instance.md` to reflect that `device_certificate_provision_type`, `message_rate`, and `use_default_server_cert` are now configurable parameters with example usage
