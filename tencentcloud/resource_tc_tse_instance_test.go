package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_instance.instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseInstance = `

resource "tencentcloud_tse_instance" "instance" {
  engine_type = "nacos"
  engine_version = "2.0.3"
  engine_product_version = "STANDARD"
  engine_region = "ap-guangzhou"
  engine_name = "nacos-test"
  trade_type = 0
  engine_resource_spec = "STANDARD"
  engine_node_num = 3
  vpc_id = "vpc-xxxxxx"
  subnet_id = "subnet-xxxxxx"
  apollo_env_params {
		name = "dev"
		engine_resource_spec = "1C2G"
		engine_node_num = 3
		storage_capacity = 20
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"
		env_desc = "dev env"

  }
  engine_tags {
		tag_key = ""
		tag_value = ""

  }
  engine_admin {
		name = "admin"
		password = "admin"
		token = "xxxxxx"

  }
  prepaid_period = 0
  prepaid_renew_flag = 1
  engine_region_infos {
		engine_region = "ap-guangzhou"
		replica = 3
		vpc_infos {
			vpc_id = "vpc-xxxxxx"
			subnet_id = "subnet-xxxxxx"
			intranet_address = ""
		}

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
