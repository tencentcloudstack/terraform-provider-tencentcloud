package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dcdb_account", &resource.Sweeper{
		Name: "tencentcloud_dcdb_account",
		F:    testSweepDCDBAccount,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dcdb_account
func testSweepDCDBAccount(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	dcdbService := DcdbService{client: cli.(*TencentCloudClient).apiV3Conn}

	account, err := dcdbService.DescribeDcdbAccount(ctx, defaultDcdbInstanceId, "")
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("dcdb account not exists. instanceId:[%s]", defaultDcdbInstanceId)
	}

	for _, v := range account.Users {
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := dcdbService.DeleteDcdbAccountById(ctx, defaultDcdbInstanceId, *v.UserName, *v.Host)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("[ERROR] delete dcdb account %s reason:[%s]", *v.UserName, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudDCDBAccountResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcdbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDcdbAccount_basic, defaultDcdbInstanceId),
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
				Config: fmt.Sprintf(testAccDcdbAccount_update, defaultDcdbInstanceId),
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
	instance_id = "%s"
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
  instance_id = "%s"
  user_name = "mysql"
  host = "127.0.0.1"
  password = "===password==="
  read_only = 0
  description = "this is a changed test account"
  max_user_connections = 10
}

`
