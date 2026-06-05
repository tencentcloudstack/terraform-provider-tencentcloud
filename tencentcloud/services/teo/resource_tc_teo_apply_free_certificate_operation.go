package teo

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoApplyFreeCertificateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoApplyFreeCertificateCreate,
		Read:   resourceTencentCloudTeoApplyFreeCertificateRead,
		Delete: resourceTencentCloudTeoApplyFreeCertificateDelete,
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
				Description: "The target domain for the free certificate application.",
			},
			"verification_method": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The verification method for the free certificate application. Valid values: `http_challenge`, `dns_challenge`.",
			},
			"dns_verification": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "DNS verification information. Returned when `verification_method` is `dns_challenge`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subdomain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host record.",
						},
						"record_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The record type.",
						},
						"record_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The record value.",
						},
					},
				},
			},
			"file_verification": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "File verification information. Returned when `verification_method` is `http_challenge`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL path for file verification.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content of the verification file.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoApplyFreeCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_apply_free_certificate.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	domain := d.Get("domain").(string)
	verificationMethod := d.Get("verification_method").(string)

	request := teov20220901.NewApplyFreeCertificateRequest()
	request.ZoneId = helper.String(zoneId)
	request.Domain = helper.String(domain)
	request.VerificationMethod = helper.String(verificationMethod)

	var response *teov20220901.ApplyFreeCertificateResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		log.Printf("[DEBUG]%s api[%s] request body [%s]\n", logId, request.GetAction(), request.ToJsonString())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ApplyFreeCertificate(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return fmt.Errorf("api[ApplyFreeCertificate] fail, reason: %s", err.Error())
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("api[ApplyFreeCertificate] response is nil")
	}

	log.Printf("[DEBUG]%s api[%s] success, response body [%s]\n", logId, request.GetAction(), response.ToJsonString())

	d.SetId(strings.Join([]string{zoneId, domain}, tccommon.FILED_SP))

	if response.Response.DnsVerification != nil {
		dnsVerificationMap := []map[string]interface{}{
			{
				"subdomain":    response.Response.DnsVerification.Subdomain,
				"record_type":  response.Response.DnsVerification.RecordType,
				"record_value": response.Response.DnsVerification.RecordValue,
			},
		}
		if err := d.Set("dns_verification", dnsVerificationMap); err != nil {
			return err
		}
	}

	if response.Response.FileVerification != nil {
		fileVerificationMap := []map[string]interface{}{
			{
				"path":    response.Response.FileVerification.Path,
				"content": response.Response.FileVerification.Content,
			},
		}
		if err := d.Set("file_verification", fileVerificationMap); err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudTeoApplyFreeCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_apply_free_certificate.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoApplyFreeCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_apply_free_certificate.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
