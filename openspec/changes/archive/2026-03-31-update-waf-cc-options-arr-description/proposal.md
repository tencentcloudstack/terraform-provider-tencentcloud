# Proposal: Update tencentcloud_waf_cc options_arr Description

## Background

The `options_arr` field in the `tencentcloud_waf_cc` resource currently has an outdated description that doesn't match the latest API documentation. According to the official API documentation (https://cloud.tencent.com/document/api/627/97646), the `OptionsArr` field now includes additional field types and match operators that are not documented in the current schema.

Additionally, the current description uses backticks (\`\`) which need to be changed to double quotes ("") for consistency with Terraform provider standards.

## Goal

Update the `options_arr` field description in `tencentcloud_waf_cc` resource to:
1. Reflect the latest API documentation (updated 2026-03-06)
2. Include all newly supported key types: `URL`, `Method`, `Post`, `Referer`, `Cookie`, `User-Agent`, `CustomHeader`, `IPLocation`, `CaptchaRisk`, `CaptchaDeviceRisk`, `CaptchaScore`
3. Document all match operators for each key type according to the latest API specification
4. Replace backticks (\`\`) with double quotes ("") in the Description field

## Scope

### In Scope
- Update `options_arr` field Description in `resource_tc_waf_cc.go` Schema definition
- Update match operator values for all key types based on latest API documentation
- Add newly supported key types: `URL` and `IPLocation`
- Format code with `go fmt` after changes

### Out of Scope
- Functional logic changes
- Other fields or resources
- Test files or documentation files

## Files to Modify

| File | Change |
|------|--------|
| `tencentcloud/services/waf/resource_tc_waf_cc.go` | Update `options_arr` field Description (lines 82-103) |

## API Reference

- **API Documentation**: https://cloud.tencent.com/document/api/627/97646
- **Last Updated**: 2026-03-06 03:48:41
- **Field**: OptionsArr

## Expected Outcome

The `options_arr` field will have an accurate, up-to-date description that:
- Matches the official API documentation
- Uses proper formatting (double quotes instead of backticks)
- Includes all supported key types and match operators
- Helps users correctly configure CC attack protection rules
