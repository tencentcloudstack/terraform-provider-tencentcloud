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

func ResourceTencentCloudElasticsearchStartLogstashPipelineOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchStartLogstashPipelineOperationCreate,
		Read:   resourceTencentCloudElasticsearchStartLogstashPipelineOperationRead,
		Delete: resourceTencentCloudElasticsearchStartLogstashPipelineOperationDelete,
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

			"pipeline_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Pipeline id.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchStartLogstashPipelineOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_start_logstash_pipeline_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = elasticsearch.NewStartLogstashPipelinesRequest()
		instanceId string
		pipelineId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("pipeline_id"); ok {
		pipelineId = v.(string)
		request.PipelineIds = []*string{&pipelineId}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().StartLogstashPipelines(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch startLogstashPipelineOperation failed, reason:%+v", logId, err)
		return err
	}

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 3*tccommon.ReadRetryTimeout, time.Second, service.ElasticsearchLogstashPipelineStateRefreshFunc(instanceId, pipelineId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId + tccommon.FILED_SP + pipelineId)

	return resourceTencentCloudElasticsearchStartLogstashPipelineOperationRead(d, meta)
}

func resourceTencentCloudElasticsearchStartLogstashPipelineOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_start_logstash_pipeline_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchStartLogstashPipelineOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_start_logstash_pipeline_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
