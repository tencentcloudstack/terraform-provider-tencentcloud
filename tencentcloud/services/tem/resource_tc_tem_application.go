package tem

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTemApplication() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemApplicationRead,
		Create: resourceTencentCloudTemApplicationCreate,
		Update: resourceTencentCloudTemApplicationUpdate,
		Delete: resourceTencentCloudTemApplicationDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "application name.",
			},

			"description": {
				Type:        schema.TypeString,
				Required:    true,
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
				Description: "repo type, 0: tcr personal, 1: tcr enterprise, 2: public repository, 3: tcr hosted by tem, 4: demo image.",
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
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "application tag list.",
			},
		},
	}
}

func resourceTencentCloudTemApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_application.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			tag := tem.Tag{
				TagKey:   helper.String(key),
				TagValue: helper.String(value.(string)),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	request.DeployMode = helper.String("IMAGE")

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().CreateApplication(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::tem:%s:uin/:application/%s", region, applicationId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTemApplicationRead(d, meta)
}

func resourceTencentCloudTemApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_application.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region
	tags, err := tagService.DescribeResourceTags(ctx, "tem", "application", region, applicationId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTemApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_application.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := tem.NewModifyApplicationInfoRequest()

	request.ApplicationId = helper.String(d.Id())

	if d.HasChange("application_name") {
		return fmt.Errorf("`application_name` do not support change now.")
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().ModifyApplicationInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("tem", "application", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTemApplicationRead(d, meta)
}

func resourceTencentCloudTemApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_application.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	applicationId := d.Id()

	if err := service.DeleteTemApplicationById(ctx, applicationId); err != nil {
		return err
	}

	return nil
}
