package apm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApmInstanceResource_basic -v
func TestAccTencentCloudApmInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApmInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_instance.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "description", "desc."),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "trace_duration", "15"),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "span_daily_counters", "0"),
				),
			},
			{
				Config: testAccApmInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_instance.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "description", "desc update."),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "trace_duration", "15"),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.example", "span_daily_counters", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_apm_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApmInstance = `
resource "tencentcloud_apm_instance" "example" {
  name                = "tf-example"
  description         = "desc."
  trace_duration      = 15
  span_daily_counters = 0
  tags = {
    createdBy = "terraform"
  }
}
`

const testAccApmInstanceUpdate = `
resource "tencentcloud_apm_instance" "example" {
  name                = "tf-example-update"
  description         = "desc update."
  trace_duration      = 15
  span_daily_counters = 0
  tags = {
    createdBy = "terraform"
  }
}
`
