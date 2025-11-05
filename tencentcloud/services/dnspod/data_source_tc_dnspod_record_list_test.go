package dnspod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodRecordListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordListDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_list.record_list")),
			},
		},
	})
}

func TestAccTencentCloudDnspodRecordListDataSource_subDomains(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordListDataSource_subDomains,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_list.subdomains"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnspod_record_list.subdomains", "record_list.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnspod_record_list.subdomains", "instance_list.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudDnspodRecordListDataSource_withoutSubDomain(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordListDataSource_withoutSubDomain,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_list.subdomains"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dnspod_record_list.subdomains", "record_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dnspod_record_list.subdomains", "instance_list.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudDnspodRecordListDataSource_filterAtNS(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordListDataSource_filterAtNSTrue,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_list.subdomains"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnspod_record_list.subdomains", "record_list.#", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnspod_record_list.subdomains", "instance_list.#", "0"),
				),
			},
			{
				Config: testAccDnspodRecordListDataSource_filterAtNSFalse,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_list.subdomains"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnspod_record_list.subdomains", "record_list.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnspod_record_list.subdomains", "instance_list.#", "2"),
				),
			},
		},
	})
}

const testAccDnspodRecordListDataSource = `

data "tencentcloud_dnspod_record_list" "record_list" {
  domain = "tencentcloud-terraform-provider.cn"
  # domain_id = 123
  # sub_domain = "www"
  record_type = ["A", "NS", "CNAME", "NS", "AAAA"]
  # record_line = [""]
  group_id = []
  keyword = ""
  sort_field = "UPDATED_ON"
  sort_type = "DESC"
  record_value = "bicycle.dnspod.net"
  record_status = ["ENABLE"]
  weight_begin = 0
  weight_end = 100
  mx_begin = 0
  mx_end = 10
  ttl_begin = 1
  ttl_end = 864000
  updated_at_begin = "2021-09-07"
  updated_at_end = "2023-12-07"
  remark = ""
  is_exact_sub_domain = true
  # project_id = -1
}

`

const testAccDnspodRecordListDataSource_subDomains = `
data "tencentcloud_dnspod_record_list" "subdomains" {
  domain              = "mikatong.xyz"
  is_exact_sub_domain = true
  sub_domains          = ["tes1029","tes103"]
}
`

const testAccDnspodRecordListDataSource_withoutSubDomain = `
data "tencentcloud_dnspod_record_list" "subdomains" {
  domain              = "mikatong.xyz"
  is_exact_sub_domain = true
}
`

const testAccDnspodRecordListDataSource_filterAtNSTrue = `
data "tencentcloud_dnspod_record_list" "subdomains" {
  domain              = "tencentcloud-terraform-provider.cn"
  sub_domains = ["@"]
  record_type = ["NS"]
  filter_at_ns = true
}
`
const testAccDnspodRecordListDataSource_filterAtNSFalse = `
data "tencentcloud_dnspod_record_list" "subdomains" {
  domain              = "tencentcloud-terraform-provider.cn"
  sub_domains = ["@"]
  record_type = ["NS"]
  filter_at_ns = false
}
`
