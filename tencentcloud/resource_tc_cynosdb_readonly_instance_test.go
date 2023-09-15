package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCynosdbReadonlyInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCynosdbReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbReadonlyInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbReadonlyInstanceExists("tencentcloud_cynosdb_readonly_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_name", "tf-cynosdb-readonly-instance"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "force_delete", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_cpu_core", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_memory_size", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_duration", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_start_time", "10800"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_weekdays.#", "7"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_readonly_instance.foo", "instance_memory_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_readonly_instance.foo", "instance_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_readonly_instance.foo", "instance_storage_size"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "vpc_id", "vpc-4owdpnwr"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "subnet_id", "subnet-m4qpx38w"),
				),
			},
			{
				Config: testAccCynosdbReadonlyInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_duration", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_start_time", "21600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_weekdays.#", "6"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_cpu_core", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_memory_size", "4"),
				),
			},
			{
				ResourceName:            "tencentcloud_cynosdb_readonly_instance.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
		},
	})
}

func testAccCheckCynosdbReadonlyInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cynosdbService := CynosdbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_readonly_instance" {
			continue
		}

		_, _, has, err := cynosdbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("cynosdb readonly instance still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbReadonlyInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb readonly instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb readonly instance id is not set")
		}
		cynosdbService := CynosdbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, _, has, err := cynosdbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("cynosdb readonly instance doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const readonlyInstanceVar = `
variable "readonly_subnet" {
  default = "subnet-m4qpx38w"
}
`

const testAccCynosdbReadonlyInstance = testAccCynosdbBasic + readonlyInstanceVar + `
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2
  param_items {
    name = "character_set_server"
    current_value = "utf8"
  }

#  tags = {
#    test = "test"
#  }

  force_delete = true

  rw_group_sg = [
    "` + defaultSecurityGroup + `",
  ]
}

resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = tencentcloud_cynosdb_cluster.foo.id
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 1
  instance_memory_size = 2
  vpc_id               = var.my_vpc
  subnet_id            = var.readonly_subnet

  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]
}
`

const testAccCynosdbReadonlyInstance_update = testAccCynosdbBasic + `
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2
  param_items {
    name = "character_set_server"
    current_value = "utf8"
  }

#  tags = {
#    test = "test"
#  }

  force_delete = true

  rw_group_sg = [
    "` + defaultSecurityGroup + `",
  ]
}

resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = tencentcloud_cynosdb_cluster.foo.id
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 2
  instance_memory_size = 4

  instance_maintain_duration   = 7200
  instance_maintain_start_time = 21600
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Thu",
    "Wed",
    "Tue",
  ]
}
`
