package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralBackupResource_basic -v
func TestAccTencentCloudSqlserverGeneralBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverGeneralBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralBackup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBackupExists("tencentcloud_sqlserver_general_backup.general_backup"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_backup.general_backup", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_backup.general_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverGeneralBackupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBackupExists("tencentcloud_sqlserver_general_backup.general_backup"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_backup.general_backup", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_backup.general_backup", "backup_name"),
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
			if *result.Response.Status == SQLSERVER_BACKUP_FAIL {
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

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
			if *result.Response.Status == SQLSERVER_BACKUP_FAIL {
				return fmt.Errorf("sqlserver general_backup %s is not found", rs.Primary.ID)
			}
			return nil
		} else {
			return fmt.Errorf("sqlserver general_backup %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverGeneralBackup = `
resource "tencentcloud_sqlserver_general_backup" "general_backup" {
  strategy = 0
  instance_id = "mssql-qelbzgwf"
  backup_name = "create_sqlserver_backup_name"
}
`

const testAccSqlserverGeneralBackupUpdate = `
resource "tencentcloud_sqlserver_general_backup" "general_backup" {
  strategy = 0
  instance_id = "mssql-qelbzgwf"
  backup_name = "update_sqlserver_backup_name"
}
`
