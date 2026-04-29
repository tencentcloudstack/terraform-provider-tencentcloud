## 1. Service Layer

- [x] 1.1 Add `DescribeTeoAliasDomainById` method to `tencentcloud/services/teo/service_tencentcloud_teo.go` - calls DescribeAliasDomains API with Filters (alias-name exact match), Limit=1000, handles pagination, uses resource.Retry(tccommon.ReadRetryTimeout, ...), returns *teo.AliasDomain or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_alias_domain.go` with schema definition including zone_id (Required, ForceNew), alias_name (Required, ForceNew), target_name (Required), cert_type (Optional), cert_id (Optional, TypeList), status/forbid_mode/created_on/modified_on (Computed), and Importer support
- [x] 2.2 Implement `resourceTencentCloudTeoAliasDomainCreate` function - calls CreateAliasDomain API with retry, sets composite ID (zone_id#alias_name), validates response not empty, calls Read at end
- [x] 2.3 Implement `resourceTencentCloudTeoAliasDomainRead` function - parses composite ID, calls service layer DescribeTeoAliasDomainById, sets computed fields (checking nil before setting), handles resource not found, does NOT set cert_type/cert_id from API
- [x] 2.4 Implement `resourceTencentCloudTeoAliasDomainUpdate` function - checks d.HasChange() for target_name, cert_type, cert_id, calls ModifyAliasDomain API with retry when mutable args change
- [x] 2.5 Implement `resourceTencentCloudTeoAliasDomainDelete` function - parses composite ID, calls DeleteAliasDomain API with retry, passes single alias_name in AliasNames

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_alias_domain` resource in `tencentcloud/provider.go` resource map, mapping to `teo.ResourceTencentCloudTeoAliasDomain()`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_alias_domain_test.go` with unit tests using gomonkey mock approach for Create, Read, Update, Delete functions

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_alias_domain.md` with resource description, example usage (including cert_type/cert_id example), and import section noting composite ID format zone_id#alias_name
- [x] 5.2 Update `tencentcloud/provider.md` to include the new resource entry
