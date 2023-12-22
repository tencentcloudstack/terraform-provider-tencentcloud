package rum

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
)

func ResourceTencentCloudRumProjectStatusConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumProjectStatusConfigCreate,
		Read:   resourceTencentCloudRumProjectStatusConfigRead,
		Update: resourceTencentCloudRumProjectStatusConfigUpdate,
		Delete: resourceTencentCloudRumProjectStatusConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},
			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`resume`, `stop`.",
			},
		},
	}
}

func resourceTencentCloudRumProjectStatusConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_rum_project_status_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var projectId int
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(int)
	}

	d.SetId(strconv.Itoa(projectId))

	return resourceTencentCloudRumProjectStatusConfigRead(d, meta)
}

func resourceTencentCloudRumProjectStatusConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_rum_project_status_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	projectId := d.Id()

	project, err := service.DescribeRumProjectStatusConfigById(ctx, projectId)
	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumProjectStatusConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if project.ID != nil {
		_ = d.Set("project_id", project.ID)
	}

	return nil
}

func resourceTencentCloudRumProjectStatusConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_rum_project_status_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	projectId, _ := strconv.ParseInt(d.Id(), 10, 64)

	var operate string
	if v, ok := d.GetOk("operate"); ok {
		operate = v.(string)
	}

	if operate == "resume" {
		request := rum.NewResumeProjectRequest()
		request.ProjectId = &projectId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRumClient().ResumeProject(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s resume rum project failed, reason:%+v", logId, err)
			return err
		}
	} else if operate == "stop" {
		request := rum.NewStopProjectRequest()
		request.ProjectId = &projectId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRumClient().StopProject(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s stop rum project failed, reason:%+v", logId, err)
			return err
		}
	} else {
		return fmt.Errorf("[CRITAL]%s Operation type error", logId)
	}

	return resourceTencentCloudRumProjectStatusConfigRead(d, meta)
}

func resourceTencentCloudRumProjectStatusConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_rum_project_status_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
