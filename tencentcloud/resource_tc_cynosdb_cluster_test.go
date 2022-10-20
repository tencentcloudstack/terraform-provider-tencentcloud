package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_cynosdb", &resource.Sweeper{
		Name: "tencentcloud_cynosdb",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := CynosdbService{client: client}

			instances, err := service.DescribeClusters(ctx, nil)

			if err != nil {
				return err
			}

			for _, v := range instances {
				id := *v.ClusterId
				name := *v.ClusterName
				status := *v.Status
				if status != "running" {
					continue
				}
				if !strings.HasPrefix(name, "tf-cynosdb") {
					continue
				}
				_, err := service.IsolateCluster(ctx, id)
				if err != nil {
					continue
				}
				if err = service.OfflineCluster(ctx, id); err != nil {
					continue
				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudCynosdbClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCynosdbClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbClusterExists("tencentcloud_cynosdb_cluster.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "available_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "vpc_id", "vpc-4owdpnwr"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "subnet_id", "subnet-qpxez62e"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "db_type", "MYSQL"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "db_version", "5.7"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "storage_limit", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "cluster_name", "tf-cynosdb"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_maintain_duration", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_maintain_start_time", "10800"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_maintain_weekdays.#", "7"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_cpu_core", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_memory_size", "2"),
					//resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "force_delete", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "rw_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "ro_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "port", "5432"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "charge_type", CYNOSDB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "instance_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "instance_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "instance_storage_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "charset", "utf8"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "cluster_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "storage_used"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "rw_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "rw_group_instances.0.instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "rw_group_instances.0.instance_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "rw_group_addr.0.ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster.foo", "rw_group_addr.0.port"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.0.name", "character_set_server"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.0.current_value", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.1.name", "time_zone"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.1.current_value", "+09:00"),
				),
			},
			{
				ResourceName:            "tencentcloud_cynosdb_cluster.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "force_delete", "storage_limit", "param_items", "ro_group_sg", "prarm_template_id"},
			},
			{
				Config: testAccCynosdbCluster_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_maintain_duration", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_maintain_start_time", "21600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_maintain_weekdays.#", "6"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_cpu_core", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "instance_memory_size", "4"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "rw_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "ro_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.0.name", "character_set_server"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.0.old_value", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.0.current_value", "utf8mb4"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.1.name", "time_zone"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.1.old_value", "+09:00"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster.foo", "param_items.1.current_value", "+08:00"),
				),
			},
		},
	})
}

func testAccCheckCynosdbClusterDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cynosdbService := CynosdbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_cluster" {
			continue
		}

		_, _, has, err := cynosdbService.DescribeClusterById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("cynosdb cluster still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster id is not set")
		}
		cynosdbService := CynosdbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, _, has, err := cynosdbService.DescribeClusterById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("cynosdb cluster doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbBasic = `
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

variable "my_vpc" {
  default = "vpc-4owdpnwr"
}

variable "my_subnet" {
  default = "subnet-qpxez62e"
}

variable "my_param_template" {
	default = "15765"
}
`

const testAccCynosdbCluster = testAccCynosdbBasic + `
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
  param_items {
    name = "time_zone"
    current_value = "+09:00"
  }

#  tags = {
#    test = "test"
#  }

  force_delete = true

  rw_group_sg = [
    "` + defaultSecurityGroup + `",
  ]
  ro_group_sg = [
    "` + defaultSecurityGroup + `",
  ]
  prarm_template_id = var.my_param_template
}
`

const testAccCynosdbCluster_update = testAccCynosdbBasic + `
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 21600
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Tue",
  ]

  instance_cpu_core    = 2
  instance_memory_size = 4

  param_items {
    name = "character_set_server"
    old_value = "utf8"
    current_value = "utf8mb4"
  }

  param_items {
    name = "time_zone"
    old_value = "+09:00"
    current_value = "+08:00"
  }

#  tags = {
#    test = "test-update"
#  }

  force_delete = true

  rw_group_sg = [
    "` + defaultSecurityGroup2 + `",
  ]
  ro_group_sg = [
    "` + defaultSecurityGroup2 + `",
  ]
}
`
