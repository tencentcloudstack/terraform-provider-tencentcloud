package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudEipsDataSource_basic -v
func TestAccTencentCloudEipsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.example", "eip_list.#", "0"),

					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example_by_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_id", "eip_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.eip_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_id", "eip_list.0.eip_name", "tf-example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.eip_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.public_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.create_time"),

					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example_by_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_name", "eip_list.0.eip_name", "tf-example"),

					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example_by_tags"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_tags", "eip_list.0.tags.test", "test"),
				),
			},
		},
	})
}

const testAccEipsDataSource = `
resource "tencentcloud_eip" "example" {
  name = "tf-example"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_eips" "example" {}

data "tencentcloud_eips" "example_by_id" {
  eip_id = tencentcloud_eip.example.id
}

data "tencentcloud_eips" "example_by_name" {
  eip_name = tencentcloud_eip.example.name
}

data "tencentcloud_eips" "example_by_tags" {
  tags = tencentcloud_eip.example.tags
}
`
