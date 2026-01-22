package emr

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEmrDeployYarnOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrDeployYarnOperationCreate,
		Read:   resourceTencentCloudEmrDeployYarnOperationRead,
		Delete: resourceTencentCloudEmrDeployYarnOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "EMR Instance ID.",
			},
		},
	}
}

func resourceTencentCloudEmrDeployYarnOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_deploy_yarn_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
	)

	instanceId = d.Get("instance_id").(string)
	deployRequest := emr.NewDeployYarnConfRequest()
	deployRequest.InstanceId = helper.String(instanceId)
	var flowId *uint64
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().DeployYarnConfWithContext(ctx, deployRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deployRequest.GetAction(), deployRequest.ToJsonString(), result.ToJsonString())
		}
		flowId = result.Response.FlowId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update emr yarn failed, reason:%+v", logId, err)
		return err
	}

	if flowId != nil {
		emrService := EMRService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, strconv.FormatUint(*flowId, 10), F_KEY_FLOW_ID, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	d.SetId(instanceId)

	return resourceTencentCloudEmrDeployYarnOperationRead(d, meta)
}

func resourceTencentCloudEmrDeployYarnOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_deploy_yarn_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEmrDeployYarnOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_deploy_yarn_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
