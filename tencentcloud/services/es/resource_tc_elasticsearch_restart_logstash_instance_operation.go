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

func ResourceTencentCloudElasticsearchRestartLogstashInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchRestartLogstashInstanceOperationCreate,
		Read:   resourceTencentCloudElasticsearchRestartLogstashInstanceOperationRead,
		Delete: resourceTencentCloudElasticsearchRestartLogstashInstanceOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Restart type, 0 full restart, 1 rolling restart.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchRestartLogstashInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_logstash_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = elasticsearch.NewRestartLogstashInstanceRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, _ := d.GetOk("type"); v != nil {
		request.Type = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().RestartLogstashInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch restartLogstashInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 3*tccommon.ReadRetryTimeout, time.Second, service.ElasticsearchLogstashStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId)

	return resourceTencentCloudElasticsearchRestartLogstashInstanceOperationRead(d, meta)
}

func resourceTencentCloudElasticsearchRestartLogstashInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_logstash_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchRestartLogstashInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_logstash_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
