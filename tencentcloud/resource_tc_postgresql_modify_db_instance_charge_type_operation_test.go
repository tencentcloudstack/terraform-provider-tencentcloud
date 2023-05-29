package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestObjectPgModifyChargeType = "tencentcloud_postgresql_modify_db_instance_charge_type_operation.modify_db_instance_charge_type_operation"

func TestAccTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlModifyDbInstanceChargeTypeOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestObjectPgModifyChargeType, "id"),
					resource.TestCheckResourceAttrSet(TestObjectPgModifyChargeType, "db_instance_id"),
					resource.TestCheckResourceAttr(TestObjectPgModifyChargeType, "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(TestObjectPgModifyChargeType, "period", "2"),
					resource.TestCheckResourceAttr(TestObjectPgModifyChargeType, "auto_renew_flag", "0"),
					resource.TestCheckResourceAttr(TestObjectPgModifyChargeType, "auto_voucher", "0"),
				),
			},
		},
	})
}

const testAccPostgresqlModifyDbInstanceChargeTypeOperation = defaultVpcSubnets + `
data "tencentcloud_availability_zones_by_product" "zone" {
	product = "postgres"
  }
  
  data "tencentcloud_security_groups" "internal" {
	name = "default"
  }
  
  locals {
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
  }

resource "tencentcloud_postgresql_instance" "prepaid" {
	name 				= "tf_postsql_pre"
	availability_zone   = var.default_az
	charge_type 		= "POSTPAID_BY_HOUR"
	vpc_id  	  		= local.vpc_id
	subnet_id 		    = local.subnet_id
	engine_version	    = "13.3"
	root_password	    = "t1qaA2k1wgvfa3?ZZZ"
	security_groups     = [local.sg_id]
	charset			= "LATIN1"
	project_id 		= 0
	memory 			= 2
	storage 		= 20
  }


resource "tencentcloud_postgresql_modify_db_instance_charge_type_operation" "modify_db_instance_charge_type_operation" {
  db_instance_id = tencentcloud_postgresql_instance.prepaid.id
  instance_charge_type = "PREPAID"
  period = 2
  auto_renew_flag = 0
  auto_voucher = 0
}

`
