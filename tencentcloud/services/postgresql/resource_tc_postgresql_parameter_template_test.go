package postgresql_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_postgresql_parameter_template
	resource.AddTestSweepers("tencentcloud_postgresql_parameter_template", &resource.Sweeper{
		Name: "tencentcloud_postgresql_parameter_template",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			postgresqlService := svcpostgresql.NewPostgresqlService(client)

			temps, err := postgresqlService.DescribePostgresqlParameterTemplatesByFilter(ctx, map[string]interface{}{})
			if err != nil {
				return err
			}

			for _, v := range temps {

				name := *v.TemplateName
				id := *v.TemplateId

				if strings.HasPrefix(name, tcacctest.KeepResource) || strings.HasPrefix(name, tcacctest.DefaultResource) {
					continue
				}

				// delete
				err = postgresqlService.DeletePostgresqlParameterTemplateById(ctx, id)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudPostgresqlParameterTemplateResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckPostgresqlParameterTemplateDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlParameterTemplate,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlParameterTemplateExists("tencentcloud_postgresql_parameter_template.parameter_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_parameter_template.parameter_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "template_name", "tf_test_pg_temp"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "db_major_version", "13"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "db_engine", "postgresql"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "template_description", "For_tf_test"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.#", "2"),
					// resource.TestCheckTypeSetElemAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.*",
					// 	map[string]string{
					// 		"name":           "lc_time",
					// 		"expected_value": "POSIX",
					// 	}),
					// resource.TestCheckTypeSetElemAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.*",
					// 	map[string]string{
					// 		"name":           "timezone",
					// 		"expected_value": "PRC",
					// 	}),
				),
			},
			{
				Config: testAccPostgresqlParameterTemplate_update_desc,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlParameterTemplateExists("tencentcloud_postgresql_parameter_template.parameter_template"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "template_description", "For_tf_test_desc_changed"),
				),
			},
			{
				Config: testAccPostgresqlParameterTemplate_update_name,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlParameterTemplateExists("tencentcloud_postgresql_parameter_template.parameter_template"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "template_name", "tf_test_pg_temp_name_changed"),
				),
			},
			{
				Config: testAccPostgresqlParameterTemplate_update_multiple_attr,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlParameterTemplateExists("tencentcloud_postgresql_parameter_template.parameter_template"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.#", "2"),
					// resource.TestCheckTypeSetElemAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.*",
					// 	map[string]string{
					// 		"name":           "timezone",
					// 		"expected_value": "UTC",
					// 	}),
					// resource.TestCheckTypeSetElemAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.*",
					// 	map[string]string{
					// 		"name":           "lock_timeout",
					// 		"expected_value": "123",
					// 	}),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "delete_param_set.#", "1"),
					// resource.TestCheckTypeSetElemAttr("tencentcloud_postgresql_parameter_template.parameter_template", "delete_param_set.*", "lc_time"),
				),
			},
			{
				ResourceName:            "tencentcloud_postgresql_parameter_template.parameter_template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_param_set"},
			},
		},
	})
}

func testAccCheckPostgresqlParameterTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_postgresql_parameter_template" {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcpostgresql.NewPostgresqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		ret, _ := service.DescribePostgresqlParameterTemplateById(ctx, rs.Primary.ID)

		if ret != nil && ret.TemplateId != nil {
			return fmt.Errorf("delete postgresql parameter template %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPostgresqlParameterTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcpostgresql.NewPostgresqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		ret, err := service.DescribePostgresqlParameterTemplateById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("postgresql parameter template %s not found", rs.Primary.ID)
		}
		return nil
	}
}

const testAccPostgresqlParameterTemplate = `

resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "tf_test_pg_temp"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "For_tf_test"

  modify_param_entry_set {
	name = "lc_time"
	expected_value = "POSIX"
  }
  modify_param_entry_set {
	name = "timezone"
	expected_value = "PRC"
  }
}

`

const testAccPostgresqlParameterTemplate_update_desc = `

resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "tf_test_pg_temp"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "For_tf_test_desc_changed"

  modify_param_entry_set {
	name = "lc_time"
	expected_value = "POSIX"
  }
  modify_param_entry_set {
	name = "timezone"
	expected_value = "PRC"
  }
}

`

const testAccPostgresqlParameterTemplate_update_name = `

resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "tf_test_pg_temp_name_changed"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "For_tf_test_desc_changed"

  modify_param_entry_set {
	name = "lc_time"
	expected_value = "POSIX"
  }
  modify_param_entry_set {
	name = "timezone"
	expected_value = "PRC"
  }
}

`

const testAccPostgresqlParameterTemplate_update_multiple_attr = `

resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "tf_test_pg_temp_name_multi_changed"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "For_tf_test_desc_multi_changed"

  modify_param_entry_set {
	name = "timezone"
	expected_value = "UTC"
  }
  modify_param_entry_set {
	name = "lock_timeout"
	expected_value = "123"
  }

  delete_param_set = ["lc_time"]
}

`
