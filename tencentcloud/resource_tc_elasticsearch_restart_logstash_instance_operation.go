package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudElasticsearchRestartLogstashInstanceOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_elasticsearch_restart_logstash_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().RestartLogstashInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch restartLogstashInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"1"}, 3*readRetryTimeout, time.Second, service.ElasticsearchLogstashStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId)

	return resourceTencentCloudElasticsearchRestartLogstashInstanceOperationRead(d, meta)
}

func resourceTencentCloudElasticsearchRestartLogstashInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_restart_logstash_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchRestartLogstashInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_restart_logstash_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
