package as_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudAsStartInstanceRefreshResource_basic -v
func TestAccTencentCloudAsStartInstanceRefreshResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsStartInstanceRefresh,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_as_start_instance_refresh.example", "id"),
				),
			},
		},
	})
}

const testAccAsStartInstanceRefresh = `
resource "tencentcloud_as_start_instance_refresh" "example" {
  auto_scaling_group_id = "asg-2l55y7u7"
  refresh_mode          = "ROLLING_UPDATE_RESET"
  refresh_settings {
    check_instance_target_health = false
    rolling_update_settings {
      batch_number = 1
      batch_pause  = "AUTOMATIC"
    }
  }
}
`
