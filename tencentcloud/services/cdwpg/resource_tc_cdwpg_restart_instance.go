package cdwpg

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwpgv20201230 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCdwpgRestartInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdwpgRestartInstanceCreate,
		Read:   resourceTencentCloudCdwpgRestartInstanceRead,
		Delete: resourceTencentCloudCdwpgRestartInstanceDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id (e.g., \"cdwpg-xxxx\").",
			},

			"node_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Node types to restart (gtm/cn/dn).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"node_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Node ids to restart (specify nodes to reboot).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudCdwpgRestartInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_restart_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
	)
	var (
		request  = cdwpgv20201230.NewRestartInstanceRequest()
		response = cdwpgv20201230.NewRestartInstanceResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("node_types"); ok {
		nodeTypesSet := v.(*schema.Set).List()
		for i := range nodeTypesSet {
			nodeTypes := nodeTypesSet[i].(string)
			request.NodeTypes = append(request.NodeTypes, helper.String(nodeTypes))
		}
	}

	if v, ok := d.GetOk("node_ids"); ok {
		nodeIdsSet := v.(*schema.Set).List()
		for i := range nodeIdsSet {
			nodeIds := nodeIdsSet[i].(string)
			request.NodeIds = append(request.NodeIds, helper.String(nodeIds))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwpgV20201230Client().RestartInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdwpg restart instance failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Serving"}, 10*tccommon.ReadRetryTimeout, time.Second, service.InstanceStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	d.SetId(instanceId)

	return resourceTencentCloudCdwpgRestartInstanceRead(d, meta)
}

func resourceTencentCloudCdwpgRestartInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_restart_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdwpgRestartInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_restart_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
