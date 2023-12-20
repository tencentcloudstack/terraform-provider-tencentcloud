package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbEncryptAttributesConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() { tcacctest.AccPreCheck(t) },
		// PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDcdbEncryptAttributesConfig, "encrypt_attributes"),
				// PreventDiskCleanup: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_encrypt_attributes_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_encrypt_attributes_config.config", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_encrypt_attributes_config.config", "encrypt_enabled", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_encrypt_attributes_config.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// const testAccDcdbConfig_common_ins = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
const testAccDcdbConfig_common_ins = tcacctest.DefaultAzVariable + `
data "tencentcloud_security_groups" "internal" {
	name = "default"
}

data "tencentcloud_vpc_instances" "vpc" {
	name ="Default-VPC"
}
	
data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

locals {
	vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}

// resource "tencentcloud_dcdb_db_instance" "post_instance" {
// 	instance_name = "test_dcdb_db_post_instance_"
// 	zones = [var.default_az]
// 	period = 1
// 	shard_memory = "2"
// 	shard_storage = "10"
// 	shard_node_count = "2"
// 	shard_count = "2"
// 	vpc_id = local.vpc_id
// 	subnet_id = local.subnet_id
// 	db_version_id = "8.0"
// 	resource_tags {
// 	  tag_key = "aaa"
// 	  tag_value = "bbb"
// 	}
// 	security_group_ids = [local.sg_id]
// }

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
	instance_name = "test_dcdb_db_hourdb_instance_%s"
	zones = [var.default_az]
	shard_memory = "2"
	shard_storage = "10"
	shard_node_count = "2"
	shard_count = "2"
	vpc_id = local.vpc_id
	subnet_id = local.subnet_id
	security_group_id = local.sg_id
	db_version_id = "8.0"
	resource_tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
}

  locals {
	// dcdb_id = tencentcloud_dcdb_db_instance.post_instance.id
	dcdb_id = tencentcloud_dcdb_hourdb_instance.hourdb_instance.id
  }

`

const testAccDcdbEncryptAttributesConfig = testAccDcdbConfig_common_ins + `

resource "tencentcloud_dcdb_encrypt_attributes_config" "config" {
  instance_id = local.dcdb_id
  encrypt_enabled = 1
}

`
