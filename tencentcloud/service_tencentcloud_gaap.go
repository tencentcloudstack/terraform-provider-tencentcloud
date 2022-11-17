package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type gaapRealserverBind struct {
	id     string
	ip     string
	port   int
	weight int
}

type gaapHttpRule struct {
	listenerId                 string
	domain                     string
	path                       string
	realserverType             string
	scheduler                  string
	healthCheck                bool
	interval                   int
	connectTimeout             int
	healthCheckPath            string
	healthCheckMethod          string
	healthCheckStatusCodes     []int
	forwardHost                string
	serverNameIndicationSwitch string
	serverNameIndication       string
}

type GaapService struct {
	client *connectivity.TencentCloudClient
}

func (me *GaapService) CreateRealserver(ctx context.Context, address, name string, projectId int) (id string, err error) {
	logId := getLogId(ctx)

	request := gaap.NewAddRealServersRequest()
	request.RealServerName = &name
	request.RealServerIP = []*string{&address}
	request.ProjectId = helper.IntUint64(projectId)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().AddRealServers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.RealServerSet) != 1 {
			err = fmt.Errorf("api[%s] return %d realservers", request.GetAction(), len(response.Response.RealServerSet))
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
			Name:   helper.String("RealServerName"),
			Values: []*string{name},
		})
	}
	for k, v := range tags {
		request.TagSet = append(request.TagSet, &gaap.TagPair{
			TagKey:   helper.String(k),
			TagValue: helper.String(v),
		})
	}
	request.ProjectId = helper.IntInt64(projectId)

	request.Limit = helper.IntUint64(50)
	offset := 0

	// run loop at least one times
	count := 50
	for count == 50 {
		request.Offset = helper.IntUint64(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseGaapClient().DescribeRealServers(request)
			if err != nil {
				count = 0

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			realservers = append(realservers, response.Response.RealServerSet...)
			count = len(response.Response.RealServerSet)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read realservers failed, reason: %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *GaapService) ModifyRealserverName(ctx context.Context, id, name string) error {
	logId := getLogId(ctx)

	request := gaap.NewModifyRealServerNameRequest()
	request.RealServerId = &id
	request.RealServerName = &name

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

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
		ratelimit.Check(request.GetAction())

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

func (me *GaapService) createCertificate(ctx context.Context, certificateType int, content, name string, key *string) (id string, err error) {
	logId := getLogId(ctx)

	request := gaap.NewCreateCertificateRequest()
	request.CertificateType = helper.IntInt64(certificateType)
	request.CertificateContent = &content
	request.CertificateAlias = &name
	request.CertificateKey = key

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

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
		ratelimit.Check(request.GetAction())

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
		ratelimit.Check(request.GetAction())

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
		ratelimit.Check(request.GetAction())

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
	createRequest.ProjectId = helper.IntInt64(projectId)
	createRequest.Bandwidth = helper.IntUint64(bandwidth)
	createRequest.Concurrent = helper.IntUint64(concurrent)
	createRequest.AccessRegion = &accessRegion
	createRequest.RealServerRegion = &realserverRegion
	for k, v := range tags {
		createRequest.TagSet = append(createRequest.TagSet, &gaap.TagPair{
			TagKey:   helper.String(k),
			TagValue: helper.String(v),
		})
	}

	if err := resource.Retry(2*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(createRequest.GetAction())
		createRequest.ClientToken = helper.String(helper.BuildToken())

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
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet

		switch len(proxies) {
		case 0:
			err := errors.New("read no proxy")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)

		default:
			err := errors.New("return more than 1 proxy")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.NonRetryableError(err)

		case 1:
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

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(enableRequest.GetAction())
		enableRequest.ClientToken = helper.String(helper.BuildToken())

		response, err := client.OpenProxies(enableRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, enableRequest.GetAction(), enableRequest.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.OperationFailedInstanceSet) > 0 {
			err := fmt.Errorf("api[%s] enable proxy failed", enableRequest.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		// proxy may be enabled
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
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet

		if len(proxies) != 1 {
			err := fmt.Errorf("api[%s] read %d proxies", describeRequest.GetAction(), len(proxies))
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		proxy := proxies[0]
		if proxy.Status == nil {
			err := fmt.Errorf("api[%s] proxy status is nil", describeRequest.GetAction())
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

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(disableRequest.GetAction())
		disableRequest.ClientToken = helper.String(helper.BuildToken())

		response, err := client.CloseProxies(disableRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, disableRequest.GetAction(), disableRequest.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.OperationFailedInstanceSet) > 0 {
			err := fmt.Errorf("api[%s] disable proxy failed", disableRequest.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		// proxy may be disabled
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
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet

		if len(proxies) != 1 {
			err := fmt.Errorf("api[%s] read %d proxies", describeRequest.GetAction(), len(proxies))
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		proxy := proxies[0]
		if proxy.Status == nil {
			err := fmt.Errorf("api[%s] proxy status is nil", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *proxy.Status != GAAP_PROXY_CLOSED {
			err := errors.New("proxy is still disabling")
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
			Name:   helper.String("ProjectId"),
			Values: []*string{helper.String(strconv.Itoa(*projectId))},
		})
	}
	if accessRegion != nil {
		request.Filters = append(request.Filters, &gaap.Filter{
			Name:   helper.String("AccessRegion"),
			Values: []*string{accessRegion},
		})
	}
	if realserverRegion != nil {
		request.Filters = append(request.Filters, &gaap.Filter{
			Name:   helper.String("RealServerRegion"),
			Values: []*string{realserverRegion},
		})
	}
	for k, v := range tags {
		request.TagSet = append(request.TagSet, &gaap.TagPair{
			TagKey:   helper.String(k),
			TagValue: helper.String(v),
		})
	}

	request.Limit = helper.IntUint64(100)
	offset := 0

	// to run loop at least one times
	count := 100
	for count == 100 {
		request.Offset = helper.IntUint64(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseGaapClient().DescribeProxies(request)
			if err != nil {
				count = 0

				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == "ResourceNotFound" {
						return nil
					}
				}

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

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		request.ClientToken = helper.String(helper.BuildToken())

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
	request.ProjectId = helper.IntInt64(projectId)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		request.ClientToken = helper.String(helper.BuildToken())

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

	modifyRequest := gaap.NewModifyProxyConfigurationRequest()
	modifyRequest.ProxyId = &id
	if bandwidth != nil {
		modifyRequest.Bandwidth = helper.IntUint64(*bandwidth)
	}
	if concurrent != nil {
		modifyRequest.Concurrent = helper.IntUint64(*concurrent)
	}

	if err := resource.Retry(2*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(modifyRequest.GetAction())
		modifyRequest.ClientToken = helper.String(helper.BuildToken())

		if _, err := client.ModifyProxyConfiguration(modifyRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, modifyRequest.GetAction(), modifyRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify proxy configuration failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeProxiesRequest()
	describeRequest.ProxyIds = []*string{&id}

	if err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		proxies := response.Response.ProxySet

		if len(proxies) != 1 {
			err := fmt.Errorf("api[%s] read %d proxies", describeRequest.GetAction(), len(proxies))
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		proxy := proxies[0]
		if proxy.Status == nil {
			err := fmt.Errorf("api[%s] proxy status is nil", describeRequest.GetAction())
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
	deleteRequest.Force = helper.IntInt64(1)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())
		deleteRequest.ClientToken = helper.String(helper.BuildToken())

		response, err := client.DestroyProxies(deleteRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.OperationFailedInstanceSet) > 0 {
			err := fmt.Errorf("api[%s] delete proxy failed", deleteRequest.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		if len(response.Response.InvalidStatusInstanceSet) > 0 {
			err := fmt.Errorf("api[%s] proxy can't be deleted", deleteRequest.GetAction())
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
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeProxies(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" {
					return nil
				}
			}

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

func (me *GaapService) CreateTCPListener(
	ctx context.Context,
	name, scheduler, realserverType, proxyId string,
	port, interval, connectTimeout, clientIPMethod int,
	healthCheck bool,
) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewCreateTCPListenersRequest()
	request.ListenerName = &name
	request.Scheduler = &scheduler
	request.RealServerType = &realserverType
	request.ProxyId = &proxyId
	request.Ports = []*uint64{helper.IntUint64(port)}
	if healthCheck {
		request.HealthCheck = helper.IntUint64(1)
	} else {
		request.HealthCheck = helper.IntUint64(0)
	}
	request.DelayLoop = helper.IntUint64(interval)
	request.ConnectTimeout = helper.IntUint64(connectTimeout)
	request.ClientIPMethod = helper.IntInt64(clientIPMethod)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.CreateTCPListeners(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.ListenerIds) != 1 {
			err := fmt.Errorf("api[%s] return %d TCP listener ids", request.GetAction(), len(response.Response.ListenerIds))
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.ListenerIds[0]
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create TCP listener failed, reason: %v", logId, err)
		return "", err
	}

	if err := waitLayer4ListenerReady(ctx, client, id, "TCP"); err != nil {
		log.Printf("[CRITAL]%s create TCP listener failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) CreateUDPListener(
	ctx context.Context,
	name, scheduler, realserverType, proxyId string,
	port int,
) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewCreateUDPListenersRequest()
	request.ListenerName = &name
	request.Scheduler = &scheduler
	request.RealServerType = &realserverType
	request.ProxyId = &proxyId
	request.Ports = []*uint64{helper.IntUint64(port)}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.CreateUDPListeners(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.ListenerIds) != 1 {
			err := fmt.Errorf("api[%s] return %d UDP listener ids", request.GetAction(), len(response.Response.ListenerIds))
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.ListenerIds[0]
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create UDP listener failed, reason: %v", logId, err)
		return "", err
	}

	if err := waitLayer4ListenerReady(ctx, client, id, "UDP"); err != nil {
		log.Printf("[CRITAL]%s create UDP listener failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) BindLayer4ListenerRealservers(ctx context.Context, id, protocol, proxyId string, realserverBinds []gaapRealserverBind) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewBindListenerRealServersRequest()
	request.ListenerId = &id
	request.RealServerBindSet = make([]*gaap.RealServerBindSetReq, 0, len(realserverBinds))
	for _, bind := range realserverBinds {
		request.RealServerBindSet = append(request.RealServerBindSet, &gaap.RealServerBindSetReq{
			RealServerId:     helper.String(bind.id),
			RealServerPort:   helper.IntUint64(bind.port),
			RealServerIP:     helper.String(bind.ip),
			RealServerWeight: helper.IntUint64(bind.weight),
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.BindListenerRealServers(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s bind realservers to layer4 listener failed, reason: %v", logId, err)
		return err
	}

	if err := waitLayer4ListenerReady(ctx, client, id, protocol); err != nil {
		log.Printf("[CRITAL]%s bind realservers to layer4 listener failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DescribeTCPListeners(ctx context.Context, proxyId, listenerId, name *string, port *int) (listeners []*gaap.TCPListener, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeTCPListenersRequest()
	request.ProxyId = proxyId
	request.ListenerId = listenerId

	if port != nil {
		request.Port = helper.IntUint64(*port)
	}

	// if port set, name can't use fuzzy search
	if name != nil {
		if port != nil {
			request.ListenerName = name
		} else {
			request.SearchValue = name
		}
	}

	request.Limit = helper.IntUint64(50)
	offset := 0

	// to run loop at least once
	count := 50
	for count == 50 {
		request.Offset = helper.IntUint64(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseGaapClient().DescribeTCPListeners(request)
			if err != nil {
				count = 0

				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			count = len(response.Response.ListenerSet)
			listeners = append(listeners, response.Response.ListenerSet...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read TCP listeners failed, reason: %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *GaapService) DescribeUDPListeners(ctx context.Context, proxyId, id, name *string, port *int) (listeners []*gaap.UDPListener, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeUDPListenersRequest()
	request.ProxyId = proxyId
	request.ListenerId = id

	if port != nil {
		request.Port = helper.IntUint64(*port)
	}

	// if port set, name can't use fuzzy search
	if name != nil {
		if port != nil {
			request.ListenerName = name
		} else {
			request.SearchValue = name
		}
	}

	request.Limit = helper.IntUint64(50)
	offset := 0

	// to run loop at least once
	count := 50
	for count == 50 {
		request.Offset = helper.IntUint64(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseGaapClient().DescribeUDPListeners(request)
			if err != nil {
				count = 0

				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err, GAAPInternalError)
			}

			count = len(response.Response.ListenerSet)
			listeners = append(listeners, response.Response.ListenerSet...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read UDP listeners failed, reason: %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *GaapService) ModifyTCPListenerAttribute(
	ctx context.Context,
	proxyId, id string,
	name, scheduler *string,
	healthCheck *bool,
	interval, connectTimeout int,
) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewModifyTCPListenerAttributeRequest()
	request.ProxyId = &proxyId
	request.ListenerId = &id
	request.ListenerName = name
	request.Scheduler = scheduler
	if healthCheck != nil {
		if *healthCheck {
			request.HealthCheck = helper.IntUint64(1)
		} else {
			request.HealthCheck = helper.IntUint64(0)
		}
	}
	request.DelayLoop = helper.IntUint64(interval)
	request.ConnectTimeout = helper.IntUint64(connectTimeout)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyTCPListenerAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify TCP listener attribute failed, reason: %v", logId, err)
		return err
	}

	if err := waitLayer4ListenerReady(ctx, client, id, "TCP"); err != nil {
		log.Printf("[CRITAL]%s modify TCP listener attribute failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) ModifyUDPListenerAttribute(
	ctx context.Context,
	proxyId, id string,
	name, scheduler *string,
) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewModifyUDPListenerAttributeRequest()
	request.ProxyId = &proxyId
	request.ListenerId = &id
	request.ListenerName = name
	request.Scheduler = scheduler

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyUDPListenerAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify UDP listener attribute failed, reason: %v", logId, err)
		return err
	}

	if err := waitLayer4ListenerReady(ctx, client, id, "UDP"); err != nil {
		log.Printf("[CRITAL]%s modify UDP listener attribute failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DeleteLayer4Listener(ctx context.Context, id, proxyId, protocol string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	deleteRequest := gaap.NewDeleteListenersRequest()
	deleteRequest.ProxyId = &proxyId
	deleteRequest.ListenerIds = []*string{&id}
	deleteRequest.Force = helper.IntUint64(0)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		response, err := client.DeleteListeners(deleteRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err, GAAPInternalError)
		}

		// listener may not exist
		if len(response.Response.InvalidStatusListenerSet) > 0 {
			return nil
		}

		// delete failed
		if len(response.Response.OperationFailedListenerSet) > 0 {
			err := fmt.Errorf("api[%s] listener delete failed", deleteRequest.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		// delete successfully
		if len(response.Response.OperationSucceedListenerSet) > 0 {
			return nil
		}

		err = fmt.Errorf("api[%s] listener delete status unknown", deleteRequest.GetAction())
		log.Printf("[CRITAL]%s %v", logId, err)
		return resource.NonRetryableError(err)
	}); err != nil {
		log.Printf("[CRITAL]%s delete listener failed, reason: %v", logId, err)
		return err
	}

	switch protocol {
	case "TCP":
		describeRequest := gaap.NewDescribeTCPListenersRequest()
		describeRequest.ListenerId = &id

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(describeRequest.GetAction())

			response, err := client.DescribeTCPListeners(describeRequest)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && sdkError.Message == fmt.Sprintf("ListenerId(%s) Not Exist.", id)) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
				return retryError(err, GAAPInternalError)
			}

			if len(response.Response.ListenerSet) > 0 {
				err := errors.New("listener still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s delete listener failed, reason: %v", logId, err)
			return err
		}

	case "UDP":
		describeRequest := gaap.NewDescribeUDPListenersRequest()
		describeRequest.ListenerId = &id

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(describeRequest.GetAction())

			response, err := client.DescribeUDPListeners(describeRequest)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && sdkError.Message == fmt.Sprintf("ListenerId(%s) Not Exist.", id)) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
				return retryError(err, GAAPInternalError)
			}

			if len(response.Response.ListenerSet) > 0 {
				err := errors.New("listener still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s delete listener failed, reason: %v", logId, err)
			return err
		}
	}

	return nil
}

func (me *GaapService) CreateSecurityPolicy(ctx context.Context, proxyId, action string) (id string, err error) {
	logId := getLogId(ctx)

	request := gaap.NewCreateSecurityPolicyRequest()
	request.ProxyId = &proxyId
	request.DefaultAction = &action

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().CreateSecurityPolicy(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.PolicyId == nil {
			err := fmt.Errorf("api[%s] security policy id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.PolicyId
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create security policy failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) EnableSecurityPolicy(ctx context.Context, proxyId, policyId string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	enableRequest := gaap.NewOpenSecurityPolicyRequest()
	enableRequest.ProxyId = &proxyId

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(enableRequest.GetAction())

		if _, err := client.OpenSecurityPolicy(enableRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, enableRequest.GetAction(), enableRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s enable security policy failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeSecurityPolicyDetailRequest()
	describeRequest.PolicyId = &policyId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeSecurityPolicyDetail(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "PolicyId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.Status == nil {
			err := fmt.Errorf("api[%s] security policy status is nil", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *response.Response.Status != GAAP_SECURITY_POLICY_BOUND {
			err := errors.New("security policy still binding")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s enable security policy failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DisableSecurityPolicy(ctx context.Context, proxyId, policyId string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	disableRequest := gaap.NewCloseSecurityPolicyRequest()
	disableRequest.ProxyId = &proxyId

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(disableRequest.GetAction())

		if _, err := client.CloseSecurityPolicy(disableRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, disableRequest.GetAction(), disableRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s disable security policy failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeSecurityPolicyDetailRequest()
	describeRequest.PolicyId = &policyId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeSecurityPolicyDetail(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "PolicyId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.Status == nil {
			err := fmt.Errorf("api[%s] security policy status is nil", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *response.Response.Status != GAAP_SECURITY_POLICY_UNBIND {
			err := errors.New("security policy still unbinding")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s disable security policy failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DescribeSecurityPolicy(ctx context.Context, id string) (proxyId, status, action string, exist bool, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeSecurityPolicyDetailRequest()
	request.PolicyId = &id

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().DescribeSecurityPolicyDetail(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "PolicyId")) {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.ProxyId == nil {
			err := fmt.Errorf("api[%s] security policy id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}
		proxyId = *response.Response.ProxyId

		if response.Response.Status == nil {
			err := fmt.Errorf("api[%s] security policy status is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}
		status = *response.Response.Status

		if response.Response.DefaultAction == nil {
			err := fmt.Errorf("api[%s] security policy action is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}
		action = *response.Response.DefaultAction

		exist = true

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read security policy failed, reason: %v", logId, err)
		return "", "", "", false, err
	}

	return
}

func (me *GaapService) DeleteSecurityPolicy(ctx context.Context, id string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	deleteRequest := gaap.NewDeleteSecurityPolicyRequest()
	deleteRequest.PolicyId = &id

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteSecurityPolicy(deleteRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete security policy failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeSecurityPolicyDetailRequest()
	describeRequest.PolicyId = &id

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		_, err := client.DescribeSecurityPolicyDetail(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "PolicyId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err, GAAPInternalError)
		}

		err = errors.New("security policy still exists")
		log.Printf("[DEBUG]%s %v", logId, err)
		return resource.RetryableError(err)
	}); err != nil {
		log.Printf("[CRITAL]%s delete security policy failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) CreateSecurityRule(
	ctx context.Context,
	policyId, name, cidrIp, port, action, protocol string,
) (id string, err error) {
	logId := getLogId(ctx)

	request := gaap.NewCreateSecurityRulesRequest()
	request.PolicyId = &policyId
	request.RuleList = []*gaap.SecurityPolicyRuleIn{
		{
			SourceCidr:    &cidrIp,
			DestPortRange: &port,
			Protocol:      &protocol,
			AliasName:     &name,
			Action:        &action,
		},
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().CreateSecurityRules(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.RuleIdList) != 1 {
			err := fmt.Errorf("api[%s] return %d rule ids", request.GetAction(), len(response.Response.RuleIdList))
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if response.Response.RuleIdList[0] == nil {
			err := fmt.Errorf("api[%s] rule id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.RuleIdList[0]
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create security rule failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) DescribeSecurityRule(ctx context.Context, id string) (securityRule *gaap.SecurityPolicyRuleOut, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeSecurityRulesRequest()
	request.SecurityRuleIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().DescribeSecurityRules(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "SecurityRuleId")) {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		switch len(response.Response.SecurityRuleSet) {
		case 0:
			return nil

		default:
			err := fmt.Errorf("api[%s] return more than 1 security rule", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)

		case 1:
		}

		securityRule = response.Response.SecurityRuleSet[0]

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read security rule failed, reason: %v", logId, err)
		return nil, err
	}

	return
}

func (me *GaapService) ModifySecurityRuleName(ctx context.Context, policyId, ruleId, name string) error {
	logId := getLogId(ctx)

	request := gaap.NewModifySecurityRuleRequest()
	request.PolicyId = &policyId
	request.RuleId = &ruleId
	request.AliasName = &name

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseGaapClient().ModifySecurityRule(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify security rule name failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DeleteSecurityRule(ctx context.Context, policyId, ruleId string) error {
	logId := getLogId(ctx)

	request := gaap.NewDeleteSecurityRulesRequest()
	request.PolicyId = &policyId
	request.RuleIdList = []*string{&ruleId}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseGaapClient().DeleteSecurityRules(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete security rule failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) CreateHTTPListener(ctx context.Context, name, proxyId string, port int) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewCreateHTTPListenerRequest()
	request.ProxyId = &proxyId
	request.ListenerName = &name
	request.Port = helper.IntUint64(port)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.CreateHTTPListener(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.ListenerId == nil {
			err := fmt.Errorf("api[%s] HTTP listener id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.ListenerId
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create HTTP listener failed, reason: %v", logId, err)
		return "", err
	}

	if err := waitLayer7ListenerReady(ctx, client, proxyId, id, "HTTP"); err != nil {
		log.Printf("[CRITAL]%s create HTTP listener failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) CreateHTTPSListener(
	ctx context.Context,
	name, certificateId, forwardProtocol, proxyId string,
	polyClientCertificateIds []string,
	port, authType int,
) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewCreateHTTPSListenerRequest()
	request.ProxyId = &proxyId
	request.CertificateId = &certificateId
	request.ForwardProtocol = &forwardProtocol
	request.ListenerName = &name
	request.Port = helper.IntUint64(port)
	request.AuthType = helper.IntUint64(authType)
	request.PolyClientCertificateIds = helper.Strings(polyClientCertificateIds)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.CreateHTTPSListener(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.ListenerId == nil {
			err := fmt.Errorf("api[%s] HTTPS listener id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.ListenerId
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create HTTPS listener failed, reason: %v", logId, err)
		return "", err
	}

	if err := waitLayer7ListenerReady(ctx, client, proxyId, id, "HTTPS"); err != nil {
		log.Printf("[CRITAL]%s create HTTPS listener failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) DescribeHTTPListeners(
	ctx context.Context,
	proxyId, id, name *string,
	port *int,
) (listeners []*gaap.HTTPListener, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeHTTPListenersRequest()
	request.ProxyId = proxyId
	request.ListenerId = id

	if port != nil {
		request.Port = helper.IntUint64(*port)
	}

	// if port set, name can't use fuzzy search
	if name != nil {
		if port != nil {
			request.ListenerName = name
		} else {
			request.SearchValue = name
		}
	}

	request.Limit = helper.IntUint64(50)
	offset := 0

	// run loop at least once
	count := 50
	for count == 50 {
		request.Offset = helper.IntUint64(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseGaapClient().DescribeHTTPListeners(request)
			if err != nil {
				count = 0

				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			count = len(response.Response.ListenerSet)
			listeners = append(listeners, response.Response.ListenerSet...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *GaapService) DescribeHTTPSListeners(
	ctx context.Context,
	proxyId, listenerId, name *string,
	port *int,
) (listeners []*gaap.HTTPSListener, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeHTTPSListenersRequest()
	request.ProxyId = proxyId
	request.ListenerId = listenerId

	if port != nil {
		request.Port = helper.IntUint64(*port)
	}

	// if port set, name can't use fuzzy search
	if name != nil {
		if port != nil {
			request.ListenerName = name
		} else {
			request.SearchValue = name
		}
	}

	request.Limit = helper.IntUint64(50)
	offset := 0

	// run loop at least once
	count := 50
	for count == 50 {
		request.Offset = helper.IntUint64(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseGaapClient().DescribeHTTPSListeners(request)
			if err != nil {
				count = 0

				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			count = len(response.Response.ListenerSet)
			listeners = append(listeners, response.Response.ListenerSet...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *GaapService) ModifyHTTPListener(ctx context.Context, id, proxyId, name string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewModifyHTTPListenerAttributeRequest()
	request.ListenerId = &id
	request.ListenerName = &name
	request.ProxyId = &proxyId

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyHTTPListenerAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify HTTP listener failed, reason: %v", logId, err)
		return err
	}

	if err := waitLayer7ListenerReady(ctx, client, proxyId, id, "HTTP"); err != nil {
		log.Printf("[CRITAL]%s modify HTTP listener failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) ModifyHTTPSListener(
	ctx context.Context,
	proxyId, id string,
	name, forwardProtocol, certificateId *string,
	polyClientCertificateIds []string,
) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewModifyHTTPSListenerAttributeRequest()
	request.ProxyId = &proxyId
	request.ListenerId = &id
	request.ListenerName = name
	request.ForwardProtocol = forwardProtocol
	request.CertificateId = certificateId
	request.PolyClientCertificateIds = helper.Strings(polyClientCertificateIds)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyHTTPSListenerAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify HTTPS listener failed, reason: %v", logId, err)
		return err
	}

	if err := waitLayer7ListenerReady(ctx, client, proxyId, id, "HTTPS"); err != nil {
		log.Printf("[CRITAL]%s modify HTTPS listener failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DeleteLayer7Listener(ctx context.Context, id, proxyId, protocol string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	deleteRequest := gaap.NewDeleteListenersRequest()
	deleteRequest.ProxyId = &proxyId
	deleteRequest.ListenerIds = []*string{helper.String(id)}
	deleteRequest.Force = helper.IntUint64(0)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		response, err := client.DeleteListeners(deleteRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}

		// listener may not exist
		if len(response.Response.InvalidStatusListenerSet) > 0 {
			return nil
		}

		// delete failed
		if len(response.Response.OperationFailedListenerSet) > 0 {
			err := fmt.Errorf("api[%s] listener delete failed", deleteRequest.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		// delete successfully
		if len(response.Response.OperationSucceedListenerSet) > 0 {
			return nil
		}

		err = fmt.Errorf("api[%s] listener delete status unknown", deleteRequest.GetAction())
		log.Printf("[CRITAL]%s %v", logId, err)
		return resource.NonRetryableError(err)
	}); err != nil {
		log.Printf("[CRITAL]%s delete listener failed, reason: %v", logId, err)
		return err
	}

	switch protocol {
	case "HTTP":
		describeRequest := gaap.NewDescribeHTTPListenersRequest()
		// don't set proxy id it may cause InternalError
		//describeRequest.ProxyId = &proxyId
		describeRequest.ListenerId = &id

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(describeRequest.GetAction())

			response, err := client.DescribeHTTPListeners(describeRequest)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && sdkError.Message == fmt.Sprintf("ListenerId(%s) Not Exist.", id)) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
				return retryError(err, GAAPInternalError)
			}

			if len(response.Response.ListenerSet) > 0 {
				err := errors.New("listener still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s delete listener failed, reason: %v", logId, err)
			return err
		}

	case "HTTPS":
		describeRequest := gaap.NewDescribeHTTPSListenersRequest()
		// don't set proxy id it may cause InternalError
		//describeRequest.ProxyId = &proxyId
		describeRequest.ListenerId = &id

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(describeRequest.GetAction())

			response, err := client.DescribeHTTPSListeners(describeRequest)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && sdkError.Message == fmt.Sprintf("ListenerId(%s) Not Exist.", id)) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
				return retryError(err, GAAPInternalError)
			}

			if len(response.Response.ListenerSet) > 0 {
				err := errors.New("listener still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s delete listener failed, reason: %v", logId, err)
			return err
		}
	}

	return nil
}

func waitLayer4ListenerReady(ctx context.Context, client *gaap.Client, id, protocol string) (err error) {
	logId := getLogId(ctx)

	switch protocol {
	case "TCP":
		request := gaap.NewDescribeTCPListenersRequest()
		request.ListenerId = &id

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := client.DescribeTCPListeners(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err, GAAPInternalError)
			}

			switch len(response.Response.ListenerSet) {
			case 0:
				err := fmt.Errorf("api[%s] return empty TCP listener set", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)

			default:
				err := fmt.Errorf("api[%s] return more than 1 TCP listener", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.NonRetryableError(err)

			case 1:
			}

			listener := response.Response.ListenerSet[0]
			if listener.ListenerStatus == nil {
				err := fmt.Errorf("api[%s] TCP listener status is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *listener.ListenerStatus != GAAP_LISTENER_RUNNING {
				err := errors.New("TCP listener is not ready")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		})

	case "UDP":
		request := gaap.NewDescribeUDPListenersRequest()
		request.ListenerId = &id

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := client.DescribeUDPListeners(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err, GAAPInternalError)
			}

			switch len(response.Response.ListenerSet) {
			case 0:
				err := fmt.Errorf("api[%s] return empty UDP listener set", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)

			default:
				err := fmt.Errorf("api[%s] return more than 1 UDP listener", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.NonRetryableError(err)

			case 1:
			}

			listener := response.Response.ListenerSet[0]
			if listener.ListenerStatus == nil {
				err := fmt.Errorf("api[%s] UDP listener status is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *listener.ListenerStatus != GAAP_LISTENER_RUNNING {
				err := errors.New("UDP listener is not ready")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		})
	}

	return
}

func waitLayer7ListenerReady(ctx context.Context, client *gaap.Client, proxyId, id, protocol string) (err error) {
	logId := getLogId(ctx)

	switch protocol {
	case "HTTP":
		request := gaap.NewDescribeHTTPListenersRequest()
		request.ProxyId = &proxyId
		request.ListenerId = &id

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := client.DescribeHTTPListeners(request)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && sdkError.Message == fmt.Sprintf("ListenerId(%s) Not Exist.", id)) {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			switch len(response.Response.ListenerSet) {
			case 0:
				err := fmt.Errorf("api[%s] return empty HTTP listener set", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)

			default:
				err := fmt.Errorf("api[%s] return more than 1 HTTP listener", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.NonRetryableError(err)

			case 1:
			}

			listener := response.Response.ListenerSet[0]
			if listener.ListenerStatus == nil {
				err := fmt.Errorf("api[%s] HTTP listener status is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *listener.ListenerStatus != GAAP_LISTENER_RUNNING {
				err := errors.New("HTTP listener is not ready")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		})

	case "HTTPS":
		request := gaap.NewDescribeHTTPSListenersRequest()
		request.ProxyId = &proxyId
		request.ListenerId = &id

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := client.DescribeHTTPSListeners(request)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && sdkError.Message == fmt.Sprintf("ListenerId(%s) Not Exist.", id)) {
						return nil
					}
				}
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			switch len(response.Response.ListenerSet) {
			case 0:
				err := fmt.Errorf("api[%s] return empty HTTPS listener set", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)

			default:
				err := fmt.Errorf("api[%s] return more than 1 HTTPS listener", request.GetAction())
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.NonRetryableError(err)

			case 1:
			}

			listener := response.Response.ListenerSet[0]
			if listener.ListenerStatus == nil {
				err := fmt.Errorf("api[%s] HTTPS listener status is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *listener.ListenerStatus != GAAP_LISTENER_RUNNING {
				err := errors.New("HTTPS listener is not ready")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			return nil
		})
	}

	return
}

func (me *GaapService) CreateHTTPDomain(ctx context.Context, listenerId, domain string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	createRequest := gaap.NewCreateDomainRequest()
	createRequest.ListenerId = &listenerId
	createRequest.Domain = &domain

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(createRequest.GetAction())

		if _, err := client.CreateDomain(createRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, createRequest.GetAction(), createRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create HTTP domain failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeRulesRequest()
	describeRequest.ListenerId = &listenerId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeRules(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		for _, rule := range response.Response.DomainRuleSet {
			if rule.Domain == nil {
				err := fmt.Errorf("api[%s] domain is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *rule.Domain == domain {
				return nil
			}
		}

		err = errors.New("domain not found")
		log.Printf("[DEBUG]%s %v", logId, err)
		return resource.RetryableError(err)
	}); err != nil {
		log.Printf("[CRITAL]%s create HTTP domain failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) CreateHTTPSDomain(
	ctx context.Context,
	listenerId, domain, certificateId string,
	polyClientCertificateIds []string,
) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	createRequest := gaap.NewCreateDomainRequest()
	createRequest.ListenerId = &listenerId
	createRequest.Domain = &domain
	createRequest.CertificateId = &certificateId

	for _, polyId := range polyClientCertificateIds {
		createRequest.PolyClientCertificateIds = append(createRequest.PolyClientCertificateIds, helper.String(polyId))
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(createRequest.GetAction())

		if _, err := client.CreateDomain(createRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, createRequest.GetAction(), createRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create HTTPS domain failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeRulesRequest()
	describeRequest.ListenerId = &listenerId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeRules(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		for _, rule := range response.Response.DomainRuleSet {
			if rule.Domain == nil {
				err := fmt.Errorf("api[%s] domain is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *rule.Domain == domain {
				return nil
			}
		}

		err = errors.New("domain not found")
		log.Printf("[DEBUG]%s %v", logId, err)
		return resource.RetryableError(err)
	}); err != nil {
		log.Printf("[CRITAL]%s create HTTPS domain failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) SetAdvancedAuth(
	ctx context.Context,
	listenerId, domain string,
	realserverAuth, basicAuth, gaapAuth bool,
	realserverCertificateIds []string,
	realserverCertificateDomain, basicAuthId, gaapAuthId *string,
) error {
	logId := getLogId(ctx)

	request := gaap.NewSetAuthenticationRequest()
	request.ListenerId = &listenerId
	request.Domain = &domain

	if realserverAuth {
		request.RealServerAuth = helper.IntInt64(1)
	} else {
		request.RealServerAuth = helper.IntInt64(0)
	}

	request.RealServerCertificateDomain = realserverCertificateDomain

	for _, id := range realserverCertificateIds {
		request.PolyRealServerCertificateIds = append(request.PolyRealServerCertificateIds, helper.String(id))
	}

	if basicAuth {
		request.BasicAuth = helper.IntInt64(1)
	} else {
		request.BasicAuth = helper.IntInt64(0)
	}

	request.BasicAuthConfId = basicAuthId

	if gaapAuth {
		request.GaapAuth = helper.IntInt64(1)
	} else {
		request.GaapAuth = helper.IntInt64(0)
	}

	request.GaapCertificateId = gaapAuthId

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseGaapClient().SetAuthentication(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s set advanced auth failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DescribeDomain(ctx context.Context, listenerId, domain string) (domainRet *gaap.DomainRuleSet, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeRulesRequest()
	request.ListenerId = &listenerId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().DescribeRules(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		for _, rule := range response.Response.DomainRuleSet {
			if rule.Domain == nil {
				err := fmt.Errorf("api[%s] domain rule domain is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *rule.Domain == domain {
				domainRet = rule
				break
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read domain failed, reason: %v", logId, err)
		return nil, err
	}

	return
}

func (me *GaapService) ModifyDomainCertificate(
	ctx context.Context,
	listenerId, domain, certificateId string,
	polyClientCertificateIds []string,
) error {
	logId := getLogId(ctx)

	request := gaap.NewModifyCertificateRequest()
	request.ListenerId = &listenerId
	request.Domain = &domain
	request.CertificateId = &certificateId

	for _, polyId := range polyClientCertificateIds {
		request.PolyClientCertificateIds = append(request.PolyClientCertificateIds, helper.String(polyId))
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseGaapClient().ModifyCertificate(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s update domain certificate failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DeleteDomain(ctx context.Context, listenerId, domain string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	deleteRequest := gaap.NewDeleteDomainRequest()
	deleteRequest.ListenerId = &listenerId
	deleteRequest.Domain = &domain
	deleteRequest.Force = helper.IntUint64(0)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteDomain(deleteRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete domain failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeRulesRequest()
	describeRequest.ListenerId = &listenerId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeRules(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err, GAAPInternalError)
		}

		for _, rule := range response.Response.DomainRuleSet {
			if rule.Domain == nil {
				err := fmt.Errorf("api[%s] domain rule domain is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *rule.Domain == domain {
				err := errors.New("domain still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete domain failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) CreateHttpRule(ctx context.Context, httpRule gaapHttpRule) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewCreateRuleRequest()
	request.ListenerId = &httpRule.listenerId
	request.Domain = &httpRule.domain
	request.Path = &httpRule.path
	request.RealServerType = &httpRule.realserverType
	request.Scheduler = &httpRule.scheduler
	if httpRule.healthCheck {
		request.HealthCheck = helper.IntUint64(1)
	} else {
		request.HealthCheck = helper.IntUint64(0)
	}

	request.CheckParams = &gaap.RuleCheckParams{
		DelayLoop:      helper.IntUint64(httpRule.interval),
		ConnectTimeout: helper.IntUint64(httpRule.connectTimeout),
		Path:           &httpRule.healthCheckPath,
		Method:         &httpRule.healthCheckMethod,
		StatusCode:     make([]*uint64, 0, len(httpRule.healthCheckStatusCodes)),
	}
	for _, code := range httpRule.healthCheckStatusCodes {
		request.CheckParams.StatusCode = append(request.CheckParams.StatusCode, helper.IntUint64(code))
	}

	request.ForwardHost = &httpRule.forwardHost
	if httpRule.serverNameIndicationSwitch != "" {
		request.ServerNameIndicationSwitch = &httpRule.serverNameIndicationSwitch
	}
	if httpRule.serverNameIndication != "" {
		request.ServerNameIndication = &httpRule.serverNameIndication
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.CreateRule(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if response.Response.RuleId == nil {
			err := fmt.Errorf("api[%s] HTTP rule id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.RuleId
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create HTTP rule failed, reason: %v", logId, err)
		return "", err
	}

	if err := waitHttpRuleReady(ctx, client, httpRule.listenerId, id); err != nil {
		log.Printf("[CRITAL]%s create HTTP rule failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *GaapService) BindHttpRuleRealservers(ctx context.Context, listenerId, ruleId string, realservers []gaapRealserverBind) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewBindRuleRealServersRequest()
	request.RuleId = &ruleId

	// if realservers is nil, request.RealServerBindSet will bi nil and remove realserver binding
	for _, realserver := range realservers {
		request.RealServerBindSet = append(request.RealServerBindSet, &gaap.RealServerBindSetReq{
			RealServerId:     helper.String(realserver.id),
			RealServerPort:   helper.IntUint64(realserver.port),
			RealServerIP:     helper.String(realserver.ip),
			RealServerWeight: helper.IntUint64(realserver.weight),
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.BindRuleRealServers(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s bind HTTP rule realservers failed, reason: %v", logId, err)
		return err
	}

	if err := waitHttpRuleReady(ctx, client, listenerId, ruleId); err != nil {
		log.Printf("[CRITAL]%s bind HTTP rule realservers failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) DescribeHttpRule(ctx context.Context, id string) (rule *gaap.RuleInfo, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeRulesByRuleIdsRequest()
	request.RuleIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().DescribeRulesByRuleIds(request)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == GAAPResourceNotFound || (sdkErr.Code == "InvalidParameter" && strings.Contains(sdkErr.Message, "ruleId")) {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		if len(response.Response.RuleSet) == 0 {
			return nil
		}

		rule = response.Response.RuleSet[0]

		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func (me *GaapService) ModifyHTTPRuleAttribute(
	ctx context.Context,
	listenerId, ruleId, healthCheckPath, healthCheckMethod, sniSwitch, sni string,
	path, scheduler *string,
	healthCheck bool,
	interval, connectTimeout int,
	healthCheckStatusCodes []int,
) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewModifyRuleAttributeRequest()
	request.ListenerId = &listenerId
	request.RuleId = &ruleId
	request.Path = path
	request.Scheduler = scheduler
	request.ServerNameIndicationSwitch = &sniSwitch
	request.ServerNameIndication = &sni

	if healthCheck {
		request.HealthCheck = helper.IntUint64(1)
	} else {
		request.HealthCheck = helper.IntUint64(0)
	}

	request.CheckParams = &gaap.RuleCheckParams{
		DelayLoop:      helper.IntUint64(interval),
		ConnectTimeout: helper.IntUint64(connectTimeout),
		Path:           &healthCheckPath,
		Method:         &healthCheckMethod,
		StatusCode:     make([]*uint64, 0, len(healthCheckStatusCodes)),
	}
	for _, code := range healthCheckStatusCodes {
		request.CheckParams.StatusCode = append(request.CheckParams.StatusCode, helper.IntUint64(code))
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyRuleAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify HTTP rule attribute failed, reason: %v", logId, err)
		return err
	}

	return waitHttpRuleReady(ctx, client, listenerId, ruleId)
}

func (me *GaapService) DeleteHttpRule(ctx context.Context, listenerId, ruleId string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	deleteRequest := gaap.NewDeleteRuleRequest()
	deleteRequest.ListenerId = &listenerId
	deleteRequest.RuleId = &ruleId
	deleteRequest.Force = helper.IntUint64(1)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteRule(deleteRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete HTTP rule failed, reason: %v", logId, err)
		return err
	}

	describeRequest := gaap.NewDescribeRulesRequest()
	describeRequest.ListenerId = &listenerId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeRules(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err, GAAPInternalError)
		}

		for _, domainRule := range response.Response.DomainRuleSet {
			for _, rule := range domainRule.RuleSet {
				if rule.RuleId == nil {
					err := fmt.Errorf("api[%s] HTTP rule id is nil", describeRequest.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				if *rule.RuleId == ruleId {
					err := errors.New("HTTP rule still exists")
					log.Printf("[DEBUG]%s %v", logId, err)
					return resource.RetryableError(err)
				}
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete HTTP rule failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func waitHttpRuleReady(ctx context.Context, client *gaap.Client, listenerId, ruleId string) error {
	logId := getLogId(ctx)

	request := gaap.NewDescribeRulesRequest()
	request.ListenerId = &listenerId

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.DescribeRules(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		for _, domainRule := range response.Response.DomainRuleSet {
			for _, rule := range domainRule.RuleSet {
				if rule.RuleId == nil {
					err := fmt.Errorf("api[%s] rule id is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				if rule.RuleStatus == nil {
					err := fmt.Errorf("api[%s] rule status is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				if *rule.RuleId == ruleId {
					if *rule.RuleStatus != GAAP_HTTP_RULE_RUNNING {
						err := errors.New("HTTP rule is not ready")
						log.Printf("[DEBUG]%s %v", logId, err)
						return resource.RetryableError(err)
					}
					return nil
				}
			}
		}

		err = fmt.Errorf("api[%s] HTTP rule not found", request.GetAction())
		log.Printf("[DEBUG]%s %v", logId, err)
		return resource.RetryableError(err)
	})
}

func (me *GaapService) DescribeDomains(ctx context.Context, listenerId, domain string) (domains []*gaap.DomainRuleSet, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeRulesRequest()
	request.ListenerId = &listenerId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().DescribeRules(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		for _, domainRule := range response.Response.DomainRuleSet {
			if domainRule.Domain == nil {
				err := fmt.Errorf("api[%s] domain rule domain is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if strings.Contains(*domainRule.Domain, domain) {
				domains = append(domains, domainRule)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read domain failed, reason: %v", logId, err)
		return nil, err
	}

	return
}

func (me *GaapService) DescribeSecurityRules(ctx context.Context, policyId string) (securityRules []*gaap.SecurityPolicyRuleOut, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeSecurityPolicyDetailRequest()
	request.PolicyId = &policyId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseGaapClient().DescribeSecurityPolicyDetail(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "PolicyId")) {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		securityRules = response.Response.RuleList

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read security rule failed, reason: %v", logId, err)
		return nil, err
	}

	return
}

func (me *GaapService) DescribeCertificates(ctx context.Context, id, name *string, certificateType *int) (certificates []*gaap.Certificate, err error) {
	logId := getLogId(ctx)

	request := gaap.NewDescribeCertificatesRequest()

	if certificateType != nil {
		request.CertificateType = helper.IntInt64(*certificateType)
	}

	request.Limit = helper.IntUint64(50)

	offset := 0
	count := 50
	// run loop at least once
	for count == 50 {
		request.Offset = helper.IntUint64(offset)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseGaapClient().DescribeCertificates(request)
			if err != nil {
				count = 0

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			count = len(response.Response.CertificateSet)

			for _, certificate := range response.Response.CertificateSet {
				if certificate.CertificateId == nil {
					err := fmt.Errorf("api[%s] certificate id is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				if id != nil && *certificate.CertificateId != *id {
					continue
				}

				if certificate.CertificateAlias == nil {
					err := fmt.Errorf("api[%s] certificate name is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				// if name set, use fuzzy search
				if name != nil && !strings.Contains(*certificate.CertificateAlias, *name) {
					continue
				}

				if certificate.CertificateType == nil {
					err := fmt.Errorf("api[%s] certificate type is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				if certificate.CreateTime == nil {
					err := fmt.Errorf("api[%s] certificate create time is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				certificates = append(certificates, certificate)
			}

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read gaap certificates failed, reason: %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *GaapService) ModifyHTTPRuleForwardHost(ctx context.Context, listenerId, ruleId, forwardHost string) error {
	logId := getLogId(ctx)
	client := me.client.UseGaapClient()

	request := gaap.NewModifyRuleAttributeRequest()
	request.ListenerId = &listenerId
	request.RuleId = &ruleId
	request.ForwardHost = &forwardHost

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyRuleAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify HTTP rule forward host failed, reason: %v", logId, err)
		return err
	}

	if err := waitHttpRuleReady(ctx, client, listenerId, ruleId); err != nil {
		log.Printf("[CRITAL]%s modify HTTP rule forward host failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *GaapService) CreateDomainErrorPageInfo(ctx context.Context,
	listenerId, domain, body string,
	newErrorCode *int64,
	errorCodes []int,
	clearHeaders []string,
	setHeaders map[string]string,
) (id string, errRet error) {
	client := me.client.UseGaapClient()

	request := gaap.NewCreateDomainErrorPageInfoRequest()
	request.ListenerId = &listenerId
	request.Domain = &domain
	request.Body = &body
	request.NewErrorNo = newErrorCode

	for _, code := range errorCodes {
		request.ErrorNos = append(request.ErrorNos, helper.IntInt64(code))
	}

	request.ClearHeaders = helper.Strings(clearHeaders)

	for k, v := range setHeaders {
		request.SetHeaders = append(request.SetHeaders, &gaap.HttpHeaderParam{
			HeaderName:  helper.String(k),
			HeaderValue: helper.String(v),
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		resp, err := client.CreateDomainErrorPageInfo(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "FailedOperation.DomainAlreadyExisted" {
				return resource.NonRetryableError(helper.WrapErrorf(err, "", sdkError.RequestId, sdkError.Message))
			}

			return retryError(err)
		}

		id = *resp.Response.ErrorPageId

		return nil
	}); err != nil {
		return "", helper.WrapErrorf(err, "", "", "create gaap domain error page info failed")
	}

	return
}

func (me *GaapService) DescribeDomainErrorPageInfo(ctx context.Context, listenerId, domain, id string) (info *gaap.DomainErrorPageInfo, err error) {
	client := me.client.UseGaapClient()

	request := gaap.NewDescribeDomainErrorPageInfoRequest()
	request.ListenerId = &listenerId
	request.Domain = &domain

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		resp, err := client.DescribeDomainErrorPageInfo(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" || (sdkError.Code == "InvalidParameter" && strings.Contains(sdkError.Message, "ListenerId")) {
					return nil
				}
			}

			return retryError(err)
		}

		for _, pageInfo := range resp.Response.ErrorPageSet {
			if *pageInfo.ErrorPageId == id {
				info = pageInfo
				break
			}
		}

		return nil
	}); err != nil {
		return nil, helper.WrapErrorf(err, id, "", "describe domain error page info failed")
	}

	return
}

func (me *GaapService) DescribeDomainErrorPageInfoList(ctx context.Context, listenerId, domain string) (list []*gaap.DomainErrorPageInfo, err error) {
	client := me.client.UseGaapClient()

	request := gaap.NewDescribeDomainErrorPageInfoRequest()
	request.ListenerId = &listenerId
	request.Domain = &domain

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		resp, err := client.DescribeDomainErrorPageInfo(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "ResourceNotFound" {
				return nil
			}

			return retryError(err)
		}

		list = resp.Response.ErrorPageSet

		return nil
	}); err != nil {
		return nil, helper.WrapErrorf(err, "", "", "describe domain error page info list failed")
	}

	return
}

func (me *GaapService) DeleteDomainErrorPageInfo(ctx context.Context, id string) error {
	client := me.client.UseGaapClient()

	request := gaap.NewDeleteDomainErrorPageInfoRequest()
	request.ErrorPageId = &id

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if _, err := client.DeleteDomainErrorPageInfo(request); err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "ResourceNotFound" {
				return nil
			}

			return retryError(err)
		}

		return nil
	}); err != nil {
		return helper.WrapErrorf(err, id, "", "delete domain error page info failed")
	}

	return nil
}
