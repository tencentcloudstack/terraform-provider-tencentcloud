package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterMaintenanceWindowAndExclusionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterMaintenanceWindowAndExclusion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.example", "id")),
			},
			{
				Config: testAccKubernetesClusterMaintenanceWindowAndExclusionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.example", "id")),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKubernetesClusterMaintenanceWindowAndExclusion = `
resource "tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion" "example" {
  cluster_id       = "cls-d2cit6no"
  maintenance_time = "01:00:00"
  duration         = 4
  day_of_week      = ["MO", "TU", "WE", "TH", "FR"]
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

const testAccKubernetesClusterMaintenanceWindowAndExclusionUpdate = `
resource "tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion" "example" {
  cluster_id       = "cls-d2cit6no"
  maintenance_time = "02:00:00"
  duration         = 2
  day_of_week      = ["MO", "TU", "WE", "TH", "FR"]
  exclusions {
    name     = "name1"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }
}
`
