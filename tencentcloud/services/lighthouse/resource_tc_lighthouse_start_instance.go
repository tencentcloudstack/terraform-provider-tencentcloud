package lighthouse

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

func ResourceTencentCloudLighthouseStartInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseStartInstanceCreate,
		Read:   resourceTencentCloudLighthouseStartInstanceRead,
		Delete: resourceTencentCloudLighthouseStartInstanceDelete,
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

func resourceTencentCloudLighthouseStartInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_start_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := lighthouse.NewStartInstancesRequest()
	instanceId := d.Get("instance_id").(string)
	request.InstanceIds = []*string{&instanceId}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().StartInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse startInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseStartInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseStartInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_start_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseStartInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_start_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
