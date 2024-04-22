package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testPostgresqlReadonlyGroupResourceKey = "tencentcloud_postgresql_readonly_group.group"

func TestAccTencentCloudPostgresqlReadonlyGroupResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
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
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "min_delay_eliminate_reserve", "1"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "net_info_list.#"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "net_info_list.0.ip"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "net_info_list.0.port", "5432"),
				),
			},
			{
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Config: testAccPostgresqlReadonlyGroupInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "name", "tf_ro_group_test_updated"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "subnet_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "net_info_list.#"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "net_info_list.0.ip"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "net_info_list.0.port", "5432"),
				),
			},
		},
	})
}

const testAccPostgresqlReadonlyGroupInstance string = tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_postgresql_instance" "test" {
	name              = "example"
	availability_zone = var.default_az
	charge_type       = "POSTPAID_BY_HOUR"
	vpc_id            = local.vpc_id
	subnet_id         = local.subnet_id
	engine_version    = "10.4"
	root_user         = "root123"
	root_password     = "Root123$"
	charset           = "UTF8"
	project_id        = 0
	memory            = 2
	storage           = 10
  
	tags = {
	  test = "tf"
	}
  }

resource "tencentcloud_postgresql_readonly_group" "group" {
	master_db_instance_id = tencentcloud_postgresql_instance.test.id
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

const testAccPostgresqlReadonlyGroupInstance_update string = tcacctest.DefaultVpcSubnets + `
  resource "tencentcloud_postgresql_instance" "test" {
	name              = "example"
	availability_zone = var.default_az
	charge_type       = "POSTPAID_BY_HOUR"
	vpc_id            = local.vpc_id
	subnet_id         = local.subnet_id
	engine_version    = "10.4"
	root_user         = "root123"
	root_password     = "Root123$"
	charset           = "UTF8"
	project_id        = 0
	memory            = 2
	storage           = 10
  
	tags = {
	  test = "tf"
	}
  }

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
	master_db_instance_id = tencentcloud_postgresql_instance.test.id
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
