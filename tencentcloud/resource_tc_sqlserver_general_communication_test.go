package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralCommunicationResource_basic -v
func TestAccTencentCloudSqlserverGeneralCommunicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverGeneralCommunicationDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCommunication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverGeneralCommunicationExists("tencentcloud_sqlserver_general_communication.general_communication"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_communication.general_communication", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_communication.general_communication",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSqlserverGeneralCommunicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_general_communication" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		instanceId := rs.Primary.ID
		result, err := service.DescribeSqlserverGeneralCommunicationById(ctx, instanceId)
		if err != nil {
			return err
		}

		if result != nil {
			return fmt.Errorf("sqlserver general communicationinstance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverGeneralCommunicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		instanceId := rs.Primary.ID
		result, err := service.DescribeSqlserverGeneralCommunicationById(ctx, instanceId)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("sqlserver general communicationinstance %s is not found", rs.Primary.ID)
		} else {
			return nil
		}
	}
}

const testAccSqlserverGeneralCommunication = `
resource "tencentcloud_sqlserver_general_communication" "general_communication" {
  instance_id = "mssql-qelbzgwf"
}
`
