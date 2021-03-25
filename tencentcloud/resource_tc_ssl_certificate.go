/*
Provides a resource to create a SSL certificate.

Example Usage

```hcl
resource "tencentcloud_ssl_certificate" "foo" {
  name       = "test-ssl-certificate"
  type       = "CA"
  project_id = 0
  cert       = "-----BEGIN CERTIFICATE-----\nMIIERzCCAq+gAwIBAgIBAjANBgkqhkiG9w0BAQsFADAoMQ0wCwYDVQQDEwR0ZXN0\nMRcwFQYDVQQKEw50ZXJyYWZvcm0gdGVzdDAeFw0xOTA4MTMwMzE5MzlaFw0yOTA4\nMTAwMzE5MzlaMC4xEzARBgNVBAMTCnNlcnZlciBzc2wxFzAVBgNVBAoTDnRlcnJh\nZm9ybS10ZXN0MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA1Ryp+DKK\nSNFKZsPtwfR+jzOnQ8YFieIKYgakV688d8YgpolenbmeEPrzT87tunFD7G9f6ALG\nND8rj7npj0AowxhOL/h/v1D9u0UsIaj5i2GWJrqNAhGLaxWiEB/hy5WOiwxDrGei\ngQqJkFM52Ep7G1Yx7PHJmKFGwN9FhIsFi1cNZfVRopZuCe/RMPNusNVZaIi+qcEf\nfsE1cmfmuSlG3Ap0RKOIyR0ajDEzqZn9/0R7VwWCF97qy8TNYk94K/1tq3zyhVzR\nZ83xOSfrTqEfb3so3AU2jyKgYdwr/FZS72VCHS8IslgnqJW4izIXZqgIKmHaRZtM\nN4jUloi6l/6lktt6Lsgh9xECecxziSJtPMaog88aC8HnMqJJ3kScGCL36GYG+Kaw\n5PnDlWXBaeiDe8z/eWK9+Rr2M+rhTNxosAVGfDJyxAXyiX49LQ0v7f9qzwc/0JiD\nbvsUv1cm6OgpoEMP9SXqqBdwGqeKbD2/2jlP48xlYP6l1SoJG3GgZ8dbAgMBAAGj\ndjB0MAwGA1UdEwEB/wQCMAAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0PAQH/\nBAUDAweAADAdBgNVHQ4EFgQULwWKBQNLL9s3cb3tTnyPVg+mpCMwHwYDVR0jBBgw\nFoAUKwfrmq791mY831S6UHARHtgYnlgwDQYJKoZIhvcNAQELBQADggGBAMo5RglS\nAHdPgaicWJvmvjjexjF/42b7Rz4pPfMjYw6uYO8He/f4UZWv5CZLrbEe7MywaK3y\n0OsfH8AhyN29pv2x8g9wbmq7omZIOZ0oCAGduEXs/A/qY/hFaCohdkz/IN8qi6JW\nVXreGli3SrpcHFchSwHTyJEXgkutcGAsOvdsOuVSmplOyrkLHc8uUe8SG4j8kGyg\nEzaszFjHkR7g1dVyDVUedc588mjkQxYeAamJgfkgIhljWKMa2XzkVMcVfQHfNpM1\nn+bu8SmqRt9Wma2bMijKRG/Blm756LoI+skY+WRZmlDnq8zj95TT0vceGP0FUWh5\nhKyiocABmpQs9OK9HMi8vgSWISP+fYgkm/bKtKup2NbZBoO5/VL2vCEPInYzUhBO\njCbLMjNjtM5KriCaR7wDARgHiG0gBEPOEW1PIjZ9UOH+LtIxbNZ4eEIIINLHnBHf\nL+doVeZtS/gJc4G4Adr5HYuaS9ZxJ0W2uy0eQlOHzjyxR6Mf/rpnilJlcQ==\n-----END CERTIFICATE-----"
}
```

Import

ssl certificate can be imported using the id, e.g.

```
  $ terraform import tencentcloud_ssl_certificate.cert GjTNRoK7
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCertificateCreate,
		Read:   resourceTencentCloudSslCertificateRead,
		Update: resourceTencentCloudSslCertificateUpdate,
		Delete: resourceTencentCloudSslCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the SSL certificate.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SSL_CERT_TYPE),
				ForceNew:     true,
				Description:  "Type of the SSL certificate. Valid values: `CA` and `SVR`.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project ID of the SSL certificate. Default is `0`.",
			},
			"cert": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Content of the SSL certificate. Not allowed newline at the start and end.",
				ValidateFunc: func(v interface{}, k string) (wss []string, errs []error) {
					value := v.(string)
					if strings.HasPrefix(value, "\n") {
						errs = append(errs, errors.New("cert can't have \\n prefix"))
						return
					}

					if strings.HasSuffix(value, "\n") {
						errs = append(errs, errors.New("cert can't have \\n suffix"))
					}
					return
				},
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Key of the SSL certificate and required when certificate type is `SVR`. Not allowed newline at the start and end.",
				ValidateFunc: func(v interface{}, k string) (wss []string, errs []error) {
					value := v.(string)
					if strings.HasPrefix(value, "\n") {
						errs = append(errs, errors.New("key can't have \\n prefix"))
						return
					}

					if strings.HasSuffix(value, "\n") {
						errs = append(errs, errors.New("key can't have \\n suffix"))
					}
					return
				},
			},

			// computed
			"product_zh_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate authority.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary domain of the SSL certificate.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the SSL certificate.",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Beginning time of the SSL certificate.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ending time of the SSL certificate.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the SSL certificate.",
			},
			"subject_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "ALL domains included in the SSL certificate. Including the primary domain name.",
			},
		},
	}
}

func resourceTencentCloudSslCertificateCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_certificate.create")()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		sslService       = SSLService{client: m.(*TencentCloudClient).apiV3Conn}
		outErr, inErr    error
		id               string
		describeResponse *ssl.DescribeCertificateDetailResponse
	)

	request := ssl.NewUploadCertificateRequest()
	request.CertificatePublicKey = helper.String(d.Get("cert").(string))
	request.CertificateType = helper.String(d.Get("type").(string))
	request.ProjectId = helper.Uint64(uint64(d.Get("project_id").(int)))
	request.Alias = helper.String(d.Get("name").(string))
	if raw, ok := d.GetOk("key"); ok {
		request.CertificatePrivateKey = helper.String(raw.(string))
	}
	if *request.CertificateType == "SVR" && (request.CertificatePrivateKey == nil || *request.CertificatePrivateKey == "") {
		return errors.New("when type is SVR, key can't be empty")
	}

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		id, inErr = sslService.UploadCertificate(ctx, request)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s create certificate failed, reason: %v", logId, outErr)
		return outErr
	}

	describeRequest := ssl.NewDescribeCertificateDetailRequest()
	describeRequest.CertificateId = &id
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
		if inErr != nil {
			return retryError(inErr)
		}
		if describeResponse == nil || describeResponse.Response == nil {
			err := fmt.Errorf("TencentCloud SDK %s return empty response", describeRequest.GetAction())
			return retryError(err)
		}
		if describeResponse.Response.Status == nil {
			err := fmt.Errorf("api[%s] certificate status is nil", describeRequest.GetAction())
			return resource.NonRetryableError(err)
		}

		if *describeResponse.Response.Status != SSL_STATUS_AVAILABLE {
			err := fmt.Errorf("certificate is not available, status is %d", *describeResponse.Response.Status)
			return resource.RetryableError(err)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s create certificate failed, reason: %v", logId, outErr)
		return outErr
	}
	d.SetId(id)

	return resourceTencentCloudSslCertificateRead(d, m)
}

func resourceTencentCloudSslCertificateRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_certificate.read")()
	defer inconsistentCheck(d, m)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		sslService       = SSLService{client: m.(*TencentCloudClient).apiV3Conn}
		outErr, inErr    error
		id               = d.Id()
		describeResponse *ssl.DescribeCertificateDetailResponse
	)

	describeRequest := ssl.NewDescribeCertificateDetailRequest()
	describeRequest.CertificateId = &id
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
		if inErr != nil {
			if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == CertificateNotFound {
					return nil
				}
			}
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s read certificate failed, reason: %v", logId, outErr)
		return outErr
	}

	if describeResponse == nil || describeResponse.Response == nil || describeResponse.Response.CertificateId == nil {
		d.SetId("")
		return nil
	}

	certificate := describeResponse.Response
	if nilNames := CheckNil(certificate, map[string]string{
		"Alias":                "name",
		"CertificateType":      "type",
		"ProjectId":            "project id",
		"CertificatePublicKey": "cert",
		"ProductZhName":        "product zh name",
		"Domain":               "domain",
		"Status":               "status",
		"CertBeginTime":        "begin time",
		"CertEndTime":          "end time",
		"InsertTime":           "create time",
	}); len(nilNames) > 0 {
		return fmt.Errorf("certificate %v are nil", nilNames)
	}

	_ = d.Set("name", certificate.Alias)
	_ = d.Set("type", certificate.CertificateType)
	projectId, err := strconv.Atoi(*certificate.ProjectId)
	if err != nil {
		return err
	}
	_ = d.Set("project_id", projectId)
	_ = d.Set("cert", strings.TrimRight(*certificate.CertificatePublicKey, "\n"))
	_ = d.Set("product_zh_name", certificate.ProductZhName)
	_ = d.Set("domain", certificate.Domain)
	_ = d.Set("status", certificate.Status)
	_ = d.Set("begin_time", certificate.CertBeginTime)
	_ = d.Set("end_time", certificate.CertEndTime)
	_ = d.Set("create_time", certificate.InsertTime)

	subjectAltNames := make([]string, 0, len(certificate.SubjectAltName))
	for _, subjectAltName := range certificate.SubjectAltName {
		subjectAltNames = append(subjectAltNames, *subjectAltName)
	}
	_ = d.Set("subject_names", subjectAltNames)

	return nil
}

func resourceTencentCloudSslCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_certificate.update")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		id         = d.Id()
		sslService = SSLService{client: m.(*TencentCloudClient).apiV3Conn}
	)

	d.Partial(true)
	if d.HasChange("name") {
		aliasRequest := ssl.NewModifyCertificateAliasRequest()
		aliasRequest.CertificateId = helper.String(id)
		_, alias := d.GetChange("name")
		aliasRequest.Alias = helper.String(alias.(string))

		if outErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if inErr := sslService.ModifyCertificateAlias(ctx, aliasRequest); inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					code := sdkErr.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkErr)
					}
				}
				return retryError(inErr)
			}
			return nil
		}); outErr != nil {
			return outErr
		}
		d.SetPartial("name")
	}
	if d.HasChange("project_id") {
		projectRequest := ssl.NewModifyCertificateProjectRequest()
		projectRequest.CertificateIdList = []*string{
			helper.String(id),
		}
		_, projectId := d.GetChange("project_id")
		projectRequest.ProjectId = helper.Uint64(uint64(projectId.(int)))

		if outErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if inErr := sslService.ModifyCertificateProject(ctx, projectRequest); inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					code := sdkErr.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkErr)
					}
				}
				return retryError(inErr)
			}
			return nil
		}); outErr != nil {
			return outErr
		}
		d.SetPartial("project_id")
	}
	d.Partial(false)
	return resourceTencentCloudSslCertificateRead(d, m)
}

func resourceTencentCloudSslCertificateDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_certificate.delete")()
	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		sslService    = SSLService{client: m.(*TencentCloudClient).apiV3Conn}
		outErr, inErr error
		id            = d.Id()
		deleteResult  bool
	)
	request := ssl.NewDeleteCertificateRequest()
	request.CertificateId = helper.String(id)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		deleteResult, inErr = sslService.DeleteCertificate(ctx, request)
		if inErr != nil {
			return retryError(inErr)
		}
		if !deleteResult {
			return resource.NonRetryableError(errors.New("failed to delete certificate"))
		}
		return nil
	})

	if outErr != nil {
		log.Printf("[CRITAL]%s delete SSL certificate failed, reason:%+v", logId, outErr)
		return outErr
	}
	return nil
}
