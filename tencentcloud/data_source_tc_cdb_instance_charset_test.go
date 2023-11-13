package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbInstanceCharsetDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbInstanceCharsetDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_instance_charset.instance_charset")),
			},
		},
	})
}

const testAccCdbInstanceCharsetDataSource = `

data "tencentcloud_cdb_instance_charset" "instance_charset" {
  instance_id = ""
  }

`
