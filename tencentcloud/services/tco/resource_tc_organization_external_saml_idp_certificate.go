package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organizationv20210331 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationExternalSamlIdpCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationExternalSamlIdpCertificateCreate,
		Read:   resourceTencentCloudOrganizationExternalSamlIdpCertificateRead,
		Delete: resourceTencentCloudOrganizationExternalSamlIdpCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Space ID.",
			},

			"x509_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "X509 certificate in PEM format, provided by the SAML identity provider.",
			},

			// computed
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate ID.",
			},

			"serial_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate serial number.",
			},

			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate issuer.",
			},

			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Certificate version.",
			},

			"signature_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate signature algorithm.",
			},

			"not_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate expiration date.",
			},

			"not_before": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate creation date.",
			},
		},
	}
}

func resourceTencentCloudOrganizationExternalSamlIdpCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_idp_certificate.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = organizationv20210331.NewAddExternalSAMLIdPCertificateRequest()
		response      = organizationv20210331.NewAddExternalSAMLIdPCertificateResponse()
		zoneId        string
		certificateId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("x509_certificate"); ok {
		request.X509Certificate = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddExternalSAMLIdPCertificateWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Add external saml idp certificate failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create organization external saml idp certificate failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.CertificateId == nil {
		return fmt.Errorf("CertificateId is nil.")
	}

	certificateId = *response.Response.CertificateId
	d.SetId(strings.Join([]string{zoneId, certificateId}, tccommon.FILED_SP))
	return resourceTencentCloudOrganizationExternalSamlIdpCertificateRead(d, meta)
}

func resourceTencentCloudOrganizationExternalSamlIdpCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_idp_certificate.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	certificateId := idSplit[1]

	respData, err := service.DescribeOrganizationExternalSamlIdpCertificateById(ctx, zoneId, certificateId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_organization_external_saml_idp_certificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.X509Certificate != nil {
		_ = d.Set("x509_certificate", respData.X509Certificate)
	}

	if respData.CertificateId != nil {
		_ = d.Set("certificate_id", respData.CertificateId)
	}

	if respData.SerialNumber != nil {
		_ = d.Set("serial_number", respData.SerialNumber)
	}

	if respData.Issuer != nil {
		_ = d.Set("issuer", respData.Issuer)
	}

	if respData.Version != nil {
		_ = d.Set("version", respData.Version)
	}

	if respData.SignatureAlgorithm != nil {
		_ = d.Set("signature_algorithm", respData.SignatureAlgorithm)
	}

	if respData.NotAfter != nil {
		_ = d.Set("not_after", respData.NotAfter)
	}

	if respData.NotBefore != nil {
		_ = d.Set("not_before", respData.NotBefore)
	}

	return nil
}

func resourceTencentCloudOrganizationExternalSamlIdpCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_idp_certificate.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = organizationv20210331.NewRemoveExternalSAMLIdPCertificateRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	certificateId := idSplit[1]

	request.ZoneId = &zoneId
	request.CertificateId = &certificateId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().RemoveExternalSAMLIdPCertificateWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete organization external saml idp certificate failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
