package dts_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	dtssvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dts"
)

func TestAccTencentCloudDtsSyncConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_access_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_access_type"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "job_name", "tf_test_sync_config"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "job_mode", "liteMode"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "run_mode", "Immediate"),
					// objects
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.mode", "Partial"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.db_name", "tf_ci_test"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.new_db_name", "tf_ci_test_new"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.db_mode", "Partial"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.table_mode", "All"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.0.table_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.0.new_table_name", "test_new"),
					// src_info dest_info
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.instance_id", "cdb-fitq5t9h"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.user", "keep_dts"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.db_name", "tf_ci_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.0.db_name", "tf_ci_test_new"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.subnet_id"),

					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "auto_retry_time_range_minutes", "0"),
				),
			},
			{
				ResourceName:            "tencentcloud_dts_sync_config.sync_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dst_info.0.password", "src_info.0.password", "job_mode"},
			},
		},
	})
}

func TestAccTencentCloudDtsSyncConfigResource_ccn(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncConfig_ccn,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_access_type", "ccn"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_access_type", "cdb"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "job_name"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "job_mode", "liteMode"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "run_mode", "Immediate"),
					// objects
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.mode", "Partial"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.db_name", "tf_ci_test"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.new_db_name", "tf_ci_test_new"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.db_mode", "Partial"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.table_mode", "All"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.0.table_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.0.new_table_name", "test_new"),
					// src_info dest_info
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.region", "ap-shanghai"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.user"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.password"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.port"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.ccn_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.database_net_env", "TencentVPC"),

					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.user"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.password"),

					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "auto_retry_time_range_minutes", "0"),
				),
			},
			{
				ResourceName: "tencentcloud_dts_sync_config.sync_config",
				ImportState:  true,
			},
		},
	})
}

const testAccDtsSyncConfig = testAccDtsMigrateJob_vpc_config + `


resource "tencentcloud_dts_sync_job" "sync_job" {
	pay_mode = "PostPay"
	src_database_type = "mysql"
	src_region = "ap-guangzhou"
	dst_database_type = "cynosdbmysql"
	dst_region = "ap-guangzhou"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
	auto_renew = 0
	instance_class = "micro"
  }

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id = tencentcloud_dts_sync_job.sync_job.job_id
  src_access_type = "cdb"
  dst_access_type = "cdb"
  
  job_name = "tf_test_sync_config"
  job_mode = "liteMode"
  run_mode = "Immediate"

  objects {
	mode = "Partial"
      databases {
	    db_name = "tf_ci_test"
			new_db_name = "tf_ci_test_new"
			db_mode = "Partial"
			table_mode = "All"
			tables {
				table_name = "test"
				new_table_name = "test_new"
			}
	  }
  }
  src_info {
		region        = "ap-guangzhou"
		instance_id   = "cdb-fitq5t9h"
		user          = "keep_dts"
		password      = "Letmein123"
		db_name       = "tf_ci_test"
		vpc_id        = local.vpc_id
		subnet_id     = local.subnet_id
  }
  dst_info {
		region        = "ap-guangzhou"
		instance_id   = "cynosdbmysql-bws8h88b"
		user          = "keep_dts"
		password      = "Letmein123"
		db_name       = "tf_ci_test_new"
		vpc_id        = "vpc-pewdpc0d"
		subnet_id     = "subnet-driddx4g"
  }
  auto_retry_time_range_minutes = 0
}

`

const testAccDtsSyncConfig_ccn = `

locals {
  vpc_id_sh    = "vpc-evtcyb3g"
  subnet_id_sh = "subnet-1t83cxkp"
  src_ip       = data.tencentcloud_mysql_instance.src_mysql.instance_list.0.intranet_ip
  src_port     = data.tencentcloud_mysql_instance.src_mysql.instance_list.0.intranet_port
  ccn_id       = data.tencentcloud_ccn_instances.ccns.instance_list.0.ccn_id
  dst_mysql_id = data.tencentcloud_mysql_instance.dst_mysql.instance_list.0.mysql_id
}

variable "src_az_sh" {
  default = "ap-shanghai"
}

variable "dst_az_gz" {
  default = "ap-guangzhou"
}

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_name = "keep_sync_config_ccn_2_cdb"
}

data "tencentcloud_ccn_instances" "ccns" {
  name = "keep-ccn-dts-sh"
}

data "tencentcloud_mysql_instance" "src_mysql" {
  instance_name = "keep_dts_mysql_src"
}

data "tencentcloud_mysql_instance" "dst_mysql" {
  instance_name = "keep_dts_mysql_src"
}

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id          = data.tencentcloud_dts_sync_jobs.sync_jobs.list.0.job_id
  src_access_type = "ccn"
  dst_access_type = "cdb"

  job_mode = "liteMode"
  run_mode = "Immediate"

  objects {
    mode = "Partial"
    databases {
      db_name     = "tf_ci_test"
      new_db_name = "tf_ci_test_new"
      db_mode     = "Partial"
      table_mode  = "All"
      tables {
        table_name     = "test"
        new_table_name = "test_new"
      }
    }
  }
  src_info { // shanghai to guangzhou via ccn
    region           = var.src_az_sh
    user             = "keep_dts"
    password         = "Letmein123"
    ip               = local.src_ip
    port             = local.src_port
    vpc_id           = local.vpc_id_sh
    subnet_id        = local.subnet_id_sh
    ccn_id           = local.ccn_id
    database_net_env = "TencentVPC"
  }
  dst_info {
    region      = var.dst_az_gz
    instance_id = local.dst_mysql_id
    user        = "keep_dts"
    password    = "Letmein123"
  }
  auto_retry_time_range_minutes = 0
}


`

