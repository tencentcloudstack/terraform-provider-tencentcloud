package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmInstanceVncUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceVncUrlDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_instance_vnc_url.instance_vnc_url")),
			},
		},
	})
}

const testAccCvmInstanceVncUrlDataSource = `

data "tencentcloud_cvm_instance_vnc_url" "instance_vnc_url" {
  instance_id = "ins-r9hr2upy"
  }

`
