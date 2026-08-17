## 1. Schema Definition

- [x] 1.1 Add `destroy_protect` schema field (TypeString, Optional, Computed) to `TencentMsyqlBasicInfo()` in `tencentcloud/services/cdb/resource_tc_mysql_instance.go` with Description explaining valid values `on`/`off`

## 2. Create Operation

- [x] 2.1 Add `DestroyProtect` parameter setting in `mysqlAllInstanceRoleSet()` function in `tencentcloud/services/cdb/resource_tc_mysql_instance.go` — read `destroy_protect` from schema with `d.GetOk()`, set on both `requestByMonth` (*cdb.CreateDBInstanceRequest) and `requestByUse` (*cdb.CreateDBInstanceHourRequest) using the existing dual-path `okByMonth` pattern

## 3. Read Operation

- [x] 3.1 Add `DestroyProtect` read in `tencentMsyqlBasicInfoRead()` function in `tencentcloud/services/cdb/resource_tc_mysql_instance.go` — check `mysqlInfo.DestroyProtect != nil` before setting state with `_ = d.Set("destroy_protect", mysqlInfo.DestroyProtect)`

## 4. Tests

- [x] 4.1 Add test case for `destroy_protect` parameter in `tencentcloud/services/cdb/resource_tc_mysql_instance_test.go` using mock (gomonkey) approach to mock the CDB API and verify business logic

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/cdb/resource_tc_mysql_instance.md` with `destroy_protect` usage example
