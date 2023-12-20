package ci_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localci "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ci"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCIGuetzli_basic(t *testing.T) {
	resourceName := "tencentcloud_ci_guetzli.basic"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudCIGuetzliConfig_basic("on"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCIGuetzliOn(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "on"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTencentCloudCIGuetzliConfig_basic("off"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "off"),
				),
			},
		},
	})
}

func testAccCheckCIGuetzliOn(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		bucket := rs.Primary.ID
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.GetCiGuetzliById(context.Background(), bucket)

		if err != nil {
			return fmt.Errorf("error getting COS bucket Guetzli set (%s): %w", bucket, err)
		}
		if res.GuetzliStatus != "on" {
			return fmt.Errorf("error setting COS bucket Guetzli (%s): status(%s)", bucket, res.GuetzliStatus)
		}
		return nil
	}
}

func testAccTencentCloudCIGuetzliConfig_basic(status string) string {
	return fmt.Sprintf(`
variable "bucket" {
	default = %[1]q
	}

resource "tencentcloud_ci_guetzli" "basic" {
	bucket = var.bucket
	status = %[2]q
}
`, tcacctest.DefaultCiBucket, status)
}
