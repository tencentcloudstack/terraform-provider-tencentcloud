package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlInstanceOperation_disisolate = "tencentcloud_postgresql_disisolate_db_instance_operation.disisolate_db_instance_operation"
const testAccPostgresqlInstanceOperation_isolate = "tencentcloud_postgresql_isolate_db_instance_operation.isolate_db_instance_operation"
const testAccPostgresqlInstanceOperation_renew = "tencentcloud_postgresql_renew_db_instance_operation.renew_db_instance_operation"
const testAccPostgresqlInstanceOperation_restart = "tencentcloud_postgresql_restart_db_instance_operation.restart_db_instance_operation"

func TestAccTencentCloudPostgresqlCommonDbInstanceOperationResource_all(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlCommonDbInstanceOperation_restart,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlInstanceOperation_restart, "id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlInstanceOperation_restart, "db_instance_id"),
				),
			},
			{
				Config: testAccPostgresqlCommonDbInstanceOperation_isolate,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlInstanceOperation_isolate, "id"),
					resource.TestCheckResourceAttr(testAccPostgresqlInstanceOperation_isolate, "db_instance_id_set.#", "1"),
				),
			},
			{
				Config: testAccPostgresqlCommonDbInstanceOperation_disisolate,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlInstanceOperation_disisolate, "id"),
					resource.TestCheckResourceAttr(testAccPostgresqlInstanceOperation_disisolate, "db_instance_id_set.#", "1"),
					resource.TestCheckResourceAttr(testAccPostgresqlInstanceOperation_disisolate, "period", "1"),
					resource.TestCheckResourceAttr(testAccPostgresqlInstanceOperation_disisolate, "auto_voucher", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudPostgresqlCommonDbInstanceOperationResource_renew(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlCommonDbInstanceOperation_renew,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlInstanceOperation_renew, "id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlInstanceOperation_renew, "db_instance_id"),
					resource.TestCheckResourceAttr(testAccPostgresqlInstanceOperation_renew, "period", "1"),
					resource.TestCheckResourceAttr(testAccPostgresqlInstanceOperation_renew, "auto_voucher", "0"),
				),
			},
		},
	})
}

const testAccPostgresqlCommonDbInstanceOperation_isolate = OperationPresetPGSQL + `

resource "tencentcloud_postgresql_isolate_db_instance_operation" "isolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
}

`

const testAccPostgresqlCommonDbInstanceOperation_disisolate = OperationPresetPGSQL + `

resource "tencentcloud_postgresql_disisolate_db_instance_operation" "disisolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
  period = 1
  auto_voucher = false
}

`

const testAccPostgresqlCommonDbInstanceOperation_restart = OperationPresetPGSQL + `

resource "tencentcloud_postgresql_restart_db_instance_operation" "restart_db_instance_operation" {
  db_instance_id = local.pgsql_id
}

`

const testAccPostgresqlCommonDbInstanceOperation_renew = testAccPostgresqlInstanceCommonBase_PREPAID + `

resource "tencentcloud_postgresql_renew_db_instance_operation" "renew_db_instance_operation" {
  db_instance_id = tencentcloud_postgresql_instance.oper_test_PREPAID.id
  period = 1
  auto_voucher = 0
}

`

const testAccPostgresqlInstanceCommonBase_PREPAID = defaultAzVariable + `

data "tencentcloud_security_groups" "sg" {
	name = "default"
  }

  data "tencentcloud_vpc_subnets" "gz3" {
	availability_zone = var.default_az
	is_default = true
  }

  locals {
	my_vpc_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
	my_subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
	my_sg_id = data.tencentcloud_security_groups.sg.security_groups.0.security_group_id
  }

resource "tencentcloud_postgresql_instance" "oper_test_PREPAID" {
	name 				= "tf_postsql_instance_common_operation"
	availability_zone   = var.default_az
	charge_type 		= "PREPAID"
	period              = 1
	vpc_id  	  		= local.my_vpc_id
	subnet_id 		    = local.my_subnet_id
	engine_version	    = "13.3"
	root_password	    = "t1qaA2k1wgvfa3?ZZZ"
	security_groups     = [local.my_sg_id]
	charset			    = "LATIN1"
	project_id 		    = 0
	memory 			    = 4
	storage 			= 20

	db_kernel_version = "v13.3_r1.1"

	tags = {
	  tf = "test"
	}
  }
`
