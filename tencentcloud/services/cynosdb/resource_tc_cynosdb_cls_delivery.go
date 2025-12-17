package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	cynosdbv20190107 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbClsDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClsDeliveryCreate,
		Read:   resourceTencentCloudCynosdbClsDeliveryRead,
		Update: resourceTencentCloudCynosdbClsDeliveryUpdate,
		Delete: resourceTencentCloudCynosdbClsDeliveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Intance ID.",
			},

			"cls_info_list": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Log shipping configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Log delivery area.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Log topic ID.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Log topic name.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Log set ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Log set name.",
						},
					},
				},
			},

			"log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Log type.",
			},

			"running_status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Delivery status. true: Enabled; false: Disabled.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClsDeliveryCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cls_delivery.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = cynosdbv20190107.NewCreateCLSDeliveryRequest()
		response   = cynosdbv20190107.NewCreateCLSDeliveryResponse()
		instanceId string
		topicId    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("cls_info_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			clsInfo := cynosdb.CLSInfo{}
			if v, ok := dMap["region"]; ok && v.(string) != "" {
				clsInfo.Region = helper.String(v.(string))
			}

			if v, ok := dMap["topic_id"]; ok && v.(string) != "" {
				clsInfo.TopicId = helper.String(v.(string))
				topicId = v.(string)

				clsInfo.TopicOperation = helper.String("reuse")
				clsInfo.GroupOperation = helper.String("reuse")
			}

			if v, ok := dMap["topic_name"]; ok && v.(string) != "" {
				clsInfo.TopicName = helper.String(v.(string))

				clsInfo.TopicOperation = helper.String("create")
				clsInfo.GroupOperation = helper.String("create")
			}

			if v, ok := dMap["group_id"]; ok && v.(string) != "" {
				clsInfo.GroupId = helper.String(v.(string))

				clsInfo.TopicOperation = helper.String("reuse")
				clsInfo.GroupOperation = helper.String("reuse")
			}

			if v, ok := dMap["group_name"]; ok && v.(string) != "" {
				clsInfo.GroupName = helper.String(v.(string))

				clsInfo.TopicOperation = helper.String("create")
				clsInfo.GroupOperation = helper.String("create")
			}

			request.CLSInfoList = append(request.CLSInfoList, &clsInfo)
		}
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().CreateCLSDeliveryWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cynosdb cls delivery failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cynosdb cls delivery failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	taskId := *response.Response.TaskId
	waitReq := cynosdb.NewDescribeTasksRequest()
	waitReq.Filters = []*cynosdb.QueryFilter{
		{
			Values:     helper.Strings([]string{helper.Int64ToStr(taskId)}),
			Names:      helper.Strings([]string{"TaskId"}),
			ExactMatch: helper.Bool(true),
		},
	}

	reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeTasks(waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskList == nil || len(result.Response.TaskList) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
		}

		task := result.Response.TaskList[0]
		if task.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Status is nil."))
		}

		if *task.Status == "success" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Tasks is still running, status is %s.", *task.Status))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe tasks failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// get topicId
	if topicId == "" {
		service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		respData, err := service.DescribeCynosdbClsDeliveryById(ctx, instanceId, topicId)
		if err != nil {
			return err
		}

		if respData == nil {
			return fmt.Errorf("Describe instance cls log delivery failed, Response is nil.")
		}

		if respData.TopicId == nil {
			return fmt.Errorf("TopicId is nil.")
		}

		topicId = *respData.TopicId
	}

	d.SetId(strings.Join([]string{instanceId, topicId}, tccommon.FILED_SP))

	// set status
	if v, ok := d.GetOkExists("running_status"); ok {
		if !v.(bool) {
			request := cynosdbv20190107.NewStopCLSDeliveryRequest()
			response := cynosdbv20190107.NewStopCLSDeliveryResponse()
			request.InstanceId = helper.String(instanceId)
			request.CLSTopicIds = helper.Strings([]string{topicId})
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StopCLSDeliveryWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.TaskId == nil {
					return resource.NonRetryableError(fmt.Errorf("Stop cynosdb cls delivery failed, Response is nil."))
				}

				response = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s stop cynosdb cls delivery failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			taskId := *response.Response.TaskId
			waitReq := cynosdb.NewDescribeTasksRequest()
			waitReq.Filters = []*cynosdb.QueryFilter{
				{
					Values:     helper.Strings([]string{helper.Int64ToStr(taskId)}),
					Names:      helper.Strings([]string{"TaskId"}),
					ExactMatch: helper.Bool(true),
				},
			}

			reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeTasks(waitReq)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.TaskList == nil || len(result.Response.TaskList) == 0 {
					return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
				}

				task := result.Response.TaskList[0]
				if task.Status == nil {
					return resource.NonRetryableError(fmt.Errorf("Status is nil."))
				}

				if *task.Status == "success" {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("Tasks is still running, status is %s.", *task.Status))
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s describe tasks failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}
	return resourceTencentCloudCynosdbClsDeliveryRead(d, meta)
}

