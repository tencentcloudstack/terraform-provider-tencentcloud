package es

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudElasticsearchRestartInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchRestartInstanceOperationCreate,
		Read:   resourceTencentCloudElasticsearchRestartInstanceOperationRead,
		Delete: resourceTencentCloudElasticsearchRestartInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"force_restart": {
				Optional: true,
				Default:  false,
				ForceNew: true,
				Type:     schema.TypeBool,
				Description: "Force restart. Valid values:\n" +
					"- true: Forced restart;\n" +
					"- false: No forced restart;\n" +
					"default false.",
			},

			"restart_mode": {
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Restart mode: 0 roll restart; 1 full restart.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchRestartInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = elasticsearch.NewRestartInstanceRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOkExists("force_restart"); ok {
		request.ForceRestart = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("restart_mode"); ok {
		request.RestartMode = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().RestartInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch RestartInstanceOperation failed, reason:%+v", logId, err)
		return err
	}
	elasticsearchService := ElasticsearchService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 10*tccommon.ReadRetryTimeout, time.Second, elasticsearchService.ElasticsearchInstanceRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	d.SetId(instanceId)

	return resourceTencentCloudElasticsearchRestartInstanceOperationRead(d, meta)
}

func resourceTencentCloudElasticsearchRestartInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchRestartInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
