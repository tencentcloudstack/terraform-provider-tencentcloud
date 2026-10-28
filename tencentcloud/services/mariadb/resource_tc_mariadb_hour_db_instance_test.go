package mariadb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcmariadb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mariadb"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbHourDbInstance_basic -v
func TestAccTencentCloudMariadbHourDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMariadbHourDbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbHourDbInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMariadbHourDbInstanceExists("tencentcloud_mariadb_hour_db_instance.basic"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "instance_name", "db-test-2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "db_version_id", "8.0"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "memory", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "node_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "storage", "10"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "subnet_id", tcacctest.DefaultMariadbSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "vpc_id", tcacctest.DefaultMariadbVpcId),
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_hour_db_instance.basic", "zones.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_hour_db_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMariadbHourDbInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmariadb.NewMariadbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mariadb_hour_db_instance" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		instance, err := service.DescribeMariadbDbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("db hour Instance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMariadbHourDbInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		service := svcmariadb.NewMariadbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMariadbDbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("db hour Instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMariadbHourDbInstanceVar = `
variable "subnet_id" {
  default = "` + tcacctest.DefaultMariadbSubnetId + `"
}

variable "vpc_id" {
  default = "` + tcacctest.DefaultMariadbVpcId + `"
}
`

const testAccMariadbHourDbInstance = testAccMariadbHourDbInstanceVar + `
resource "tencentcloud_mariadb_hour_db_instance" "basic" {
  db_version_id = "8.0"
  instance_name = "db-test-2"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = var.subnet_id
  vpc_id        = var.vpc_id
  zones         = ["ap-guangzhou-6","ap-guangzhou-7"]
  tags          = {
	createdBy   = "terraform"
  }
}
`

// mockMetaMariadbHourDbInstance implements tccommon.ProviderMeta
type mockMetaMariadbHourDbInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaMariadbHourDbInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaMariadbHourDbInstance{}

func ptrMariadbHourDbInstanceString(s string) *string {
	return &s
}

func ptrMariadbHourDbInstanceInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/mariadb/ -run "TestMariadbHourDbInstanceInitParams" -v -count=1 -gcflags="all=-l"

// TestMariadbHourDbInstanceInitParams_Create tests that init_params is correctly passed to CreateHourDBInstance request
func TestMariadbHourDbInstanceInitParams_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mariadbClient := &mariadb.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMariadbClient", mariadbClient)

	var capturedRequest *mariadb.CreateHourDBInstanceRequest

	// Patch CreateHourDBInstance to capture the request and return a mock response
	patches.ApplyMethodFunc(mariadbClient, "CreateHourDBInstance", func(request *mariadb.CreateHourDBInstanceRequest) (*mariadb.CreateHourDBInstanceResponse, error) {
		capturedRequest = request
		resp := &mariadb.CreateHourDBInstanceResponse{}
		resp.Response = &mariadb.CreateHourDBInstanceResponseParams{
			DealName:    ptrMariadbHourDbInstanceString("test-deal"),
			InstanceIds: []*string{ptrMariadbHourDbInstanceString("tdsql-test-12345")},
			FlowId:      ptrMariadbHourDbInstanceInt64(12345),
			RequestId:   ptrMariadbHourDbInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeDBInstanceDetail for the Read call
	patches.ApplyMethodFunc(mariadbClient, "DescribeDBInstanceDetail", func(request *mariadb.DescribeDBInstanceDetailRequest) (*mariadb.DescribeDBInstanceDetailResponse, error) {
		resp := &mariadb.DescribeDBInstanceDetailResponse{}
		resp.Response = &mariadb.DescribeDBInstanceDetailResponseParams{
			InstanceId:   ptrMariadbHourDbInstanceString("tdsql-test-12345"),
			InstanceName: ptrMariadbHourDbInstanceString("test-instance"),
			Zone:         ptrMariadbHourDbInstanceString("ap-guangzhou-6"),
			SlaveZones:   []*string{ptrMariadbHourDbInstanceString("ap-guangzhou-7")},
			NodeCount:    ptrMariadbHourDbInstanceInt64(2),
			Memory:       ptrMariadbHourDbInstanceInt64(2),
			Storage:      ptrMariadbHourDbInstanceInt64(10),
			ProjectId:    ptrMariadbHourDbInstanceInt64(0),
			VpcId:        ptrMariadbHourDbInstanceString("vpc-test"),
			SubnetId:     ptrMariadbHourDbInstanceString("subnet-test"),
			Vip:          ptrMariadbHourDbInstanceString("10.0.0.1"),
			DbVersionId:  ptrMariadbHourDbInstanceString("8.0"),
			RequestId:    ptrMariadbHourDbInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := &mockMetaMariadbHourDbInstance{client: mockClient}
	res := svcmariadb.ResourceTencentCloudMariadbHourDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zones":         []interface{}{"ap-guangzhou-6", "ap-guangzhou-7"},
		"node_count":    2,
		"memory":        2,
		"storage":       10,
		"db_version_id": "8.0",
		"instance_name": "test-instance",
		"vpc_id":        "vpc-test",
		"subnet_id":     "subnet-test",
		"init_params": []interface{}{
			map[string]interface{}{
				"param": "character_set_server",
				"value": "utf8",
			},
			map[string]interface{}{
				"param": "lower_case_table_names",
				"value": "0",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tdsql-test-12345", d.Id())

	// Verify init_params were passed to the CreateHourDBInstance request
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, 2, len(capturedRequest.InitParams))
	assert.Equal(t, "character_set_server", *capturedRequest.InitParams[0].Param)
	assert.Equal(t, "utf8", *capturedRequest.InitParams[0].Value)
	assert.Equal(t, "lower_case_table_names", *capturedRequest.InitParams[1].Param)
	assert.Equal(t, "0", *capturedRequest.InitParams[1].Value)
}

// TestMariadbHourDbInstanceInitParams_CreateWithoutInitParams tests that when init_params is not provided,
// the hardcoded InitDBInstances call is used
func TestMariadbHourDbInstanceInitParams_CreateWithoutInitParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mariadbClient := &mariadb.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMariadbClient", mariadbClient)

	var capturedRequest *mariadb.CreateHourDBInstanceRequest
	var initDBCalled bool

	// Patch CreateHourDBInstance
	patches.ApplyMethodFunc(mariadbClient, "CreateHourDBInstance", func(request *mariadb.CreateHourDBInstanceRequest) (*mariadb.CreateHourDBInstanceResponse, error) {
		capturedRequest = request
		resp := &mariadb.CreateHourDBInstanceResponse{}
		resp.Response = &mariadb.CreateHourDBInstanceResponseParams{
			DealName:    ptrMariadbHourDbInstanceString("test-deal"),
			InstanceIds: []*string{ptrMariadbHourDbInstanceString("tdsql-test-67890")},
			FlowId:      ptrMariadbHourDbInstanceInt64(12345),
			RequestId:   ptrMariadbHourDbInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch InitDbInstance on MariadbService to simulate successful initialization
	patches.ApplyMethodFunc(&svcmariadb.MariadbService{}, "InitDbInstance", func(_ context.Context, _ string, _ []*mariadb.DBParamValue) (bool, error) {
		initDBCalled = true
		return true, nil
	})

	// Patch DescribeDBInstanceDetail for the Read call
	patches.ApplyMethodFunc(mariadbClient, "DescribeDBInstanceDetail", func(request *mariadb.DescribeDBInstanceDetailRequest) (*mariadb.DescribeDBInstanceDetailResponse, error) {
		resp := &mariadb.DescribeDBInstanceDetailResponse{}
		resp.Response = &mariadb.DescribeDBInstanceDetailResponseParams{
			InstanceId:   ptrMariadbHourDbInstanceString("tdsql-test-67890"),
			InstanceName: ptrMariadbHourDbInstanceString("test-instance"),
			Zone:         ptrMariadbHourDbInstanceString("ap-guangzhou-6"),
			SlaveZones:   []*string{ptrMariadbHourDbInstanceString("ap-guangzhou-7")},
			NodeCount:    ptrMariadbHourDbInstanceInt64(2),
			Memory:       ptrMariadbHourDbInstanceInt64(2),
			Storage:      ptrMariadbHourDbInstanceInt64(10),
			ProjectId:    ptrMariadbHourDbInstanceInt64(0),
			VpcId:        ptrMariadbHourDbInstanceString("vpc-test"),
			SubnetId:     ptrMariadbHourDbInstanceString("subnet-test"),
			Vip:          ptrMariadbHourDbInstanceString("10.0.0.1"),
			DbVersionId:  ptrMariadbHourDbInstanceString("8.0"),
			RequestId:    ptrMariadbHourDbInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := &mockMetaMariadbHourDbInstance{client: mockClient}
	res := svcmariadb.ResourceTencentCloudMariadbHourDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zones":         []interface{}{"ap-guangzhou-6", "ap-guangzhou-7"},
		"node_count":    2,
		"memory":        2,
		"storage":       10,
		"db_version_id": "8.0",
		"instance_name": "test-instance",
		"vpc_id":        "vpc-test",
		"subnet_id":     "subnet-test",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tdsql-test-67890", d.Id())

	// Verify InitParams was NOT set on the CreateHourDBInstance request
	assert.NotNil(t, capturedRequest)
	assert.Nil(t, capturedRequest.InitParams)

	// Verify InitDBInstances was called (hardcoded defaults path)
	assert.True(t, initDBCalled)
}

// TestMariadbHourDbInstanceInitParams_Schema tests the schema definition of init_params
func TestMariadbHourDbInstanceInitParams_Schema(t *testing.T) {
	res := svcmariadb.ResourceTencentCloudMariadbHourDbInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "init_params")

	initParams := res.Schema["init_params"]
	assert.Equal(t, schema.TypeList, initParams.Type)
	assert.True(t, initParams.Optional)
	assert.True(t, initParams.ForceNew)

	elemResource := initParams.Elem.(*schema.Resource)
	assert.Contains(t, elemResource.Schema, "param")
	assert.Contains(t, elemResource.Schema, "value")
	assert.Equal(t, schema.TypeString, elemResource.Schema["param"].Type)
	assert.True(t, elemResource.Schema["param"].Required)
	assert.Equal(t, schema.TypeString, elemResource.Schema["value"].Type)
	assert.True(t, elemResource.Schema["value"].Required)
}
