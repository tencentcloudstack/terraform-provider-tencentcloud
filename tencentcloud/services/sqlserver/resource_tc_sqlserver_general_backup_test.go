package sqlserver_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcsqlserver "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralBackupResource_basic -v -timeout=0
func TestAccTencentCloudSqlserverGeneralBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSqlserverGeneralBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralBackup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBackupExists("tencentcloud_sqlserver_general_backup.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_backup.example", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_backup.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverGeneralBackupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBackupExists("tencentcloud_sqlserver_general_backup.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_backup.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_backup.example", "backup_name"),
				),
			},
		},
	})
}

func testAccCheckSqlserverGeneralBackupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_general_backup" {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 6 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		flowId := idSplit[2]
		result, err := service.DescribeBackupByFlowId(ctx, instanceId, flowId)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "InvalidParameter.InputIllegal" || sdkerr.Code == "ResourceNotFound.BackupNotFound" || sdkerr.Code == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}

			return err
		}

		if result != nil {
			if *result.Response.Status == svcsqlserver.SQLSERVER_BACKUP_FAIL {
				return nil
			}

			return fmt.Errorf("sqlserver general_backup %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverBackupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 6 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		flowId := idSplit[2]
		result, err := service.DescribeBackupByFlowId(ctx, instanceId, flowId)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "InvalidParameter.InputIllegal" || sdkerr.Code == "ResourceNotFound.BackupNotFound" || sdkerr.Code == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}

			return err
		}

		if result != nil {
			if *result.Response.Status == svcsqlserver.SQLSERVER_BACKUP_FAIL {
				return fmt.Errorf("sqlserver general_backup %s is not found", rs.Primary.ID)
			}
			return nil
		} else {
			return fmt.Errorf("sqlserver general_backup %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverGeneralBackup = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = "tf_example_backup"
  strategy    = 0
}
`

const testAccSqlserverGeneralBackupUpdate = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = "tf_example_backup_update"
  strategy    = 0
}
`
