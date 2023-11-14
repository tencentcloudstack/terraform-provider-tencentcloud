package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsPersonSamplesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsPersonSamplesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_person_samples.person_samples")),
			},
		},
	})
}

const testAccMpsPersonSamplesDataSource = `

data "tencentcloud_mps_person_samples" "person_samples" {
  type = &lt;nil&gt;
  person_ids = &lt;nil&gt;
  names = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  person_set {
		person_id = &lt;nil&gt;
		name = &lt;nil&gt;
		description = &lt;nil&gt;
		face_info_set {
			face_id = &lt;nil&gt;
			url = &lt;nil&gt;
		}
		tag_set = &lt;nil&gt;
		usage_set = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;

  }
}

`
