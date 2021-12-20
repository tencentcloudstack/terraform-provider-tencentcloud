package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testSqlserverAccountResourceName = "tencentcloud_sqlserver_account"
var testSqlserverAccountResourceKey = testSqlserverAccountResourceName + ".test"

func TestAccTencentCloudSqlserverAccountResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountExists(testSqlserverAccountResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "password", "testt123"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "update_time"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "is_admin", "false"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "status"),
				),
			},
			{
				ResourceName:            testSqlserverAccountResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "is_admin"},
			},

			{
				Config: testAccSqlserverAccountUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountExists(testSqlserverAccountResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "password", "test1233"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "remark", "testt"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "update_time"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "is_admin", "false"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "status"),
				),
			},
		},
	})
}

func testAccCheckSqlserverAccountDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testSqlserverAccountResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 2 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		name := idStrs[1]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverAccountById(ctx, instanceId, name)

		if err != nil {
			return err
		}

		if !has {
			return nil
		} else {
			return fmt.Errorf("delete SQL Server account %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 2 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		name := idStrs[1]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverAccountById(ctx, instanceId, name)
		if err != nil {
			_, has, err = service.DescribeSqlserverAccountById(ctx, instanceId, name)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("SQL Server account %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverAccount string = testAccSqlserverInstance + `
resource "tencentcloud_sqlserver_account" "test" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name = "tf_sqlserver_account"
  password = "testt123"
}
`

const testAccSqlserverAccountUpdate string = testAccSqlserverInstance + `
resource "tencentcloud_sqlserver_account" "test" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name = "tf_sqlserver_account"
  password = "test1233"
  remark = "testt"
}
`
