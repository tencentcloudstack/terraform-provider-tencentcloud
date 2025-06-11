package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMonitorPolicyBindingObjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorPolicyBindingObjectBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_binding_object.binding_object", "dimensions.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_policy_binding_object.binding_object",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudMonitorPolicyBindingObjectResource_multiRegion(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorPolicyBindingObjectMultiRegion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_binding_object.binding_multi_region_object", "dimensions.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_policy_binding_object.binding_multi_region_object",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudMonitorPolicyBindingObjectResource_noRegion(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorPolicyBindingNoRegionObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_binding_object.binding_multi_region_object", "dimensions.#", "1"),
				),
			},
		},
	})
}

const testAccMonitorPolicyBindingObjectBasic string = `
resource "tencentcloud_monitor_policy_binding_object" "binding_object" {
  policy_id = "policy-dkfebnac"
  dimensions {
    dimensions_json = jsonencode(
      {
        unInstanceId = "ins-df0kqv1o"
      }
    )
  }
}
`

const testAccMonitorPolicyBindingObjectMultiRegion string = `
resource "tencentcloud_monitor_policy_binding_object" "binding_multi_region_object" {
  policy_id = "policy-dkfebnac"
  dimensions {
    dimensions_json = jsonencode(
      {
        unInstanceId = "ins-2ferl9zf"
      }
    )
    region = "ap-shanghai"
  }
  dimensions {
    dimensions_json = jsonencode(
      {
        unInstanceId = "ins-df0kqv1o"
      }
    )
    region = "ap-guangzhou"
  }
}
`

const testAccMonitorPolicyBindingNoRegionObject string = `
resource "tencentcloud_monitor_policy_binding_object" "binding_multi_region_object" {
  policy_id = "policy-wt2kvmmq"
  dimensions {
    dimensions_json = jsonencode(
      {
        domain = "keep.tencentcloud-terraform-provider.cn"
      }
    )
  }
}
`
