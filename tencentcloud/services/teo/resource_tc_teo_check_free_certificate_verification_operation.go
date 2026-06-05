package teo

import (
	"fmt"
	"log"
	"strings"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoCheckFreeCertificateVerificationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoCheckFreeCertificateVerificationOperationCreate,
		Read:   resourceTencentCloudTeoCheckFreeCertificateVerificationOperationRead,
		Delete: resourceTencentCloudTeoCheckFreeCertificateVerificationOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The domain name to verify, which is the domain used when applying for a free certificate.",
			},
			"common_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain name to which the certificate is issued when the free certificate is successfully applied.",
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The signature algorithm used by the certificate when the free certificate is successfully applied.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration time of the certificate when the free certificate is successfully applied, in ISO 8601 format.",
			},
		},
	}
}

func resourceTencentCloudTeoCheckFreeCertificateVerificationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_check_free_certificate_verification.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	domain := d.Get("domain").(string)

	request := teov20220901.NewCheckFreeCertificateVerificationRequest()
	request.ZoneId = helper.String(zoneId)
	request.Domain = helper.String(domain)

	var response *teov20220901.CheckFreeCertificateVerificationResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CheckFreeCertificateVerification(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("CheckFreeCertificateVerification API returned empty response")
	}

	d.SetId(strings.Join([]string{zoneId, domain}, tccommon.FILED_SP))

	if response.Response.CommonName != nil {
		_ = d.Set("common_name", *response.Response.CommonName)
	}
	if response.Response.SignatureAlgorithm != nil {
		_ = d.Set("signature_algorithm", *response.Response.SignatureAlgorithm)
	}
	if response.Response.ExpireTime != nil {
		_ = d.Set("expire_time", *response.Response.ExpireTime)
	}

	return resourceTencentCloudTeoCheckFreeCertificateVerificationOperationRead(d, meta)
}

func resourceTencentCloudTeoCheckFreeCertificateVerificationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_check_free_certificate_verification.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoCheckFreeCertificateVerificationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_check_free_certificate_verification.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
