package clb

import (
	"context"
	"fmt"
	"log"
	"math"

	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SSLService struct {
	client *connectivity.TencentCloudClient
}

func (me *SSLService) ApplyCertificate(ctx context.Context, request *ssl.ApplyCertificateRequest) (id string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSSLCertificateClient().ApplyCertificate(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response.CertificateId != nil {
		id = *response.Response.CertificateId
	} else {
		errRet = fmt.Errorf("[%s] error, no certificate id response: %s", request.GetAction(), response.ToJsonString())
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SSLService) CreateCertificate(ctx context.Context, request *ssl.CreateCertificateRequest) (certificateId, dealId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err := client.CreateCertificate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		if len(response.Response.CertificateIds) > 0 {
			certificateId = *response.Response.CertificateIds[0]
		}
		if len(response.Response.DealIds) > 0 {
			dealId = *response.Response.DealIds[0]
		}
		return
	}
	errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	return
}

func (me *SSLService) CommitCertificateInformation(ctx context.Context, request *ssl.CommitCertificateInformationRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err := client.CommitCertificateInformation(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *SSLService) DescribeCertificateDetail(ctx context.Context, request *ssl.DescribeCertificateDetailRequest) (response *ssl.DescribeCertificateDetailResponse, err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err = client.DescribeCertificateDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *SSLService) ModifyCertificateAlias(ctx context.Context, request *ssl.ModifyCertificateAliasRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.ModifyCertificateAliasResponse

	response, err = client.ModifyCertificateAlias(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) ModifyCertificateProject(ctx context.Context, request *ssl.ModifyCertificateProjectRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.ModifyCertificateProjectResponse

	response, err = client.ModifyCertificateProject(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	for _, v := range response.Response.FailCertificates {
		if *v == *request.CertificateIdList[0] {
			err = fmt.Errorf("failed to modify the project. certificateId=%s", *request.CertificateIdList[0])
			return
		}
	}
	return
}

func (me *SSLService) DeleteCertificate(ctx context.Context, request *ssl.DeleteCertificateRequest) (deleteResult bool, err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.DeleteCertificateResponse

	response, err = client.DeleteCertificate(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	deleteResult = *response.Response.DeleteResult
	return
}

func (me *SSLService) CancelCertificateOrder(ctx context.Context, request *ssl.CancelCertificateOrderRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.CancelCertificateOrderResponse

	response, err = client.CancelCertificateOrder(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) SubmitCertificateInformation(ctx context.Context, request *ssl.SubmitCertificateInformationRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.SubmitCertificateInformationResponse

	response, err = client.SubmitCertificateInformation(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) UploadConfirmLetter(ctx context.Context, request *ssl.UploadConfirmLetterRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.UploadConfirmLetterResponse

	response, err = client.UploadConfirmLetter(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) UploadCertificate(ctx context.Context, request *ssl.UploadCertificateRequest) (id string, err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.UploadCertificateResponse
	response, err = client.UploadCertificate(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return
	}

	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	if response.Response.CertificateId == nil {
		err = fmt.Errorf("api[%s] return id is nil", request.GetAction())
		log.Printf("[CRITAL]%s %v", logId, err)
		return
	}

	id = *response.Response.CertificateId
	return
}

func (me *SSLService) DescribeCertificates(ctx context.Context, request *ssl.DescribeCertificatesRequest) (certificateList []*ssl.Certificates, err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()

	offset := 0
	pageSize := 100
	certificateList = make([]*ssl.Certificates, 0)
	var response *ssl.DescribeCertificatesResponse
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err = client.DescribeCertificates(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.Certificates) == 0 {
			break
		}
		certificateList = append(certificateList, response.Response.Certificates...)
		if len(response.Response.Certificates) < pageSize {
			break
		}
		offset += pageSize
	}

	return
}

func (me *SSLService) checkCertificateType(ctx context.Context, certId string, checkType string) (bool, error) {

	//get certificate by id

	request := ssl.NewDescribeCertificateDetailRequest()
	request.CertificateId = helper.String(certId)
	certificate, err := me.DescribeCertificateDetail(ctx, request)
	if err != nil {
		return false, err
	}

	if certificate != nil && certificate.Response != nil && *certificate.Response.CertificateType == checkType {
		return true, nil
	} else {
		if certificate == nil || certificate.Response == nil || certificate.Response.CertificateId == nil {
			return false, fmt.Errorf("certificate id %s is not found", certId)
		}
		return false, nil
	}

}
func (me *SSLService) ModifyCertificateResubmit(ctx context.Context, request *ssl.ModifyCertificateResubmitRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err := client.ModifyCertificateResubmit(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil || response.Response.CertificateId == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if *response.Response.CertificateId != *request.CertificateId {
		err = fmt.Errorf("TencentCloud SDK %s eertificates are inconsistent, request[%s], response[%s]",
			request.GetAction(), *request.CertificateId, *response.Response.CertificateId)
		return
	}
	return
}
func (me *SSLService) CancelAuditCertificate(ctx context.Context, request *ssl.CancelAuditCertificateRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseSSLCertificateClient()

	response, err := client.CancelAuditCertificate(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		err = fmt.Errorf("TencentCloud SDK %s CancelAudit failed", request.GetAction())
		return err
	}

	return
}
func (me *SSLService) getCertificateStatus(ctx context.Context, certificateId string) (uint64, error) {
	describeRequest := ssl.NewDescribeCertificateDetailRequest()
	describeRequest.CertificateId = &certificateId

	describeResponse, err := me.DescribeCertificateDetail(ctx, describeRequest)
	if err != nil {
		return math.MaxUint64, err
	}
	if describeResponse == nil || describeResponse.Response == nil {
		err := fmt.Errorf("TencentCloud SDK %s return empty response", describeRequest.GetAction())
		return math.MaxUint64, err
	}
	if describeResponse.Response.Status == nil {
		err := fmt.Errorf("api[%s] certificate status is nil", describeRequest.GetAction())
		return math.MaxUint64, err
	}

	return *describeResponse.Response.Status, nil
}
func (me *SslService) DescribeSslDescribeCertificateByID(ctx context.Context, certificateId string) (describeCertificate *ssl.DescribeCertificateResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeCertificateRequest()
	)
	request.CertificateId = &certificateId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeCertificate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	describeCertificate = response.Response
	return
}
func (me *SslService) DescribeSslDescribeCertificateBindResourceTaskDetailByFilter(ctx context.Context, param map[string]interface{}) (describeCertificateBindResourceTaskDetail *ssl.DescribeCertificateBindResourceTaskDetailResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeCertificateBindResourceTaskDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
		if k == "ResourceTypes" {
			request.ResourceTypes = v.([]*string)
		}
		if k == "Regions" {
			request.Regions = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var offset, limit = 0, 20

	for {
		request.Offset = helper.String(helper.IntToStr(offset))
		request.Limit = helper.String(helper.IntToStr(limit))
		response, err := me.client.UseSSLCertificateClient().DescribeCertificateBindResourceTaskDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || checkResult(1, response.Response) {
			break
		}

		describeCertificateBindResourceTaskDetail.CLB = append(describeCertificateBindResourceTaskDetail.CLB, response.Response.CLB...)
		describeCertificateBindResourceTaskDetail.CDN = append(describeCertificateBindResourceTaskDetail.CDN, response.Response.CDN...)
		describeCertificateBindResourceTaskDetail.WAF = append(describeCertificateBindResourceTaskDetail.WAF, response.Response.WAF...)
		describeCertificateBindResourceTaskDetail.DDOS = append(describeCertificateBindResourceTaskDetail.DDOS, response.Response.DDOS...)
		describeCertificateBindResourceTaskDetail.LIVE = append(describeCertificateBindResourceTaskDetail.LIVE, response.Response.LIVE...)
		describeCertificateBindResourceTaskDetail.VOD = append(describeCertificateBindResourceTaskDetail.VOD, response.Response.VOD...)
		describeCertificateBindResourceTaskDetail.TKE = append(describeCertificateBindResourceTaskDetail.TKE, response.Response.TKE...)
		describeCertificateBindResourceTaskDetail.APIGATEWAY = append(describeCertificateBindResourceTaskDetail.APIGATEWAY, response.Response.APIGATEWAY...)
		describeCertificateBindResourceTaskDetail.TCB = append(describeCertificateBindResourceTaskDetail.TCB, response.Response.TCB...)
		describeCertificateBindResourceTaskDetail.TEO = append(describeCertificateBindResourceTaskDetail.TEO, response.Response.TEO...)

		if checkResult(limit, response.Response) {
			break
		}

		offset += limit
	}
	return
}
func checkResult(num int, result *ssl.DescribeCertificateBindResourceTaskDetailResponseParams) bool {
	return len(result.CLB) < num && len(result.CDN) < num &&
		len(result.WAF) < num && len(result.DDOS) < num &&
		len(result.LIVE) < num && len(result.VOD) < num &&
		len(result.TKE) < num && len(result.APIGATEWAY) < num &&
		len(result.TCB) < num && len(result.TEO) < num
}

func (me *SslService) DescribeSslDescribeCompaniesByFilter(ctx context.Context, param map[string]interface{}) (describeCompanies []*ssl.CompanyInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeCompaniesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CompanyId" {
			request.CompanyId = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSSLCertificateClient().DescribeCompanies(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.Companies) < 1 {
			break
		}
		describeCompanies = append(describeCompanies, response.Response.Companies...)
		if len(response.Response.Companies) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
func (me *SslService) DescribeSslDescribeHostApiGatewayInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostApiGatewayInstanceList []*ssl.ApiGatewayInstanceDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostApiGatewayInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeHostApiGatewayInstanceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.InstanceList) < 1 {
		return
	}

	describeHostApiGatewayInstanceList = response.Response.InstanceList
	return
}
func (me *SslService) DescribeSslDescribeHostCdnInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostCdnInstanceList *ssl.DescribeHostCdnInstanceListResponseParams, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = ssl.NewDescribeHostCdnInstanceListRequest()
		response = ssl.NewDescribeHostCdnInstanceListResponse()
	)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
		if k == "AsyncCache" {
			request.AsyncCache = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
		list   []*ssl.CdnInstanceDetail
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err = me.client.UseSSLCertificateClient().DescribeHostCdnInstanceList(request)
		if err != nil {
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		list = append(list, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}
		offset += limit
	}

	if response == nil || response.Response == nil || response.Response.InstanceList == nil {
		return
	}

	describeHostCdnInstanceList = &ssl.DescribeHostCdnInstanceListResponseParams{
		InstanceList: list,
	}
	if response.Response.AsyncOffset != nil {
		describeHostCdnInstanceList.AsyncOffset = response.Response.AsyncOffset
	}
	if response.Response.AsyncCacheTime != nil {
		describeHostCdnInstanceList.AsyncCacheTime = response.Response.AsyncCacheTime
	}
	if response.Response.AsyncTotalNum != nil {
		describeHostCdnInstanceList.AsyncOffset = response.Response.AsyncTotalNum
	}

	return
}
func (me *SslService) DescribeSslDescribeHostClbInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostClbInstanceList *ssl.DescribeHostClbInstanceListResponseParams, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = ssl.NewDescribeHostClbInstanceListRequest()
		response = ssl.NewDescribeHostClbInstanceListResponse()
	)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "AsyncCache" {
			request.AsyncCache = v.(*int64)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
		list   []*ssl.ClbInstanceDetail
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err = me.client.UseSSLCertificateClient().DescribeHostClbInstanceList(request)
		if err != nil {
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		list = append(list, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}
	if response == nil || response.Response == nil || response.Response.InstanceList == nil {
		return
	}

	describeHostClbInstanceList = &ssl.DescribeHostClbInstanceListResponseParams{
		InstanceList: list,
	}
	if response.Response.AsyncOffset != nil {
		describeHostClbInstanceList.AsyncOffset = response.Response.AsyncOffset
	}
	if response.Response.AsyncCacheTime != nil {
		describeHostClbInstanceList.AsyncCacheTime = response.Response.AsyncCacheTime
	}
	if response.Response.AsyncTotalNum != nil {
		describeHostClbInstanceList.AsyncOffset = response.Response.AsyncTotalNum
	}

	return
}
func (me *SslService) DescribeSslDescribeHostCosInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostCosInstanceList *ssl.DescribeHostCosInstanceListResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostCosInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeHostCosInstanceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.InstanceList) < 1 {
		return
	}

	describeHostCosInstanceList = &ssl.DescribeHostCosInstanceListResponseParams{
		InstanceList: response.Response.InstanceList,
	}
	if response.Response.AsyncOffset != nil {
		describeHostCosInstanceList.AsyncOffset = response.Response.AsyncOffset
	}
	if response.Response.AsyncCacheTime != nil {
		describeHostCosInstanceList.AsyncCacheTime = response.Response.AsyncCacheTime
	}
	if response.Response.AsyncTotalNum != nil {
		describeHostCosInstanceList.AsyncOffset = response.Response.AsyncTotalNum
	}
	return
}
func (me *SslService) DescribeSslDescribeHostDdosInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostDdosInstanceList []*ssl.DdosInstanceDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostDdosInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeHostDdosInstanceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.InstanceList) < 1 {
		return
	}
	describeHostDdosInstanceList = response.Response.InstanceList
	return
}
func (me *SslService) DescribeSslDescribeHostDeployRecordByFilter(ctx context.Context, param map[string]interface{}) (describeHostDeployRecord []*ssl.DeployRecordInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostDeployRecordRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSSLCertificateClient().DescribeHostDeployRecord(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DeployRecordList) < 1 {
			break
		}
		describeHostDeployRecord = append(describeHostDeployRecord, response.Response.DeployRecordList...)
		if len(response.Response.DeployRecordList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
func (me *SslService) DescribeSslDescribeHostDeployRecordDetailByFilter(
	ctx context.Context, param map[string]interface{}) (
	describeHostDeployRecordDetail []*ssl.DeployRecordDetail,
	successTotalCount *int64, failedTotalCount *int64, runningTotalCount *int64,
	errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostDeployRecordDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DeployRecordId" {
			request.DeployRecordId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSSLCertificateClient().DescribeHostDeployRecordDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DeployRecordDetailList) < 1 {
			break
		}
		describeHostDeployRecordDetail = append(describeHostDeployRecordDetail, response.Response.DeployRecordDetailList...)
		if response.Response != nil {
			if response.Response.SuccessTotalCount != nil {
				successTotalCount = response.Response.SuccessTotalCount
			}
			if response.Response.SuccessTotalCount != nil {
				failedTotalCount = response.Response.FailedTotalCount
			}
			if response.Response.SuccessTotalCount != nil {
				runningTotalCount = response.Response.SuccessTotalCount
			}
		}
		if len(response.Response.DeployRecordDetailList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
func (me *SslService) DescribeSslDescribeHostLighthouseInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostLighthouseInstanceList []*ssl.LighthouseInstanceDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostLighthouseInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeHostLighthouseInstanceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.InstanceList) < 1 {
		return
	}
	describeHostLighthouseInstanceList = response.Response.InstanceList
	return
}
func (me *SslService) DescribeSslDescribeHostLiveInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostLiveInstanceList []*ssl.LiveInstanceDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostLiveInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeHostLiveInstanceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.InstanceList) < 1 {
		return
	}
	describeHostLiveInstanceList = response.Response.InstanceList

	return
}
func (me *SslService) DescribeSslDescribeHostTeoInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostTeoInstanceList []*ssl.TeoInstanceDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostTeoInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSSLCertificateClient().DescribeHostTeoInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostTeoInstanceList = append(describeHostTeoInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
func (me *SslService) DescribeSslDescribeHostTkeInstanceListByFilter(
	ctx context.Context, param map[string]interface{}) (
	describeHostTkeInstanceList []*ssl.TkeInstanceDetail,
	asyncTotalNum *int64, asyncOffset *int64, asyncCacheTime *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostTkeInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "AsyncCache" {
			request.AsyncCache = v.(*int64)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSSLCertificateClient().DescribeHostTkeInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostTkeInstanceList = append(describeHostTkeInstanceList, response.Response.InstanceList...)
		if response.Response != nil {
			if response.Response.AsyncTotalNum != nil {
				asyncTotalNum = response.Response.AsyncTotalNum
			}
			if response.Response.AsyncOffset != nil {
				asyncOffset = response.Response.AsyncOffset
			}
			if response.Response.AsyncCacheTime != nil {
				asyncCacheTime = response.Response.AsyncCacheTime
			}
		}
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostUpdateRecordByFilter(ctx context.Context, param map[string]interface{}) (describeHostUpdateRecord []*ssl.UpdateRecordInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostUpdateRecordRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSSLCertificateClient().DescribeHostUpdateRecord(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DeployRecordList) < 1 {
			break
		}
		describeHostUpdateRecord = append(describeHostUpdateRecord, response.Response.DeployRecordList...)
		if len(response.Response.DeployRecordList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostUpdateRecordDetailByFilter(
	ctx context.Context, param map[string]interface{}) (
	describeHostUpdateRecordDetail []*ssl.UpdateRecordDetails,
	successTotalCount *int64, failedTotalCount *int64, runningTotalCount *int64, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostUpdateRecordDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DeployRecordId" {
			request.DeployRecordId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = helper.Int64ToStrPoint(offset)
		request.Limit = helper.Int64ToStrPoint(limit)
		response, err := me.client.UseSSLCertificateClient().DescribeHostUpdateRecordDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RecordDetailList) < 1 {
			break
		}
		describeHostUpdateRecordDetail = append(describeHostUpdateRecordDetail, response.Response.RecordDetailList...)
		if response.Response != nil {
			if response.Response.SuccessTotalCount != nil {
				successTotalCount = response.Response.SuccessTotalCount
			}
			if response.Response.SuccessTotalCount != nil {
				failedTotalCount = response.Response.FailedTotalCount
			}
			if response.Response.SuccessTotalCount != nil {
				runningTotalCount = response.Response.SuccessTotalCount
			}
		}
		if len(response.Response.RecordDetailList) < int(limit) {
			break
		}
		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostVodInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostVodInstanceList []*ssl.VodInstanceDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostVodInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeHostVodInstanceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.InstanceList) < 1 {
		return
	}

	describeHostVodInstanceList = response.Response.InstanceList
	return
}

func (me *SslService) DescribeSslDescribeHostWafInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostWafInstanceList []*ssl.LiveInstanceDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeHostWafInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CertificateId" {
			request.CertificateId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "IsCache" {
			request.IsCache = v.(*uint64)
		}
		if k == "Filters" {
			request.Filters = v.([]*ssl.Filter)
		}
		if k == "OldCertificateId" {
			request.OldCertificateId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeHostWafInstanceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.InstanceList) < 1 {
		return
	}

	describeHostWafInstanceList = response.Response.InstanceList
	return
}

func (me *SslService) DescribeSslDescribeManagerDetailByFilter(ctx context.Context, managerId string) (describeManagerDetailResponse *ssl.DescribeManagerDetailResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeManagerDetailRequest()
	)
	request.ManagerId = helper.StrToInt64Point(managerId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeManagerDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	describeManagerDetailResponse = response.Response
	return
}

func (me *SslService) DescribeSslDescribeManagersByFilter(ctx context.Context, param map[string]interface{}) (describeManagers []*ssl.ManagerInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeManagersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CompanyId" {
			request.CompanyId = v.(*int64)
		}
		if k == "ManagerName" {
			request.ManagerName = v.(*string)
		}
		if k == "ManagerMail" {
			request.ManagerMail = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSSLCertificateClient().DescribeManagers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Managers) < 1 {
			break
		}
		describeManagers = append(describeManagers, response.Response.Managers...)
		if len(response.Response.Managers) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeCertificateBindResourceTaskResultByFilter(ctx context.Context, param map[string]interface{}) (describeCertificateBindResourceTaskResult []*ssl.SyncTaskBindResourceResult, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ssl.NewDescribeCertificateBindResourceTaskResultRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TaskIds" {
			request.TaskIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSSLCertificateClient().DescribeCertificateBindResourceTaskResult(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	describeCertificateBindResourceTaskResult = response.Response.SyncTaskBindResourceResult

	return
}
