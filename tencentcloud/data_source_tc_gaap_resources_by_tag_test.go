package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapResourcesByTagDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapResourcesByTagDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_resources_by_tag.resources_by_tag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_resources_by_tag.resources_by_tag", "resource_set.#"),
				),
			},
		},
	})
}

const testAccGaapResourcesByTagDataSource = `
data "tencentcloud_gaap_resources_by_tag" "resources_by_tag" {
  tag_key = "tagKey"
  tag_value = "tagValue"
}
`
