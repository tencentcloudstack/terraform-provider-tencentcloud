package cdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -test.run TestAccTencentNeedFixCloudMysqlDrInstanceResource_basic -v
func TestAccTencentNeedFixCloudMysqlDrInstanceResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-shanghai")
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlDrInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDrInstance,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlDrInstanceExists("tencentcloud_mysql_dr_instance.mysql_dr"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "instance_name", "mysql-dr-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "mem_size", "8000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "volume_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "intranet_port", "3360"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "tags.test", "test-tf"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_dr_instance.mysql_dr",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMysqlDrInstanceUp,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlDrInstanceExists("tencentcloud_mysql_dr_instance.mysql_dr"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "instance_name", "mysql-dr-test-up"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "mem_size", "8000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "volume_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "intranet_port", "3360"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_dr_instance.mysql_dr", "tags.test", "test-tf"),
				),
			},
		},
	})
}

func testAccCheckMysqlDrInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_dr_instance" {
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

func testAccCheckMysqlDrInstanceExists(n string) resource.TestCheckFunc {
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

const testAccMysqlDrInstance = `
resource "tencentcloud_mysql_dr_instance" "mysql_dr" {
  master_instance_id = "cdb-adjdu3t5"
  master_region      = "ap-guangzhou"
  auto_renew_flag   = 0
  availability_zone = "ap-shanghai-3"
  charge_type       = "POSTPAID"
  cpu               = 4
  device_type       = "UNIVERSAL"
  first_slave_zone  = "ap-shanghai-4"
  instance_name     = "mysql-dr-test"
  mem_size          = 8000
  prepaid_period    = 1
  project_id        = 0
  security_groups = [
    "sg-q4d821qk",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-5vfntba5"
  volume_size       = 100
  vpc_id            = "vpc-h6s1s3aa"
  intranet_port      = 3360
  tags = {
    test = "test-tf"
  }
}

`

const testAccMysqlDrInstanceUp = `
resource "tencentcloud_mysql_dr_instance" "mysql_dr" {
  master_instance_id = "cdb-adjdu3t5"
  master_region      = "ap-guangzhou"
  auto_renew_flag   = 0
  availability_zone = "ap-shanghai-3"
  charge_type       = "POSTPAID"
  cpu               = 4
  device_type       = "UNIVERSAL"
  first_slave_zone  = "ap-shanghai-4"
  instance_name     = "mysql-dr-test-up"
  mem_size          = 8000
  prepaid_period    = 1
  project_id        = 0
  security_groups = [
    "sg-q4d821qk",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-5vfntba5"
  volume_size       = 100
  vpc_id            = "vpc-h6s1s3aa"
  intranet_port      = 3360
  tags = {
    test = "test-tf"
  }
}

`
