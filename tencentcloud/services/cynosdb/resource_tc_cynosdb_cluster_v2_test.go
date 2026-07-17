package cynosdb_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"
)

func init() {
	resource.AddTestSweepers("tencentcloud_cynosdb_v2", &resource.Sweeper{
		Name: "tencentcloud_cynosdb_v2",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svccynosdb.NewCynosdbService(client)

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

// go test -i; go test -test.run TestAccTencentCloudCynosdbClusterV2ResourceBasic -v
func TestAccTencentCloudCynosdbClusterV2ResourceBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCynosdbClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbV2Cluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbClusterV2Exists("tencentcloud_cynosdb_cluster_v2.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "available_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "vpc_id", "vpc-m0d2dbnn"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "subnet_id", "subnet-j10lsueq"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "db_type", "MYSQL"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "db_version", "5.7"),
					// resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "storage_limit", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "cluster_name", "tf-cynosdb"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_maintain_duration", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_maintain_start_time", "10800"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_maintain_weekdays.#", "7"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_cpu_core", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_memory_size", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "force_delete", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "rw_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "ro_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "port", "5432"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "charge_type", svccynosdb.CYNOSDB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "instance_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "instance_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "instance_storage_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "charset", "utf8"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "cluster_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "storage_used"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "rw_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "rw_group_instances.0.instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "rw_group_instances.0.instance_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "rw_group_addr.0.ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_v2.foo", "rw_group_addr.0.port"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.0.name", "character_set_server"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.0.current_value", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.1.name", "time_zone"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.1.current_value", "+09:00"),
				),
			},
			{
				ResourceName:            "tencentcloud_cynosdb_cluster_v2.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "force_delete", "storage_limit", "param_items", "ro_group_sg", "prarm_template_id"},
			},
			{
				Config: testAccCynosdbV2Cluster_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "vpc_id", "vpc-m0d2dbnn"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "subnet_id", "subnet-j10lsueq"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_maintain_duration", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_maintain_start_time", "10800"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_maintain_weekdays.#", "7"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_cpu_core", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "instance_memory_size", "4"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "rw_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "ro_group_sg.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.0.name", "character_set_server"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.0.old_value", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.0.current_value", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.1.name", "time_zone"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.1.old_value", "+09:00"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "param_items.1.current_value", "+09:00"),
				),
			},
		},
	})
}
func TestAccTencentCloudCynosdbClusterV2ResourceServerless(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCynosdbClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbV2ClusterServerless,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbClusterV2Exists("tencentcloud_cynosdb_cluster_v2.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "db_mode", "SERVERLESS"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_v2.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"force_delete",
					"storage_limit",
					"min_cpu",
					"max_cpu",
					"auto_pause",
					"auto_pause_delay",
				},
			},
			{
				Config: testAccCynosdbV2ClusterServerlessPause,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbClusterV2Exists("tencentcloud_cynosdb_cluster_v2.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "db_mode", "SERVERLESS"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "serverless_status", "pause"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "serverless_status_flag", "pause"),
				),
			},
			{
				Config: testAccCynosdbV2ClusterServerlessResume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbClusterV2Exists("tencentcloud_cynosdb_cluster_v2.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "db_mode", "SERVERLESS"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "serverless_status", "resume"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_v2.foo", "serverless_status_flag", "resume"),
				),
			},
		},
	})
}

func testAccCheckCynosdbClusterV2Destroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_cluster_v2" {
			continue
		}

		_, _, has, err := cynosdbService.DescribeClusterById(ctx, rs.Primary.ID)
		if err != nil {
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Message == "record not found" {
					return nil
				}
			}
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("cynosdb cluster still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbClusterV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster id is not set")
		}
		cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccCynosdbV2Basic = `
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

variable "my_vpc" {
  default = "vpc-m0d2dbnn"
}

variable "my_subnet" {
  default = "subnet-j10lsueq"
}

variable "my_param_template" {
	default = "15765"
}

variable "rw_group_sg" {
	default = "sg-05f7wnhn"
}
`

const testAccCynosdbV2Cluster = testAccCynosdbV2Basic + `
resource "tencentcloud_cynosdb_cluster_v2" "foo" {
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

  tags = {
    test = "test"
  }

  force_delete = true

  rw_group_sg = [
    var.rw_group_sg
  ]
  ro_group_sg = [
    var.rw_group_sg
  ]
#  prarm_template_id = var.my_param_template
}
`

const testAccCynosdbV2Cluster_update = testAccCynosdbV2Basic + `
resource "tencentcloud_cynosdb_cluster_v2" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb-update"
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

  instance_cpu_core    = 2
  instance_memory_size = 4
  param_items {
    name = "character_set_server"
	old_value = "utf8"
    current_value = "utf8"
  }
  param_items {
    name = "time_zone"
	old_value = "+09:00"
    current_value = "+09:00"
  }

  tags = {
    test = "test"
  }

  force_delete = true

  rw_group_sg = [
    var.rw_group_sg
  ]
  ro_group_sg = [
    var.rw_group_sg
  ]
}
`

const testAccCynosdbV2ClusterServerless = testAccCynosdbV2Basic + `
resource "tencentcloud_cynosdb_cluster_v2" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  cluster_name                 = "tf-cynosdb-s"
  password                     = "cynos@123"
  db_mode                      = "SERVERLESS"
  min_cpu 					   = 0.25
  max_cpu 					   = 1
  auto_pause 				   = "yes"
  auto_pause_delay 			   = 1000
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

  force_delete = true
}`
const testAccCynosdbV2ClusterServerlessPause = testAccCynosdbV2Basic + `
resource "tencentcloud_cynosdb_cluster_v2" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  cluster_name                 = "tf-cynosdb-s"
  password                     = "cynos@123"
  db_mode                      = "SERVERLESS"
  min_cpu 					   = 0.25
  max_cpu 					   = 1
  auto_pause 				   = "yes"
  auto_pause_delay 			   = 1000
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
  serverless_status_flag       = "pause"
  force_delete = true
}`
const testAccCynosdbV2ClusterServerlessResume = testAccCynosdbV2Basic + `
resource "tencentcloud_cynosdb_cluster_v2" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  cluster_name                 = "tf-cynosdb-s"
  password                     = "cynos@123"
  db_mode                      = "SERVERLESS"
  min_cpu 					   = 0.25
  max_cpu 					   = 1
  auto_pause 				   = "yes"
  auto_pause_delay 			   = 1000
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
  serverless_status_flag       = "resume"
  force_delete = true
}`
