## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_domain_shared_cname_attachment.go` with schema definition (`zone_id` ForceNew, `shared_cname` ForceNew, `domain_names` mutable), and CRUD functions (Create calls `BindSharedCNAMEWithContext` with BindType=bind, Read calls `DescribeSharedCNAMEWithContext` to populate domain_names, Update diffs domain_names and calls bind/unbind accordingly, Delete calls `BindSharedCNAMEWithContext` with BindType=unbind)
- [x] 1.2 Register the resource `tencentcloud_teo_domain_shared_cname_attachment` in `tencentcloud/provider.go`
- [x] 1.3 Add the resource entry to `tencentcloud/provider.md`

## 2. Documentation

- [x] 2.1 Create resource example file `tencentcloud/services/teo/resource_tc_teo_domain_shared_cname_attachment.md` with Example Usage and Import sections

## 3. Testing

- [x] 3.1 Create unit test file `tencentcloud/services/teo/resource_tc_teo_domain_shared_cname_attachment_test.go` using gomonkey to mock cloud API calls (WithContext variants), and verify with `go test -gcflags=all=-l`
