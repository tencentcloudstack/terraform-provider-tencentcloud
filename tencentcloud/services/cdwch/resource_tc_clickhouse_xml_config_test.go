package cdwch_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseXmlConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseXmlConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_xml_config.xml_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "instance_id", "cdwch-pcap78rz"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.0.file_name", "metrika.xml"),
					resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.0.new_conf_value"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.0.file_path", "/etc/clickhouse-server"),
				),
			},
			{
				Config: testAccClickhouseXmlConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_xml_config.xml_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "instance_id", "cdwch-pcap78rz"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.0.file_name", "metrika.xml"),
					resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.0.new_conf_value"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_xml_config.xml_config", "modify_conf_context.0.file_path", "/etc/clickhouse-server"),
				),
			},
			{
				ResourceName:      "tencentcloud_clickhouse_xml_config.xml_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClickhouseXmlConfig = `
resource "tencentcloud_clickhouse_xml_config" "xml_config" {
  instance_id = "cdwch-pcap78rz"
  modify_conf_context {
    file_name      = "metrika.xml"
    new_conf_value = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHlhbmRleD4KICAgIDxjbGlja2hvdXNlX3JlbW90ZV9zZXJ2ZXJzPgogICAgICAgIDxkZWZhdWx0X2NsdXN0ZXI+CiAgICAgICAgICAgIDxzaGFyZD4KICAgICAgICAgICAgICAgIDxpbnRlcm5hbF9yZXBsaWNhdGlvbj50cnVlPC9pbnRlcm5hbF9yZXBsaWNhdGlvbj4KICAgICAgICAgICAgICAgIDxyZXBsaWNhPgogICAgICAgICAgICAgICAgICAgIDxob3N0PjkuMC4xNi4xNjwvaG9zdD4KICAgICAgICAgICAgICAgICAgICA8cG9ydD45MDAwPC9wb3J0PgogICAgICAgICAgICAgICAgPC9yZXBsaWNhPgogICAgICAgICAgICAgICAgPHJlcGxpY2E+CiAgICAgICAgICAgICAgICAgICAgPGhvc3Q+OS4wLjE2LjI8L2hvc3Q+CiAgICAgICAgICAgICAgICAgICAgPHBvcnQ+OTAwMDwvcG9ydD4KICAgICAgICAgICAgICAgIDwvcmVwbGljYT4KICAgICAgICAgICAgPC9zaGFyZD4KICAgICAgICA8L2RlZmF1bHRfY2x1c3Rlcj4KICAgIDwvY2xpY2tob3VzZV9yZW1vdGVfc2VydmVycz4KICAgIDx6b29rZWVwZXItc2VydmVycz4KICAgICAgICA8bm9kZT4KICAgICAgICAgICAgPGhvc3Q+OS4wLjE2LjY8L2hvc3Q+CiAgICAgICAgICAgIDxwb3J0PjIxODE8L3BvcnQ+CiAgICAgICAgPC9ub2RlPgogICAgICAgIDxub2RlPgogICAgICAgICAgICA8aG9zdD45LjAuMTYuMTM8L2hvc3Q+CiAgICAgICAgICAgIDxwb3J0PjIxODE8L3BvcnQ+CiAgICAgICAgPC9ub2RlPgogICAgPC96b29rZWVwZXItc2VydmVycz4KPC95YW5kZXg+Cg=="
    file_path      = "/etc/clickhouse-server"
  }
}
`

const testAccClickhouseXmlConfigUpdate = `
resource "tencentcloud_clickhouse_xml_config" "xml_config" {
  instance_id = "cdwch-pcap78rz"
  modify_conf_context {
    file_name      = "metrika.xml"
    new_conf_value = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHlhbmRleD4KICAgIDxjbGlja2hvdXNlX3JlbW90ZV9zZXJ2ZXJzPgogICAgICAgIDxkZWZhdWx0X2NsdXN0ZXI+CiAgICAgICAgICAgIDxzaGFyZD4KICAgICAgICAgICAgICAgIDxpbnRlcm5hbF9yZXBsaWNhdGlvbj50cnVlPC9pbnRlcm5hbF9yZXBsaWNhdGlvbj4KICAgICAgICAgICAgICAgIDxyZXBsaWNhPgogICAgICAgICAgICAgICAgICAgIDxob3N0PjkuMC4xNi4xNjwvaG9zdD4KICAgICAgICAgICAgICAgICAgICA8cG9ydD45MDAwPC9wb3J0PgogICAgICAgICAgICAgICAgPC9yZXBsaWNhPgogICAgICAgICAgICAgICAgPHJlcGxpY2E+CiAgICAgICAgICAgICAgICAgICAgPGhvc3Q+OS4wLjE2LjI8L2hvc3Q+CiAgICAgICAgICAgICAgICAgICAgPHBvcnQ+OTAwMDwvcG9ydD4KICAgICAgICAgICAgICAgIDwvcmVwbGljYT4KICAgICAgICAgICAgPC9zaGFyZD4KICAgICAgICA8L2RlZmF1bHRfY2x1c3Rlcj4KICAgIDwvY2xpY2tob3VzZV9yZW1vdGVfc2VydmVycz4KICAgIDx6b29rZWVwZXItc2VydmVycz4KICAgICAgICA8bm9kZT4KICAgICAgICAgICAgPGhvc3Q+OS4wLjE2LjY8L2hvc3Q+CiAgICAgICAgICAgIDxwb3J0PjIxODE8L3BvcnQ+CiAgICAgICAgPC9ub2RlPgogICAgICAgIDxub2RlPgogICAgICAgICAgICA8aG9zdD45LjAuMTYuMTM8L2hvc3Q+CiAgICAgICAgICAgIDxwb3J0PjIxODE8L3BvcnQ+CiAgICAgICAgPC9ub2RlPgogICAgICAgIDxub2RlPgogICAgICAgICAgICA8aG9zdD45LjAuMTYuMTU8L2hvc3Q+CiAgICAgICAgICAgIDxwb3J0PjIxODE8L3BvcnQ+CiAgICAgICAgPC9ub2RlPgogICAgPC96b29rZWVwZXItc2VydmVycz4KPC95YW5kZXg+Cg=="
    file_path      = "/etc/clickhouse-server"
  }
}
`
