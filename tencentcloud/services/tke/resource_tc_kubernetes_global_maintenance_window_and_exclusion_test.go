package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesGlobalMaintenanceWindowAndExclusionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesGlobalMaintenanceWindowAndExclusion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_global_maintenance_window_and_exclusion.example", "id"),
				),
			},
			{
				Config: testAccKubernetesGlobalMaintenanceWindowAndExclusionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_global_maintenance_window_and_exclusion.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_global_maintenance_window_and_exclusion.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKubernetesGlobalMaintenanceWindowAndExclusion = `
resource "tencentcloud_kubernetes_global_maintenance_window_and_exclusion" "example" {
  maintenance_time = "02:00:00"
  duration         = 4
  day_of_week      = ["MO", "TU", "WE", "TH", "FR"]
  target_regions   = ["ap-guangzhou"]
  exclusions {
    name     = "name1"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }

  exclusions {
    name     = "name2"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }

  exclusions {
    name     = "name3"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }
}
`

const testAccKubernetesGlobalMaintenanceWindowAndExclusionUpdate = `
resource "tencentcloud_kubernetes_global_maintenance_window_and_exclusion" "example" {
  maintenance_time = "04:00:00"
  duration         = 2
  day_of_week      = ["MO", "TU", "WE", "TH", "FR"]
  target_regions   = ["ap-guangzhou"]
  exclusions {
    name     = "name1"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }
}
`
