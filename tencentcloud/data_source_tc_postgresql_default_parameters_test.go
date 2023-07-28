package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlDefaultParametersDataSource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDefaultParametersDataSource,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_default_parameters.default_parameters"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_default_parameters.default_parameters", "db_major_version", "13"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "db_engine"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.param_value_type"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.default_value"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.current_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.max"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.min"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.enum_value.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.param_description_ch"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.param_description_en"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.need_reboot"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.classification_cn"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.classification_en"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_related"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.advanced"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.last_modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.standby_related"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set.0.name"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set.0.db_kernel_version"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set.0.value"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set.0.unit"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set.0.max"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set.0.min"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.version_relation_set.0.enum_value.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set.0.name"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set.0.memory"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set.0.value"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set.0.unit"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set.0.max"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set.0.min"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_default_parameters.default_parameters", "param_info_set.0.spec_relation_set.0.enum_value.#"),
				),
			},
		},
	})
}

const testAccPostgresqlDefaultParametersDataSource = CommonPresetPGSQL + `

data "tencentcloud_postgresql_default_parameters" "default_parameters" {
  db_major_version = "13"
  db_engine = "postgresql"
}

`