// ---- gomonkey-based unit tests for the schema_mode field ----
//
// These tests mock the DTS cloud API client methods so they can run without
// real cloud credentials. They verify that the `schema_mode` field is correctly
// mapped into the ConfigureSyncJob request (write path) and read back from the
// DescribeSyncJobs response (read path).

type mockMetaDtsSyncConfig struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDtsSyncConfig) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDtsSyncConfig{}

func newMockMetaDtsSyncConfig() *mockMetaDtsSyncConfig {
	return &mockMetaDtsSyncConfig{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDtsSyncConfig(s string) *string { return &s }

// TestDtsSyncConfig_Update_SchemaMode verifies that setting `schema_mode` in
// the configuration causes the value to be sent on the ConfigureSyncJob
// request's Objects.Databases[].SchemaMode field (write path).
func TestDtsSyncConfig_Update_SchemaMode(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncConfig().client, "UseDtsClient", dtsClient)

	// Capture and assert the ConfigureSyncJob request carries SchemaMode.
	patches.ApplyMethodFunc(dtsClient, "ConfigureSyncJob", func(request *dts.ConfigureSyncJobRequest) (*dts.ConfigureSyncJobResponse, error) {
		assert.NotNil(t, request.Objects)
		assert.NotNil(t, request.Objects.Databases)
		assert.Len(t, request.Objects.Databases, 1)
		assert.NotNil(t, request.Objects.Databases[0].SchemaMode)
		assert.Equal(t, "Partial", *request.Objects.Databases[0].SchemaMode)
		assert.Equal(t, "tf_ci_test", *request.Objects.Databases[0].DbName)

		resp := dts.NewConfigureSyncJobResponse()
		resp.Response = &dts.ConfigureSyncJobResponseParams{
			RequestId: ptrStringDtsSyncConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSyncJobs for both the state refresh polling and the final Read.
	var describeCallCount int
	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		describeCallCount++
		resp := dts.NewDescribeSyncJobsResponse()
		resp.Response = &dts.DescribeSyncJobsResponseParams{
			JobList: []*dts.SyncJobInfo{
				{
					JobId:         ptrStringDtsSyncConfig("sync-test-schema-mode"),
					JobName:       ptrStringDtsSyncConfig("tf_test_sync_config"),
					SrcAccessType: ptrStringDtsSyncConfig("cdb"),
					DstAccessType: ptrStringDtsSyncConfig("cdb"),
					RunMode:       ptrStringDtsSyncConfig("Immediate"),
					Status:        ptrStringDtsSyncConfig("Initialized"),
					TradeStatus:   ptrStringDtsSyncConfig("Normal"),
					Objects: &dts.Objects{
						Mode: ptrStringDtsSyncConfig("Partial"),
						Databases: []*dts.Database{
							{
								DbName:     ptrStringDtsSyncConfig("tf_ci_test"),
								NewDbName:  ptrStringDtsSyncConfig("tf_ci_test_new"),
								DbMode:     ptrStringDtsSyncConfig("Partial"),
								SchemaMode: ptrStringDtsSyncConfig("Partial"),
								TableMode:  ptrStringDtsSyncConfig("All"),
							},
						},
					},
					SrcInfo: &dts.Endpoint{
						Region:     ptrStringDtsSyncConfig("ap-guangzhou"),
						InstanceId: ptrStringDtsSyncConfig("cdb-fitq5t9h"),
						User:       ptrStringDtsSyncConfig("keep_dts"),
						Password:   ptrStringDtsSyncConfig("Letmein123"),
						DbName:     ptrStringDtsSyncConfig("tf_ci_test"),
						VpcId:      ptrStringDtsSyncConfig("vpc-i5yyodl9"),
						SubnetId:   ptrStringDtsSyncConfig("subnet-hhi88a58"),
					},
					DstInfo: &dts.Endpoint{
						Region:     ptrStringDtsSyncConfig("ap-guangzhou"),
						InstanceId: ptrStringDtsSyncConfig("cynosdbmysql-bws8h88b"),
						User:       ptrStringDtsSyncConfig("keep_dts"),
						Password:   ptrStringDtsSyncConfig("Letmein123"),
						DbName:     ptrStringDtsSyncConfig("tf_ci_test_new"),
						VpcId:      ptrStringDtsSyncConfig("vpc-i5yyodl9"),
						SubnetId:   ptrStringDtsSyncConfig("subnet-hhi88a58"),
					},
					AutoRetryTimeRangeMinutes: func() *int64 { v := int64(0); return &v }(),
				},
			},
			RequestId: ptrStringDtsSyncConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDtsSyncConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncConfig()
	rawConfig := map[string]interface{}{
		"job_id":          "sync-test-schema-mode",
		"src_access_type": "cdb",
		"dst_access_type": "cdb",
		"job_name":        "tf_test_sync_config",
		"job_mode":        "liteMode",
		"run_mode":        "Immediate",
		"objects": []interface{}{
			map[string]interface{}{
				"mode": "Partial",
				"databases": []interface{}{
					map[string]interface{}{
						"db_name":     "tf_ci_test",
						"new_db_name": "tf_ci_test_new",
						"db_mode":     "Partial",
						"schema_mode": "Partial",
						"table_mode":  "All",
					},
				},
			},
		},
		"src_info": []interface{}{
			map[string]interface{}{
				"region":      "ap-guangzhou",
				"instance_id": "cdb-fitq5t9h",
				"user":        "keep_dts",
				"password":    "Letmein123",
				"db_name":     "tf_ci_test",
				"vpc_id":      "vpc-i5yyodl9",
				"subnet_id":   "subnet-hhi88a58",
			},
		},
		"dst_info": []interface{}{
			map[string]interface{}{
				"region":      "ap-guangzhou",
				"instance_id": "cynosdbmysql-bws8h88b",
				"user":        "keep_dts",
				"password":    "Letmein123",
				"db_name":     "tf_ci_test_new",
				"vpc_id":      "vpc-i5yyodl9",
				"subnet_id":   "subnet-hhi88a58",
			},
		},
		"auto_retry_time_range_minutes": 0,
	}
	d := schema.TestResourceDataRaw(t, res.Schema, rawConfig)
	d.SetId("sync-test-schema-mode")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "Partial", d.Get("objects.0.databases.0.schema_mode").(string))
	assert.GreaterOrEqual(t, describeCallCount, 1)
}

// TestDtsSyncConfig_Read_SchemaMode verifies that the `schema_mode` value
// returned by DescribeSyncJobs is read back into Terraform state (read path).
func TestDtsSyncConfig_Read_SchemaMode(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncConfig().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		resp := dts.NewDescribeSyncJobsResponse()
		resp.Response = &dts.DescribeSyncJobsResponseParams{
			JobList: []*dts.SyncJobInfo{
				{
					JobId:         ptrStringDtsSyncConfig("sync-test-schema-mode"),
					JobName:       ptrStringDtsSyncConfig("tf_test_sync_config"),
					SrcAccessType: ptrStringDtsSyncConfig("cdb"),
					DstAccessType: ptrStringDtsSyncConfig("cdb"),
					RunMode:       ptrStringDtsSyncConfig("Immediate"),
					Objects: &dts.Objects{
						Mode: ptrStringDtsSyncConfig("Partial"),
						Databases: []*dts.Database{
							{
								DbName:     ptrStringDtsSyncConfig("tf_ci_test"),
								NewDbName:  ptrStringDtsSyncConfig("tf_ci_test_new"),
								DbMode:     ptrStringDtsSyncConfig("Partial"),
								SchemaMode: ptrStringDtsSyncConfig("All"),
								TableMode:  ptrStringDtsSyncConfig("All"),
							},
						},
					},
					SrcInfo: &dts.Endpoint{
						Region:     ptrStringDtsSyncConfig("ap-guangzhou"),
						InstanceId: ptrStringDtsSyncConfig("cdb-fitq5t9h"),
						User:       ptrStringDtsSyncConfig("keep_dts"),
						DbName:     ptrStringDtsSyncConfig("tf_ci_test"),
					},
					DstInfo: &dts.Endpoint{
						Region:     ptrStringDtsSyncConfig("ap-guangzhou"),
						InstanceId: ptrStringDtsSyncConfig("cynosdbmysql-bws8h88b"),
						User:       ptrStringDtsSyncConfig("keep_dts"),
						DbName:     ptrStringDtsSyncConfig("tf_ci_test_new"),
					},
				},
			},
			RequestId: ptrStringDtsSyncConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDtsSyncConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sync-test-schema-mode")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "All", d.Get("objects.0.databases.0.schema_mode").(string))
	assert.Equal(t, "Partial", d.Get("objects.0.databases.0.db_mode").(string))
}

// keep context import for compile safety / future use.
var _ = context.TODO
