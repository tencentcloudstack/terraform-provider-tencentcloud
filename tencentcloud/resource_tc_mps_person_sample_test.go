package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMpsPersonSampleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsPersonSample,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_person_sample.person_sample", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_person_sample.person_sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsPersonSample = `

resource "tencentcloud_mps_person_sample" "person_sample" {
  name          = "test"
  usages        = [
    "Review.Face"
  ]
  description   = "test"
  face_contents = [
    filebase64("./person.png")
  ]
}

`
