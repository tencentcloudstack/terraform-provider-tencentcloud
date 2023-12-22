package sqlserver_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcsqlserver "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_db
	resource.AddTestSweepers("tencentcloud_sqlserver_db", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_db",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svcsqlserver.NewSqlserverService(client)

			instances, err := service.DescribeSqlserverInstances(ctx, "", "", -1, "", "", -1)

			if err != nil {
				return err
			}

			var (
				insId    string
				subInsId string
			)

			for _, v := range instances {
				if *v.Name == tcacctest.DefaultSQLServerName {
					insId = *v.InstanceId
				}
				if *v.Name == tcacctest.DefaultSubSQLServerName {
					subInsId = *v.InstanceId
				}
			}

			dbs, err := service.DescribeDBsOfInstance(ctx, insId)

			if err != nil {
				return err
			}

			for i := range dbs {
				db := dbs[i]
				if !strings.HasPrefix(*db.Name, "test") {
					continue
				}
				err := service.DeleteSqlserverDB(ctx, insId, []*string{db.Name})
				if err != nil {
					continue
				}
			}

			// Clear sub instance db
			subDbs, err := service.DescribeDBsOfInstance(ctx, subInsId)

			for i := range subDbs {
				db := subDbs[i]
				if *db.Name == tcacctest.DefaultSQLServerPubSubDB {
					err = service.DeleteSqlserverDB(ctx, subInsId, []*string{db.Name})
					break
				}
			}

			if err != nil {
				log.Printf("Delete sub instance DB fail: %s", err.Error())
			}
			return nil
		},
	})
}

func TestAccTencentCloudSqlserverDB_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSqlserverDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDB_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.example", "name", "tf_example_db"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.example", "charset", "Chinese_PRC_BIN"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.example", "remark", "test-remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.example", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.example", "instance_id"),
				),
				Destroy: false,
			},
			{
				ResourceName:      "tencentcloud_sqlserver_db.mysqlserver_db",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverDB_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverDBExists("tencentcloud_sqlserver_db.example"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.example", "remark", "test-remark-update"),
				),
			},
		},
	})
}

func testAccCheckSqlserverDBDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	sqlserverService := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_db" {
			continue
		}
		_, has, err := sqlserverService.DescribeDBDetailsById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("SQL Server DB still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckSqlserverDBExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("SQL Server DB %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SQL Server DB id is not set")
		}

		sqlserverService := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := sqlserverService.DescribeDBDetailsById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("SQL Server DB %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccSqlserverDB_basic = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}
`

const testAccSqlserverDB_basic_update_remark = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db_update"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark-update"
}
`
