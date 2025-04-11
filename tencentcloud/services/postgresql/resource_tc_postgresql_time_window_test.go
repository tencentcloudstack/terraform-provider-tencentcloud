package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlTimeWindowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlTimeWindow,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_time_window.postgresql_time_window", "id"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_time_window.postgresql_time_window", "maintain_duration", "2"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_time_window.postgresql_time_window", "maintain_start_time", "04:00"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_time_window.postgresql_time_window", "maintain_week_days.#", "7"),
				),
			},
			{
				ResourceName:      "tencentcloud_postgresql_time_window.postgresql_time_window",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccPostgresqlTimeWindowUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_time_window.postgresql_time_window", "id"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_time_window.postgresql_time_window", "maintain_duration", "3"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_time_window.postgresql_time_window", "maintain_start_time", "05:00"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_time_window.postgresql_time_window", "maintain_week_days.#", "6"),
				),
			},
		},
	})
}

const testAccPostgresqlTimeWindow = testAccPostgresqlInstance + `

resource "tencentcloud_postgresql_time_window" "postgresql_time_window" {
    db_instance_id      = tencentcloud_postgresql_instance.test.id
    maintain_duration   = 2
    maintain_start_time = "04:00"
    maintain_week_days  = [
        "friday",
        "monday",
        "saturday",
        "sunday",
        "thursday",
        "tuesday",
        "wednesday",
    ]
}
`

const testAccPostgresqlTimeWindowUp = testAccPostgresqlInstance + `

resource "tencentcloud_postgresql_time_window" "postgresql_time_window" {
    db_instance_id      = tencentcloud_postgresql_instance.test.id
    maintain_duration   = 3
    maintain_start_time = "05:00"
    maintain_week_days  = [
        "friday",
        "monday",
        "saturday",
        "sunday",
        "thursday",
        "tuesday",
    ]
}
`
