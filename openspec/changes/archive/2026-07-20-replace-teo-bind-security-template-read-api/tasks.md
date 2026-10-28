## 1. Service layer

- [x] 1.1 Replace `DescribeTeoBindSecurityTemplateById` in `tencentcloud/services/teo/service_tencentcloud_teo.go` to drop `DescribeSecurityTemplateBindings` and use `DescribeZones` + `DescribeWebSecurityTemplates`
- [x] 1.2 Add `describeTeoAllZoneIds` helper paging `DescribeZones` with `Limit=100` via `UseTeoV20220901Client()`
- [x] 1.3 Batch `DescribeWebSecurityTemplates` by at most 100 zone IDs per request, filter by `TemplateId` and `BindDomains.Domain`, return synthesized `EntityStatus{Entity, Status}`
- [x] 1.4 Wrap `DescribeZones` / `DescribeWebSecurityTemplates` calls with `resource.Retry(tccommon.ReadRetryTimeout, ...)` using `tccommon.RetryError`

## 2. Resource and extension code

- [x] 2.1 Update `resourceTencentCloudTeoBindSecurityTemplateRead` to log `[CRUD] teo_bind_security_template id=%s` with `d.Id()` before `d.SetId("")`
- [x] 2.2 Clean up `resourceTeoBindSecurityTemplateCreateStateRefreshFunc_0_0` in `resource_tc_teo_bind_security_template_extension.go`: remove unused `DescribeSecurityTemplateBindingsRequest` field and guard nil `resp.Status`

## 3. Tests

- [x] 3.1 Add gomonkey-based unit tests in `resource_tc_teo_bind_security_template_test.go` for: read success, read not-found, read no-zone, read >100 zones (batching), and schema validation
- [x] 3.2 Run `go test ./tencentcloud/services/teo/ -run "TestTeoBindSecurityTemplate_" -v -count=1 -gcflags="all=-l"` and confirm all pass

## 4. Docs and changelog

- [x] 4.1 Verify `resource_tc_teo_bind_security_template.md` example/import content is consistent (no schema change; no edit required beyond the existing docs commit)
- [ ] 4.2 Update `.changelog/4261.txt` (or add a new changelog entry in the finalize phase) to describe the read-path API replacement as a bug-fix/enhancement
