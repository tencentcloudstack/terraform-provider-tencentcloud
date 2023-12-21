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

func ResourceTencentCloudElasticsearchRestartKibanaOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchRestartKibanaOperationCreate,
		Read:   resourceTencentCloudElasticsearchRestartKibanaOperationRead,
		Delete: resourceTencentCloudElasticsearchRestartKibanaOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchRestartKibanaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_kibana_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	var (
		request    = elasticsearch.NewRestartKibanaRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().RestartKibana(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch RestartKibanaOperation failed, reason:%+v", logId, err)
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

	return resourceTencentCloudElasticsearchRestartKibanaOperationRead(d, meta)
}

func resourceTencentCloudElasticsearchRestartKibanaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_kibana_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchRestartKibanaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_kibana_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
