package tencentcloud

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsEventResource_basic(t *testing.T) {
	t.Parallel()
	randIns := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := randIns.Intn(1000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsEvent, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_event.event", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_event.event", "event_name"),
					resource.TestCheckResourceAttr("tencentcloud_mps_event.event", "description", "tf test mps event description"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsEvent_update, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_event.event", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_event.event", "event_name", fmt.Sprintf("tf_test_event_%d_changed", randomNum)),
					resource.TestCheckResourceAttr("tencentcloud_mps_event.event", "description", "tf test mps event description changed"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_event.event",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsEvent = `

resource "tencentcloud_mps_event" "event" {
  event_name = "tf_test_event_%d"
  description = "tf test mps event description"
}

`

const testAccMpsEvent_update = `

resource "tencentcloud_mps_event" "event" {
  event_name = "tf_test_event_%d_changed"
  description = "tf test mps event description changed"
}

`
