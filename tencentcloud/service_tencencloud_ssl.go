package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss/v20180426"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
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
	createRequest.ProjectId = stringToPointer(strconv.Itoa(projectId))
	createRequest.ModuleType = stringToPointer(SSL_MODULE_TYPE)
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
	describeRequest.ModuleType = stringToPointer(SSL_MODULE_TYPE)
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
	request.ModuleType = stringToPointer(SSL_MODULE_TYPE)
	request.SearchKey = name
	request.Id = id
	request.CertType = certType
	request.WithCert = stringToPointer(SSL_WITH_CERT)

	var offset uint64

	request.Offset = &offset
	request.Limit = intToPointer(20)

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
	deleteRequest.ModuleType = stringToPointer(SSL_MODULE_TYPE)
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
	describeRequest.ModuleType = stringToPointer(SSL_MODULE_TYPE)
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
