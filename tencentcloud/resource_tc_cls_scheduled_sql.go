/*
Provides a resource to create a cls scheduled_sql

Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-example-logset"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "tf-example-topic"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_scheduled_sql" "scheduled_sql" {
  src_topic_id = tencentcloud_cls_topic.topic.id
  name = "tf-example-task"
  enable_flag = 1
  dst_resource {
    topic_id = tencentcloud_cls_topic.topic.id
    region = "ap-guangzhou"
    biz_type = 0
    metric_name = "test"

  }
  scheduled_sql_content = "xxx"
  process_start_time = 1690515360000
  process_type = 1
  process_period = 10
  process_time_window = "@m-15m,@m"
  process_delay = 5
  src_topic_region = "ap-guangzhou"
  process_end_time = 1690515360000
  syntax_rule = 0
}
```

Import

cls scheduled_sql can be imported using the id, e.g.

```
terraform import tencentcloud_cls_scheduled_sql.scheduled_sql scheduled_sql_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsScheduledSql() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsScheduledSqlCreate,
		Read:   resourceTencentCloudClsScheduledSqlRead,
		Update: resourceTencentCloudClsScheduledSqlUpdate,
		Delete: resourceTencentCloudClsScheduledSqlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"src_topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "src topic id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "task name.",
			},

			"enable_flag": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "task enable flag.",
			},

			"dst_resource": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "scheduled slq dst resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "dst topic id.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "topic region.",
						},
						"biz_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "topic type.",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "metric name.",
						},
					},
				},
			},

			"scheduled_sql_content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "scheduled sql content.",
			},

			"process_start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "process start timestamp.",
			},

			"process_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "process type.",
			},

			"process_period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "process period.",
			},

			"process_time_window": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "process time window.",
			},

			"process_delay": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "process delay.",
			},

			"src_topic_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "src topic region.",
			},

			"process_end_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "process end timestamp.",
			},

			"syntax_rule": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "syntax rule.",
			},
		},
	}
}

func resourceTencentCloudClsScheduledSqlCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_scheduled_sql.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateScheduledSqlRequest()
		response = cls.NewCreateScheduledSqlResponse()
		taskId   string
	)
	if v, ok := d.GetOk("src_topic_id"); ok {
		request.SrcTopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_flag"); ok {
		request.EnableFlag = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "dst_resource"); ok {
		scheduledSqlResouceInfo := cls.ScheduledSqlResouceInfo{}
		if v, ok := dMap["topic_id"]; ok {
			scheduledSqlResouceInfo.TopicId = helper.String(v.(string))
		}
		if v, ok := dMap["region"]; ok {
			scheduledSqlResouceInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["biz_type"]; ok {
			scheduledSqlResouceInfo.BizType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["metric_name"]; ok {
			scheduledSqlResouceInfo.MetricName = helper.String(v.(string))
		}
		request.DstResource = &scheduledSqlResouceInfo
	}

	if v, ok := d.GetOk("scheduled_sql_content"); ok {
		request.ScheduledSqlContent = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("process_start_time"); ok {
		request.ProcessStartTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("process_type"); ok {
		request.ProcessType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("process_period"); ok {
		request.ProcessPeriod = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("process_time_window"); ok {
		request.ProcessTimeWindow = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("process_delay"); ok {
		request.ProcessDelay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("src_topic_region"); ok {
		request.SrcTopicRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("process_end_time"); ok {
		request.ProcessEndTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("syntax_rule"); ok {
		request.SyntaxRule = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateScheduledSql(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls scheduledSql failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudClsScheduledSqlRead(d, meta)
}

func resourceTencentCloudClsScheduledSqlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_scheduled_sql.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	scheduledSqlId := d.Id()

	scheduledSql, err := service.DescribeClsScheduledSqlById(ctx, scheduledSqlId)
	if err != nil {
		return err
	}

	if scheduledSql == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsScheduledSql` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if scheduledSql.SrcTopicId != nil {
		_ = d.Set("src_topic_id", scheduledSql.SrcTopicId)
	}

	if scheduledSql.Name != nil {
		_ = d.Set("name", scheduledSql.Name)
	}

	if scheduledSql.EnableFlag != nil {
		_ = d.Set("enable_flag", scheduledSql.EnableFlag)
	}

	if scheduledSql.DstResource != nil {
		dstResourceMap := map[string]interface{}{}

		if scheduledSql.DstResource.TopicId != nil {
			dstResourceMap["topic_id"] = scheduledSql.DstResource.TopicId
		}

		if scheduledSql.DstResource.Region != nil {
			dstResourceMap["region"] = scheduledSql.DstResource.Region
		}

		if scheduledSql.DstResource.BizType != nil {
			dstResourceMap["biz_type"] = scheduledSql.DstResource.BizType
		}

		if scheduledSql.DstResource.MetricName != nil {
			dstResourceMap["metric_name"] = scheduledSql.DstResource.MetricName
		}

		_ = d.Set("dst_resource", []interface{}{dstResourceMap})
	}

	if scheduledSql.ScheduledSqlContent != nil {
		_ = d.Set("scheduled_sql_content", scheduledSql.ScheduledSqlContent)
	}

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", *scheduledSql.ProcessStartTime, location)
	if err != nil {
		return err
	}

	if scheduledSql.ProcessStartTime != nil {
		_ = d.Set("process_start_time", startTime.UnixNano()/int64(time.Millisecond))
	}

	if scheduledSql.ProcessType != nil {
		_ = d.Set("process_type", scheduledSql.ProcessType)
	}

	if scheduledSql.ProcessPeriod != nil {
		_ = d.Set("process_period", scheduledSql.ProcessPeriod)
	}

	if scheduledSql.ProcessTimeWindow != nil {
		_ = d.Set("process_time_window", scheduledSql.ProcessTimeWindow)
	}

	if scheduledSql.ProcessDelay != nil {
		_ = d.Set("process_delay", scheduledSql.ProcessDelay)
	}

	if scheduledSql.SrcTopicRegion != nil {
		_ = d.Set("src_topic_region", scheduledSql.SrcTopicRegion)
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", *scheduledSql.ProcessStartTime)
	if err != nil {
		return err
	}
	if scheduledSql.ProcessEndTime != nil {
		_ = d.Set("process_end_time", endTime.Unix())
	}

	if scheduledSql.SyntaxRule != nil {
		_ = d.Set("syntax_rule", scheduledSql.SyntaxRule)
	}

	return nil
}

func resourceTencentCloudClsScheduledSqlUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_scheduled_sql.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cls.NewModifyScheduledSqlRequest()

	scheduledSqlId := d.Id()

	request.TaskId = &scheduledSqlId

	immutableArgs := []string{"src_topic_id", "name", "enable_flag", "dst_resource", "scheduled_sql_content", "process_start_time", "process_type", "process_period", "process_time_window", "process_delay", "src_topic_region", "process_end_time", "syntax_rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("src_topic_id") {
		if v, ok := d.GetOk("src_topic_id"); ok {
			request.SrcTopicId = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_flag") {
		if v, ok := d.GetOkExists("enable_flag"); ok {
			request.EnableFlag = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("dst_resource") {
		if dMap, ok := helper.InterfacesHeadMap(d, "dst_resource"); ok {
			scheduledSqlResouceInfo := cls.ScheduledSqlResouceInfo{}
			if v, ok := dMap["topic_id"]; ok {
				scheduledSqlResouceInfo.TopicId = helper.String(v.(string))
			}
			if v, ok := dMap["region"]; ok {
				scheduledSqlResouceInfo.Region = helper.String(v.(string))
			}
			if v, ok := dMap["biz_type"]; ok {
				scheduledSqlResouceInfo.BizType = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["metric_name"]; ok {
				scheduledSqlResouceInfo.MetricName = helper.String(v.(string))
			}
			request.DstResource = &scheduledSqlResouceInfo
		}
	}

	if d.HasChange("scheduled_sql_content") {
		if v, ok := d.GetOk("scheduled_sql_content"); ok {
			request.ScheduledSqlContent = helper.String(v.(string))
		}
	}

	if d.HasChange("process_period") {
		if v, ok := d.GetOkExists("process_period"); ok {
			request.ProcessPeriod = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("process_time_window") {
		if v, ok := d.GetOk("process_time_window"); ok {
			request.ProcessTimeWindow = helper.String(v.(string))
		}
	}

	if d.HasChange("process_delay") {
		if v, ok := d.GetOkExists("process_delay"); ok {
			request.ProcessDelay = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("src_topic_region") {
		if v, ok := d.GetOk("src_topic_region"); ok {
			request.SrcTopicRegion = helper.String(v.(string))
		}
	}

	if d.HasChange("syntax_rule") {
		if v, ok := d.GetOkExists("syntax_rule"); ok {
			request.SyntaxRule = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyScheduledSql(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cls scheduledSql failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsScheduledSqlRead(d, meta)
}

func resourceTencentCloudClsScheduledSqlDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_scheduled_sql.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	scheduledSqlId := d.Id()
	var srcTopicId string
	if v, ok := d.GetOk("src_topic_id"); ok {
		srcTopicId = v.(string)
	}
	if err := service.DeleteClsScheduledSqlById(ctx, scheduledSqlId, srcTopicId); err != nil {
		return err
	}

	return nil
}
