package cvm_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmInstanceVncUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceVncUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_instance_vnc_url.instance_vnc_url"),
					resource.TestCheckResourceAttrWith("data.tencentcloud_cvm_instance_vnc_url.instance_vnc_url", "instance_vnc_url", func(value string) error {
						if value == "" {
							return fmt.Errorf("vnc url is empty")
						}
						return nil
					}),
				),
			},
		},
	})
}

const testAccCvmInstanceVncUrlDataSource = `
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    filter {
        name = "instance-family"
        values = ["SA2","SA3","SA4","SA5","S2","S3"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}

resource "tencentcloud_instance" "test_cvm" {
	image_id = data.tencentcloud_images.default.images.0.image_id
	availability_zone = "ap-guangzhou-7"
	instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    orderly_security_groups = ["sg-5275dorp"]
	instance_charge_type = "POSTPAID_BY_HOUR"
}

data "tencentcloud_cvm_instance_vnc_url" "instance_vnc_url" {
  instance_id = tencentcloud_instance.test_cvm.id
}
`
