package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_instance.migrate_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_instance.migrate_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateInstance = `

resource "tencentcloud_dts_migrate_instance" "migrate_instance" {
  src_database_type = &lt;nil&gt;
  dst_database_type = &lt;nil&gt;
  src_region = &lt;nil&gt;
  dst_region = &lt;nil&gt;
  instance_class = &lt;nil&gt;
  count = &lt;nil&gt;
  job_name = &lt;nil&gt;
  tags {
		tag_key = &lt;nil&gt;
		tag_value = &lt;nil&gt;

  }
              complete_mode = &lt;nil&gt;
}

`
