package controlcenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudControlcenterAccountFactoryBaselineConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccControlcenterAccountFactoryBaselineConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_controlcenter_account_factory_baseline_config.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_controlcenter_account_factory_baseline_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccControlcenterAccountFactoryBaselineConfig = `
resource "tencentcloud_controlcenter_account_factory_baseline_config" "example" {
  name = "default"
  baseline_config_items {
    identifier = "TCC-AF_VPC_SUBNET"
    configuration = jsonencode({
      "VpcName" : "tf-example",
      "CidrBlock" : "10.0.0.0/16",
      "Region" : "1",
      "RegionName" : "ap-guangzhou",
      "Subnets" : [
        {
          "CidrBlock" : "10.0.0.0/24",
          "SubnetName" : "abc",
          "Zone" : "ap-guangzhou-6"
        }
      ]
    })
  }

  baseline_config_items {
    identifier    = "TCC-AF_PRESET_TAG"
    configuration = "{\"TagValuePairs\":[{\"Key\":\"key\",\"Values\":[\"value\"]}]}"
  }
}
`
