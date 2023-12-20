package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEipsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.data_eips"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.data_eips", "eip_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.eip_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.data_eips", "eip_list.0.eip_name", "tf-test-eip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.eip_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.public_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.create_time"),

					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.tags"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.tags", "eip_list.0.tags.test", "test"),
				),
			},
		},
	})
}

const testAccEipsDataSource = `
resource "tencentcloud_eip" "eip" {
  name = "tf-test-eip"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_eips" "data_eips" {
  eip_id = tencentcloud_eip.eip.id
}

data "tencentcloud_eips" "tags" {
  tags = tencentcloud_eip.eip.tags
}
`
