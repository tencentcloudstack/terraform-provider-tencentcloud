package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	sslv20191205 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func resourceTencentCloudSslCheckCertificateDomainVerificationOperationCreatePreHandleResponse0(ctx context.Context, resp *sslv20191205.CheckCertificateDomainVerificationResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)

	response := resp.Response.VerificationResults
	tmpList := make([]map[string]interface{}, 0, len(response))
	for _, item := range response {
		tmpObject := make(map[string]interface{})
		if item.Domain != nil {
			tmpObject["domain"] = item.Domain
		}

		if item.VerifyType != nil {
			tmpObject["verify_type"] = item.VerifyType
		}

		if item.LocalCheck != nil {
			tmpObject["local_check"] = item.LocalCheck
		}

		if item.CaCheck != nil {
			tmpObject["ca_check"] = item.CaCheck
		}

		if item.LocalCheckFailReason != nil {
			tmpObject["local_check_fail_reason"] = item.LocalCheckFailReason
		}

		if item.CheckValue != nil {
			tmpValueList := make([]string, 0, len(item.CheckValue))
			for _, v := range item.CheckValue {
				tmpValueList = append(tmpValueList, *v)
			}

			tmpObject["check_value"] = tmpValueList
		}

		if item.Frequently != nil {
			tmpObject["frequently"] = item.Frequently
		}

		if item.Issued != nil {
			tmpObject["issued"] = item.Issued
		}

		tmpList = append(tmpList, tmpObject)
	}

	_ = d.Set("verification_results", tmpList)

	return nil
}
