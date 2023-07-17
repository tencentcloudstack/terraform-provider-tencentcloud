package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testPostgresqlReadonlyGroupResourceKey = "tencentcloud_postgresql_readonly_group.group"

func TestAccTencentCloudPostgresqlReadonlyGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Config: testAccPostgresqlReadonlyGroupInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "name", "tf_ro_group_test"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "replay_lag_eliminate", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "replay_latency_eliminate", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "max_replay_lag", "100"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "max_replay_latency", "512"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "min_delay_eliminate_reserve", "1"),
				),
			},
			{
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Config: testAccPostgresqlReadonlyGroupInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "name", "tf_ro_group_test_updated"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "subnet_id"),
				),
			},
		},
	})
}

const testAccPostgresqlReadonlyGroupInstance string = CommonPresetPGSQL + defaultVpcSubnets + `
resource "tencentcloud_postgresql_readonly_group" "group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_ro_group_test"
	project_id = 0
	vpc_id  = local.vpc_id
	subnet_id 	= local.subnet_id
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
  }
`

const testAccPostgresqlReadonlyGroupInstance_update string = CommonPresetPGSQL + defaultVpcSubnets + `
resource "tencentcloud_vpc" "vpc" {
	cidr_block = "172.18.111.0/24"
	name       = "test-pg-rogroup-network-vpc"
  }
  
  resource "tencentcloud_subnet" "subnet" {
	availability_zone = var.default_az
	cidr_block        = "172.18.111.0/24"
	name              = "test-pg-rogroup-network-sub1"
	vpc_id            = tencentcloud_vpc.vpc.id
  }

  locals {
	new_vpc_id = tencentcloud_subnet.subnet.vpc_id
	new_subnet_id = tencentcloud_subnet.subnet.id
  }

resource "tencentcloud_postgresql_readonly_group" "group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_ro_group_test_updated"
	project_id = 0
	vpc_id  	  		= local.new_vpc_id
	subnet_id 		= local.new_subnet_id
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
  }
`
