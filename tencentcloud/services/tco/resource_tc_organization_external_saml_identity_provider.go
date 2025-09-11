package tco

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organizationv20210331 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationExternalSamlIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationExternalSamlIdentityProviderCreate,
		Read:   resourceTencentCloudOrganizationExternalSamlIdentityProviderRead,
		Delete: resourceTencentCloudOrganizationExternalSamlIdentityProviderDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Space ID.",
			},

			"encoded_metadata_document": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IdP metadata document (Base64 encoded). Provided by an IdP that supports the SAML 2.0 protocol.",
			},

			"sso_status": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "SSO enabling status. Valid values: Enabled, Disabled (default).",
			},

			"entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IdP identifier.",
			},

			"login_url": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IdP login URL.",
			},

			"x509_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.",
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},
		},
	}
}

func resourceTencentCloudOrganizationExternalSamlIdentityProviderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_identity_provider.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = organizationv20210331.NewSetExternalSAMLIdentityProviderRequest()
		zoneId  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("encoded_metadata_document"); ok {
		request.EncodedMetadataDocument = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sso_status"); ok {
		request.SSOStatus = helper.String(v.(string))
	}

	if v, ok := d.GetOk("entity_id"); ok {
		request.EntityId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("login_url"); ok {
		request.LoginUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("x509_certificate"); ok {
		request.X509Certificate = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().SetExternalSAMLIdentityProviderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create organization external saml identity provider failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(zoneId)
	return resourceTencentCloudOrganizationExternalSamlIdentityProviderRead(d, meta)
}

func resourceTencentCloudOrganizationExternalSamlIdentityProviderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_identity_provider.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	respData, err := service.DescribeOrganizationExternalSamlIdentityProviderById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_organization_external_saml_identity_provider` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.EncodedMetadataDocument != nil {
		_ = d.Set("encoded_metadata_document", respData.EncodedMetadataDocument)
	}

	if respData.SSOStatus != nil {
		_ = d.Set("sso_status", respData.SSOStatus)
	}

	if respData.EntityId != nil {
		_ = d.Set("entity_id", respData.EntityId)
	}

	if respData.LoginUrl != nil {
		_ = d.Set("login_url", respData.LoginUrl)
	}

	if respData.CertificateIds != nil {
		_ = d.Set("certificate_ids", respData.CertificateIds)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudOrganizationExternalSamlIdentityProviderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_identity_provider.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = organizationv20210331.NewClearExternalSAMLIdentityProviderRequest()
		zoneId  = d.Id()
	)

	request.ZoneId = &zoneId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().ClearExternalSAMLIdentityProviderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete organization external saml identity provider failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
