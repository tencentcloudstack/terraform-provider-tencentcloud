package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationAttributeDataSource_basic -v
func TestAccTencentCloudTsfApplicationAttributeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationAttributeDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application_attribute.application_attribute"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_attribute.application_attribute", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_attribute.application_attribute", "result.0.group_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_attribute.application_attribute", "result.0.instance_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_attribute.application_attribute", "result.0.run_instance_count"),
				),
			},
		},
	})
}

const testAccTsfApplicationAttributeDataSource = `

data "tencentcloud_tsf_application_attribute" "application_attribute" {
  application_id = "application-a24x29xv"
}

`
