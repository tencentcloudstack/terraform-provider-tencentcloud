package sqlserver_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcsqlserver "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testSqlserverAccountResourceName = "tencentcloud_sqlserver_account"
var testSqlserverAccountResourceKey = testSqlserverAccountResourceName + ".example"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_account
	resource.AddTestSweepers("tencentcloud_sqlserver_account", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_account",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()

			service := svcsqlserver.NewSqlserverService(client)

			db, err := service.DescribeSqlserverInstances(ctx, "", tcacctest.DefaultSQLServerName, -1, "", "", -1)

			if err != nil {
				return err
			}

			if len(db) == 0 {
				return fmt.Errorf("%s not exists", tcacctest.DefaultSQLServerName)
			}

			instanceId := *db[0].InstanceId

			accounts, _ := service.DescribeSqlserverAccounts(ctx, instanceId)

			for i := range accounts {
				account := accounts[i]
				name := *account.Name
				created, err := time.Parse("2006-01-02 15:04:05", *account.CreateTime)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist(name, &created) {
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
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSqlserverAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountExists(testSqlserverAccountResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "name", "tf_example_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "password", "Qwer@234"),
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
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "name", "tf_example_db_update"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "password", "Qwer@234Update"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "remark", "test-remark-update"),
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, tccommon.FILED_SP)
		if len(idStrs) != 2 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		name := idStrs[1]

		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, tccommon.FILED_SP)
		if len(idStrs) != 2 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		name := idStrs[1]

		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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

const testAccSqlserverAccount string = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234"
  remark      = "test-remark"
}
`

const testAccSqlserverAccountUpdate string = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234Update"
  remark      = "test-remark-update"
}
`
