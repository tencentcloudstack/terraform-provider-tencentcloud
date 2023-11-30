package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudProjectCreate,
		Read:   resourceTencentCloudProjectRead,
		Update: resourceTencentCloudProjectUpdate,
		Delete: resourceTencentCloudProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of project.",
			},

			"info": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description of project.",
			},

			"disable": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "If disable project. 1 means disable, 0 means enable. Default 0.",
			},

			"creator_uin": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Uin of creator.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create time.",
			},
		},
	}
}

func resourceTencentCloudProjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_project.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tag.NewAddProjectRequest()
		response  = tag.NewAddProjectResponse()
		projectId uint64
	)
	if v, ok := d.GetOk("project_name"); ok {
		request.ProjectName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("info"); ok {
		request.Info = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTagClient().AddProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tag project failed, reason:%+v", logId, err)
		return err
	}

	projectId = *response.Response.ProjectId
	d.SetId(helper.UInt64ToStr(projectId))

	if v, ok := d.GetOkExists("disable"); ok {
		if v.(int) == 1 {
			ctx := context.WithValue(context.TODO(), logIdKey, logId)

			service := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
			projectId := helper.StrToUInt64(d.Id())

			if err := service.DisableProjectById(ctx, projectId); err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudProjectRead(d, meta)
}

func resourceTencentCloudProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_project.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectId := helper.StrToUInt64(d.Id())

	project, disable, err := service.DescribeProjectById(ctx, projectId)
	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TagProject` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if project.ProjectName != nil {
		_ = d.Set("project_name", project.ProjectName)
	}

	if project.ProjectInfo != nil {
		_ = d.Set("info", project.ProjectInfo)
	}

	if disable != nil {
		_ = d.Set("disable", disable)
	}

	if project.CreatorUin != nil {
		_ = d.Set("creator_uin", project.CreatorUin)
	}

	if project.CreateTime != nil {
		_ = d.Set("create_time", project.CreateTime)
	}

	return nil
}

func resourceTencentCloudProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_project.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tag.NewUpdateProjectRequest()

	projectId := helper.StrToUInt64(d.Id())

	request.ProjectId = &projectId

	mutableArgs := []string{"project_name", "info", "disable"}

	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("project_name"); ok {
			request.ProjectName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("info"); ok {
			request.Info = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("disable"); ok {
			request.Disable = helper.IntInt64(v.(int))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTagClient().UpdateProject(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update tag project failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudProjectRead(d, meta)
}

func resourceTencentCloudProjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_project.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	projectId := helper.StrToUInt64(d.Id())

	if err := service.DisableProjectById(ctx, projectId); err != nil {
		return err
	}

	return nil
}
