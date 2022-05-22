package tencentcloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudDnspodRecord(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnspodRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDnspodRecord,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodRecordExists("tencentcloud_dnspod_record.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "domain", "terraform.com"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "value", "1.2.3.9"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "sub_domain", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "status", "ENABLE"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record.demo", "record_line", "默认"),
				),
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
		items := strings.Split(rs.Primary.ID, FILED_SP)
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

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, e := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn.UseDnsPodClient().DescribeRecord(request)
			if e != nil {
				return retryError(e)
			}

			if response.Response.RecordInfo != nil {
				return nil
			}
			return retryError(fmt.Errorf("Dnspod record is null!"))
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

		items := strings.Split(rs.Primary.ID, FILED_SP)
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

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, e := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn.UseDnsPodClient().DescribeRecord(request)
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "InvalidParameter.RecordIdInvalid" {
					return nil
				}
				return retryError(e)
			}
			if response.Response.RecordInfo == nil {
				return nil
			}
			return retryError(fmt.Errorf("Dnspod record still exist!"))
		})
		if err != nil {
			return err
		}
	}
	return nil
}

const testAccTencentCloudDnspodRecord = `
resource "tencentcloud_dnspod_record" "demo" {
  domain="terraform.com"
  record_type="A"
  record_line="默认"
  value="1.2.3.9"
  sub_domain="demo"
}
`
