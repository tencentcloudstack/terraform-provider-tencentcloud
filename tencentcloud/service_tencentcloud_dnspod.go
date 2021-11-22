package tencentcloud

import (
	"context"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

type DnspodService struct {
	client *connectivity.TencentCloudClient
}

// ////////api
func (me *DnspodService) ModifyDnsPodDomainStatus(ctx context.Context, domain string, status string) (errRet error) {
	logId := getLogId(ctx)
	request := dnspod.NewModifyDomainStatusRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.Domain = helper.String(domain)
	request.Status = &status

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseDnsPodClient().ModifyDomainStatus(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify dnspod domain status failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *DnspodService) ModifyDnsPodDomainRemark(ctx context.Context, domain string, remark string) (errRet error) {
	logId := getLogId(ctx)
	request := dnspod.NewModifyDomainRemarkRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.Domain = helper.String(domain)
	request.Remark = &remark

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseDnsPodClient().ModifyDomainRemark(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify dnspod domain remark failed, reason: %v", logId, err)
		return err
	}
	return
}
