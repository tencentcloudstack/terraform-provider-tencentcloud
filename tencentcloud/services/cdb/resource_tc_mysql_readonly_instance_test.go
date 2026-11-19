package cdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"

	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cdb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tag_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlReadonlyInstanceResource_basic -v
func TestAccTencentCloudMysqlReadonlyInstanceResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "2000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "ro_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.test", "test-tf"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_readonly_instance.mysql_readonly",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(tcacctest.CommonPresetMysql, "read"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "read"),
				),
			},
			// update tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(tcacctest.CommonPresetMysql, "readonly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "readonly"),
				),
			},
			// remove tag
			{
				Config: testAccMysqlReadonlyInstance(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckNoResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role"),
				),
			},
			// update instance_name
			{
				Config: testAccMysqlReadonlyInstance_update(tcacctest.CommonPresetMysql, "mysql-readonly-update", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-update"),
				),
			},
			// update mem_size
			{
				Config: testAccMysqlReadonlyInstance_memSize(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "1000"),
				),
			},
			{
				Config: testAccMysqlReadonlyInstance_roGroup(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly_ro_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly_ro_group", "ro_group_id"),
				),
			},
			// // update intranet_port
			// {
			// 	Config: testAccMysqlReadonlyInstance_update(CommonPresetMysql, "mysql-readonly-update", "3361"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
			// 		resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3361"),
			// 	),
			// },
		},
	})
}

func testAccCheckMysqlReadonlyInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_readonly_instance" {
			continue
		}
		instance, err := mysqlService.DescribeRunningDBInstanceById(ctx, rs.Primary.ID)
		if instance != nil {
			return fmt.Errorf("mysql instance still exist")
		}
		if err != nil {
			sdkErr, ok := err.(*errors.TencentCloudSDKError)
			if ok && sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
				continue
			}
			return err
		}
	}
	return nil
}

func testAccCheckMysqlInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mysql instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mysql instance id is not set")
		}

		mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := mysqlService.DescribeDBInstanceById(ctx, rs.Primary.ID)
		if instance == nil {
			return fmt.Errorf("mysql instance %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccMysqlReadonlyInstance(mysqlTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase)
}

func testAccMysqlReadonlyInstance_multiTags(mysqlTestCase, value string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  cpu                = 1
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
    role = "%s"
  }
}
	`, mysqlTestCase, value)
}

func testAccMysqlReadonlyInstance_update(mysqlTestCase, instance_name, instranet_port string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  cpu                = 1
  volume_size        = 200
  instance_name      = "%s"
  intranet_port      = %s 
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase, instance_name, instranet_port)
}

func testAccMysqlReadonlyInstance_memSize(mysqlTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 1000
  cpu                = 1
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase)
}

func testAccMysqlReadonlyInstance_roGroup(mysqlTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly_ro_group" {
  master_instance_id = local.mysql_id
  ro_group_id 		 = tencentcloud_mysql_readonly_instance.mysql_readonly.ro_group_id
  mem_size           = 1000
  cpu                = 1
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase)
}

// ---- gomonkey-based unit tests for the client-token create retry mechanism ----

type mockMetaForMysqlReadonlyInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForMysqlReadonlyInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForMysqlReadonlyInstance{}

func newMockMetaForMysqlReadonlyInstance() *mockMetaForMysqlReadonlyInstance {
	return &mockMetaForMysqlReadonlyInstance{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStrMysqlReadonly(s string) *string    { return &s }
func ptrInt64MysqlReadonly(v int64) *int64    { return &v }
func ptrUint64MysqlReadonly(v uint64) *uint64 { return &v }

// patchReadonlyInstanceReadApis mocks all SDK client methods invoked during the
// mysql readonly instance Read flow (which Create calls at the end) so that the
// Create path completes without touching the real cloud API.
func patchReadonlyInstanceReadApis(patches *gomonkey.Patches, cdbClient *cdb_sdk.Client, tagClient *tag_sdk.Client, roInstanceId, masterInstanceId string) {
	// DescribeDBInstances: returns one running readonly instance bound to the master
	patches.ApplyMethodFunc(cdbClient, "DescribeDBInstances", func(request *cdb_sdk.DescribeDBInstancesRequest) (*cdb_sdk.DescribeDBInstancesResponse, error) {
		status := int64(localcdb.MYSQL_STATUS_RUNNING)
		payType := int64(localcdb.MysqlPayByUse)
		autoRenew := int64(localcdb.MYSQL_RENEW_CLOSE)
		wanStatus := int64(0)
		taskStatus := int64(0)
		cdbError := int64(0)
		port := int64(3306)
		cpu := int64(1)
		mem := int64(1000)
		vol := int64(50)
		project := int64(0)
		instanceName := "mysql-readonly-test"
		engineVersion := "5.7"
		vip := "10.0.0.1"
		zone := "ap-guangzhou-3"
		deviceType := "UNIVERSAL"
		resp := cdb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &cdb_sdk.DescribeDBInstancesResponseParams{
			TotalCount: ptrInt64MysqlReadonly(1),
			Items: []*cdb_sdk.InstanceInfo{
				{
					InstanceId:    ptrStrMysqlReadonly(roInstanceId),
					InstanceName:  ptrStrMysqlReadonly(instanceName),
					Status:        &status,
					PayType:       &payType,
					AutoRenew:     &autoRenew,
					Memory:        &mem,
					Cpu:           &cpu,
					Volume:        &vol,
					ProjectId:     &project,
					WanStatus:     &wanStatus,
					TaskStatus:    &taskStatus,
					CdbError:      &cdbError,
					Vip:           ptrStrMysqlReadonly(vip),
					Vport:         &port,
					EngineVersion: ptrStrMysqlReadonly(engineVersion),
					Zone:          ptrStrMysqlReadonly(zone),
					DeviceType:    ptrStrMysqlReadonly(deviceType),
					MasterInfo: &cdb_sdk.MasterInfo{
						Region:     ptrStrMysqlReadonly("ap-guangzhou"),
						Zone:       ptrStrMysqlReadonly("ap-guangzhou-3"),
						InstanceId: ptrStrMysqlReadonly(masterInstanceId),
					},
				},
			},
			RequestId: ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})

	// DescribeDBSecurityGroups: empty list
	patches.ApplyMethodFunc(cdbClient, "DescribeDBSecurityGroups", func(request *cdb_sdk.DescribeDBSecurityGroupsRequest) (*cdb_sdk.DescribeDBSecurityGroupsResponse, error) {
		resp := cdb_sdk.NewDescribeDBSecurityGroupsResponse()
		resp.Response = &cdb_sdk.DescribeDBSecurityGroupsResponseParams{
			Groups:    []*cdb_sdk.SecurityGroup{},
			RequestId: ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})

	// DescribeRoGroups: returns a ro group containing the readonly instance id
	patches.ApplyMethodFunc(cdbClient, "DescribeRoGroups", func(request *cdb_sdk.DescribeRoGroupsRequest) (*cdb_sdk.DescribeRoGroupsResponse, error) {
		resp := cdb_sdk.NewDescribeRoGroupsResponse()
		resp.Response = &cdb_sdk.DescribeRoGroupsResponseParams{
			RoGroups: []*cdb_sdk.RoGroup{
				{
					RoGroupId:   ptrStrMysqlReadonly("ro-group-fake-id"),
					RoGroupMode: ptrStrMysqlReadonly("allinone"),
					RoInstances: []*cdb_sdk.RoInstanceInfo{
						{
							InstanceId: ptrStrMysqlReadonly(roInstanceId),
						},
					},
				},
			},
			RequestId: ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})

	// tag client DescribeResourceTagsByResourceIds: empty tags
	patches.ApplyMethodFunc(tagClient, "DescribeResourceTagsByResourceIds", func(request *tag_sdk.DescribeResourceTagsByResourceIdsRequest) (*tag_sdk.DescribeResourceTagsByResourceIdsResponse, error) {
		resp := tag_sdk.NewDescribeResourceTagsByResourceIdsResponse()
		resp.Response = &tag_sdk.DescribeResourceTagsByResourceIdsResponseParams{
			TotalCount: ptrUint64MysqlReadonly(0),
			Tags:       []*tag_sdk.TagResource{},
			RequestId:  ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})
}

// patchReadonlyInstanceBackupApis mocks DescribeBackups so the create-flow
// backup-waiting retry passes immediately.
func patchReadonlyInstanceBackupApis(patches *gomonkey.Patches, cdbClient *cdb_sdk.Client) {
	patches.ApplyMethodFunc(cdbClient, "DescribeBackups", func(request *cdb_sdk.DescribeBackupsRequest) (*cdb_sdk.DescribeBackupsResponse, error) {
		resp := cdb_sdk.NewDescribeBackupsResponse()
		resp.Response = &cdb_sdk.DescribeBackupsResponseParams{
			TotalCount: ptrInt64MysqlReadonly(1),
			Items: []*cdb_sdk.BackupInfo{
				{
					BackupId: ptrInt64MysqlReadonly(1),
				},
			},
			RequestId: ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})
}

// TestMysqlReadonlyInstance_Create_Prepaid_SetsClientToken verifies that the
// prepaid (包年包月) create path sets a non-empty ClientToken generated by
// helper.BuildToken() on the CreateDBInstance request.
func TestMysqlReadonlyInstance_Create_Prepaid_SetsClientToken(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	tagClient := &tag_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseMysqlClient", cdbClient)
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseMysqlClientRegion", cdbClient)
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseTagClient", tagClient)

	instanceId := "cdb-ro-prepaid-001"
	masterInstanceId := "cdb-master-prepaid-001"

	var capturedClientToken *string
	patches.ApplyMethodFunc(cdbClient, "CreateDBInstance", func(request *cdb_sdk.CreateDBInstanceRequest) (*cdb_sdk.CreateDBInstanceResponse, error) {
		if request.ClientToken != nil {
			val := *request.ClientToken
			capturedClientToken = &val
		}
		resp := cdb_sdk.NewCreateDBInstanceResponse()
		resp.Response = &cdb_sdk.CreateDBInstanceResponseParams{
			InstanceIds: []*string{ptrStrMysqlReadonly(instanceId)},
			RequestId:   ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})

	patchReadonlyInstanceBackupApis(patches, cdbClient)
	patchReadonlyInstanceReadApis(patches, cdbClient, tagClient, instanceId, masterInstanceId)

	meta := newMockMetaForMysqlReadonlyInstance()
	res := localcdb.ResourceTencentCloudMysqlReadonlyInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"master_instance_id": masterInstanceId,
		"instance_name":      "mysql-readonly-test",
		"mem_size":           1000,
		"volume_size":        50,
		"charge_type":        "PREPAID",
		"prepaid_period":     1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())

	assert.NotNil(t, capturedClientToken, "ClientToken should be set on CreateDBInstance request")
	assert.NotEmpty(t, *capturedClientToken, "ClientToken should be a non-empty string")
}

// TestMysqlReadonlyInstance_Create_Postpaid_SetsClientToken verifies that the
// postpaid (按量计费) create path sets a non-empty ClientToken generated by
// helper.BuildToken() on the CreateDBInstanceHour request.
func TestMysqlReadonlyInstance_Create_Postpaid_SetsClientToken(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	tagClient := &tag_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseMysqlClient", cdbClient)
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseMysqlClientRegion", cdbClient)
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseTagClient", tagClient)

	instanceId := "cdb-ro-postpaid-002"
	masterInstanceId := "cdb-master-postpaid-002"

	var capturedClientToken *string
	patches.ApplyMethodFunc(cdbClient, "CreateDBInstanceHour", func(request *cdb_sdk.CreateDBInstanceHourRequest) (*cdb_sdk.CreateDBInstanceHourResponse, error) {
		if request.ClientToken != nil {
			val := *request.ClientToken
			capturedClientToken = &val
		}
		resp := cdb_sdk.NewCreateDBInstanceHourResponse()
		resp.Response = &cdb_sdk.CreateDBInstanceHourResponseParams{
			InstanceIds: []*string{ptrStrMysqlReadonly(instanceId)},
			RequestId:   ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})

	patchReadonlyInstanceBackupApis(patches, cdbClient)
	patchReadonlyInstanceReadApis(patches, cdbClient, tagClient, instanceId, masterInstanceId)

	meta := newMockMetaForMysqlReadonlyInstance()
	res := localcdb.ResourceTencentCloudMysqlReadonlyInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"master_instance_id": masterInstanceId,
		"instance_name":      "mysql-readonly-test",
		"mem_size":           1000,
		"volume_size":        50,
		"charge_type":        "POSTPAID",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())

	assert.NotNil(t, capturedClientToken, "ClientToken should be set on CreateDBInstanceHour request")
	assert.NotEmpty(t, *capturedClientToken, "ClientToken should be a non-empty string")
}

// TestMysqlReadonlyInstance_Create_RetriesOnEmptyInstanceIds verifies that when
// the create API responds with an empty InstanceIds list (while a client token
// was provided), the create path retries and ultimately succeeds once the API
// returns a valid instance id.
func TestMysqlReadonlyInstance_Create_RetriesOnEmptyInstanceIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	tagClient := &tag_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseMysqlClient", cdbClient)
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseMysqlClientRegion", cdbClient)
	patches.ApplyMethodReturn(newMockMetaForMysqlReadonlyInstance().client, "UseTagClient", tagClient)

	instanceId := "cdb-ro-retry-003"
	masterInstanceId := "cdb-master-retry-003"

	var capturedClientToken *string
	callCount := 0
	patches.ApplyMethodFunc(cdbClient, "CreateDBInstanceHour", func(request *cdb_sdk.CreateDBInstanceHourRequest) (*cdb_sdk.CreateDBInstanceHourResponse, error) {
		if request.ClientToken != nil {
			val := *request.ClientToken
			capturedClientToken = &val
		}
		callCount++
		resp := cdb_sdk.NewCreateDBInstanceHourResponse()
		// first call returns an empty InstanceIds list -> should trigger retry
		if callCount == 1 {
			resp.Response = &cdb_sdk.CreateDBInstanceHourResponseParams{
				InstanceIds: []*string{},
				RequestId:   ptrStrMysqlReadonly("fake-request-id"),
			}
			return resp, nil
		}
		// subsequent call returns the instance id -> success
		resp.Response = &cdb_sdk.CreateDBInstanceHourResponseParams{
			InstanceIds: []*string{ptrStrMysqlReadonly(instanceId)},
			RequestId:   ptrStrMysqlReadonly("fake-request-id"),
		}
		return resp, nil
	})

	patchReadonlyInstanceBackupApis(patches, cdbClient)
	patchReadonlyInstanceReadApis(patches, cdbClient, tagClient, instanceId, masterInstanceId)

	meta := newMockMetaForMysqlReadonlyInstance()
	res := localcdb.ResourceTencentCloudMysqlReadonlyInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"master_instance_id": masterInstanceId,
		"instance_name":      "mysql-readonly-test",
		"mem_size":           1000,
		"volume_size":        50,
		"charge_type":        "POSTPAID",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())

	// the create API must have been called more than once (retry happened)
	assert.Greater(t, callCount, 1, "create API should be retried after an empty InstanceIds response")
	assert.NotNil(t, capturedClientToken, "ClientToken should be set on the create request")
	assert.NotEmpty(t, *capturedClientToken, "ClientToken should be a non-empty string")
}
