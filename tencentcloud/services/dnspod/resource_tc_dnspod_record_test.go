package dnspod_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudDnspodRecordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDnspodRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDnspodRecord,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodRecordExists("tencentcloud_dnspod_record.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "value", "1.2.3.9"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "sub_domain", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "status", "ENABLE"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_line", "默认"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "remark", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "weight", "100"),
				),
			},
			{
				Config: testAccTencentCloudDnspodRecordRemarkUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodRecordExists("tencentcloud_dnspod_record.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "value", "1.2.3.9"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "sub_domain", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "status", "ENABLE"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_line", "默认"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "remark", "terraform-test1"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "weight", "100"),
				),
			},
			{
				Config: testAccTencentCloudDnspodRecordValueUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodRecordExists("tencentcloud_dnspod_record.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "value", "1.2.3.10"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "sub_domain", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "status", "ENABLE"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_line", "默认"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "remark", "terraform-test1"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "weight", "100"),
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_record.demo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudDnspodRecordResource_MX(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDnspodRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDnspodRecordMXValueWithOutDot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodRecordExists("tencentcloud_dnspod_record.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_type", "MX"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_line", "默认"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "value", "1.2.3.9"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "sub_domain", "@"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "mx", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "ttl", "86400"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "status", "ENABLE"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "remark", "terraform-test"),
				),
			},
			{
				Config: testAccTencentCloudDnspodRecordMXValueWithDot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodRecordExists("tencentcloud_dnspod_record.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_type", "MX"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_line", "默认"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "value", "1.2.3.9."),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "sub_domain", "@"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "mx", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "ttl", "86400"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "status", "ENABLE"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "remark", "terraform-test"),
				),
			},
			{
				Config: testAccTencentCloudDnspodRecordMxUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodRecordExists("tencentcloud_dnspod_record.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_type", "MX"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_line", "默认"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "value", "1.2.3.19."),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "sub_domain", "@"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "mx", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "ttl", "86400"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "status", "ENABLE"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "remark", "terraform-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_record.demo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDnspodRecordExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("dnspod record %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dnspod record id is not set")
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) < 2 {
			return nil
		}
		request := dnspod.NewDescribeRecordRequest()
		request.Domain = helper.String(items[0])
		recordId, err := strconv.Atoi(items[1])
		if err != nil {
			return err
		}
		request.RecordId = helper.IntUint64(recordId)

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			response, e := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DescribeRecord(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if response.Response.RecordInfo != nil {
				return nil
			}
			return tccommon.RetryError(fmt.Errorf("Dnspod record is null!"))
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckDnspodRecordDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dnspod_record" {
			continue
		}

		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) < 2 {
			return nil
		}
		request := dnspod.NewDescribeRecordRequest()
		request.Domain = helper.String(items[0])
		recordId, err := strconv.Atoi(items[1])
		if err != nil {
			return err
		}
		request.RecordId = helper.IntUint64(recordId)

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			response, e := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DescribeRecord(request)
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "InvalidParameter.RecordIdInvalid" {
					return nil
				}
				return tccommon.RetryError(e)
			}
			if response.Response.RecordInfo == nil {
				return nil
			}
			return tccommon.RetryError(fmt.Errorf("Dnspod record still exist!"))
		})
		if err != nil {
			return err
		}
	}
	return nil
}

const testAccTencentCloudDnspodRecord = `
resource "tencentcloud_dnspod_record" "demo" {
	domain="iac-tf.cloud"
	record_type="A"
	record_line="默认"
	value="1.2.3.9"
	sub_domain="demo"
	remark="terraform-test"
	weight=100
}
`
const testAccTencentCloudDnspodRecordRemarkUp = `
resource "tencentcloud_dnspod_record" "demo" {
	domain="iac-tf.cloud"
	record_type="A"
	record_line="默认"
	value="1.2.3.9"
	sub_domain="demo"
	remark="terraform-test1"
	weight=100
}
`
const testAccTencentCloudDnspodRecordValueUpdate = `
resource "tencentcloud_dnspod_record" "demo" {
	domain="iac-tf.cloud"
	record_type="A"
	record_line="默认"
	value="1.2.3.10"
	sub_domain="demo"
	remark="terraform-test1"
	weight=100
}
`

const testAccTencentCloudDnspodRecordMXValueWithOutDot = `
resource "tencentcloud_dnspod_record" "demo" {
  domain      = "iac-tf.cloud"
  record_type = "MX"
  record_line = "默认"
  value       = "1.2.3.9"
  sub_domain  = "@"
  mx          = 10
  ttl         = 86400
  status      = "ENABLE"
  remark      = "terraform-test"
}
`
const testAccTencentCloudDnspodRecordMXValueWithDot = `
resource "tencentcloud_dnspod_record" "demo" {
  domain      = "iac-tf.cloud"
  record_type = "MX"
  record_line = "默认"
  value       = "1.2.3.9."
  sub_domain  = "@"
  mx          = 10
  ttl         = 86400
  status      = "ENABLE"
  remark      = "terraform-test"
}
`

const testAccTencentCloudDnspodRecordMxUpdate = `
resource "tencentcloud_dnspod_record" "demo" {
  domain      = "iac-tf.cloud"
  record_type = "MX"
  record_line = "默认"
  value       = "1.2.3.19."
  sub_domain  = "@"
  mx          = 10
  ttl         = 86400
  status      = "ENABLE"
  remark      = "terraform-test"
}
`
