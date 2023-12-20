package css

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCssPullStreamTaskRestart() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssPullStreamTaskRestartCreate,
		Read:   resourceTencentCloudCssPullStreamTaskRestartRead,
		Delete: resourceTencentCloudCssPullStreamTaskRestartDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Task Id.",
			},

			"operator": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Task operator.",
			},
		},
	}
}

func resourceTencentCloudCssPullStreamTaskRestartCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_pull_stream_task_restart.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = css.NewRestartLivePullStreamTaskRequest()
		taskId  string
	)
	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
		request.TaskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("operator"); ok {
		request.Operator = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().RestartLivePullStreamTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate css restartPushTask failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(taskId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"active"}, 6*tccommon.ReadRetryTimeout, time.Second, service.CssRestartPushTaskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCssPullStreamTaskRestartRead(d, meta)
}

func resourceTencentCloudCssPullStreamTaskRestartRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_pull_stream_task_restart.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCssPullStreamTaskRestartDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_pull_stream_task_restart.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
