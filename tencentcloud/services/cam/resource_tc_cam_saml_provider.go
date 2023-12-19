package cam

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamSAMLProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamSAMLProviderCreate,
		Read:   resourceTencentCloudCamSAMLProviderRead,
		Update: resourceTencentCloudCamSAMLProviderUpdate,
		Delete: resourceTencentCloudCamSAMLProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of CAM SAML provider.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the CAM SAML provider.",
			},
			"meta_data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The meta data document of the CAM SAML provider.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the CAM SAML provider.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the CAM SAML provider.",
			},
			"provider_arn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ARN of the CAM SAML provider.",
			},
		},
	}
}

func resourceTencentCloudCamSAMLProviderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_saml_provider.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cam.NewCreateSAMLProviderRequest()
	request.Name = helper.String(d.Get("name").(string))
	request.SAMLMetadataDocument = helper.String(d.Get("meta_data").(string))
	//special check function
	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	var response *cam.CreateSAMLProviderResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().CreateSAMLProvider(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "IdentityNameInUse") {
					return resource.NonRetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM SAML provider failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response.Response.ProviderArn == nil {
		return fmt.Errorf("CAM SAML provider id is nil")
	}

	d.SetId(d.Get("name").(string))
	_ = d.Set("provider_arn", *response.Response.ProviderArn)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	samlProviderId := d.Id()
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeSAMLProviderById(ctx, samlProviderId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM SAML provider failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)

	return resourceTencentCloudCamSAMLProviderRead(d, meta)
}

func resourceTencentCloudCamSAMLProviderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_saml_provider.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	samlProviderId := d.Id()
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *cam.GetSAMLProviderResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeSAMLProviderById(ctx, samlProviderId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM SAML provider failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil || instance.Response == nil || instance.Response.Name == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", *instance.Response.Name)
	_ = d.Set("create_time", *instance.Response.CreateTime)
	_ = d.Set("update_time", *instance.Response.ModifyTime)
	_ = d.Set("meta_data", *instance.Response.SAMLMetadata)
	if instance.Response.Description != nil {
		_ = d.Set("description", *instance.Response.Description)
	}
	return nil
}

func resourceTencentCloudCamSAMLProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_saml_provider.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	SAMLProviderId := d.Id()
	request := cam.NewUpdateSAMLProviderRequest()
	request.Name = &SAMLProviderId
	changeFlag := false

	if d.HasChange("description") {
		request.Description = helper.String(d.Get("description").(string))
		changeFlag = true

	}
	if d.HasChange("meta_data") {
		request.SAMLMetadataDocument = helper.String(d.Get("meta_data").(string))
		changeFlag = true
	}

	if changeFlag {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().UpdateSAMLProvider(request)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM SAML provider description failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCamSAMLProviderRead(d, meta)
}

func resourceTencentCloudCamSAMLProviderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_saml_provider.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	SAMLProviderId := d.Id()
	request := cam.NewDeleteSAMLProviderRequest()
	request.Name = &SAMLProviderId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().DeleteSAMLProvider(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM SAML provider failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
