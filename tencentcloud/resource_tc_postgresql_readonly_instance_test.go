package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testPostgresqlReadonlyInstanceResourceKey = "tencentcloud_postgresql_readonly_instance.instance"

func TestAccTencentCloudPostgresqlReadonlyInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccPostgresqlReadonlyInstanceInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "name", "tf_ro_instance_test"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "db_version", "13.3"),
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
		},
	})
}

const testAccPostgresqlReadonlyInstanceInstance string = CommonPresetPGSQL + defaultVpcSubnets + defaultSecurityGroupData + `
  resource "tencentcloud_postgresql_readonly_instance" "instance" {
	auto_renew_flag       = 0
	db_version            = "13.3"
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