func resourceTencentCloudCynosdbClsDeliveryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cls_delivery.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	topicId := idSplit[1]

	respData, err := service.DescribeCynosdbClsDeliveryById(ctx, instanceId, topicId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cynosdb_cls_delivery` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", *respData.InstanceId)
	}

	dMap := make(map[string]interface{}, 0)
	if respData.Region != nil {
		dMap["region"] = *respData.Region
	}

	if respData.TopicId != nil {
		dMap["topic_id"] = *respData.TopicId
	}

	if respData.TopicName != nil {
		dMap["topic_name"] = trimPrefixSuffix(*respData.TopicName, "cloud_cynos_", "_topic")
	}

	if respData.GroupId != nil {
		dMap["group_id"] = *respData.GroupId
	}

	if respData.GroupName != nil {
		dMap["group_name"] = trimPrefixSuffix(*respData.GroupName, "cloud_cynos_", "_logset")
	}

	_ = d.Set("cls_info_list", []interface{}{dMap})

	if respData.LogType != nil {
		_ = d.Set("log_type", *respData.LogType)
	}

	if respData.Status != nil {
		if *respData.Status == "running" {
			_ = d.Set("running_status", true)
		} else if *respData.Status == "offlined" {
			_ = d.Set("running_status", false)
		}
	}

	return nil
}

func resourceTencentCloudCynosdbClsDeliveryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cls_delivery.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	topicId := idSplit[1]

	if d.HasChange("running_status") {
		if v, ok := d.GetOkExists("running_status"); ok {
			if v.(bool) {
				request := cynosdbv20190107.NewStartCLSDeliveryRequest()
				response := cynosdbv20190107.NewStartCLSDeliveryResponse()
				request.InstanceId = helper.String(instanceId)
				request.CLSTopicIds = helper.Strings([]string{topicId})
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StartCLSDeliveryWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.TaskId == nil {
						return resource.NonRetryableError(fmt.Errorf("Start cynosdb cls delivery failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s start cynosdb cls delivery failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				taskId := *response.Response.TaskId
				waitReq := cynosdb.NewDescribeTasksRequest()
				waitReq.Filters = []*cynosdb.QueryFilter{
					{
						Values:     helper.Strings([]string{helper.Int64ToStr(taskId)}),
						Names:      helper.Strings([]string{"TaskId"}),
						ExactMatch: helper.Bool(true),
					},
				}

				reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeTasks(waitReq)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.TaskList == nil || len(result.Response.TaskList) == 0 {
						return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
					}

					task := result.Response.TaskList[0]
					if task.Status == nil {
						return resource.NonRetryableError(fmt.Errorf("Status is nil."))
					}

					if *task.Status == "success" {
						return nil
					}

					return resource.RetryableError(fmt.Errorf("Tasks is still running, status is %s.", *task.Status))
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s describe tasks failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			} else {
				request := cynosdbv20190107.NewStopCLSDeliveryRequest()
				response := cynosdbv20190107.NewStopCLSDeliveryResponse()
				request.InstanceId = helper.String(instanceId)
				request.CLSTopicIds = helper.Strings([]string{topicId})
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StopCLSDeliveryWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.TaskId == nil {
						return resource.NonRetryableError(fmt.Errorf("Stop cynosdb cls delivery failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s stop cynosdb cls delivery failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				taskId := *response.Response.TaskId
				waitReq := cynosdb.NewDescribeTasksRequest()
				waitReq.Filters = []*cynosdb.QueryFilter{
					{
						Values:     helper.Strings([]string{helper.Int64ToStr(taskId)}),
						Names:      helper.Strings([]string{"TaskId"}),
						ExactMatch: helper.Bool(true),
					},
				}

				reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeTasks(waitReq)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.TaskList == nil || len(result.Response.TaskList) == 0 {
						return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
					}

					task := result.Response.TaskList[0]
					if task.Status == nil {
						return resource.NonRetryableError(fmt.Errorf("Status is nil."))
					}

					if *task.Status == "success" {
						return nil
					}

					return resource.RetryableError(fmt.Errorf("Tasks is still running, status is %s.", *task.Status))
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s describe tasks failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	return resourceTencentCloudCynosdbClsDeliveryRead(d, meta)
}

func resourceTencentCloudCynosdbClsDeliveryDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cls_delivery.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = cynosdbv20190107.NewDeleteCLSDeliveryRequest()
		response = cynosdbv20190107.NewDeleteCLSDeliveryResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	topicId := idSplit[1]

	// get status
	respData, err := service.DescribeCynosdbClsDeliveryById(ctx, instanceId, topicId)
	if err != nil {
		return err
	}

	if respData == nil {
		return fmt.Errorf("Describe instance cls log delivery failed, Response is nil.")
	}

	if respData.Status == nil {
		return fmt.Errorf("Status is nil.")
	}

	// stop first
	if *respData.Status == "running" {
		stopReq := cynosdbv20190107.NewStopCLSDeliveryRequest()
		stopResp := cynosdbv20190107.NewStopCLSDeliveryResponse()
		stopReq.InstanceId = helper.String(instanceId)
		stopReq.CLSTopicIds = helper.Strings([]string{topicId})
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StopCLSDeliveryWithContext(ctx, stopReq)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, stopReq.GetAction(), stopReq.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.TaskId == nil {
				return resource.NonRetryableError(fmt.Errorf("Stop cynosdb cls delivery failed, Response is nil."))
			}

			stopResp = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s stop cynosdb cls delivery failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait stop
		taskId := *stopResp.Response.TaskId
		waitReq := cynosdb.NewDescribeTasksRequest()
		waitReq.Filters = []*cynosdb.QueryFilter{
			{
				Values:     helper.Strings([]string{helper.Int64ToStr(taskId)}),
				Names:      helper.Strings([]string{"TaskId"}),
				ExactMatch: helper.Bool(true),
			},
		}

		reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeTasks(waitReq)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.TaskList == nil || len(result.Response.TaskList) == 0 {
				return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
			}

			task := result.Response.TaskList[0]
			if task.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Status is nil."))
			}

			if *task.Status == "success" {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("Tasks is still running, status is %s.", *task.Status))
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s describe tasks failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	// delete
	request.InstanceId = &instanceId
	request.CLSTopicIds = helper.Strings([]string{topicId})
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DeleteCLSDeliveryWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete cynosdb cls delivery failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cynosdb cls delivery failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait delete
	taskId := *response.Response.TaskId
	waitReq := cynosdb.NewDescribeTasksRequest()
	waitReq.Filters = []*cynosdb.QueryFilter{
		{
			Values:     helper.Strings([]string{helper.Int64ToStr(taskId)}),
			Names:      helper.Strings([]string{"TaskId"}),
			ExactMatch: helper.Bool(true),
		},
	}

	reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeTasks(waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskList == nil || len(result.Response.TaskList) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
		}

		task := result.Response.TaskList[0]
		if task.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Status is nil."))
		}

		if *task.Status == "success" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Tasks is still running, status is %s.", *task.Status))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe tasks failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func trimPrefixSuffix(s, prefix, suffix string) string {
	trimmed := strings.TrimPrefix(s, prefix)
	trimmed = strings.TrimSuffix(trimmed, suffix)
	return trimmed
}
