package dcdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcdb"

	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dcdb_account", &resource.Sweeper{
		Name: "tencentcloud_dcdb_account",
		F:    testSweepDCDBAccount,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dcdb_account
func testSweepDCDBAccount(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	dcdbService := svcdcdb.NewDcdbService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	account, err := dcdbService.DescribeDcdbAccount(ctx, tcacctest.DefaultDcdbInstanceId, "")
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("dcdb account not exists. instanceId:[%s]", tcacctest.DefaultDcdbInstanceId)
	}

	for _, v := range account.Users {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := dcdbService.DeleteDcdbAccountById(ctx, tcacctest.DefaultDcdbInstanceId, *v.UserName, *v.Host)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("[ERROR] delete dcdb account %s reason:[%s]", *v.UserName, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudDcdbAccountResource_basic(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Nanosecond()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDcdbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDcdbAccount_basic, tcacctest.DefaultDcdbInstanceId, timestamp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcdbAccountExists("tencentcloud_dcdb_account.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.basic", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.basic", "user_name"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "host", "127.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "password", "===password==="),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "description", "this is a test account"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "max_user_connections", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "read_only", "0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDcdbAccount_update, tcacctest.DefaultDcdbInstanceId, timestamp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcdbAccountExists("tencentcloud_dcdb_account.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.basic", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.basic", "user_name"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "host", "127.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account.basic", "password", "===password===updated==="),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dcdbService := svcdcdb.NewDcdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dcdb_account" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		instanceId := idSplit[0]
		userName := idSplit[1]

		account, err := dcdbService.DescribeDcdbAccount(ctx, instanceId, userName)
		if err != nil {
			return err
		}
		if account != nil && len(account.Users) > 0 {
			return fmt.Errorf("dcdb account still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDcdbAccountExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("dcdb account %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dcdb account id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		instanceId := idSplit[0]
		userName := idSplit[1]

		dcdbService := svcdcdb.NewDcdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
	user_name = "mysql_%d"
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
  user_name = "mysql_%d"
  host = "127.0.0.1"
  password = "===password===updated==="
  read_only = 0
  description = "this is a changed test account"
  max_user_connections = 10
}

`
