package tco

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	organizationv20210331 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationExternalSamlIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationExternalSamlIdentityProviderCreate,
		Read:   resourceTencentCloudOrganizationExternalSamlIdentityProviderRead,
		Update: resourceTencentCloudOrganizationExternalSamlIdentityProviderUpdate,
		Delete: resourceTencentCloudOrganizationExternalSamlIdentityProviderDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Space ID.",
			},

			"encoded_metadata_document": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"x509_certificate"},
				Description:   "IdP metadata document (Base64 encoded). Provided by an IdP that supports the SAML 2.0 protocol.",
			},

			"sso_status": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "SSO enabling status. Valid values: Enabled, Disabled (default).",
			},

			"entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "IdP identifier.",
			},

			"login_url": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "IdP login URL.",
			},

			"x509_certificate": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"encoded_metadata_document"},
				Description:   "X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.",
			},

			"another_x509_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Another X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.",
			},

			// computed
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate ID.",
			},

			"another_certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Another certificate ID.",
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

func resourceTencentCloudOrganizationExternalSamlIdentityProviderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_identity_provider.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = organizationv20210331.NewSetExternalSAMLIdentityProviderRequest()
		response = organizationv20210331.NewSetExternalSAMLIdentityProviderResponse()
		zoneId   string
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

		if result == nil || result.Response == nil || result.Response.CertificateIds == nil {
			return resource.NonRetryableError(fmt.Errorf("Create organization external saml identity provider failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create organization external saml identity provider failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.CertificateIds) == 0 {
		return fmt.Errorf("CertificateIds is nil.")
	}

	// set main certificate id
	_ = d.Set("certificate_id", response.Response.CertificateIds[0])
	d.SetId(zoneId)

	// another certificate
	if v, ok := d.GetOk("another_x509_certificate"); ok {
		request := organization.NewAddExternalSAMLIdPCertificateRequest()
		response := organization.NewAddExternalSAMLIdPCertificateResponse()
		request.ZoneId = &zoneId
		request.X509Certificate = helper.String(v.(string))
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddExternalSAMLIdPCertificateWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Create another organization external saml identity provider failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s create another organization external saml IdP certificate failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if response.Response.CertificateId == nil {
			return fmt.Errorf("Another certificateId is nil.")
		}

		// set another certificate id
		_ = d.Set("another_certificate_id", response.Response.CertificateId)
	}

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

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudOrganizationExternalSamlIdentityProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_external_saml_identity_provider.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId = d.Id()
	)

	if d.HasChange("encoded_metadata_document") || d.HasChange("x509_certificate") || d.HasChange("another_x509_certificate") {
		if d.Get("encoded_metadata_document").(string) == "" && d.Get("x509_certificate").(string) == "" && d.Get("another_x509_certificate").(string) == "" {
			return fmt.Errorf("At least one certificate must be retained.")
		}

		if d.HasChange("encoded_metadata_document") || d.HasChange("x509_certificate") {
			oldEmdInterface, newEmdInterface := d.GetChange("encoded_metadata_document")
			oldEmd := oldEmdInterface.(string)
			newEmd := newEmdInterface.(string)
			if newEmd != "" {
				return fmt.Errorf("Currently, `encoded_metadata_document` does not support adding new value.")
			}

			oldX509CertInterface, newX509CertInterface := d.GetChange("x509_certificate")
			oldX509Cert := oldX509CertInterface.(string)
			newX509Cert := newX509CertInterface.(string)

			// delete first
			if oldEmd != "" || oldX509Cert != "" {
				request := organization.NewRemoveExternalSAMLIdPCertificateRequest()
				tmpCertificateId := d.Get("certificate_id").(string)

				request.ZoneId = &zoneId
				request.CertificateId = helper.String(tmpCertificateId)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().RemoveExternalSAMLIdPCertificateWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Remove organization external saml identity provider failed, Response is nil."))
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s remove organization external saml identity provider failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				// Clear certificate_id
				_ = d.Set("certificate_id", "")
			}

			// add new
			if newX509Cert != "" {
				request := organization.NewAddExternalSAMLIdPCertificateRequest()
				response := organization.NewAddExternalSAMLIdPCertificateResponse()
				tmpX509Certificate := d.Get("x509_certificate").(string)

				request.ZoneId = &zoneId
				request.X509Certificate = helper.String(tmpX509Certificate)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddExternalSAMLIdPCertificateWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Add organization external saml identity provider failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s add organization external saml identity provider failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				if response.Response.CertificateId == nil {
					return fmt.Errorf("CertificateId is nil.")
				}

				_ = d.Set("certificate_id", response.Response.CertificateId)
			}
		}

		if d.HasChange("another_x509_certificate") {
			oldAnotherX509CertInterface, newAnotherX509CertInterface := d.GetChange("another_x509_certificate")
			oldAnotherX509Cert := oldAnotherX509CertInterface.(string)
			newAnotherX509Cert := newAnotherX509CertInterface.(string)

			// delete first
			if oldAnotherX509Cert != "" {
				request := organization.NewRemoveExternalSAMLIdPCertificateRequest()
				tmpCertificateId := d.Get("another_certificate_id").(string)

				request.ZoneId = &zoneId
				request.CertificateId = helper.String(tmpCertificateId)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().RemoveExternalSAMLIdPCertificateWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Remove another organization external saml identity provider failed, Response is nil."))
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s remove another organization external saml identity provider failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				// Clear certificate_id
				_ = d.Set("another_certificate_id", "")
			}

			// add new
			if newAnotherX509Cert != "" {
				request := organization.NewAddExternalSAMLIdPCertificateRequest()
				response := organization.NewAddExternalSAMLIdPCertificateResponse()

				request.ZoneId = &zoneId
				request.X509Certificate = helper.String(newAnotherX509Cert)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddExternalSAMLIdPCertificateWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Add another organization external saml identity provider failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s add another organization external saml identity provider failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				if response.Response.CertificateId == nil {
					return fmt.Errorf("CertificateId is nil.")
				}

				_ = d.Set("another_certificate_id", response.Response.CertificateId)
			}
		}
	}

	return resourceTencentCloudOrganizationExternalSamlIdentityProviderRead(d, meta)
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
