package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseInstanceVncUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceVncUrlDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_instance_vnc_url.instance_vnc_url")),
			},
		},
	})
}

const testAccLighthouseInstanceVncUrlDataSource = `

data "tencentcloud_lighthouse_instance_vnc_url" "instance_vnc_url" {
  instance_id = "lhins-123456"
}

`
