package captcha

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	captchaintl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CaptchaService struct {
	client *connectivity.TencentCloudClient
}

func (me *CaptchaService) DescribeCaptchaInfoInternationalById(ctx context.Context, captchaAppId string) (ret *captchaintl.DescribeCaptchaConsoleSubDataInternational, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := captchaintl.NewDescribeCaptchaInfoListInternationalRequest()
	response := captchaintl.NewDescribeCaptchaInfoListInternationalResponse()
	request.CaptchaAppId = helper.String(captchaAppId)
	request.PageIndex = helper.IntInt64(0)
	request.PageSize = helper.IntInt64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCaptchaClient().DescribeCaptchaInfoListInternational(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe captcha info list international failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if response.Response.Data == nil || len(response.Response.Data.DataList) < 1 {
		return
	}

	ret = response.Response.Data.DataList[0]
	return
}

func (me *CaptchaService) DescribeCaptchaIpWhiteListInternationalById(ctx context.Context, captchaAppid, id int64) (ret *captchaintl.DescribeCaptchaWhiteListItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := captchaintl.NewDescribeIpWhiteListInternationalRequest()
	response := captchaintl.NewDescribeIpWhiteListInternationalResponse()
	request.CaptchaAppid = helper.Int64(captchaAppid)
	request.PageIndex = helper.IntInt64(0)
	request.PageSize = helper.IntInt64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCaptchaClient().DescribeIpWhiteListInternational(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ip white list international failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if response.Response.Data == nil || len(response.Response.Data.DataList) < 1 {
		return
	}

	for _, item := range response.Response.Data.DataList {
		if item != nil && item.Id != nil && *item.Id == id {
			ret = item
			return
		}
	}

	return
}
