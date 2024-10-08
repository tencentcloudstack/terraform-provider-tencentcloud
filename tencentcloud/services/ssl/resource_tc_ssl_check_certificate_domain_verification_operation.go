package ssl

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslCheckCertificateDomainVerificationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCheckCertificateDomainVerificationCreate,
		Read:   resourceTencentCloudSslCheckCertificateDomainVerificationRead,
		Delete: resourceTencentCloudSslCheckCertificateDomainVerificationDelete,

		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The certificate ID.",
			},
			// computed
			"verification_results": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Domain name verification results.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name.",
						},
						"verify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain Verify Type.",
						},
						"local_check": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Local inspection results.",
						},
						"ca_check": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CA inspection results.",
						},
						"local_check_fail_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Check the reason for the failure.",
						},
						"check_value": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detected values.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"frequently": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether frequent requests.",
						},
						"issued": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether issued.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSslCheckCertificateDomainVerificationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_check_certificate_domain_verification_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = ssl.NewCheckCertificateDomainVerificationRequest()
		response      = []*ssl.DomainValidationResult{}
		certificateId string
	)

	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
		certificateId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().CheckCertificateDomainVerification(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			err := fmt.Errorf("[DEBUG]%s Check certificate domain verification failed.\n", logId)
			return tccommon.RetryError(err)
		}

		response = result.Response.VerificationResults
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s Check certificate domain verification failed, reason: %+v", logId, err)
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(response))
	for _, item := range response {
		tmpObject := make(map[string]interface{})
		if item.Domain != nil {
			tmpObject["domain"] = item.Domain
		}

		if item.VerifyType != nil {
			tmpObject["verify_type"] = item.VerifyType
		}

		if item.LocalCheck != nil {
			tmpObject["local_check"] = item.LocalCheck
		}

		if item.CaCheck != nil {
			tmpObject["ca_check"] = item.CaCheck
		}

		if item.LocalCheckFailReason != nil {
			tmpObject["local_check_fail_reason"] = item.LocalCheckFailReason
		}

		if item.CheckValue != nil {
			tmpValueList := make([]string, 0, len(item.CheckValue))
			for _, v := range item.CheckValue {
				tmpValueList = append(tmpValueList, *v)
			}

			tmpObject["check_value"] = tmpValueList
		}

		if item.Frequently != nil {
			tmpObject["frequently"] = item.Frequently
		}

		if item.Issued != nil {
			tmpObject["issued"] = item.Issued
		}

		tmpList = append(tmpList, tmpObject)
	}

	_ = d.Set("verification_results", tmpList)

	d.SetId(certificateId)

	return resourceTencentCloudSslCheckCertificateDomainVerificationRead(d, meta)
}

func resourceTencentCloudSslCheckCertificateDomainVerificationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_check_certificate_domain_verification_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslCheckCertificateDomainVerificationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_check_certificate_domain_verification_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
