/*
Provides a resource to create a cam user_saml_config

Example Usage

```hcl
resource "tencentcloud_cam_user_saml_config" "user_saml_config" {
  saml_metadata_document = "./metadataDocument.xml"
  # saml_metadata_document  = <<-EOT
  # <?xml version="1.0" encoding="utf-8"?></EntityDescriptor>
  # EOT
}
```

Import

cam user_saml_config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_saml_config.user_saml_config user_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamUserSamlConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserSamlConfigCreate,
		Read:   resourceTencentCloudCamUserSamlConfigRead,
		Update: resourceTencentCloudCamUserSamlConfigUpdate,
		Delete: resourceTencentCloudCamUserSamlConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"saml_metadata_document": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SAML metadata document, xml format, support string content or file path.",
				StateFunc: func(v interface{}) string {
					saml := v.(string)
					if saml != "" {
						b := strings.HasSuffix(saml, ".xml")
						if b {
							metadata, _ := ReadFromFile(saml)
							return string(metadata)
						}
					}
					return saml
				},
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Status: `0`: not set, `11`: enabled, `2`: disabled.",
			},

			"metadata_document_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path used to save the samlMetadat file.",
			},
		},
	}
}

func resourceTencentCloudCamUserSamlConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_saml_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewCreateUserSAMLConfigRequest()
		response = cam.NewCreateUserSAMLConfigResponse()
	)
	if v, ok := d.GetOk("saml_metadata_document"); ok {
		saml := v.(string)
		b := strings.HasSuffix(saml, ".xml")
		if b {
			metadata, err := ReadFromFile(v.(string))
			if err != nil {
				return err
			}
			saml = string(metadata)
		}
		request.SAMLMetadataDocument = helper.String(StringToBase64(saml))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateUserSAMLConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam userSamlConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(*response.Response.RequestId)
	return resourceTencentCloudCamUserSamlConfigRead(d, meta)
}

func resourceTencentCloudCamUserSamlConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_saml_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	samlConfig, err := service.DescribeCamUserSamlConfigById(ctx)
	if err != nil {
		return err
	}

	if samlConfig == nil || samlConfig.Response == nil || *samlConfig.Response.Status == 2 {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamUserSamlConfig` status is closed, please check if it has been closed.", logId)
		return nil
	}
	userSamlConfig := samlConfig.Response

	if userSamlConfig.SAMLMetadata != nil {
		metadata, err := Base64ToString(*userSamlConfig.SAMLMetadata)
		if err != nil {
			return fmt.Errorf("`SamlConfig.SAMLMetadata` %s does not be decoded to xml", *userSamlConfig.SAMLMetadata)
		}
		_ = d.Set("saml_metadata_document", metadata)
	}

	if userSamlConfig.Status != nil {
		_ = d.Set("status", userSamlConfig.Status)
	}

	output, ok := d.GetOk("metadata_document_files")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), d.Get("saml_metadata_document")); err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudCamUserSamlConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_saml_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cam.NewUpdateUserSAMLConfigRequest()

	if d.HasChange("saml_metadata_document") {
		if v, ok := d.GetOk("saml_metadata_document"); ok {
			saml := v.(string)
			b := strings.HasSuffix(saml, ".xml")
			if b {
				metadata, err := ReadFromFile(v.(string))
				if err != nil {
					return err
				}
				saml = string(metadata)
			}
			request.SAMLMetadataDocument = helper.String(StringToBase64(saml))
		}
	}

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	samlConfig, describeErr := service.DescribeCamUserSamlConfigById(ctx)
	if describeErr != nil {
		return describeErr
	}
	if *samlConfig.Response.Status == 2 {
		request.Operate = helper.String("enable")
	} else {
		request.Operate = helper.String("updateSAML")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateUserSAMLConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cam userSamlConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCamUserSamlConfigRead(d, meta)
}

func resourceTencentCloudCamUserSamlConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_saml_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	if err := service.DeleteCamUserSamlConfigById(ctx); err != nil {
		return err
	}

	return nil
}
