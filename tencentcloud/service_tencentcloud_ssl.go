package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss/v20180426"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SslService struct {
	client *connectivity.TencentCloudClient
}

func (me *SslService) CreateCertificate(ctx context.Context, certType, cert, name string, projectId int, key *string) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseSslClient()

	createRequest := ssl.NewUploadCertRequest()
	createRequest.Cert = &cert
	createRequest.CertType = &certType
	createRequest.ProjectId = helper.String(strconv.Itoa(projectId))
	createRequest.ModuleType = helper.String(SSL_MODULE_TYPE)
	createRequest.Alias = &name
	createRequest.Key = key

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(createRequest.GetAction())

		response, err := client.UploadCert(createRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, createRequest.GetAction(), createRequest.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.Id == nil {
			err := fmt.Errorf("api[%s] return id is nil", createRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.Id
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create certificate failed, reason: %v", logId, err)
		return "", err
	}

	describeRequest := ssl.NewDescribeCertListRequest()
	describeRequest.ModuleType = helper.String(SSL_MODULE_TYPE)
	describeRequest.Id = &id

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeCertList(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		var certificate *ssl.SSLCertificate
		for _, c := range response.Response.CertificateSet {
			if c.Id == nil {
				err := fmt.Errorf("api[%s] certificate id is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *c.Id == id {
				certificate = c
				break
			}
		}

		if certificate == nil {
			err := fmt.Errorf("api[%s] certificate not found", describeRequest.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		if certificate.Status == nil {
			err := fmt.Errorf("api[%s] certificate status is nil", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *certificate.Status != SSL_STATUS_AVAILABLE {
			err := fmt.Errorf("certificate is not available, status is %d", *certificate.Status)
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create certificate failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *SslService) DescribeCertificates(ctx context.Context, id, name, certType *string) (certificates []*ssl.SSLCertificate, err error) {
	logId := getLogId(ctx)

	request := ssl.NewDescribeCertListRequest()
	request.ModuleType = helper.String(SSL_MODULE_TYPE)
	request.SearchKey = name
	request.Id = id
	request.CertType = certType
	request.WithCert = helper.String(SSL_WITH_CERT)

	var offset uint64

	request.Offset = &offset
	request.Limit = helper.IntUint64(20)

	// run loop at least once
	count := 20
	for count == 20 {
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseSslClient().DescribeCertList(request)
			if err != nil {
				count = 0

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			count = len(response.Response.CertificateSet)
			certificates = append(certificates, response.Response.CertificateSet...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read certificates failed, reason: %v", logId, err)
			return nil, err
		}

		offset += uint64(count)
	}

	return
}

func (me *SslService) DeleteCertificate(ctx context.Context, id string) error {
	logId := getLogId(ctx)
	client := me.client.UseSslClient()

	deleteRequest := ssl.NewDeleteCertRequest()
	deleteRequest.ModuleType = helper.String(SSL_MODULE_TYPE)
	deleteRequest.Id = &id

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteCert(deleteRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete certificate failed, reason: %v", logId, err)
		return err
	}

	describeRequest := ssl.NewDescribeCertListRequest()
	describeRequest.ModuleType = helper.String(SSL_MODULE_TYPE)
	describeRequest.Id = &id

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeCertList(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		for _, c := range response.Response.CertificateSet {
			if c.Id == nil {
				err := fmt.Errorf("api[%s] certificate id is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *c.Id == id {
				err := errors.New("certificate still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete certificate failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *SslService) CheckCertificateType(ctx context.Context, certId string, checkType string) (bool, error) {

	//get certificate by id

	certificates, err := me.DescribeCertificates(ctx, &certId, nil, nil)
	if err != nil {
		return false, err
	}

	var certificate *ssl.SSLCertificate
	for _, c := range certificates {
		if c.Id == nil {
			return false, errors.New("certificate id is nil")
		}

		if *c.Id == certId {
			certificate = c
			break
		}
	}

	if certificate != nil && *certificate.CertType == checkType {
		return true, nil
	} else {
		if certificate == nil {
			return false, fmt.Errorf("certificate id %s is not found", certId)
		}
		return false, nil
	}

}

func (me *SslService) DescribeSslDescribeCertificateByFilter(ctx context.Context, param map[string]interface{}) (describeCertificate []*ssl.DescribeCertificateResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ssl.NewDescribeCertificateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeCertificate(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.OwnerUin) < 1 {
			break
		}
		describeCertificate = append(describeCertificate, response.Response.OwnerUin...)
		if len(response.Response.OwnerUin) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeCompaniesByFilter(ctx context.Context, param map[string]interface{}) (describeCompanies []*ssl.CompanyInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		response, err := me.client.UseSslClient().DescribeCompanies(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Companies) < 1 {
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
		logId   = getLogId(ctx)
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostApiGatewayInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostApiGatewayInstanceList = append(describeHostApiGatewayInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostCdnInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostCdnInstanceList []*ssl.CdnInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ssl.NewDescribeHostCdnInstanceListRequest()
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
		if k == "AsyncCache" {
			request.AsyncCache = v.(*int64)
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
		response, err := me.client.UseSslClient().DescribeHostCdnInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostCdnInstanceList = append(describeHostCdnInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostClbInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostClbInstanceList []*ssl.ClbInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ssl.NewDescribeHostClbInstanceListRequest()
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
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostClbInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostClbInstanceList = append(describeHostClbInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostCosInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostCosInstanceList []*ssl.CosInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostCosInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostCosInstanceList = append(describeHostCosInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostDdosInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostDdosInstanceList []*ssl.DdosInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostDdosInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostDdosInstanceList = append(describeHostDdosInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostDeployRecordByFilter(ctx context.Context, param map[string]interface{}) (describeHostDeployRecord []*ssl.DeployRecordInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostDeployRecord(request)
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

func (me *SslService) DescribeSslDescribeHostDeployRecordDetailByFilter(ctx context.Context, param map[string]interface{}) (describeHostDeployRecordDetail []*ssl.DeployRecordDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostDeployRecordDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DeployRecordDetailList) < 1 {
			break
		}
		describeHostDeployRecordDetail = append(describeHostDeployRecordDetail, response.Response.DeployRecordDetailList...)
		if len(response.Response.DeployRecordDetailList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostLiveInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostLiveInstanceList []*ssl.LiveInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostLiveInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostLiveInstanceList = append(describeHostLiveInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostTeoInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostTeoInstanceList []*ssl.TeoInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostTeoInstanceList(request)
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

func (me *SslService) DescribeSslDescribeHostTkeInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostTkeInstanceList []*ssl.TkeInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostTkeInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostTkeInstanceList = append(describeHostTkeInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostUpdateRecordByFilter(ctx context.Context, param map[string]interface{}) (describeHostUpdateRecord []*ssl.UpdateRecordInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostUpdateRecord(request)
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

func (me *SslService) DescribeSslDescribeHostUpdateRecordDetailByFilter(ctx context.Context, param map[string]interface{}) (describeHostUpdateRecordDetail []*ssl.UpdateRecordDetails, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostUpdateRecordDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RecordDetailList) < 1 {
			break
		}
		describeHostUpdateRecordDetail = append(describeHostUpdateRecordDetail, response.Response.RecordDetailList...)
		if len(response.Response.RecordDetailList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostVodInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostVodInstanceList []*ssl.VodInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostVodInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostVodInstanceList = append(describeHostVodInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeHostWafInstanceListByFilter(ctx context.Context, param map[string]interface{}) (describeHostWafInstanceList []*ssl.LiveInstanceDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeHostWafInstanceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		describeHostWafInstanceList = append(describeHostWafInstanceList, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeManagerDetailByFilter(ctx context.Context, param map[string]interface{}) (describeManagerDetail []*ssl.DescribeManagerDetailResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ssl.NewDescribeManagerDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSslClient().DescribeManagerDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Status) < 1 {
			break
		}
		describeManagerDetail = append(describeManagerDetail, response.Response.Status...)
		if len(response.Response.Status) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SslService) DescribeSslDescribeManagersByFilter(ctx context.Context, param map[string]interface{}) (describeManagers []*ssl.ManagerInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		response, err := me.client.UseSslClient().DescribeManagers(request)
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
