/*
Provides a resource to create a CAM-ROLE-SSO (Only support OIDC).

Example Usage

```hcl
resource "tencentcloud_cam_role_sso" "foo" {
	name="tf_cam_role_sso"
	identity_url="https://login.microsoftonline.com/.../v2.0"
	identity_key="..."
	client_ids=["..."]
	description="this is a description"
}
```

Import

CAM-ROLE-SSO can be imported using the `name`, e.g.

```
$ terraform import tencentcloud_cam_role_sso.foo "test"
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamRoleSSO() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRoleSSOCreate,
		Read:   resourceTencentCloudCamRoleSSORead,
		Update: resourceTencentCloudCamRoleSSOUpdate,
		Delete: resourceTencentCloudCamRoleSSODelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of resource.",
			},
			"identity_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identity provider URL.",
			},
			"identity_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sign the public key.",
			},
			"client_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Client ids.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of resource.",
			},
		},
	}
}

func resourceTencentCloudCamRoleSSOCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_sso.create")()

	logId := getLogId(contextNil)

	request := cam.NewCreateOIDCConfigRequest()
	request.Name = helper.String(d.Get("name").(string))
	request.IdentityUrl = helper.String(d.Get("identity_url").(string))
	request.IdentityKey = helper.String(d.Get("identity_key").(string))
	request.Description = helper.String(d.Get("description").(string))
	request.ClientId = helper.InterfacesStringsPoint(d.Get("client_ids").(*schema.Set).List())

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateOIDCConfig(request)
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
		log.Printf("[CRITAL]%s create CAM role SSO failed, reason:%s\n", logId, err.Error())
		return err
	}
	d.SetId(d.Get("name").(string))
	return resourceTencentCloudCamRoleSSORead(d, meta)
}

func resourceTencentCloudCamRoleSSORead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_sso.read")()

	logId := getLogId(contextNil)
	request := cam.NewDescribeOIDCConfigRequest()
	request.Name = helper.String(d.Id())
	var response *cam.DescribeOIDCConfigResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DescribeOIDCConfig(request)
		if e != nil {
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role SSO failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("identity_key", *response.Response.IdentityKey)
	_ = d.Set("identity_url", *response.Response.IdentityUrl)
	_ = d.Set("name", *response.Response.Name)
	_ = d.Set("description", *response.Response.Description)
	clientIds := make([]string, 0)
	for _, clientId := range response.Response.ClientId {
		clientIds = append(clientIds, *clientId)
	}
	_ = d.Set("client_ids", clientIds)

	return nil
}

func resourceTencentCloudCamRoleSSOUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_sso.update")()
	logId := getLogId(contextNil)
	request := cam.NewUpdateOIDCConfigRequest()
	if d.HasChange("name") {
		return fmt.Errorf("not support change name")
	}
	request.Name = helper.String(d.Id())
	if d.HasChange("identity_url") || d.HasChange("identity_key") || d.HasChange("description") || d.HasChange("client_ids") {
		request.IdentityKey = helper.String(d.Get("identity_key").(string))
		request.IdentityUrl = helper.String(d.Get("identity_url").(string))
		request.Description = helper.String(d.Get("description").(string))
		request.ClientId = helper.InterfacesStringsPoint(d.Get("client_ids").(*schema.Set).List())
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateOIDCConfig(request)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update CAM Role SSO failed, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCamRoleSSORead(d, meta)
}

func resourceTencentCloudCamRoleSSODelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_sso.delete")()
	logId := getLogId(contextNil)
	request := cam.NewDeleteOIDCConfigRequest()
	name := d.Id()
	request.Name = helper.String(name)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DeleteOIDCConfig(request)
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
