package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchDiagnoseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchDiagnoseInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_diagnose_instance.diagnose_instance", "id")),
			},
		},
	})
}

const testAccElasticsearchDiagnoseInstance = `
resource "tencentcloud_elasticsearch_diagnose_instance" "diagnose_instance" {
	instance_id = "es-5wn36he6"
	diagnose_jobs = ["cluster_health"]
	diagnose_indices = "*"
  }
`
