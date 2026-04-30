## 1. Schema Definition

- [x] 1.1 Add `default_deny_security_action_parameters` field to the `security_policy` schema in `resource_tc_teo_web_security_template.go` (TypeList, Optional, Computed, MaxItems: 1)
- [x] 1.2 Add `managed_rules` sub-block (TypeList, Optional, Computed, MaxItems: 1) with DenyActionParameters fields: block_ip, block_ip_duration, return_custom_page, response_code, error_page_id, stall (all TypeString, Optional)
- [x] 1.3 Add `other_modules` sub-block (TypeList, Optional, Computed, MaxItems: 1) with DenyActionParameters fields: block_ip, block_ip_duration, return_custom_page, response_code, error_page_id, stall (all TypeString, Optional)

## 2. Create Operation

- [x] 2.1 Add expand function for `default_deny_security_action_parameters` to convert Terraform schema to cloud API `DefaultDenySecurityActionParameters` struct in Create operation
- [x] 2.2 Set `SecurityPolicy.DefaultDenySecurityActionParameters` in the Create request when the field is specified

## 3. Read Operation

- [x] 3.1 Add flatten function for `DefaultDenySecurityActionParameters` to convert cloud API response to Terraform state in Read operation
- [x] 3.2 Handle nil checks for `DefaultDenySecurityActionParameters`, `ManagedRules`, and `OtherModules` in the flatten logic

## 4. Update Operation

- [x] 4.1 Add expand function for `default_deny_security_action_parameters` in Update operation to set `SecurityPolicy.DefaultDenySecurityActionParameters` in the Modify request

## 5. Unit Tests

- [x] 5.1 Add unit test cases for expand/flatten functions of `default_deny_security_action_parameters` in `resource_tc_teo_web_security_template_test.go`
- [x] 5.2 Add mock test for Create operation with `default_deny_security_action_parameters`
- [x] 5.3 Add mock test for Read operation with `DefaultDenySecurityActionParameters` response
- [x] 5.4 Add mock test for Update operation with `default_deny_security_action_parameters`
- [x] 5.5 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 6. Documentation

- [x] 6.1 Update `resource_tc_teo_web_security_template.md` to add `default_deny_security_action_parameters` example usage in the HCL example
