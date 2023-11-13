package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmProgramFpgaImageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
  instance_id = "ins-r8hr2upy"
  f_p_g_a_url = "fpga-test-123456.cos.ap-guangzhou.myqcloud.com/test.xclbin"
  d_b_d_fs = 
  dry_run = false
}

`
