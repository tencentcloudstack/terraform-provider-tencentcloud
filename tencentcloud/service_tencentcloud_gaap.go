package tencentcloud

import (
	"context"
	"fmt"
	"log"

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
