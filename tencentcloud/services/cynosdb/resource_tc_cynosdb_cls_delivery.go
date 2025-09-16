package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
				Description: "Instance ID.",
			},

			"cls_info_list": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Log delivery configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log topic operations: Optional: create or reuse. create: Creates a new log topic using TopicName. reuse: Reuses an existing log topic using TopicId. Reusing an existing log topic and creating a new log set is not allowed.",
						},
						"group_operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log set operations: Optional: create or reuse. create: Creates a new log set, using the GroupName. reuse: Reuses an existing log set, using the GroupId. The combination of reusing an existing log topic and creating a new log set is not allowed.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log delivery region.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Log topic ID.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Log topic name.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Log set ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Log set name.",
						},
					},
				},
			},

			"log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Log type.",
			},

			"enable_cls_delivery": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable CLS delivery. Default value: true (enabled).",
			},

			// computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Delivery status. running, offlined.",
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
		service    = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = cynosdbv20190107.NewCreateCLSDeliveryRequest()
		response   = cynosdbv20190107.NewCreateCLSDeliveryResponse()
		instanceId string
		groupId    string
		topicId    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("cls_info_list"); ok {
		for _, item := range v.([]interface{}) {
			cLSInfoListMap := item.(map[string]interface{})
			cLSInfo := cynosdbv20190107.CLSInfo{}
			if v, ok := cLSInfoListMap["topic_operation"].(string); ok && v != "" {
				cLSInfo.TopicOperation = helper.String(v)
			}

			if v, ok := cLSInfoListMap["group_operation"].(string); ok && v != "" {
				cLSInfo.GroupOperation = helper.String(v)
			}

			if v, ok := cLSInfoListMap["region"].(string); ok && v != "" {
				cLSInfo.Region = helper.String(v)
			}

			if v, ok := cLSInfoListMap["topic_id"].(string); ok && v != "" {
				cLSInfo.TopicId = helper.String(v)
			}

			if v, ok := cLSInfoListMap["topic_name"].(string); ok && v != "" {
				cLSInfo.TopicName = helper.String(v)
			}

			if v, ok := cLSInfoListMap["group_id"].(string); ok && v != "" {
				cLSInfo.GroupId = helper.String(v)
			}

			if v, ok := cLSInfoListMap["group_name"].(string); ok && v != "" {
				cLSInfo.GroupName = helper.String(v)
			}

			request.CLSInfoList = append(request.CLSInfoList, &cLSInfo)
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

		if result == nil || result.Response == nil {
			return tccommon.RetryError(fmt.Errorf("Create cynosdb cls delivery failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cynosdb cls delivery failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}

	// wait for create
	taskId := response.Response.TaskId
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(*taskId, 10), []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// get cls info
	clsReq := cynosdbv20190107.NewDescribeTasksRequest()
	clsResp := cynosdbv20190107.NewDescribeTasksResponse()
	clsReq.Filters = []*cynosdb.QueryFilter{
		{
			ExactMatch: helper.Bool(true),
			Names:      helper.Strings([]string{"TaskId"}),
			Values:     helper.Strings([]string{strconv.FormatInt(*taskId, 10)}),
		},
	}
	reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeTasksWithContext(ctx, clsReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, clsReq.GetAction(), clsReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return tccommon.RetryError(fmt.Errorf("Describe cynosdb cls delivery failed, Response is nil."))
		}

		clsResp = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe cynosdb cls delivery failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if clsResp.Response.TaskList == nil || len(clsResp.Response.TaskList) == 0 {
		return fmt.Errorf("TaskList is nil.")
	}

	if len(clsResp.Response.TaskList) != 1 {
		return fmt.Errorf("TaskList length is not 1.")
	}

	if clsResp.Response.TaskList[0].InstanceCLSDeliveryInfos == nil || len(clsResp.Response.TaskList[0].InstanceCLSDeliveryInfos) == 0 {
		return fmt.Errorf("InstanceCLSDeliveryInfos is nil.")
	}

	instanceCLSDeliveryInfo := clsResp.Response.TaskList[0].InstanceCLSDeliveryInfos[0]
	if instanceCLSDeliveryInfo.GroupId == nil || instanceCLSDeliveryInfo.TopicId == nil {
		return fmt.Errorf("GroupId or TopicId is nil.")
	}

	groupId = *instanceCLSDeliveryInfo.GroupId
	topicId = *instanceCLSDeliveryInfo.TopicId
	d.SetId(strings.Join([]string{instanceId, groupId, topicId}, tccommon.FILED_SP))

	// operate cls delivery
	if v, ok := d.GetOkExists("enable_cls_delivery"); ok {
		if !v.(bool) {
			// stop cls delivery
			request := cynosdb.NewStopCLSDeliveryRequest()
			response := cynosdb.NewStopCLSDeliveryResponse()
			request.InstanceId = helper.String(instanceId)
			request.CLSTopicIds = helper.Strings([]string{topicId})
			if v, ok := d.GetOk("log_type"); ok {
				request.LogType = helper.String(v.(string))
			}
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StopCLSDeliveryWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return tccommon.RetryError(fmt.Errorf("Stop cynosdb cls delivery failed, Response is nil."))
				}

				response = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s stop cynosdb cls delivery failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if response.Response.TaskId == nil {
				return fmt.Errorf("TaskId is nil.")
			}

			// wait for start
			taskId := response.Response.TaskId
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(*taskId, 10), []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
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
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	groupId := idSplit[1]
	topicId := idSplit[2]

	respData, err := service.DescribeCynosdbClsDeliveryById(ctx, instanceId, groupId, topicId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `cynosdb_cls_delivery` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	tmpList := make([]map[string]interface{}, 0, 1)
	tmpMap := make(map[string]interface{}, 0)
	if respData.Region != nil {
		tmpMap["region"] = respData.Region
	}

	if respData.TopicId != nil {
		tmpMap["topic_id"] = respData.TopicId
	}

	if respData.TopicName != nil {
		tmpMap["topic_name"] = respData.TopicName
	}

	if respData.GroupId != nil {
		tmpMap["group_id"] = respData.GroupId
	}

	if respData.GroupName != nil {
		tmpMap["group_name"] = respData.GroupName
	}

	tmpList = append(tmpList, tmpMap)
	_ = d.Set("cls_info_list", tmpList)

	if respData.LogType != nil {
		_ = d.Set("log_type", respData.LogType)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	return nil
}

func resourceTencentCloudCynosdbClsDeliveryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cls_delivery.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	topicId := idSplit[2]

	if d.HasChange("enable_cls_delivery") {
		if v, ok := d.GetOkExists("enable_cls_delivery"); ok {
			if v.(bool) {
				// start cls delivery
				request := cynosdb.NewStartCLSDeliveryRequest()
				response := cynosdb.NewStartCLSDeliveryResponse()
				request.InstanceId = helper.String(instanceId)
				request.CLSTopicIds = helper.Strings([]string{topicId})
				if v, ok := d.GetOk("log_type"); ok {
					request.LogType = helper.String(v.(string))
				}
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StartCLSDeliveryWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return tccommon.RetryError(fmt.Errorf("Start cynosdb cls delivery failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s start cynosdb cls delivery failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				if response.Response.TaskId == nil {
					return fmt.Errorf("TaskId is nil.")
				}

				// wait for start
				taskId := response.Response.TaskId
				conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(*taskId, 10), []string{}))
				if _, e := conf.WaitForState(); e != nil {
					return e
				}
			} else {
				// stop cls delivery
				request := cynosdb.NewStopCLSDeliveryRequest()
				response := cynosdb.NewStopCLSDeliveryResponse()
				request.InstanceId = helper.String(instanceId)
				request.CLSTopicIds = helper.Strings([]string{topicId})
				if v, ok := d.GetOk("log_type"); ok {
					request.LogType = helper.String(v.(string))
				}
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StopCLSDeliveryWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return tccommon.RetryError(fmt.Errorf("Stop cynosdb cls delivery failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s stop cynosdb cls delivery failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				if response.Response.TaskId == nil {
					return fmt.Errorf("TaskId is nil.")
				}

				// wait for start
				taskId := response.Response.TaskId
				conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(*taskId, 10), []string{}))
				if _, e := conf.WaitForState(); e != nil {
					return e
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
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	groupId := idSplit[1]
	topicId := idSplit[2]

	respData, err := service.DescribeCynosdbClsDeliveryById(ctx, instanceId, groupId, topicId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `cynosdb_cls_delivery` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	// check status
	status := respData.Status
	if status != nil && *status == "running" {
		// stop cls delivery first
		request := cynosdb.NewStopCLSDeliveryRequest()
		response := cynosdb.NewStopCLSDeliveryResponse()
		request.InstanceId = helper.String(instanceId)
		request.CLSTopicIds = helper.Strings([]string{topicId})
		if v, ok := d.GetOk("log_type"); ok {
			request.LogType = helper.String(v.(string))
		}
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().StopCLSDeliveryWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return tccommon.RetryError(fmt.Errorf("Stop cynosdb cls delivery failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s stop cynosdb cls delivery failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if response.Response.TaskId == nil {
			return fmt.Errorf("TaskId is nil.")
		}

		// wait for start
		taskId := response.Response.TaskId
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(*taskId, 10), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	request := cynosdbv20190107.NewDeleteCLSDeliveryRequest()
	response := cynosdbv20190107.NewDeleteCLSDeliveryResponse()
	request.InstanceId = helper.String(instanceId)
	request.CLSTopicIds = helper.Strings([]string{topicId})
	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DeleteCLSDeliveryWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return tccommon.RetryError(fmt.Errorf("Delete cynosdb cls delivery failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cynosdb cls delivery failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}

	// wait for delete
	taskId := response.Response.TaskId
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(*taskId, 10), []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
