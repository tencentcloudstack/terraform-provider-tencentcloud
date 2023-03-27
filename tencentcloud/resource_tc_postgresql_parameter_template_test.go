package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_postgresql_parameter_template
	resource.AddTestSweepers("tencentcloud_postgresql_parameter_template", &resource.Sweeper{
		Name: "tencentcloud_postgresql_parameter_template",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			postgresqlService := PostgresqlService{client: client}

			temps, err := postgresqlService.DescribePostgresqlParameterTemplatesByFilter(ctx, map[string]interface{}{})
			if err != nil {
				return err
			}

			for _, v := range temps {

				name := *v.TemplateName
				id := *v.TemplateId

				if strings.HasPrefix(name, keepResource) || strings.HasPrefix(name, defaultResource) {
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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckPostgresqlParameterTemplateDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlParameterTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlParameterTemplateExists("tencentcloud_postgresql_parameter_template.parameter_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_parameter_template.parameter_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "template_name", "tf_test_pg_temp"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "db_major_version", "13"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "db_engine", "postgresql"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.0.name", "lc_time"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.0.expected_value", "POSIX"),
				),
			},
			{
				Config: testAccPostgresqlParameterTemplate_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlParameterTemplateExists("tencentcloud_postgresql_parameter_template.parameter_template"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.0.name", "timezone"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.0.expected_value", "UTC"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.1.name", "lock_timeout"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "modify_param_entry_set.1.expected_value", "123"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameter_template.parameter_template", "delete_param_set.#", "1"),
				),
			},
			{
				ResourceName:            "tencentcloud_postgresql_parameter_template.parameter_template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_param_entry_set", "delete_param_set"},
			},
		},
	})
}

func testAccCheckPostgresqlParameterTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_postgresql_parameter_template" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

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

const testAccPostgresqlParameterTemplate_update = `

resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "tf_test_pg_temp"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "For_tf_test"

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
