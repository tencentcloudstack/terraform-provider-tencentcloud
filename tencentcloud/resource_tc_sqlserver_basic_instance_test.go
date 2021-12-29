package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testSqlserverBasicInstanceResourceName = "tencentcloud_sqlserver_basic_instance"
var testSqlserverBasicInstanceResourceKey = testSqlserverBasicInstanceResourceName + ".test"

func TestAccTencentCloudSqlserverBasicInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverBasicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBasicInstancePostpaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBasicInstanceExists(testSqlserverBasicInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "name", "tf_sqlserver_basic_instance"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "memory", "2"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "storage", "20"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "cpu", "1"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "machine_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "vip"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "vport"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "status"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "auto_renew"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "tags.test", "test"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "maintenance_start_time", "09:00"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "maintenance_time_span", "3"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "maintenance_week_set.#", "3"),
				),
			},
			{
				Config: testAccSqlserverBasicInstancePostpaidUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBasicInstanceExists(testSqlserverBasicInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "name", "tf_sqlserver_basic_instance_update"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "maintenance_start_time", "08:00"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "maintenance_time_span", "4"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "maintenance_week_set.#", "4"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "storage", "100"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "cpu", "1"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "machine_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "vip"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "vport"),
					resource.TestCheckResourceAttrSet(testSqlserverBasicInstanceResourceKey, "status"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "security_groups.#", "1"),
					resource.TestCheckNoResourceAttr(testSqlserverBasicInstanceResourceKey, "tags.test"),
					resource.TestCheckResourceAttr(testSqlserverBasicInstanceResourceKey, "tags.abc", "abc"),
				),
			},
			{
				ResourceName:            testSqlserverBasicInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_voucher", "period"},
			},
		},
	})
}

func testAccCheckSqlserverBasicInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testSqlserverBasicInstanceResourceName {
			continue
		}

		_, has, err := service.DescribeSqlserverInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has {
			return fmt.Errorf("delete SQL Server Basic instance %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverBasicInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("SQL Server Basic instance %s is not found", rs.Primary.ID)
		}
		return nil
	}
}

const testAccSqlserverBasicInstanceBasic = `
variable "availability_zone"{
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
	name       = "tf-sqlserver-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
	availability_zone = "ap-guangzhou-3"
	name              = "tf-sqlserver-subnet"
	vpc_id            = tencentcloud_vpc.foo.id
	cidr_block        = "10.0.0.0/16"
	is_multicast      = false
}
`

const testAccSqlserverBasicInstancePostpaid string = testAccSqlserverBasicInstanceBasic + `
resource "tencentcloud_sqlserver_basic_instance" "test" {
	name                    = "tf_sqlserver_basic_instance"
	availability_zone       = var.availability_zone
	charge_type             = "POSTPAID_BY_HOUR"
	vpc_id                  = tencentcloud_vpc.foo.id
	subnet_id               = tencentcloud_subnet.foo.id
	security_groups         = ["` + defaultSecurityGroup + `"]
	project_id              = 0
	memory                  = 2
	storage                 = 20
	cpu                     = 1
	machine_type            = "CLOUD_PREMIUM"
	maintenance_week_set    = [1,2,3]
	maintenance_start_time  = "09:00"
	maintenance_time_span   = 3

	tags = {
		"test" = "test"
	}
}
`

const testAccSqlserverBasicInstancePostpaidUpdate string = testAccSqlserverBasicInstanceBasic + `
resource "tencentcloud_sqlserver_basic_instance" "test" {
	name                    = "tf_sqlserver_basic_instance_update"
	availability_zone       = var.availability_zone
	charge_type             = "POSTPAID_BY_HOUR"
	vpc_id                  = tencentcloud_vpc.foo.id
	subnet_id               = tencentcloud_subnet.foo.id
	security_groups         = ["` + defaultSecurityGroup + `"]
	project_id              = 0
	memory                  = 4
	storage                 = 100
	cpu                     = 1
	machine_type            = "CLOUD_PREMIUM"
	maintenance_week_set    = [1,2,3,4]
	maintenance_start_time  = "08:00"
	maintenance_time_span   = 4

	tags = {
		abc = "abc"
	}
}
`
