## 1. Resource Fix

- [x] 1.1 Add `"strings"` to imports in `resource_tc_mongodb_instance_srv_connection.go`
- [x] 1.2 In `resourceTencentCloudMongodbInstanceSrvConnectionRead`, truncate `*domain` at the first `.` before calling `d.Set("domain", ...)`
