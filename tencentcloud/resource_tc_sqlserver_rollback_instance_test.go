package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRollbackInstanceResource_basic -v
func TestAccTencentCloudSqlserverRollbackInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverRollbackDBDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRollbackInstance(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_rollback_instance.rollback_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_rollback_instance.rollback_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSqlserverRollbackDBDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_rollback_instance" {
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

const testAccSqlserverRollbackInstanceString = `
resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-qelbzgwf"
  time = "%s"
  rename_restore {
    old_name = "keep_pubsub_db2"
	new_name = "rollback_pubsub_db2"
  }
}
`

func testAccSqlserverRollbackInstance() string {
	var currentTime = time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	return fmt.Sprintf(testAccSqlserverRollbackInstanceString, currentTime)
}
