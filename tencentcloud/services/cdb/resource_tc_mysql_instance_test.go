package cdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"

	"context"
	"fmt"
	"testing"

	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

const TestAccTencentCloudMysqlMasterInstance_availability_zone = "ap-guangzhou-3"
const TestAccTencentCloudMysqlMasterInstance_availability_zone_4 = "ap-guangzhou-4"
const TestAccTencentCloudMysqlInstanceName = "testAccMysql"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_mysql_instance
	resource.AddTestSweepers("tencentcloud_mysql_instance", &resource.Sweeper{
		Name: "tencentcloud_mysql_instance",
		F:    testSweepMySQLInstance,
	})
}

func testSweepMySQLInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return err
	}
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := localcdb.NewMysqlService(client)

	request := cdb.NewDescribeDBInstancesRequest()
	request.Limit = helper.IntUint64(2000)

	response, err := client.UseMysqlClient().DescribeDBInstances(request)
	if err != nil {
		return err
	}

	if len(response.Response.Items) == 0 {
		return nil
	}

	for _, v := range response.Response.Items {
		id := *v.InstanceId
		name := *v.InstanceName
		if tcacctest.IsResourcePersist(name, nil) {
			continue
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err := service.IsolateDBInstance(ctx, id)
			if err != nil {
				//for the pay order wait
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})
		if err != nil {
			return err
		}

		_ = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			mysqlInfo, err := service.DescribeDBInstanceById(ctx, id)

			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if mysqlInfo == nil {
				return nil
			}
			if *mysqlInfo.Status == localcdb.MYSQL_STATUS_ISOLATING || *mysqlInfo.Status == localcdb.MYSQL_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("mysql isolating."))
			}
			if *mysqlInfo.Status == localcdb.MYSQL_STATUS_ISOLATED {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("after IsolateDBInstance mysql Status is %d", *mysqlInfo.Status))
		})

		err = service.OfflineIsolatedInstances(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceResource_prepaid -v
func TestAccTencentCloudMysqlInstanceResource_prepaid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMySQLPrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.prepaid"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.prepaid", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "charge_type", "PREPAID"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.prepaid", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.prepaid", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.prepaid", "first_slave_zone"),
				),
			},
			{
				ResourceName:            "tencentcloud_mysql_instance.prepaid",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "prepaid_period", "force_delete", "param_template_id", "fast_upgrade"},
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceResource_DeviceType -v
func TestAccTencentCloudMysqlInstanceResource_DeviceType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMySQLDeviceType,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_exclusive"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_exclusive", "device_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_exclusive", "first_slave_zone"),
				),
			},
			{
				ResourceName:            "tencentcloud_mysql_instance.mysql_exclusive",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "prepaid_period", "force_delete", "param_template_id", "fast_upgrade"},
			},
			{
				Config: testAccMySQLDeviceTypeUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_exclusive"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_exclusive", "device_type", "EXCLUSIVE"),
				),
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceResource_MasterInstance_fullslave -v
func TestAccTencentCloudMysqlInstanceResource_MasterInstance_fullslave(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlMasterInstance_fullslave(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", TestAccTencentCloudMysqlInstanceName),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "second_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
				),
			},
			{
				Config: testAccMysqlMasterInstanceUp_fullslave(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", TestAccTencentCloudMysqlInstanceName),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone_4),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "second_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone_4),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlInstanceResource_MasterInstance_basic_and_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", TestAccTencentCloudMysqlInstanceName),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "intranet_port", "3360"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "engine_version", "5.7"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),

					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "task_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "gtid"),
				),
			},
			{
				ResourceName:            "tencentcloud_mysql_instance.mysql_master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "prepaid_period", "first_slave_zone", "force_delete"},
			},
			// add tag
			{
				Config: testAccMysqlMasterInstance_multiTags("master"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "tags.role", "master"),
				),
			},
			// update tag
			{
				Config: testAccMysqlMasterInstance_multiTags("master-version2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "tags.role", "master-version2"),
				),
			},
			// remove tag
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckNoResourceAttr("tencentcloud_mysql_instance.mysql_master", "tags.role"),
				),
			},

			// open internet service
			{
				Config: testAccMysqlMasterInstance_internet_service(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "internet_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "internet_port"),
				),
			},

			//close internet  service
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "0")),
			},

			//modify  parameters
			{
				Config: testAccMysqlMasterInstance_parameters(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "parameters.max_connections", "1000")),
			},
			//remove parameters and  restore
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master")),
			},
			// update instance_name
			{
				Config: testAccMysqlMasterInstance_update("testAccMysql-version1", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", "testAccMysql-version1"),
				),
			},
			// update intranet_port
			{
				Config: testAccMysqlMasterInstance_update("testAccMysql-version1", "3361"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "intranet_port", "3361"),
				),
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceResource_mysql8 -v
func TestAccTencentCloudMysqlInstanceResource_mysql8(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlMasterInstance_mysql8("utf8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql8"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql8", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "charge_type", "POSTPAID"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "engine_version", "8.0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "internet_service", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "first_slave_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "second_slave_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "slave_sync_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "availability_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "instance_name", "myTestMysql"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "volume_size", "25"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "intranet_port", "3306"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "tags.createdBy", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "parameters.character_set_server", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "parameters.lower_case_table_names", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "parameters.max_connections", "1000"),
				),
			},
			{
				Config: testAccMysqlMasterInstance_mysql8("gbk"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql8"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql8", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "charge_type", "POSTPAID"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "engine_version", "8.0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "internet_service", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "first_slave_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "second_slave_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "slave_sync_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "availability_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "instance_name", "myTestMysql"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "volume_size", "25"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "intranet_port", "3306"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "tags.createdBy", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "parameters.character_set_server", "gbk"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "parameters.lower_case_table_names", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql8", "parameters.max_connections", "1000"),
				),
			},
		},
	})
}

func testAccCheckMysqlMasterInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_instance" {
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

func testAccCheckMysqlMasterInstanceExists(n string) resource.TestCheckFunc {
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

const testAccMySQLPrepaid = `
resource "tencentcloud_mysql_instance" "prepaid" {
  charge_type       = "PREPAID"
  prepaid_period    = 1
  mem_size          = 2000
  volume_size       = 30
  instance_name     = "testAccMysqlPrepaid"
  engine_version    = "5.7"
  intranet_port     = 3360
  root_password     = "test1234"
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete      = true
}
`

const testAccMySQLDeviceType = `
variable "temporary_param_tmpl_id" {
	default = 19656
}

resource "tencentcloud_mysql_instance" "mysql_exclusive" {
  charge_type       = "POSTPAID"
  mem_size          = 16000
  cpu               = 2
  volume_size       = 50
  instance_name     = "testAccMysqlBasic"
  engine_version    = "5.7"
  intranet_port     = 3360
  root_password     = "test1234"
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete      = true
  device_type       = "EXCLUSIVE"
  param_template_id = var.temporary_param_tmpl_id
}
`

const testAccMySQLDeviceTypeUpdate = `
variable "temporary_param_tmpl_id" {
	default = 19656
}

resource "tencentcloud_mysql_instance" "mysql_exclusive" {
  charge_type       = "POSTPAID"
  mem_size          = 16000
  cpu               = 2
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  intranet_port     = 3360
  root_password     = "test1234"
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete      = true
  device_type       = "EXCLUSIVE"
  fast_upgrade      = 1
  param_template_id = var.temporary_param_tmpl_id
}
`

func testAccMysqlMasterInstance_basic() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete      = true
}`
}

func testAccMysqlMasterInstance_fullslave() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  slave_deploy_mode = 0
  first_slave_zone  = "ap-guangzhou-3"
  second_slave_zone = "ap-guangzhou-3"
  slave_sync_mode   = 2
  force_delete      = true
}`
}

func testAccMysqlMasterInstanceUp_fullslave() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  slave_deploy_mode = 1
  first_slave_zone  = "ap-guangzhou-4"
  second_slave_zone = "ap-guangzhou-4"
  slave_sync_mode   = 1
  force_delete      = true
}`
}

func testAccMysqlMasterInstance_internet_service(open bool) string {
	tag := "0"
	if open {
		tag = "1"
	}
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  internet_service  = ` + tag + `
  force_delete      = true
}`

}

func testAccMysqlMasterInstance_parameters() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete      = true
  
  parameters = {
    max_connections = "1000"
  }
}`
}

func testAccMysqlMasterInstance_multiTags(value string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete      = true
  tags = {
    test = "test-tf"
    role = "%s"
  }
}
	`, value)
}

func testAccMysqlMasterInstance_update(instance_name, instranet_port string) string {
	tpl := `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "%s"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = %s
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete      = true
}`
	return fmt.Sprintf(tpl, instance_name, instranet_port)
}

func testAccMysqlMasterInstance_mysql8(value string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_instance" "mysql8" {
	internet_service = 1
	engine_version   = "8.0"
	charge_type = "POSTPAID"
	root_password     = "password123"
	slave_deploy_mode = 0
	first_slave_zone  = "ap-guangzhou-4"
	second_slave_zone = "ap-guangzhou-4"
	slave_sync_mode   = 1
	availability_zone = "ap-guangzhou-4"
	project_id        = 0
	instance_name     = "myTestMysql"
	mem_size          = 1000
	volume_size       = 25
	intranet_port     = 3306
	security_groups   = ["sg-ngx2bo7j"]
  
	tags = {
	  createdBy = "terraform"
	}
  
	parameters = {
	  character_set_server = "%s"
	  lower_case_table_names = 0
	  max_connections = "1000"
	}
}`, value)
}
