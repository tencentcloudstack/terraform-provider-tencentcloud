package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AuditService struct {
	client *connectivity.TencentCloudClient
}

func (me *AuditService) DescribeAuditById(ctx context.Context, name string) (auditInfo *audit.DescribeAuditResponse, has bool, errRet error) {
	logId := getLogId(ctx)
	request := audit.NewDescribeAuditRequest()
	request.AuditName = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *audit.DescribeAuditResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAuditClient().DescribeAudit(request)
		if e != nil {
			log.Printf("[CRITAL]%s %s fail, reason:%s\n", logId, request.GetAction(), e.Error())
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound.AuditNotExist") {
					return nil
				}
			}
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil && *response.Response.AuditName != "" {
		has = true
		auditInfo = response
	}
	return
}

func (me *AuditService) DescribeAuditCosRegions(ctx context.Context) (regions []*audit.CosRegionInfo, errRet error) {
	logId := getLogId(ctx)
	request := audit.NewListCosEnableRegionRequest()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAuditClient().ListCosEnableRegion(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	regions = response.Response.EnableRegions
	return
}

func (me *AuditService) DescribeAuditCmqRegions(ctx context.Context) (regions []*audit.CmqRegionInfo, errRet error) {
	logId := getLogId(ctx)
	request := audit.NewListCmqEnableRegionRequest()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAuditClient().ListCmqEnableRegion(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	regions = response.Response.EnableRegions
	return
}

func (me *AuditService) DescribeKeyAlias(ctx context.Context, region string) (keyMetadatas []*audit.KeyMetadata, errRet error) {
	logId := getLogId(ctx)
	request := audit.NewListKeyAliasByRegionRequest()
	request.KmsRegion = &region
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAuditClient().ListKeyAliasByRegion(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	keyMetadatas = response.Response.KeyMetadatas
	return
}

func (me *AuditService) DescribeAuditTrackById(ctx context.Context, trackId string) (track *audit.DescribeAuditTrackResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = audit.NewDescribeAuditTrackRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.TrackId = helper.StrToUint64Point(trackId)

	response, err := me.client.UseAuditClient().DescribeAuditTrack(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	track = response.Response
	return
}

func (me *AuditService) DeleteAuditTrackById(ctx context.Context, trackId string) (errRet error) {
	logId := getLogId(ctx)

	request := audit.NewDeleteAuditTrackRequest()

	request.TrackId = helper.StrToUint64Point(trackId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAuditClient().DeleteAuditTrack(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
