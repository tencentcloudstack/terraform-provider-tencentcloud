## 1. Resource Schema Changes

- [x] 1.1 Add `message` computed field (`TypeString`, `Computed: true`) to the resource schema in `resource_tc_teo_bind_security_template.go`.
- [x] 1.2 In the `resourceTencentCloudTeoBindSecurityTemplateRead` function, set `message` from `respData.Message` when `respData.Message` is not nil.

## 2. Documentation

- [x] 2.1 Update `resource_tc_teo_bind_security_template.md` example usage to include the `message` attribute in the output.

## 3. Testing

- [x] 3.1 Add test cases in `resource_tc_teo_bind_security_template_test.go` to verify the `message` schema field exists and is `Computed`.
- [x] 3.2 Add test case to verify the `message` attribute remains empty when the service layer does not populate `Message`.

## 4. Validation

- [x] 4.1 Verify the code compiles successfully with `go build ./...`. (Handled by tfpacer-finalize skill per workspace rules)
- [x] 4.2 Run `make doc` to regenerate website documentation. (Handled by tfpacer-finalize skill per workspace rules)
