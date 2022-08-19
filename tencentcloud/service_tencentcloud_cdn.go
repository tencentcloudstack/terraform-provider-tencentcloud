package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CdnService struct {
	client *connectivity.TencentCloudClient
}

type CdnVerifyRecordResponse struct {
	SubDomain  *string `json:"SubDomain,omitempty" name:"SubDomain"`
	Record     *string `json:"Record,omitempty" name:"Record"`
	RecordType *string `json:"RecordType,omitempty" name:"RecordType"`
	RequestId  *string `json:"RequestId,omitempty" name:"RequestId"`
}

func (me *CdnService) DescribeDomainsConfigByDomain(ctx context.Context, domain string) (domainConfig *cdn.DetailDomain, errRet error) {
	logId := getLogId(ctx)
	request := cdn.NewDescribeDomainsConfigRequest()
	request.Filters = make([]*cdn.DomainFilter, 0, 1)
	filter := &cdn.DomainFilter{
		Name:  helper.String("domain"),
		Value: []*string{&domain},
	}
	request.Filters = append(request.Filters, filter)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().DescribeDomainsConfig(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			if sdkErr.Code == CDN_HOST_NOT_FOUND {
				return
			}
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if len(response.Response.Domains) > 0 {
		domainConfig = response.Response.Domains[0]
	}
	return
}

func (me *CdnService) DeleteDomain(ctx context.Context, domain string) error {
	logId := getLogId(ctx)
	request := cdn.NewDeleteCdnDomainRequest()
	request.Domain = &domain

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCdnClient().DeleteCdnDomain(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CdnService) StopDomain(ctx context.Context, domain string) error {
	logId := getLogId(ctx)
	request := cdn.NewStopCdnDomainRequest()
	request.Domain = &domain

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCdnClient().StopCdnDomain(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CdnService) StartDomain(ctx context.Context, domain string) error {
	logId := getLogId(ctx)
	request := cdn.NewStartCdnDomainRequest()
	request.Domain = &domain

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCdnClient().StartCdnDomain(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CdnService) DescribeDomainsConfigByFilters(ctx context.Context,
	filterMap map[string]interface{}) (domainConfig []*cdn.DetailDomain, errRet error) {

	var (
		logId         = getLogId(ctx)
		request       = cdn.NewDescribeDomainsConfigRequest()
		err           error
		offset, limit int64 = 0, 100
		response      *cdn.DescribeDomainsConfigResponse
	)

	request.Filters = make([]*cdn.DomainFilter, 0, len(filterMap))

	for k, v := range filterMap {
		value := v.(string)
		filter := &cdn.DomainFilter{
			Name:  helper.String(k),
			Value: []*string{&value},
		}
		request.Filters = append(request.Filters, filter)
	}

	request.Limit = &limit
	request.Offset = &offset

	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseCdnClient().DescribeDomainsConfig(request)

			if err != nil {
				if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
					if sdkErr.Code == CDN_HOST_NOT_FOUND {
						return nil
					}
				}
				return retryError(err, InternalError)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			errRet = err
			return
		}

		domainConfig = append(domainConfig, response.Response.Domains...)
		if len(response.Response.Domains) < int(limit) {
			break
		}
		offset += limit
	}
	return
}

func (me *CdnService) VerifyDomainRecord(ctx context.Context, domain string) (result bool, errRet error) {
	logId := getLogId(ctx)
	request := cdn.NewVerifyDomainRecordRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().VerifyDomainRecord(request)

	if err != nil {
		return
	}

	if response.Response.Result != nil {
		result = *response.Response.Result
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdnService) CreateVerifyRecord(ctx context.Context, domain string) (resp *cdn.CreateVerifyRecordResponseParams, errRet error) {
	logId := getLogId(ctx)
	request := cdn.NewCreateVerifyRecordRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.Domain = &domain

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().CreateVerifyRecord(request)

	if err != nil {
		errRet = err
		return
	}

	resp = response.Response

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdnService) DescribePurgeTasks(ctx context.Context, request *cdn.DescribePurgeTasksRequest) (task []*cdn.PurgeTask, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().DescribePurgeTasks(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response != nil {
		task = response.Response.PurgeLogs
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdnService) DescribePushTasks(ctx context.Context, request *cdn.DescribePushTasksRequest) (task []*cdn.PushTask, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().DescribePushTasks(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response != nil {
		task = response.Response.PushLogs
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdnService) PurgeUrlsCache(ctx context.Context, request *cdn.PurgeUrlsCacheRequest) (taskId string, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().PurgeUrlsCache(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response.TaskId != nil {
		taskId = *response.Response.TaskId
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdnService) PushUrlsCache(ctx context.Context, request *cdn.PushUrlsCacheRequest) (taskId string, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().PushUrlsCache(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response.TaskId != nil {
		taskId = *response.Response.TaskId
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func GetUrlsHash(urls []string) string {
	return hashcode.Strings(urls)
}
