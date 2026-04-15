## 1. Schema Definition

- [x] 1.1 Define resource schema in resource_tc_teo_dns_record_v9.go
- [x] 1.2 Add resource identifier using composite ID format (zone_id#record_id)
- [x] 1.3 Add required fields: zone_id, name, type, content
- [x] 1.4 Add optional fields: location, ttl, weight, priority
- [x] 1.5 Add computed fields: status, created_on, modified_on
- [x] 1.6 Add validation functions for field ranges (ttl: 60-86400, weight: -1 to 100, priority: 0-50)
- [x] 1.7 Add DiffSuppressFunc for type-specific fields (location/weight for A/AAAA/CNAME, priority for MX)
- [x] 1.8 Add Timeouts configuration block for async operations (create, update, default)

## 2. Service Layer Implementation

- [x] 2.1 Add CreateDnsRecord API call in service_tencentcloud_teo.go
- [x] 2.2 Add DescribeDnsRecords API call with filtering support in service_tencentcloud_teo.go
- [x] 2.3 Add ModifyDnsRecords API call in service_tencentcloud_teo.go
- [x] 2.4 Add DeleteDnsRecords API call in service_tencentcloud_teo.go
- [x] 2.5 Implement helper functions for parsing and generating composite resource IDs
- [x] 2.6 Implement retry mechanism for async operations using helper.Retry()

## 3. CRUD Functions Implementation

- [x] 3.1 Implement Create function in resource_tc_teo_dns_record_v9.go
- [x] 3.2 Implement Read function with retry logic for eventual consistency in resource_tc_teo_dns_record_v9.go
- [x] 3.3 Implement Update function with field merging and retry logic in resource_tc_teo_dns_record_v9.go
- [x] 3.4 Implement Delete function with retry logic in resource_tc_teo_dns_record_v9.go
- [x] 3.5 Implement error handling and validation in CRUD functions
- [x] 3.6 Register the resource in the provider's ResourcesMap

## 4. Unit Tests Implementation

- [x] 4.1 Create test file resource_tc_teo_dns_record_v9_test.go
- [x] 4.2 Mock CreateDnsRecord API for create operations testing
- [x] 4.3 Mock DescribeDnsRecords API for read operations testing
- [x] 4.4 Mock ModifyDnsRecords API for update operations testing
- [x] 4.5 Mock DeleteDnsRecords API for delete operations testing
- [x] 4.6 Write test case for creating A record with basic parameters
- [x] 4.7 Write test case for creating CNAME record
- [x] 4.8 Write test case for creating MX record with priority
- [x] 4.9 Write test case for creating record with location and weight
- [x] 4.10 Write test case for reading existing record
- [x] 4.11 Write test case for updating record content
- [x] 4.12 Write test case for updating multiple fields
- [x] 4.13 Write test case for deleting existing record
- [x] 4.14 Write test case for invalid resource ID format
- [x] 4.15 Write test case for missing required fields
- [x] 4.16 Write test case for field value validation (ttl, weight, priority ranges)
- [x] 4.17 Write test case for type-specific field validation

## 5. Documentation

- [x] 5.1 Create resource example file resource_tc_teo_dns_record_v9.md with usage examples
- [x] 5.2 Add example for creating A record
- [x] 5.3 Add example for creating CNAME record with location and weight
- [x] 5.4 Add example for creating MX record with priority
- [x] 5.5 Add example for updating record fields
- [x] 5.6 Add example for deleting record
- [x] 5.7 Add timeout configuration examples
- [ ] 5.8 Generate website documentation using `make doc` command

## 6. Code Quality Verification

- [ ] 6.1 Run `go fmt` to format the code (deferred to finalization)
- [ ] 6.2 Run `go vet` to check for common errors (deferred to finalization)
- [x] 6.3 Run unit tests using `go test` to verify CRUD operations
- [x] 6.4 Verify the resource can be imported in provider without errors
- [ ] 6.5 Run acceptance tests with TF_ACC=1 (if environment available) (deferred to finalization)

## 7. Final Checks

- [x] 7.1 Verify all artifacts are created: proposal.md, design.md, specs/teo-dns-record-v9/spec.md, tasks.md
- [x] 7.2 Verify resource code follows project conventions and best practices
- [x] 7.3 Verify error messages are clear and user-friendly
- [x] 7.4 Verify timeout configuration works correctly for async operations
- [x] 7.5 Verify computed fields are properly set and not editable
- [x] 7.6 Verify type-specific field validation works as expected
- [x] 7.7 Verify retry mechanism handles eventual consistency correctly
