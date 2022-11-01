package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudDcdbAccount_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcdbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcdbAccountExists("tencentcloud_dcdb_account.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.basic", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "user_name", "mysql"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "host", "127.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "password", "===password==="),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "description", "this is a test account"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "max_user_connections", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "read_only", "0"),
				),
			},
			{
				Config: testAccDcdbAccount_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcdbAccountExists("tencentcloud_dcdb_account.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.basic", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "user_name", "mysql"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "host", "127.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "description", "this is a changed test account"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "read_only", "0"),
				),
			},
			{
				ResourceName: "tencentcloud_dcdb_account.basic",
				ImportState:  true,
			},
		},
	})
}

func testAccCheckDcdbAccountDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dcdbService := DcdbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dcdb_account" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		instanceId := idSplit[0]
		userName := idSplit[1]

		account, err := dcdbService.DescribeDcdbAccount(ctx, instanceId, userName)
		if err != nil {
			return err
		}
		if account.Users != nil && len(account.Users) > 0 {
			return fmt.Errorf("dcdb account still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDcdbAccountExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("dcdb account %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dcdb account id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		instanceId := idSplit[0]
		userName := idSplit[1]

		dcdbService := DcdbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		account, err := dcdbService.DescribeDcdbAccount(ctx, instanceId, userName)
		if err != nil {
			return err
		}
		if account == nil {
			return fmt.Errorf("dcdb account not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDcdbAccount_basic = `

resource "tencentcloud_dcdb_account" "basic" {
	instance_id = "tdsqlshard-lgz66iqr" # use the hard code before the dcdb_instance resource is ready.
	user_name = "mysql"
	host = "127.0.0.1"
	password = "===password==="
	read_only = 0
	description = "this is a test account"
	max_user_connections = 10
}

`
const testAccDcdbAccount_update = `

resource "tencentcloud_dcdb_account" "basic" {
  instance_id = "tdsqlshard-lgz66iqr" # use the hard code before the dcdb_instance resource is ready.
  user_name = "mysql"
  host = "127.0.0.1"
  password = "===password==="
  read_only = 0
  description = "this is a changed test account"
  max_user_connections = 10
}

`
