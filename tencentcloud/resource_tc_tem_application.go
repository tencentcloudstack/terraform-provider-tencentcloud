/*
Provides a resource to create a tem application

Example Usage

```hcl
resource "tencentcloud_tem_application" "application" {
  application_name = "xxx"
  description = "xxx"
  coding_language = "JAVA"
  use_default_image_service = 1
  repo_type = 0
  repo_server = &lt;nil&gt;
  repo_name = &lt;nil&gt;
  instance_id = &lt;nil&gt;
  tags {
		tag_key = "key"
		tag_value = "tag value"

  }
}
```

Import

tem application can be imported using the id, e.g.

```
terraform import tencentcloud_tem_application.application application_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTemApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemApplicationCreate,
		Read:   resourceTencentCloudTemApplicationRead,
		Update: resourceTencentCloudTemApplicationUpdate,
		Delete: resourceTencentCloudTemApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application name.",
			},

			"description": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application description.",
			},

			"coding_language": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Program language, like JAVA.",
			},

			"use_default_image_service": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Create image repo or not.",
			},

			"repo_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Repo type, 0: tcr personal, 1: tcr enterprise, 2: public repository, 3: tcr hosted by tem, 4: demo image.",
			},

			"repo_server": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Registry address.",
			},

			"repo_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Repository name.",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Tcr instance id.",
			},

			"tags": {
				Optional:    true,
				Description: "Application tag list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewCreateApplicationRequest()
		response      = tem.NewCreateApplicationResponse()
		applicationId string
	)
	if v, ok := d.GetOk("application_name"); ok {
		request.ApplicationName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("coding_language"); ok {
		request.CodingLanguage = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("use_default_image_service"); ok {
		request.UseDefaultImageService = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("repo_type"); ok {
		request.RepoType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("repo_server"); ok {
		request.RepoServer = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repo_name"); ok {
		request.RepoName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("tags"); v != nil {
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateApplication(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem application failed, reason:%+v", logId, err)
		return err
	}

	applicationId = *response.Response.ApplicationId
	d.SetId(applicationId)

	return resourceTencentCloudTemApplicationRead(d, meta)
}

func resourceTencentCloudTemApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationId := d.Id()

	application, err := service.DescribeTemApplicationById(ctx, applicationId)
	if err != nil {
		return err
	}

	if application == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemApplication` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if application.ApplicationName != nil {
		_ = d.Set("application_name", application.ApplicationName)
	}

	if application.Description != nil {
		_ = d.Set("description", application.Description)
	}

	if application.CodingLanguage != nil {
		_ = d.Set("coding_language", application.CodingLanguage)
	}

	if application.UseDefaultImageService != nil {
		_ = d.Set("use_default_image_service", application.UseDefaultImageService)
	}

	if application.RepoType != nil {
		_ = d.Set("repo_type", application.RepoType)
	}

	if application.RepoServer != nil {
		_ = d.Set("repo_server", application.RepoServer)
	}

	if application.RepoName != nil {
		_ = d.Set("repo_name", application.RepoName)
	}

	if application.InstanceId != nil {
		_ = d.Set("instance_id", application.InstanceId)
	}

	if application.tags != nil {
	}

	return nil
}

func resourceTencentCloudTemApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyApplicationInfoRequest()

	applicationId := d.Id()

	request.ApplicationId = &applicationId

	immutableArgs := []string{"application_name", "description", "coding_language", "use_default_image_service", "repo_type", "repo_server", "repo_name", "instance_id", "tags"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyApplicationInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem application failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTemApplicationRead(d, meta)
}

func resourceTencentCloudTemApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	applicationId := d.Id()

	if err := service.DeleteTemApplicationById(ctx, applicationId); err != nil {
		return err
	}

	return nil
}
