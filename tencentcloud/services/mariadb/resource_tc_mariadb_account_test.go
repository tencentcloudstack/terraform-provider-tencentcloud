package mariadb_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcmariadb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mariadb"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbAccount_basic -v
func TestAccTencentCloudMariadbAccount_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMariadbHourDbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMariadbHourDbAccountExists("tencentcloud_mariadb_account.account"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "user_name", "account-test"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "host", "10.101.202.22"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "read_only", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "description", "desc"),
				),
			},
			{
				ResourceName:            "tencentcloud_mariadb_account.account",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccMariadbAccountUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMariadbHourDbAccountExists("tencentcloud_mariadb_account.account"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "user_name", "account-test"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "host", "10.101.202.22"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "read_only", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_account.account", "description", "desc"),
				),
			},
		},
	})
}

func testAccCheckMariadbHourDbAccountDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmariadb.NewMariadbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mariadb_account" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		userName := idSplit[1]
		host := idSplit[2]

		account, err := service.DescribeMariadbAccount(ctx, instanceId, userName, host)
		if err != nil {
			return err
		}

		if account != nil {
			return fmt.Errorf("db account %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMariadbHourDbAccountExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		userName := idSplit[1]
		host := idSplit[2]

		service := svcmariadb.NewMariadbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		account, err := service.DescribeMariadbAccount(ctx, instanceId, userName, host)
		if err != nil {
			return err
		}

		if account == nil {
			return fmt.Errorf("db account %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMariadbAccount = testAccMariadbHourDbInstance + `

resource "tencentcloud_mariadb_account" "account" {
	instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
	user_name   = "account-test"
	host        = "10.101.202.22"
	password    = "Password123."
	read_only   = 0
	description = "desc"
  }

`

const testAccMariadbAccountUpdate = testAccMariadbHourDbInstance + `

resource "tencentcloud_mariadb_account" "account" {
	instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
	user_name   = "account-test"
	host        = "10.101.202.22"
	password    = "Update123."
	read_only   = 0
	description = "desc"
  }

`
