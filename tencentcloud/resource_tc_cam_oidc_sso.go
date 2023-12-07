package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamOIDCSSO() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamOIDCSSOCreate,
		Read:   resourceTencentCloudCamOIDCSSORead,
		Update: resourceTencentCloudCamOIDCSSOUpdate,
		Delete: resourceTencentCloudCamOIDCSSODelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"authorization_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authorization request Endpoint, OpenID Connect identity provider authorization address. Corresponds to the value of the `authorization_endpoint` field in the Openid-configuration provided by the Enterprise IdP.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client ID, the client ID registered with the OpenID Connect identity provider.",
			},
			"identity_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The signature public key requires base64_encode. Verify the public key signed by the OpenID Connect identity provider ID Token. For the security of your account, we recommend that you rotate the signed public key regularly.",
			},
			"identity_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identity provider URL. OpenID Connect identity provider identity.Corresponds to the value of the `issuer` field in the Openid-configuration provided by the Enterprise IdP.",
			},
			"mapping_filed": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Map field names. Which field in the IdP's id_token maps to the user name of the subuser, usually the sub or name field.",
			},
			"response_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authorize the request Forsonse mode. Authorization request return mode, form_post and frogment two optional modes, recommended to select form_post mode.",
			},
			"response_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authorization requests The Response type, with a fixed value id_token.",
			},
			"scope": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Authorize the request Scope. openid; email; profile; Authorization request information scope. The default is required openid.",
			},
		},
	}
}

func resourceTencentCloudCamOIDCSSOCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_oidc_sso.create")()

	logId := getLogId(contextNil)

	request := cam.NewCreateUserOIDCConfigRequest()
	request.IdentityUrl = helper.String(d.Get("identity_url").(string))
	request.IdentityKey = helper.String(d.Get("identity_key").(string))
	request.ClientId = helper.String(d.Get("client_id").(string))
	request.AuthorizationEndpoint = helper.String(d.Get("authorization_endpoint").(string))
	request.ResponseType = helper.String(d.Get("response_type").(string))
	request.ResponseMode = helper.String(d.Get("response_mode").(string))
	request.MappingFiled = helper.String(d.Get("mapping_filed").(string))
	if v, ok := d.GetOk("scope"); ok {
		request.Scope = helper.InterfacesStringsPoint(v.(*schema.Set).List())
	} else {
		request.Scope = helper.InterfacesStringsPoint([]interface{}{"openid"})
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateUserOIDCConfig(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM SSO failed, reason:%s\n", logId, err.Error())
		return err
	}
	d.SetId(d.Get("client_id").(string))
	return resourceTencentCloudCamOIDCSSORead(d, meta)
}

func resourceTencentCloudCamOIDCSSORead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_oidc_sso.read")()

	logId := getLogId(contextNil)
	request := cam.NewDescribeUserOIDCConfigRequest()
	var response *cam.DescribeUserOIDCConfigResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DescribeUserOIDCConfig(request)
		if e != nil {
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM SSO failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("authorization_endpoint", *response.Response.AuthorizationEndpoint)
	_ = d.Set("client_id", *response.Response.ClientId)
	_ = d.Set("identity_key", *response.Response.IdentityKey)
	_ = d.Set("identity_url", *response.Response.IdentityUrl)
	_ = d.Set("mapping_filed", *response.Response.MappingFiled)
	_ = d.Set("response_mode", *response.Response.ResponseMode)
	_ = d.Set("response_type", *response.Response.ResponseType)
	scope := make([]string, 0)
	for _, s := range response.Response.Scope {
		scope = append(scope, *s)
	}
	_ = d.Set("scope", scope)

	return nil
}

func resourceTencentCloudCamOIDCSSOUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_oidc_sso.update")()
	logId := getLogId(contextNil)
	request := cam.NewUpdateUserOIDCConfigRequest()
	if d.HasChange("authorization_endpoint") || d.HasChange("client_id") || d.HasChange("identity_key") || d.HasChange("identity_url") || d.HasChange("mapping_filed") || d.HasChange("response_mode") || d.HasChange("response_type") || d.HasChange("scope") {
		request.AuthorizationEndpoint = helper.String(d.Get("authorization_endpoint").(string))
		request.ClientId = helper.String(d.Get("client_id").(string))
		request.IdentityKey = helper.String(d.Get("identity_key").(string))
		request.IdentityUrl = helper.String(d.Get("identity_url").(string))
		request.MappingFiled = helper.String(d.Get("mapping_filed").(string))
		request.ResponseMode = helper.String(d.Get("response_mode").(string))
		request.ResponseType = helper.String(d.Get("response_type").(string))
		if v, ok := d.GetOk("scope"); ok {
			request.Scope = helper.InterfacesStringsPoint(v.(*schema.Set).List())
		} else {
			request.Scope = helper.InterfacesStringsPoint([]interface{}{"openid"})
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateUserOIDCConfig(request)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update CAM OIDC SSO failed, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCamOIDCSSORead(d, meta)
}

func resourceTencentCloudCamOIDCSSODelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_oidc_sso.delete")()
	logId := getLogId(contextNil)
	request := cam.NewDisableUserSSORequest()
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DisableUserSSO(request)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s disable cam sso failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
