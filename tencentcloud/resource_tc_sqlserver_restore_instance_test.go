package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRestoreInstanceResource_basic -v
func TestAccTencentCloudSqlserverRestoreInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverRestoreDBDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRestoreInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_restore_instance.restore_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_restore_instance.restore_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSqlserverRestoreDBDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_restore_instance" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		newNameListStr := idSplit[3]
		newNameList := strings.Split(newNameListStr, COMMA_SP)

		for _, name := range newNameList {
			result, err := sqlserverService.DescribeSqlserverDBS(ctx, instanceId, name)
			if err != nil {
				return err
			}

			if result != nil {
				return fmt.Errorf("SQL Server DB still exists")
			}
		}
	}

	return nil
}

const testAccSqlserverRestoreInstance = `
resource "tencentcloud_sqlserver_restore_instance" "restore_instance" {
  instance_id = "mssql-qelbzgwf"
  backup_id   = 3482091273
  rename_restore {
    old_name = "keep_pubsub_db2"
    new_name = "restore_keep_pubsub_db2"
  }
  rename_restore {
    old_name = "keep_pubsub_db"
    new_name = "restore_keep_pubsub_db"
  }
}
`
