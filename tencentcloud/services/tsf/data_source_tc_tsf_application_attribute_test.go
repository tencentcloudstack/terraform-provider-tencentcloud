package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationAttributeDataSource_basic -v
func TestAccTencentCloudTsfApplicationAttributeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationAttributeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application_attribute.application_attribute"),
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
