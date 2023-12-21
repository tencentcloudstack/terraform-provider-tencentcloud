package lighthouse

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

func ResourceTencentCloudLighthouseStopInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseStopInstanceCreate,
		Read:   resourceTencentCloudLighthouseStopInstanceRead,
		Delete: resourceTencentCloudLighthouseStopInstanceDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudLighthouseStopInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_stop_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := lighthouse.NewStopInstancesRequest()
	instanceId := d.Get("instance_id").(string)
	request.InstanceIds = []*string{&instanceId}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().StopInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse stopInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseStopInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseStopInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_stop_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseStopInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_stop_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
