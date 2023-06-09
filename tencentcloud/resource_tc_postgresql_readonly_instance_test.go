package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testPostgresqlReadonlyInstanceResourceKey = "tencentcloud_postgresql_readonly_instance.instance"

func TestAccTencentCloudPostgresqlReadonlyInstanceResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccPostgresqlReadonlyInstanceInstance_basic_without_rogroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "name", "tf_ro_instance_test"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "db_version", "15.1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "zone"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "security_groups_ids.#", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "storage", "20"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_port"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccPostgresqlReadonlyInstanceInstance_basic_update_rogroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "name", "tf_ro_instance_test_updated"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "db_version", "15.1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "zone"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "security_groups_ids.#", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "storage", "20"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "read_only_group_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudPostgresqlReadonlyInstanceResource_update_ro_group(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccPostgresqlReadonlyInstanceInstance_with_rogroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "name", "tf_ro_instance_test_rog"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "db_version", "15.1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "zone"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "security_groups_ids.#", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "storage", "20"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_port"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccPostgresqlReadonlyInstanceInstance_update_rogroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "name", "tf_ro_instance_test_rog_updated"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "db_version", "15.1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "zone"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "security_groups_ids.#", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "storage", "20"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "read_only_group_id"),
				),
			},
		},
	})
}

const testAccPostgresqlReadonlyInstanceInstance_basic_without_rogroup string = OperationPresetPGSQL + defaultVpcSubnets + defaultSecurityGroupData + `
  resource "tencentcloud_postgresql_readonly_instance" "instance" {
	auto_renew_flag       = 0
	db_version            = "15.1"
	instance_charge_type  = "POSTPAID_BY_HOUR"
	master_db_instance_id = local.pgsql_id
	memory                = 4
	name                  = "tf_ro_instance_test"
	need_support_ipv6     = 0
	project_id            = 0
	security_groups_ids   = [
	  local.sg_id,
	]
	storage               = 20
	vpc_id                = local.vpc_id
	subnet_id 	          = local.subnet_id
	zone                  = var.default_az
  }
`

const testAccPostgresqlReadonlyInstanceInstance_basic_update_rogroup string = OperationPresetPGSQL + defaultVpcSubnets + defaultSecurityGroupData + `
  resource "tencentcloud_postgresql_readonly_instance" "instance" {
	auto_renew_flag       = 0
	db_version            = "15.1"
	instance_charge_type  = "POSTPAID_BY_HOUR"
	master_db_instance_id = local.pgsql_id
	memory                = 4
	name                  = "tf_ro_instance_test_updated"
	need_support_ipv6     = 0
	project_id            = 0
	security_groups_ids   = [
	  local.sg_id,
	]
	storage               = 20
	vpc_id                = local.vpc_id
	subnet_id 	          = local.subnet_id
	zone                  = var.default_az
	read_only_group_id    = tencentcloud_postgresql_readonly_group.new_ro_group.id
  }

  resource "tencentcloud_postgresql_readonly_group" "new_ro_group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_ro_group_test_new"
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

const testAccPostgresqlReadonlyInstanceInstance_with_rogroup string = OperationPresetPGSQL + defaultVpcSubnets + defaultSecurityGroupData + `
  resource "tencentcloud_postgresql_readonly_instance" "instance" {
	auto_renew_flag       = 0
	db_version            = "15.1"
	instance_charge_type  = "POSTPAID_BY_HOUR"
	master_db_instance_id = local.pgsql_id
	memory                = 4
	name                  = "tf_ro_instance_test_rog"
	need_support_ipv6     = 0
	project_id            = 0
	security_groups_ids   = [
	  local.sg_id,
	]
	storage               = 20
	vpc_id                = local.vpc_id
	subnet_id 	          = local.subnet_id
	zone                  = var.default_az
	read_only_group_id    = local.pgrogroup_id
  }
`

const testAccPostgresqlReadonlyInstanceInstance_update_rogroup string = OperationPresetPGSQL + defaultVpcSubnets + defaultSecurityGroupData + `
  resource "tencentcloud_postgresql_readonly_instance" "instance" {
	auto_renew_flag       = 0
	db_version            = "15.1"
	instance_charge_type  = "POSTPAID_BY_HOUR"
	master_db_instance_id = local.pgsql_id
	memory                = 4
	name                  = "tf_ro_instance_test_rog_updated"
	need_support_ipv6     = 0
	project_id            = 0
	security_groups_ids   = [
	  local.sg_id,
	]
	storage               = 20
	vpc_id                = local.vpc_id
	subnet_id 	          = local.subnet_id
	zone                  = var.default_az
	read_only_group_id    = tencentcloud_postgresql_readonly_group.new_ro_group.id
  }

  resource "tencentcloud_postgresql_readonly_group" "new_ro_group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_ro_group_test_new"
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
