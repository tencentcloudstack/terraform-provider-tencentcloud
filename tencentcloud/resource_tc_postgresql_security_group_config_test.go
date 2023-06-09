package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestAccPostgresqlSecurityGroupConfigObject = "tencentcloud_postgresql_security_group_config.security_group_config"

func TestAccTencentCloudPostgresqlSecurityGroupConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlSecurityGroupConfig_ins,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlSecurityGroupConfigObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlSecurityGroupConfigObject, "security_group_id_set.#", "2"),
					resource.TestCheckResourceAttrSet(TestAccPostgresqlSecurityGroupConfigObject, "db_instance_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudPostgresqlSecurityGroupConfigResource_ro(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlSecurityGroupConfig_ro,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlSecurityGroupConfigObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlSecurityGroupConfigObject, "security_group_id_set.#", "2"),
					resource.TestCheckResourceAttrSet(TestAccPostgresqlSecurityGroupConfigObject, "read_only_group_id"),
				),
			},
		},
	})
}

const testAccPostgresqlSecurityGroupConfig_ins = CommonPresetPGSQL + defaultSecurityGroupData + `

resource "tencentcloud_postgresql_security_group_config" "security_group_config" {
  security_group_id_set = [local.sg_id, local.sg_id2]
  db_instance_id = local.pgsql_id
}

`

const testAccPostgresqlSecurityGroupConfig_ro = CommonPresetPGSQL + defaultVpcSubnets + defaultSecurityGroupData + `
resource "tencentcloud_postgresql_readonly_group" "group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_test_ro_sg"
	project_id = 0
	subnet_id             = local.subnet_id
	vpc_id                = local.vpc_id
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
  }

resource "tencentcloud_postgresql_security_group_config" "security_group_config" {
  security_group_id_set = [local.sg_id, local.sg_id2]
  read_only_group_id = tencentcloud_postgresql_readonly_group.group.id
}

`
