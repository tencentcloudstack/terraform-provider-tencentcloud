package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testSqlserverAccountDBAttachmentResourceName = "tencentcloud_sqlserver_account_db_attachment"
var testSqlserverAccountDBAttachmentResourceKey = testSqlserverAccountDBAttachmentResourceName + ".example"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_account_db_attachment
	resource.AddTestSweepers(testSqlserverAccountDBAttachmentResourceName, &resource.Sweeper{
		Name: testSqlserverAccountDBAttachmentResourceName,
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

			records, err := service.DescribeAccountDBAttachments(ctx, instanceId, defaultSQLServerAccount, defaultSQLServerDB)
			if err != nil {
				return err
			}

			if len(records) > 0 {
				err = service.ModifyAccountDBAttachment(ctx, instanceId, defaultSQLServerAccount, defaultSQLServerDB, "Delete")
			}

			if err != nil {
				return err
			}

			return nil
		},
	})
}

func TestAccTencentCloudSqlserverAccountDBAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverAccountDBAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverAccountDBAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountDBAttachmentExists(testSqlserverAccountDBAttachmentResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountDBAttachmentResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "account_name", "tf_example_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "db_name", "tf_example_db"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "privilege", "ReadOnly"),
				),
			},
			{
				ResourceName:      testSqlserverAccountDBAttachmentResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverAccountDBAttachmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountDBAttachmentExists(testSqlserverAccountDBAttachmentResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountDBAttachmentResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "account_name", "tf_example_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "db_name", "tf_example_db"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "privilege", "ReadWrite"),
				),
			},
		},
	})
}

func testAccCheckSqlserverAccountDBAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testSqlserverAccountDBAttachmentResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 3 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		accountName := idStrs[1]
		dbName := idStrs[2]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
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

func testAccCheckSqlserverAccountDBAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 3 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		accountName := idStrs[1]
		dbName := idStrs[2]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
		if err != nil {
			_, has, err = service.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
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

const testAccSqlserverAccountDBAttachment string = defaultVpcSubnets + defaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_account_db_attachment" "example" {
  instance_id  = tencentcloud_sqlserver_basic_instance.example.id
  account_name = tencentcloud_sqlserver_account.example.name
  db_name      = tencentcloud_sqlserver_db.example.name
  privilege    = "ReadOnly"
}
`

const testAccSqlserverAccountDBAttachmentUpdate string = defaultVpcSubnets + defaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_account_db_attachment" "example" {
  instance_id  = tencentcloud_sqlserver_basic_instance.example.id
  account_name = tencentcloud_sqlserver_account.example.name
  db_name      = tencentcloud_sqlserver_db.example.name
  privilege    = "ReadWrite"
}
`
