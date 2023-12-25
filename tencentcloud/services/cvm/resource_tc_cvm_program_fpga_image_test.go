package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCvmProgramFpgaImageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmProgramFpgaImage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_program_fpga_image.program_fpga_image", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_program_fpga_image.program_fpga_image",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmProgramFpgaImage = `

resource "tencentcloud_cvm_program_fpga_image" "program_fpga_image" {
  instance_id = "ins-xxxxxx"
  fpga_url = ""
  dbd_fs = ""
}

`
