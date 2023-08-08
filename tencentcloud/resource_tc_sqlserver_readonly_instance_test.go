package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testReadonlySqlserverInstanceResourceName = "tencentcloud_sqlserver_readonly_instance"
var testReadonlySqlserverInstanceResourceKey = testReadonlySqlserverInstanceResourceName + ".example"

func TestAccTencentCloudReadonlySqlserverInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckReadonlySqlserverInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReadonlySqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReadonlySqlserverInstanceExists(testReadonlySqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "name", "tf_example"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "storage", "20"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "readonly_group_id"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "status"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "tags.test", "test"),
				),
			},
			{
				ResourceName:            testReadonlySqlserverInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_upgrade", "readonly_group_type"},
			},
			{
				Config: testAccReadonlySqlserverInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReadonlySqlserverInstanceExists(testReadonlySqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "name", "tf_example_update"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "memory", "8"),
					resource.TestCheckResourceAttr(testReadonlySqlserverInstanceResourceKey, "storage", "40"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "readonly_group_id"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testReadonlySqlserverInstanceResourceKey, "status"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "tags.update", "update"),
				),
			},
		},
	})
}

func testAccCheckReadonlySqlserverInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testReadonlySqlserverInstanceResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete SQL Server instance %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckReadonlySqlserverInstanceExists(n string) resource.TestCheckFunc {
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
			_, has, err = service.DescribeSqlserverInstanceById(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("SQL Server instance %s is not found", rs.Primary.ID)
		}
	}
}

const testAccReadonlySqlserverInstance string = defaultVpcSubnets + defaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_readonly_instance" "example" {
  name                = "tf_example"
  availability_zone   = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type         = "POSTPAID_BY_HOUR"
  vpc_id              = local.vpc_id
  subnet_id           = local.subnet_id
  memory              = 4
  storage             = 20
  master_instance_id  = tencentcloud_sqlserver_basic_instance.example.id
  readonly_group_type = 1
  force_upgrade       = true

  tags = {
    "test" = "test"
  }
}
`

const testAccReadonlySqlserverInstanceUpdate string = defaultVpcSubnets + defaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_readonly_instance" "example" {
  name                = "tf_example_update"
  availability_zone   = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type         = "POSTPAID_BY_HOUR"
  vpc_id              = local.vpc_id
  subnet_id           = local.subnet_id
  memory              = 8
  storage             = 40
  master_instance_id  = tencentcloud_sqlserver_basic_instance.example.id
  readonly_group_type = 1
  force_upgrade       = true

  tags = {
    "update" = "update"
  }
}
`
