/*
Provides a resource to create a tem application

Example Usage

```hcl
resource "tencentcloud_tem_application" "application" {
  application_name = "demo"
  description = "demo for test"
  coding_language = "JAVA"
  use_default_image_service = 0
  repo_type = 2
  repo_name = "qcloud/nginx"
  repo_server = "ccr.ccs.tencentyun.com"
}
```
Import

tem application can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_application.application application_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTemApplication() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemApplicationRead,
		Create: resourceTencentCloudTemApplicationCreate,
		Update: resourceTencentCloudTemApplicationUpdate,
		Delete: resourceTencentCloudTemApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "application name.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "application description.",
			},

			"coding_language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "program language, like JAVA.",
			},

			"use_default_image_service": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "create image repo or not.",
			},

			"repo_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "repo type, 0: tcr personal, 1: tcr enterprise, 2: public repository, 4: demo image.",
			},

			"repo_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "registry address.",
			},

			"repo_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "repository name.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "tcr instance id.",
			},
		},
	}
}

func resourceTencentCloudTemApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tem.NewCreateApplicationRequest()
		response *tem.CreateApplicationResponse
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

	request.UseDefaultImageService = helper.IntInt64(d.Get("use_default_image_service").(int))

	if v, ok := d.GetOk("repo_type"); ok {
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

	request.DeployMode = helper.String("IMAGE")

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateApplication(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tem application failed, reason:%+v", logId, err)
		return err
	}

	applicationId := *response.Response.Result

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

	applications, err := service.DescribeTemApplication(ctx, applicationId)

	if len(applications.Result.Records) != 1 {
		d.SetId("")
		return nil
	}
	application := applications.Result.Records[0]

	if err != nil {
		return err
	}

	if application == nil {
		d.SetId("")
		return fmt.Errorf("resource `application` %s does not exist", applicationId)
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

	if application.RepoType != nil {
		_ = d.Set("repo_type", application.RepoType)
	}

	if application.RepoName != nil {
		_ = d.Set("repo_name", application.RepoName)
	}

	if application.InstanceId != nil {
		_ = d.Set("instance_id", application.InstanceId)
	}

	return nil
}

func resourceTencentCloudTemApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyApplicationInfoRequest()

	request.ApplicationId = helper.String(d.Id())

	if d.HasChange("application_name") {
		return fmt.Errorf("`application_name` do not support change now.")
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("coding_language") {
		return fmt.Errorf("`coding_language` do not support change now.")
	}

	if d.HasChange("use_default_image_service") {
		return fmt.Errorf("`use_default_image_service` do not support change now.")
	}

	if d.HasChange("repo_type") {
		return fmt.Errorf("`repo_type` do not support change now.")
	}

	if d.HasChange("repo_server") {
		return fmt.Errorf("`repo_server` do not support change now.")
	}

	if d.HasChange("repo_name") {
		return fmt.Errorf("`repo_name` do not support change now.")
	}

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyApplicationInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
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
