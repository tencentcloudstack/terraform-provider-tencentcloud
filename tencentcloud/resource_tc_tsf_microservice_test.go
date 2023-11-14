package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfMicroserviceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMicroservice,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_microservice.microservice", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_microservice.microservice",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfMicroservice = `

resource "tencentcloud_tsf_microservice" "microservice" {
  namespace_id = ""
  microservice_name = ""
  microservice_desc = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
