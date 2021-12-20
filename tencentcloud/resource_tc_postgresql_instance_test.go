package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testPostgresqlInstanceResourceName = "tencentcloud_postgresql_instance"
var testPostgresqlInstanceResourceKey = testPostgresqlInstanceResourceName + ".test"

func TestAccTencentCloudPostgresqlInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPostgresqlInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "name", "tf_postsql_instance"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "storage", "100"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "public_access_switch", "false"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "root_password", "1qaA2k1wgvfa3ZZZ"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "tags.tf", "test"),
				),
			},
			{
				ResourceName:            testPostgresqlInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "spec_code", "public_access_switch", "charset"},
			},

			{
				Config: testAccPostgresqlInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "name", "tf_postsql_instance_update"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "storage", "250"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "project_id", "1154137"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "public_access_switch", "true"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "root_password", "1qaA2k1wgvfa3ZZZZ"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "public_access_host"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "public_access_port"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "tags.tf", "teest"),
				),
			},
		},
	})
}

func testAccCheckPostgresqlInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testPostgresqlInstanceResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)

		if !has {
			return nil
		} else {
			if err != nil {
				return err
			}
			return fmt.Errorf("delete postgresql instance %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPostgresqlInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("postgresql instance %s is not found", rs.Primary.ID)
		}
	}
}

const testAccPostgresqlInstanceBasic = `
variable "availability_zone"{
default = "ap-guangzhou-2"
}
`

const testAccPostgresqlInstance string = testAccPostgresqlInstanceBasic + `
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance"
  availability_zone = var.availability_zone
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id                   = "` + defaultVpcId + `"
  subnet_id = "subnet-pyio7yog"
  engine_version		= "10.4"
  root_password                 = "1qaA2k1wgvfa3ZZZ"
  charset = "LATIN1"
  project_id = 0
  memory = 4
  storage = 100

	tags = {
		tf = "test"
	}
}
`

const testAccPostgresqlInstanceUpdate string = testAccPostgresqlInstanceBasic + `
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance_update"
  availability_zone = var.availability_zone
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id                   = "` + defaultVpcId + `"
  subnet_id = "subnet-pyio7yog"
  engine_version		= "10.4"
  root_password                 = "1qaA2k1wgvfa3ZZZZ"
  charset = "LATIN1"
  project_id = 1154137
  public_access_switch = true
  memory = 4
  storage = 250

	tags = {
		tf = "teest"
	}
}
`
