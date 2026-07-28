## 1. Service Layer Changes

- [x] 1.1 Rewrite `DescribeTeoBindSecurityTemplateById` in `service_tencentcloud_teo.go` to use `DescribeSecurityTemplateBindings` API instead of `DescribeZones` + `DescribeWebSecurityTemplates`. The new implementation calls `DescribeSecurityTemplateBindings` with the given `zoneId` and `templateId`, iterates through `response.SecurityTemplate[].TemplateScope[].EntityStatus[]` to find the matching `entity`, and returns the `EntityStatus` (including `Message`).
- [x] 1.2 Remove the now-unused `describeTeoAllZoneIds` helper function from `service_tencentcloud_teo.go`.

## 2. Resource Schema Changes

- [x] 2.1 Add `message` computed field (`TypeString`, `Computed: true`) to the resource schema in `resource_tc_teo_bind_security_template.go`.
- [x] 2.2 In the `resourceTencentCloudTeoBindSecurityTemplateRead` function, set `message` from `respData.Message` when `respData.Message` is not nil.

## 3. Documentation

- [x] 3.1 Update `resource_tc_teo_bind_security_template.md` example usage to include the `message` attribute in the output.

## 4. Testing

- [x] 4.1 Add test cases in `resource_tc_teo_bind_security_template_test.go` to verify the `message` attribute is set correctly during Read when the API returns a non-nil `Message` value.
- [x] 4.2 Add test cases to verify the `message` attribute is not set when the API returns a nil `Message` value.

## 5. Validation

- [x] 5.1 Verify the code compiles successfully with `go build ./...`. (Handled by tfpacer-finalize skill per workspace rules)
- [x] 5.2 Run `make doc` to regenerate website documentation. (Handled by tfpacer-finalize skill per workspace rules)