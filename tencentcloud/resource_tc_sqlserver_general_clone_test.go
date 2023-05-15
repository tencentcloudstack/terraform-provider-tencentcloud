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

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralCloneResource_basic -v
func TestAccTencentCloudSqlserverGeneralCloneResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverGeneralCloneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralClone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverCloneExists("tencentcloud_sqlserver_general_clone.general_clone"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_clone.general_clone", "instance_id"),
				),
			},
			{
				Config: testAccSqlserverGeneralCloneUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverCloneExists("tencentcloud_sqlserver_general_clone.general_clone"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_clone.general_clone", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckSqlserverGeneralCloneDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_general_clone" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		dbName := idSplit[1]

		result, err := service.DescribeSqlserverGeneralCloneById(ctx, instanceId)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}

			return err
		}

		for _, v := range result {
			if *v.Name == dbName {
				return fmt.Errorf("sqlserver general_clone %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckSqlserverCloneExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		dbName := idSplit[1]

		result, err := service.DescribeSqlserverGeneralCloneById(ctx, instanceId)
		if err != nil {
			return err
		}

		if result != nil {
			for _, v := range result {
				if *v.Name == dbName {
					return nil
				}
			}
			return fmt.Errorf("sqlserver general_clone %s is not found", rs.Primary.ID)
		} else {
			return fmt.Errorf("sqlserver general_clone %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverGeneralClone = `
resource "tencentcloud_sqlserver_general_clone" "general_clone" {
  instance_id = "mssql-qelbzgwf"
  old_name    = "keep_pubsub_db"
  new_name    = "keep_pubsub_db_new_name"
}
`

const testAccSqlserverGeneralCloneUpdate = `
resource "tencentcloud_sqlserver_general_clone" "general_clone" {
  instance_id = "mssql-qelbzgwf"
  old_name    = "keep_pubsub_db_new_name"
  new_name    = "update_keep_pubsub_db_new_name"
}
`
