package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testSqlserverAccountResourceName = "tencentcloud_sqlserver_account"
var testSqlserverAccountResourceKey = testSqlserverAccountResourceName + ".test"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_account
	resource.AddTestSweepers("tencentcloud_sqlserver_account", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_account",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := SqlserverService{client}

			db, err := service.DescribeSqlserverInstances(ctx, "", defaultSQLServerName, -1, "", "", -1)

			if err != nil {
				return err
			}

			if len(db) == 0 {
				return fmt.Errorf("%s not exists", defaultSQLServerName)
			}

			instanceId := *db[0].InstanceId

			accounts, err := service.DescribeSqlserverAccounts(ctx, instanceId)

			for i := range accounts {
				account := accounts[i]
				name := *account.Name
				created, err := time.Parse("2006-01-02 15:04:05", *account.CreateTime)
				if err != nil {
					created = time.Time{}
				}
				if isResourcePersist(name, &created) {
					continue
				}
				err = service.DeleteSqlserverAccount(ctx, instanceId, name)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudSqlserverAccountResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
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

const testAccSqlserverAccount string = CommonPresetSQLServer + `
resource "tencentcloud_sqlserver_account" "test" {
  instance_id = local.sqlserver_id
  name = "tf_sqlserver_account"
  password = "testt123"
}
`

const testAccSqlserverAccountUpdate string = CommonPresetSQLServer + `
resource "tencentcloud_sqlserver_account" "test" {
  instance_id = local.sqlserver_id
  name = "tf_sqlserver_account"
  password = "test1233"
  remark = "testt"
}
`
