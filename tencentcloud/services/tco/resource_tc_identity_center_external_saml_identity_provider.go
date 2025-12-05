package tco

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIdentityCenterExternalSamlIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterExternalSamlIdentityProviderCreate,
		Read:   resourceTencentCloudIdentityCenterExternalSamlIdentityProviderRead,
		Update: resourceTencentCloudIdentityCenterExternalSamlIdentityProviderUpdate,
		Delete: resourceTencentCloudIdentityCenterExternalSamlIdentityProviderDelete,
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

			"encoded_metadata_document": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IdP metadata document (Base64 encoded). Provided by an IdP that supports the SAML 2.0 protocol.",
			},

			"sso_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSO enabling status. Valid values: Enabled, Disabled (default).",
			},

			"entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IdP identifier.",
			},

			"login_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IdP login URL.",
			},

			"x509_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.",
			},

			"acs_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acs url.",
			},

			"certificate_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Certificate ids.",
			},

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

func resourceTencentCloudIdentityCenterExternalSamlIdentityProviderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_external_saml_identity_provider.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)

	return resourceTencentCloudIdentityCenterExternalSamlIdentityProviderUpdate(d, meta)
}

func resourceTencentCloudIdentityCenterExternalSamlIdentityProviderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_external_saml_identity_provider.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	respData, err := service.DescribeIdentityCenterExternalSamlIdentityProviderById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_identity_center_external_saml_identity_provider` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.EntityId != nil {
		_ = d.Set("entity_id", respData.EntityId)
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.EncodedMetadataDocument != nil {
		_ = d.Set("encoded_metadata_document", respData.EncodedMetadataDocument)
	}

	if respData.AcsUrl != nil {
		_ = d.Set("acs_url", respData.AcsUrl)
	}

	respData1, err := service.DescribeIdentityCenterExternalSamlIdentityProviderById1(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData1 == nil {
		log.Printf("[WARN]%s resource `identity_center_external_saml_identity_provider` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData1.EntityId != nil {
		_ = d.Set("entity_id", respData1.EntityId)
	}

	if respData1.SSOStatus != nil {
		_ = d.Set("sso_status", respData1.SSOStatus)
	}

	if respData1.EncodedMetadataDocument != nil {
		_ = d.Set("encoded_metadata_document", respData1.EncodedMetadataDocument)
	}

	if respData1.CertificateIds != nil {
		_ = d.Set("certificate_ids", respData1.CertificateIds)
	}

	if respData1.LoginUrl != nil {
		_ = d.Set("login_url", respData1.LoginUrl)
	}

	if respData1.CreateTime != nil {
		_ = d.Set("create_time", respData1.CreateTime)
	}

	if respData1.UpdateTime != nil {
		_ = d.Set("update_time", respData1.UpdateTime)
	}

	return nil
}

func resourceTencentCloudIdentityCenterExternalSamlIdentityProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_external_saml_identity_provider.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId = d.Id()
	)

	if d.HasChange("encoded_metadata_document") {
		request := organization.NewSetExternalSAMLIdentityProviderRequest()
		if v, ok := d.GetOk("encoded_metadata_document"); ok {
			request.EncodedMetadataDocument = helper.String(v.(string))
		}

		request.ZoneId = &zoneId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().SetExternalSAMLIdentityProviderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update identity center external saml identity provider failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange := false
	mutableArgs := []string{"entity_id", "login_url", "x509_certificate"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := organization.NewSetExternalSAMLIdentityProviderRequest()
		if v, ok := d.GetOk("entity_id"); ok {
			request.EntityId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("login_url"); ok {
			request.LoginUrl = helper.String(v.(string))
		}

		if v, ok := d.GetOk("x509_certificate"); ok {
			request.X509Certificate = helper.String(v.(string))
		}

		request.ZoneId = &zoneId
		request.EncodedMetadataDocument = helper.String("")
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().SetExternalSAMLIdentityProviderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update identity center external saml identity provider failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("sso_status") {
		request := organization.NewSetExternalSAMLIdentityProviderRequest()
		if v, ok := d.GetOk("sso_status"); ok {
			request.SSOStatus = helper.String(v.(string))
		}

		request.ZoneId = &zoneId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().SetExternalSAMLIdentityProviderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update identity center external saml identity provider failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudIdentityCenterExternalSamlIdentityProviderRead(d, meta)
}

func resourceTencentCloudIdentityCenterExternalSamlIdentityProviderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_external_saml_identity_provider.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	respData1, err := service.DescribeIdentityCenterExternalSamlIdentityProviderById1(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData1.SSOStatus != nil && *respData1.SSOStatus == "Enabled" {
		request := organization.NewSetExternalSAMLIdentityProviderRequest()
		request.ZoneId = helper.String(zoneId)
		request.SSOStatus = helper.String("Disabled")
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().SetExternalSAMLIdentityProviderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update identity center external saml identity provider failed, reason:%+v", logId, err)
			return err
		}
	}

	request := organization.NewClearExternalSAMLIdentityProviderRequest()
	request.ZoneId = &zoneId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().ClearExternalSAMLIdentityProviderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete identity center external saml identity provider failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
