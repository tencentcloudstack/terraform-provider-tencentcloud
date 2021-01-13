/*
Provide a resource to create a payment SSL.

Example Usage

```hcl
resource "tencentcloud_ssl_pay_certificate" "ssl" {
    product_id = 33
    domain_num = 1
    alias      = "test-ssl"
    project_id = 0
    information {
        csr_type              = "online"
        certificate_domain    = "www.domain.com"
        organization_name     = "test"
        organization_division = "test"
        organization_address  = "test"
        organization_country  = "CN"
        organization_city     = "test"
        organization_region   = "test"
        postal_code           = "0755"
        phone_area_code       = "0755"
        phone_number          = "12345678901"
        verify_type           = "DNS"
        admin_first_name      = "test"
        admin_last_name       = "test"
        admin_phone_num       = "12345678901"
        admin_email           = "test@tencent.com"
        admin_position        = "dev"
        contact_first_name    = "test"
        contact_last_name     = "test"
        contact_email         = "test@tencent.com"
        contact_number        = "12345678901"
        contact_position      = "dev"
    }
}
```

Import

payment SSL instance can be imported, e.g.

```
$ terraform import tencentcloud_ssl_pay_certificate.ssl iPQNn61x#33#1#1
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSSLInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSSLInstanceCreate,
		Read:   resourceTencentCloudSSLInstanceRead,
		Update: resourceTencentCloudSSLInstanceUpdate,
		Delete: resourceTencentCloudSSLInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(3, 42),
				Description: "Certificate commodity ID. Valid value ranges: (3~42). `3` means SecureSite Enhanced Enterprise Edition (EV Pro), " +
					"`4` means SecureSite Enhanced (EV), `5` means SecureSite Enterprise Professional Edition (OV Pro), " +
					"`6` means SecureSite Enterprise (OV), `7` means SecureSite Enterprise Type (OV) wildcard, " +
					"`8` means Geotrust enhanced (EV), `9` means Geotrust enterprise (OV), " +
					"`10` means Geotrust enterprise (OV) wildcard, `11` means TrustAsia domain type multi-domain SSL certificate, " +
					"`12` means TrustAsia domain type ( DV) wildcard, `13` means TrustAsia enterprise wildcard (OV) SSL certificate (D3), " +
					"`14` means TrustAsia enterprise (OV) SSL certificate (D3), `15` means TrustAsia enterprise multi-domain (OV) SSL certificate (D3), " +
					"`16` means TrustAsia Enhanced (EV) SSL Certificate (D3), `17` means TrustAsia Enhanced Multiple Domain (EV) SSL Certificate (D3), " +
					"`18` means GlobalSign Enterprise (OV) SSL Certificate, `19` means GlobalSign Enterprise Wildcard (OV) SSL Certificate, " +
					"`20` means GlobalSign Enhanced (EV) SSL Certificate, `21` means TrustAsia Enterprise Wildcard Multiple Domain (OV) SSL Certificate (D3), " +
					"`22` means GlobalSign Enterprise Multiple Domain (OV) SSL Certificate, `23` means GlobalSign Enterprise Multiple Wildcard Domain name (OV) SSL certificate, " +
					"`24` means GlobalSign enhanced multi-domain (EV) SSL certificate, `25` means Wotrus domain type certificate, " +
					"`26` means Wotrus domain type multi-domain certificate, `27` means Wotrus domain type wildcard certificate, " +
					"`28` means Wotrus enterprise type certificate, `29` means Wotrus enterprise multi-domain certificate, " +
					"`30` means Wotrus enterprise wildcard certificate, `31` means Wotrus enhanced certificate, " +
					"`32` means Wotrus enhanced multi-domain certificate, `33` means DNSPod national secret domain name certificate, " +
					"`34` means DNSPod national secret domain name certificate Multi-domain certificate, `35` means DNSPod national secret domain name wildcard certificate, " +
					"`37` means DNSPod national secret enterprise certificate, `38` means DNSPod national secret enterprise multi-domain certificate, " +
					"`39` means DNSPod national secret enterprise wildcard certificate, `40` means DNSPod national secret increase Strong certificate, " +
					"`41` means DNSPod national secret enhanced multi-domain certificate, `42` means TrustAsia domain-type wildcard multi-domain certificate.",
			},
			"domain_num": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Number of domain names included in the certificate.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Certificate period, currently only supports 1 year certificate purchase.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of project.",
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Remark name.",
			},
			// ssl information
			"information": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "Certificate information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"csr_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CsrTypeOnline,
							ForceNew:     true,
							ValidateFunc: validateAllowedStringValue(CsrTypeArr),
							Description: "CSR generation method. Valid values: `online`, `parse`. " +
								"`online` means online generation, `parse` means manual upload.",
						},
						"certificate_domain": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Domain name for binding certificate.",
						},
						"organization_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Company name.",
						},
						"organization_division": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Department name.",
						},
						"organization_address": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Company address.",
						},
						"organization_country": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Country name, such as China: CN.",
						},
						"organization_city": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Company city.",
						},
						"organization_region": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The province where the company is located.",
						},
						"postal_code": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Company postal code.",
						},
						"phone_area_code": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Company landline area code.",
						},
						"phone_number": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Company landline number.",
						},
						"verify_type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateAllowedStringValue(VerifyType),
							Description: "Certificate verification method. Valid values: `DNS_AUTO`, `DNS`, `FILE`. " +
								"`DNS_AUTO` means automatic DNS verification, this verification type is only supported for " +
								"domain names resolved by Tencent Cloud and the resolution status is normal, " +
								"`DNS` means manual DNS verification, `FILE` means file verification.",
						},
						"admin_first_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The first name of the administrator.",
						},
						"admin_last_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The last name of the administrator.",
						},
						"admin_phone_num": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Manager mobile phone number.",
						},
						"admin_email": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The administrator's email address.",
						},
						"admin_position": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Manager position.",
						},
						"contact_first_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Contact first name.",
						},
						"contact_last_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Contact last name.",
						},
						"contact_email": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Contact email address.",
						},
						"contact_number": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Contact phone number.",
						},
						"contact_position": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Contact position.",
						},
						"csr_content": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "CSR content uploaded.",
						},
						"domain_list": {
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Array of uploaded domain names, multi-domain certificates can be uploaded.",
						},
						"key_password": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Private key password.",
						},
					},
				},
			},
			// computed
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returned certificate ID.",
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Order ID returned.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "SSL certificate status.",
			},
		},
	}
}

func resourceTencentCloudSSLInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_pay_certificate.create")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		sslService    = SSLService{client: meta.(*TencentCloudClient).apiV3Conn}
		certificateId string
		err           error
	)

	request := ssl.NewCreateCertificateRequest()
	request.ProductId = helper.Int64(int64(d.Get("product_id").(int)))
	request.DomainNum = helper.Int64(int64(d.Get("domain_num").(int)))
	request.TimeSpan = helper.Int64(int64(d.Get("time_span").(int)))

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		certificateId, _, err = sslService.CreateCertificate(ctx, request)
		if err != nil {
			if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
				code := sdkError.GetCode()
				if code == InvalidParam || code == InvalidParameter || code == InvalidParameterValue {
					return resource.NonRetryableError(sdkError)
				}
			}
			return retryError(err)
		}
		return nil
	}); err != nil {
		return err
	}
	d.SetId(strings.Join([]string{certificateId,
		strconv.FormatInt(*request.ProductId, 10),
		strconv.FormatInt(*request.DomainNum, 10),
		strconv.FormatInt(*request.TimeSpan, 10)}, FILED_SP))

	if alias, ok := d.GetOk("alias"); ok {
		aliasRequest := ssl.NewModifyCertificateAliasRequest()
		aliasRequest.CertificateId = helper.String(certificateId)
		aliasRequest.Alias = helper.String(alias.(string))
		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err := sslService.ModifyCertificateAlias(ctx, aliasRequest); err != nil {
				if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
					code := sdkError.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkError)
					}
				}
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if projectId, ok := d.GetOk("project_id"); ok {
		projectRequest := ssl.NewModifyCertificateProjectRequest()
		projectRequest.CertificateIdList = []*string{
			helper.String(certificateId),
		}
		projectRequest.ProjectId = helper.Uint64(uint64(projectId.(int)))

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err := sslService.ModifyCertificateProject(ctx, projectRequest); err != nil {
				if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
					code := sdkError.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkError)
					}
				}
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	infoRequest := getSubmitInfoRequest(d)
	infoRequest.CertificateId = helper.String(certificateId)
	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err = sslService.SubmitCertificateInformation(ctx, infoRequest); err != nil {
			if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
				code := sdkError.GetCode()
				if code == InvalidParam || code == CertificateNotFound {
					return resource.NonRetryableError(sdkError)
				}
			}
			return retryError(err)
		}
		return nil
	}); err != nil {
		return err
	}

	commitInfoRequest := ssl.NewCommitCertificateInformationRequest()
	commitInfoRequest.CertificateId = helper.String(certificateId)
	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err = sslService.CommitCertificateInformation(ctx, commitInfoRequest); err != nil {
			if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
				code := sdkError.GetCode()
				if code == InvalidParam || code == CertificateNotFound || code == CertificateInvalid {
					return resource.NonRetryableError(sdkError)
				}
			}
			return retryError(err)
		}
		return nil
	}); err != nil {
		return err
	}

	return resourceTencentCloudSSLInstanceRead(d, meta)
}

func resourceTencentCloudSSLInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_pay_certificate.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		sslService = SSLService{client: meta.(*TencentCloudClient).apiV3Conn}
		id         = d.Id()
		err        error
		response   *ssl.DescribeCertificateDetailResponse
	)
	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("ids param is error. id:  %s", id)
	}
	request := ssl.NewDescribeCertificateDetailRequest()
	request.CertificateId = helper.String(ids[0])

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err = sslService.DescribeCertificateDetail(ctx, request)
		if err != nil {
			return retryError(err, InternalError)
		}
		if response == nil || response.Response == nil {
			err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return retryError(err)
		}
		return nil
	}); err != nil {
		return err
	}
	if response.Response.CertificateId == nil || response.Response.ProjectId == nil {
		d.SetId("")
		return nil
	}

	var productId, domainNum, timeSpan, projectId int64
	if productId, err = strconv.ParseInt(ids[1], 10, 64); err != nil {
		return err
	}
	if domainNum, err = strconv.ParseInt(ids[2], 10, 64); err != nil {
		return err
	}
	if timeSpan, err = strconv.ParseInt(ids[3], 10, 64); err != nil {
		return err
	}
	if projectId, err = strconv.ParseInt(*response.Response.ProjectId, 10, 64); err != nil {
		return err
	}

	_ = d.Set("product_id", productId)
	_ = d.Set("domain_num", domainNum)
	_ = d.Set("time_span", timeSpan)
	_ = d.Set("project_id", projectId)
	_ = d.Set("alias", response.Response.Alias)
	_ = d.Set("certificate_id", response.Response.CertificateId)
	_ = d.Set("order_id", response.Response.OrderId)
	if response.Response.Status != nil {
		_ = d.Set("status", response.Response.Status)
	}
	if response.Response.SubmittedData != nil {
		setSubmitInfo(d, response.Response.SubmittedData)
	}

	return nil
}

func resourceTencentCloudSSLInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_pay_certificate.update")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		id         = d.Id()
		sslService = SSLService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("ids param is error. id:  %s", id)
	}

	d.Partial(true)
	if d.HasChange("alias") {
		aliasRequest := ssl.NewModifyCertificateAliasRequest()
		aliasRequest.CertificateId = helper.String(ids[0])
		_, alias := d.GetChange("alias")
		aliasRequest.Alias = helper.String(alias.(string))

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err := sslService.ModifyCertificateAlias(ctx, aliasRequest); err != nil {
				if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
					code := sdkError.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkError)
					}
				}
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
		d.SetPartial("alias")
	}
	if d.HasChange("project_id") {
		projectRequest := ssl.NewModifyCertificateProjectRequest()
		projectRequest.CertificateIdList = []*string{
			helper.String(ids[0]),
		}
		_, projectId := d.GetChange("project_id")
		projectRequest.ProjectId = helper.Uint64(uint64(projectId.(int)))

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err := sslService.ModifyCertificateProject(ctx, projectRequest); err != nil {
				if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
					code := sdkError.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkError)
					}
				}
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
		d.SetPartial("project_id")
	}
	d.Partial(false)

	return resourceTencentCloudSSLInstanceRead(d, meta)
}

func resourceTencentCloudSSLInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_pay_certificate.delete")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		id         = d.Id()
		sslService = SSLService{client: meta.(*TencentCloudClient).apiV3Conn}
		err        error
	)
	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("ids param is error. id:  %s", id)
	}

	request := ssl.NewCancelCertificateOrderRequest()
	request.CertificateId = helper.String(ids[0])

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err = sslService.CancelCertificateOrder(ctx, request); err != nil {
			if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
				code := sdkError.GetCode()
				if code == InvalidParam || code == CertificateNotFound {
					return resource.NonRetryableError(sdkError)
				}
			}
			return retryError(err)
		}
		return nil
	})
	return err
}

func getSubmitInfoRequest(d *schema.ResourceData) *ssl.SubmitCertificateInformationRequest {
	infos := d.Get("information").([]interface{})
	request := ssl.NewSubmitCertificateInformationRequest()

	for _, v := range infos {
		info := v.(map[string]interface{})
		if csrType, ok := info["csr_type"]; ok {
			request.CsrType = helper.String(csrType.(string))
		}
		request.OrganizationName = helper.String(info["organization_name"].(string))
		request.OrganizationDivision = helper.String(info["organization_division"].(string))
		request.OrganizationAddress = helper.String(info["organization_address"].(string))
		request.OrganizationCountry = helper.String(info["organization_country"].(string))
		request.OrganizationCity = helper.String(info["organization_city"].(string))
		request.OrganizationRegion = helper.String(info["organization_region"].(string))
		request.PostalCode = helper.String(info["postal_code"].(string))
		request.VerifyType = helper.String(info["verify_type"].(string))
		request.PhoneNumber = helper.String(info["phone_number"].(string))
		request.PhoneAreaCode = helper.String(info["phone_area_code"].(string))
		request.AdminFirstName = helper.String(info["admin_first_name"].(string))
		request.AdminLastName = helper.String(info["admin_last_name"].(string))
		request.AdminPhoneNum = helper.String(info["admin_phone_num"].(string))
		request.AdminEmail = helper.String(info["admin_email"].(string))
		request.AdminPosition = helper.String(info["admin_position"].(string))
		request.ContactFirstName = helper.String(info["contact_first_name"].(string))
		request.ContactLastName = helper.String(info["contact_last_name"].(string))
		request.ContactEmail = helper.String(info["contact_email"].(string))
		request.ContactNumber = helper.String(info["contact_number"].(string))
		request.ContactPosition = helper.String(info["contact_position"].(string))
		request.CertificateDomain = helper.String(info["certificate_domain"].(string))
		if csrContent, ok := info["csr_content"]; ok {
			request.CsrContent = helper.String(csrContent.(string))
		}
		if domainSet, ok := info["domain_list"]; ok {
			domains := domainSet.(*schema.Set).List()
			domainList := make([]*string, len(domains))
			for index, domain := range domains {
				domainList[index] = helper.String(domain.(string))
			}
			request.DomainList = domainList
		}
		if keyPassword, ok := info["key_password"]; ok {
			request.KeyPassword = helper.String(keyPassword.(string))
		}
	}
	return request
}

func setSubmitInfo(d *schema.ResourceData, info *ssl.SubmittedData) {
	infos := make([]map[string]interface{}, 1)
	infos[1] = map[string]interface{}{
		"csr_type":              info.CsrType,
		"organization_name":     info.OrganizationName,
		"organization_division": info.OrganizationDivision,
		"organization_address":  info.OrganizationAddress,
		"organization_country":  info.OrganizationCountry,
		"organization_city":     info.OrganizationCity,
		"organization_region":   info.OrganizationRegion,
		"postal_code":           info.PostalCode,
		"phone_area_code":       info.PhoneAreaCode,
		"phone_number":          info.PhoneNumber,
		"verify_type":           info.VerifyType,
		"admin_first_name":      info.AdminFirstName,
		"admin_last_name":       info.AdminLastName,
		"admin_phone_num":       info.AdminPhoneNum,
		"admin_email":           info.AdminEmail,
		"admin_position":        info.AdminPosition,
		"contact_first_name":    info.ContactFirstName,
		"contact_last_name":     info.ContactLastName,
		"contact_email":         info.ContactEmail,
		"contact_number":        info.ContactNumber,
		"contact_position":      info.ContactPosition,
		"csr_content":           info.CsrContent,
		"certificate_domain":    info.CertificateDomain,
		"domain_list":           info.DomainList,
		"key_password":          info.KeyPassword,
	}
	_ = d.Set("information", info)
}
