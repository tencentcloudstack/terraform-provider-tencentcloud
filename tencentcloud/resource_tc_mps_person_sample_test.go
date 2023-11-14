package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsPersonSampleResource_basic(t *testing.T) {
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
  name = &lt;nil&gt;
  usages = &lt;nil&gt;
  description = &lt;nil&gt;
  face_contents = &lt;nil&gt;
}

`
