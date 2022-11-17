package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

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

func (me *DnspodService) DescribeDomain(ctx context.Context, domain string) (ret *dnspod.DescribeDomainResponse, errRet error) {

	logId := getLogId(ctx)
	request := dnspod.NewDescribeDomainRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.Domain = helper.String(domain)

	response, err := me.client.UseDnsPodClient().DescribeDomain(request)

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return response, nil
}

func (me *DnspodService) DeleteDomain(ctx context.Context, domain string) (errRet error) {

	logId := getLogId(ctx)
	request := dnspod.NewDeleteDomainRequest()
	ratelimit.Check(request.GetAction())
	request.Domain = helper.String(domain)

	response, err := me.client.UseDnsPodClient().DeleteDomain(request)

	defer func() {
		if errRet != nil {
			responseStr := ""
			if response != nil {
				responseStr = response.ToJsonString()
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
				logId,
				request.GetAction(),
				request.ToJsonString(),
				responseStr,
				errRet.Error())
		}
	}()

	errRet = err
	return
}

func (me *DnspodService) DescribeRecordList(ctx context.Context, request *dnspod.DescribeRecordListRequest) (list []*dnspod.RecordListItem, info *dnspod.RecordCountInfo, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDnsPodClient().DescribeRecordList(request)

	if err != nil {
		errRet = err
		return
	}

	list = response.Response.RecordList
	info = response.Response.RecordCountInfo

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
