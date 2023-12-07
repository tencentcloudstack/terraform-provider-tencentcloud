package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudElasticsearchLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchLogstashPipelineCreate,
		Read:   resourceTencentCloudElasticsearchLogstashPipelineRead,
		Update: resourceTencentCloudElasticsearchLogstashPipelineUpdate,
		Delete: resourceTencentCloudElasticsearchLogstashPipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Logstash instance id.",
			},
			"pipeline": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Pipeline information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pipeline id.",
						},
						"pipeline_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pipeline description information.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pipeline configuration content.",
						},
						"workers": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of Worker of pipe.",
						},
						"batch_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Pipe batch size.",
						},
						"batch_delay": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Pipeline batch processing delay.",
						},
						"queue_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pipeline buffer queue type.",
						},
						"queue_max_bytes": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pipeline buffer queue size.",
						},
						"queue_check_point_writes": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of pipeline buffer queue checkpoint writes.",
						},
					},
				},
			},

			"op_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Operation type. 1: save only; 2: save and deploy.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchLogstashPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_logstash_pipeline.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = elasticsearch.NewSaveAndDeployLogstashPipelineRequest()
		instanceId string
		pipelineId string
		opType     int
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "pipeline"); ok {
		logstashPipeline := elasticsearch.LogstashPipeline{}
		if v, ok := dMap["pipeline_id"]; ok {
			pipelineId = v.(string)
			logstashPipeline.PipelineId = helper.String(pipelineId)
		}
		if v, ok := dMap["pipeline_desc"]; ok {
			logstashPipeline.PipelineDesc = helper.String(v.(string))
		}
		if v, ok := dMap["config"]; ok {
			logstashPipeline.Config = helper.String(StringToBase64(v.(string)))
		}
		if v, ok := dMap["workers"]; ok {
			logstashPipeline.Workers = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["batch_size"]; ok {
			logstashPipeline.BatchSize = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["batch_delay"]; ok {
			logstashPipeline.BatchDelay = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["queue_type"]; ok {
			logstashPipeline.QueueType = helper.String(v.(string))
		}
		if v, ok := dMap["queue_max_bytes"]; ok {
			logstashPipeline.QueueMaxBytes = helper.String(v.(string))
		}
		if v, ok := dMap["queue_check_point_writes"]; ok {
			logstashPipeline.QueueCheckPointWrites = helper.IntUint64(v.(int))
		}
		request.Pipeline = &logstashPipeline
	}

	opType = d.Get("op_type").(int)
	request.OpType = helper.IntUint64(opType)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().SaveAndDeployLogstashPipeline(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create elasticsearch logstashPipeline failed, reason:%+v", logId, err)
		return err
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}
	targetStatue := "2"
	if opType == 1 {
		targetStatue = "0"
	}
	conf := BuildStateChangeConf([]string{}, []string{targetStatue}, 3*readRetryTimeout, time.Second, service.ElasticsearchLogstashPipelineStateRefreshFunc(instanceId, pipelineId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId + FILED_SP + pipelineId)

	return resourceTencentCloudElasticsearchLogstashPipelineRead(d, meta)
}

func resourceTencentCloudElasticsearchLogstashPipelineRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_logstash_pipeline.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	pipelineId := idSplit[1]
	_ = d.Set("instance_id", instanceId)
	logstashPipeline, err := service.DescribeElasticsearchLogstashPipelineById(ctx, instanceId, pipelineId)
	if err != nil {
		return err
	}

	if logstashPipeline != nil {
		pipelineMap := map[string]interface{}{}

		if logstashPipeline.PipelineId != nil {
			pipelineMap["pipeline_id"] = logstashPipeline.PipelineId
		}

		if logstashPipeline.PipelineDesc != nil {
			pipelineMap["pipeline_desc"] = logstashPipeline.PipelineDesc
		}

		if logstashPipeline.Config != nil {
			pipelineMap["config"] = logstashPipeline.Config
		}

		if logstashPipeline.Workers != nil {
			pipelineMap["workers"] = logstashPipeline.Workers
		}

		if logstashPipeline.BatchSize != nil {
			pipelineMap["batch_size"] = logstashPipeline.BatchSize
		}

		if logstashPipeline.BatchDelay != nil {
			pipelineMap["batch_delay"] = logstashPipeline.BatchDelay
		}

		if logstashPipeline.QueueType != nil {
			pipelineMap["queue_type"] = logstashPipeline.QueueType
		}

		if logstashPipeline.QueueMaxBytes != nil {
			pipelineMap["queue_max_bytes"] = logstashPipeline.QueueMaxBytes
		}

		if logstashPipeline.QueueCheckPointWrites != nil {
			pipelineMap["queue_check_point_writes"] = logstashPipeline.QueueCheckPointWrites
		}

		_ = d.Set("pipeline", []interface{}{pipelineMap})
	}

	return nil
}

func resourceTencentCloudElasticsearchLogstashPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_logstash_pipeline.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := elasticsearch.NewSaveAndDeployLogstashPipelineRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	pipelineId := idSplit[1]
	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("pipeline") {
		if dMap, ok := helper.InterfacesHeadMap(d, "pipeline"); ok {
			logstashPipeline := elasticsearch.LogstashPipeline{}
			if v, ok := dMap["pipeline_id"]; ok {
				logstashPipeline.PipelineId = helper.String(v.(string))
			}
			if v, ok := dMap["pipeline_desc"]; ok {
				logstashPipeline.PipelineDesc = helper.String(v.(string))
			}
			if v, ok := dMap["config"]; ok {
				logstashPipeline.Config = helper.String(StringToBase64(v.(string)))
			}
			if v, ok := dMap["workers"]; ok {
				logstashPipeline.Workers = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["batch_size"]; ok {
				logstashPipeline.BatchSize = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["batch_delay"]; ok {
				logstashPipeline.BatchDelay = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["queue_type"]; ok {
				logstashPipeline.QueueType = helper.String(v.(string))
			}
			if v, ok := dMap["queue_max_bytes"]; ok {
				logstashPipeline.QueueMaxBytes = helper.String(v.(string))
			}
			if v, ok := dMap["queue_check_point_writes"]; ok {
				logstashPipeline.QueueCheckPointWrites = helper.IntUint64(v.(int))
			}
			request.Pipeline = &logstashPipeline
		}
	}

	opType := d.Get("op_type").(int)
	request.OpType = helper.IntUint64(opType)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().SaveAndDeployLogstashPipeline(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update elasticsearch logstashPipeline failed, reason:%+v", logId, err)
		return err
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}
	targetStatue := "2"
	if opType == 1 {
		targetStatue = "0"
	}
	conf := BuildStateChangeConf([]string{}, []string{targetStatue}, 3*readRetryTimeout, time.Second, service.ElasticsearchLogstashPipelineStateRefreshFunc(instanceId, pipelineId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudElasticsearchLogstashPipelineRead(d, meta)
}

func resourceTencentCloudElasticsearchLogstashPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_logstash_pipeline.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	pipelineId := idSplit[1]

	if err := service.DeleteElasticsearchLogstashPipelineById(ctx, instanceId, pipelineId); err != nil {
		return err
	}
	conf := BuildStateChangeConf([]string{}, []string{"-99"}, 3*readRetryTimeout, time.Second, service.ElasticsearchLogstashPipelineStateRefreshFunc(instanceId, pipelineId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}
