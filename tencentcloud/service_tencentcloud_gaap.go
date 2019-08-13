package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type GaapService struct {
	client *connectivity.TencentCloudClient
}

func (me *GaapService) CreateRealserver(ctx context.Context, address, name string, projectId int, tags map[string]string) (id string, err error) {
	logId := getLogId(ctx)

	request := gaap.NewAddRealServersRequest()
	request.RealServerName = &name
	request.RealServerIP = common.StringPtrs([]string{address})
	request.ProjectId = intToPointer(projectId)
	for k, v := range tags {
		request.TagSet = append(request.TagSet, &gaap.TagPair{
			TagKey:   common.StringPtr(k),
			TagValue: stringToPointer(v),
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseGaapClient().AddRealServers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.RealServerSet) == 0 {
			err = fmt.Errorf("api[%s] return empty realserver set", request.GetAction())
			log.Printf("[CRITAL]%s, %v", logId, err)
			return resource.NonRetryableError(err)
		}

		realserver := response.Response.RealServerSet[0]

		if realserver.RealServerId == nil {
			err = fmt.Errorf("api[%s] return realserver id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if realserver.RealServerIP == nil {
			err = fmt.Errorf("api[%s] return realserver ip or domain is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *realserver.RealServerId

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create realserver failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) DescribeRealservers(ctx context.Context, address, name *string, tags map[string]string, projectId int) (realservers []*gaap.BindRealServerInfo, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeRealServersRequest()
	request.SearchValue = address
	if name != nil {
		request.Filters = append(request.Filters, &gaap.Filter{
			Name:   stringToPointer("RealServerName"),
			Values: []*string{name},
		})
	}
	for k, v := range tags {
		request.TagSet = append(request.TagSet, &gaap.TagPair{
			TagKey:   stringToPointer(k),
			TagValue: stringToPointer(v),
		})
	}
	request.ProjectId = common.Int64Ptr(int64(projectId))

	request.Limit = intToPointer(50)
	offset := 0

	// run loop at least one times
	count := 50
	for count == 50 {
		request.Offset = intToPointer(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err := me.client.UseGaapClient().DescribeRealServers(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			realservers = append(realservers, response.Response.RealServerSet...)
			count = len(response.Response.RealServerSet)
			offset += count

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read realservers failed, reason: %v", logId, err)
			return nil, err
		}
	}

	return
}

func (me *GaapService) ModifyRealserverName(ctx context.Context, id, name string) error {
	logId := getLogId(ctx)

	request := gaap.NewModifyRealServerNameRequest()
	request.RealServerId = &id
	request.RealServerName = &name

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := me.client.UseGaapClient().ModifyRealServerName(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify realserver %s name failed, reason: %v", logId, id, err)
		return err
	}
	return nil
}

func (me *GaapService) DeleteRealserver(ctx context.Context, id string) error {
	logId := getLogId(ctx)

	request := gaap.NewRemoveRealServersRequest()
	request.RealServerIds = []*string{&id}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := me.client.UseGaapClient().RemoveRealServers(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete realserver %s failed, reason: %v", logId, id, err)
		return err
	}

	return nil
}

func (me *GaapService) createCertificate(ctx context.Context, certificateType int, content string, name, key *string) (id string, err error) {
	logId := getLogId(ctx)

	request := gaap.NewCreateCertificateRequest()
	request.CertificateType = common.Int64Ptr(int64(certificateType))
	request.CertificateContent = &content
	request.CertificateAlias = name
	request.CertificateKey = key

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseGaapClient().CreateCertificate(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		if response.Response.CertificateId == nil {
			err := fmt.Errorf("certificate id is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.CertificateId

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create certiciate failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) DescribeCertificateById(ctx context.Context, id string) (certificate *gaap.CertificateDetail, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeCertificateDetailRequest()
	request.CertificateId = &id

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseGaapClient().DescribeCertificateDetail(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Message == "CertificateId not found" {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		certificate = response.Response.CertificateDetail
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read certificate failed, reason: %v", logId, err)
		return nil, err
	}

	return
}

func (me *GaapService) ModifyCertificateName(ctx context.Context, id, name string) error {
	logId := getLogId(ctx)

	request := gaap.NewModifyCertificateAttributesRequest()
	request.CertificateId = &id
	request.CertificateAlias = &name

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := me.client.UseGaapClient().ModifyCertificateAttributes(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify certificate name failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DeleteCertificate(ctx context.Context, id string) error {
	logId := getLogId(ctx)

	request := gaap.NewDeleteCertificateRequest()
	request.CertificateId = &id

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := me.client.UseGaapClient().DeleteCertificate(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete certificate failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) CreateProxy(
	ctx context.Context,
	name, accessRegion, realserverRegion string,
	bandwidth, concurrent, projectId int,
	tags map[string]string,
) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	createRequest := gaap.NewCreateProxyRequest()
	createRequest.ProxyName = &name
	createRequest.ProjectId = common.Int64Ptr(int64(projectId))
	createRequest.Bandwidth = intToPointer(bandwidth)
	createRequest.Concurrent = intToPointer(concurrent)
	createRequest.AccessRegion = &accessRegion
	createRequest.RealServerRegion = &realserverRegion
	for k, v := range tags {
		createRequest.TagSet = append(createRequest.TagSet, &gaap.TagPair{
			TagKey:   stringToPointer(k),
			TagValue: stringToPointer(v),
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, err := client.CreateProxy(createRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, createRequest.GetAction(), createRequest.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.InstanceId == nil {
			err := errors.New("proxy id is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.InstanceId
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create proxy failed, reason: %v", logId, err)
		return "", err
	}

	describeRequest := gaap.NewDescribeProxiesRequest()
	describeRequest.ProxyIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet
		if len(proxies) == 0 {
			err := errors.New("read no proxy")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		proxy := proxies[0]
		if proxy.Status == nil {
			err := errors.New("proxy status is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *proxy.Status != GAAP_PROXY_RUNNING {
			err := errors.New("proxy is still creating")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create proxy failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) EnableProxy(ctx context.Context, id string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	enableRequest := gaap.NewOpenProxiesRequest()
	enableRequest.ProxyIds = []*string{&id}
	enableRequest.ClientToken = stringToPointer(buildToken())

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, err := client.OpenProxies(enableRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, enableRequest.GetAction(), enableRequest.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.OperationFailedInstanceSet) > 0 {
			err := errors.New("enable proxy failed")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		// proxy is enabled
		if len(response.Response.InvalidStatusInstanceSet) > 0 {
			return nil
		}

		// enable proxy successfully
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s enable proxy failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeProxiesRequest()
	describeRequest.ProxyIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet
		if len(proxies) == 0 {
			err := errors.New("read no proxy")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		proxy := proxies[0]
		if proxy.Status == nil {
			err := errors.New("proxy status is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *proxy.Status != GAAP_PROXY_RUNNING {
			err := errors.New("proxy is still enabling")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s enable proxy failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DisableProxy(ctx context.Context, id string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	disableRequest := gaap.NewCloseProxiesRequest()
	disableRequest.ProxyIds = []*string{&id}
	disableRequest.ClientToken = stringToPointer(buildToken())

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, err := client.CloseProxies(disableRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, disableRequest.GetAction(), disableRequest.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.OperationFailedInstanceSet) > 0 {
			err := errors.New("disable proxy failed")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		// proxy is disabled
		if len(response.Response.InvalidStatusInstanceSet) > 0 {
			return nil
		}
		// disable proxy successfully
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s disable proxy failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeProxiesRequest()
	describeRequest.ProxyIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet
		if len(proxies) == 0 {
			err := errors.New("read no proxy")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		proxy := proxies[0]
		if proxy.Status == nil {
			err := errors.New("proxy status is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *proxy.Status != GAAP_PROXY_CLOSED {
			err := errors.New("proxy is still enabling")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s disable proxy failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DescribeProxies(
	ctx context.Context,
	ids []string,
	projectId *int,
	accessRegion, realserverRegion *string,
	tags map[string]string,
) (proxies []*gaap.ProxyInfo, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeProxiesRequest()
	if len(ids) > 0 {
		request.ProxyIds = common.StringPtrs(ids)
	}
	if projectId != nil {
		request.Filters = append(request.Filters, &gaap.Filter{
			Name:   stringToPointer("ProjectId"),
			Values: []*string{stringToPointer(strconv.Itoa(*projectId))},
		})
	}
	if accessRegion != nil {
		request.Filters = append(request.Filters, &gaap.Filter{
			Name:   stringToPointer("AccessRegion"),
			Values: []*string{accessRegion},
		})
	}
	if accessRegion != nil {
		request.Filters = append(request.Filters, &gaap.Filter{
			Name:   stringToPointer("RealServerRegion"),
			Values: []*string{realserverRegion},
		})
	}
	for k, v := range tags {
		request.TagSet = append(request.TagSet, &gaap.TagPair{
			TagKey:   stringToPointer(k),
			TagValue: stringToPointer(v),
		})
	}

	request.Limit = intToPointer(100)
	offset := 0

	// to run loop at least one times
	count := 100
	for count == 100 {
		request.Offset = intToPointer(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err := me.client.UseGaapClient().DescribeProxies(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			count = len(response.Response.ProxySet)
			proxies = append(proxies, response.Response.ProxySet...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read proxies failed, reason: %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *GaapService) ModifyProxyName(ctx context.Context, id, name string) error {
	logId := getLogId(ctx)

	request := gaap.NewModifyProxiesAttributeRequest()
	request.ProxyIds = []*string{&id}
	request.ProxyName = &name
	request.ClientToken = stringToPointer(buildToken())

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := me.client.UseGaapClient().ModifyProxiesAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify proxy name failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) ModifyProxyProjectId(ctx context.Context, id string, projectId int) error {
	logId := getLogId(ctx)

	request := gaap.NewModifyProxiesProjectRequest()
	request.ProxyIds = []*string{&id}
	request.ProjectId = common.Int64Ptr(int64(projectId))
	request.ClientToken = stringToPointer(buildToken())

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := me.client.UseGaapClient().ModifyProxiesProject(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify proxy project id failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) ModifyProxyConfiguration(ctx context.Context, id string, bandwidth, concurrent *int) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	nodifyRequest := gaap.NewModifyProxyConfigurationRequest()
	nodifyRequest.ProxyId = &id
	if bandwidth != nil {
		nodifyRequest.Bandwidth = intToPointer(*bandwidth)
	}
	if concurrent != nil {
		nodifyRequest.Concurrent = intToPointer(*concurrent)
	}
	nodifyRequest.ClientToken = stringToPointer(buildToken())

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := client.ModifyProxyConfiguration(nodifyRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, nodifyRequest.GetAction(), nodifyRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify proxy configuration failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeProxiesRequest()
	describeRequest.ProxyIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet
		if len(proxies) == 0 {
			err := errors.New("read no proxy")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		proxy := proxies[0]
		if proxy.Status == nil {
			err := errors.New("proxy status is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *proxy.Status == GAAP_PROXY_ADJUSTING {
			err := errors.New("proxy is still modifying")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify proxy configuration failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DeleteProxy(ctx context.Context, id string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	deleteRequest := gaap.NewDestroyProxiesRequest()
	deleteRequest.ProxyIds = []*string{&id}
	deleteRequest.Force = common.Int64Ptr(0)
	deleteRequest.ClientToken = stringToPointer(buildToken())

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, err := client.DestroyProxies(deleteRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.OperationFailedInstanceSet) > 0 {
			err := errors.New("delete proxy failed")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if len(response.Response.InvalidStatusInstanceSet) > 0 {
			err := errors.New("proxy can't be deleted")
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete proxy failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeProxiesRequest()
	describeRequest.ProxyIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.ProxySet) > 0 {
			err := errors.New("proxy still exist")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete proxy failed, reason: %v", logId, err)
		return err
	}

	return nil
}
