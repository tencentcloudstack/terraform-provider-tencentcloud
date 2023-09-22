/*
Provides a resource to create a cat task_set

Example Usage

```hcl
resource "tencentcloud_cat_task_set" "task_set" {
  batch_tasks {
    name           = "demo"
    target_address = "http://www.baidu.com"
  }
  task_type     = 5
  nodes         = ["12136", "12137", "12138", "12141", "12144"]
  interval      = 6
  parameters    = jsonencode(
  {
    "ipType"            = 0,
    "grabBag"           = 0,
    "filterIp"          = 0,
    "netIcmpOn"         = 1,
    "netIcmpActivex"    = 0,
    "netIcmpTimeout"    = 20,
    "netIcmpInterval"   = 0.5,
    "netIcmpNum"        = 20,
    "netIcmpSize"       = 32,
    "netIcmpDataCut"    = 1,
    "netDnsOn"          = 1,
    "netDnsTimeout"     = 5,
    "netDnsQuerymethod" = 1,
    "netDnsNs"          = "",
    "netDigOn"          = 1,
    "netDnsServer"      = 2,
    "netTracertOn"      = 1,
    "netTracertTimeout" = 60,
    "netTracertNum"     = 30,
    "whiteList"         = "",
    "blackList"         = "",
    "netIcmpActivexStr" = ""
  }
  )
  task_category = 1
  cron          = "* 0-6 * * *"
  tags          = {
    "createdBy" = "terraform"
  }
}
```
Import

cat task_set can be imported using the id, e.g.
```
$ terraform import tencentcloud_cat_task_set.task_set taskSet_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCatTaskSet() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCatTaskSetRead,
		Create: resourceTencentCloudCatTaskSetCreate,
		Update: resourceTencentCloudCatTaskSetUpdate,
		Delete: resourceTencentCloudCatTaskSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"batch_tasks": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Batch task name address.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task name.",
						},
						"target_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target address.",
						},
					},
				},
			},

			"task_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Task Type 1:Page Performance, 2:File upload,3:File Download,4:Port performance 5:Audio and video.",
			},

			"nodes": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Task Nodes.",
			},

			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task Id.",
			},

			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Task interval minutes in (1,5,10,15,30,60,120,240).",
			},

			"parameters": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "tasks parameters.",
			},

			"task_category": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Task category,1:PC,2:Mobile.",
			},

			"cron": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Timer task cron expression.",
			},

			"operate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The input is valid when the parameter is modified, `suspend`/`resume`, used to suspend/resume the dial test task.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Task status 1:TaskPending, 2:TaskRunning,3:TaskRunException,4:TaskSuspending 5:TaskSuspendException,6:TaskSuspendException,7:TaskSuspended,9:TaskDeleted.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudCatTaskSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cat_task_set.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cat.NewCreateProbeTasksRequest()
		response *cat.CreateProbeTasksResponse
		taskId   string
	)

	if v, ok := d.GetOk("batch_tasks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			probeTaskBasicConfiguration := cat.ProbeTaskBasicConfiguration{}
			if v, ok := dMap["name"]; ok {
				probeTaskBasicConfiguration.Name = helper.String(v.(string))
			}
			if v, ok := dMap["target_address"]; ok {
				probeTaskBasicConfiguration.TargetAddress = helper.String(v.(string))
			}
			request.BatchTasks = append(request.BatchTasks, &probeTaskBasicConfiguration)
		}
	}

	if v, ok := d.GetOk("task_type"); ok {
		request.TaskType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("nodes"); ok {
		nodesSet := v.(*schema.Set).List()
		for i := range nodesSet {
			nodes := nodesSet[i].(string)
			request.Nodes = append(request.Nodes, &nodes)
		}
	}

	if v, ok := d.GetOk("interval"); ok {
		request.Interval = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("parameters"); ok {
		request.Parameters = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_category"); ok {
		request.TaskCategory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cron"); ok {
		request.Cron = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCatClient().CreateProbeTasks(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cat taskSet failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskIDs[0]

	service := CatService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cat:%s:uin/:TaskId/%s", region, taskId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	err = resource.Retry(1*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeCatTaskSet(ctx, taskId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == 2 || *instance.Status == 10 {
			return nil
		}
		if *instance.Status == 3 {
			return resource.NonRetryableError(fmt.Errorf("taskSet status is %v, operate failed.", *instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("taskSet status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	d.SetId(taskId)

	return resourceTencentCloudCatTaskSetRead(d, meta)
}

func resourceTencentCloudCatTaskSetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cat_task_set.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CatService{client: meta.(*TencentCloudClient).apiV3Conn}

	taskId := d.Id()

	taskSet, err := service.DescribeCatTaskSet(ctx, taskId)

	if err != nil {
		return err
	}

	if taskSet == nil {
		d.SetId("")
		return fmt.Errorf("resource `taskSet` %s does not exist", taskId)
	}

	if taskSet != nil {
		batchTasksList := []interface{}{}
		batchTasksMap := map[string]interface{}{}
		if taskSet.Name != nil {
			batchTasksMap["name"] = taskSet.Name
		}
		if taskSet.TargetAddress != nil {
			batchTasksMap["target_address"] = taskSet.TargetAddress
		}

		batchTasksList = append(batchTasksList, batchTasksMap)

		_ = d.Set("batch_tasks", batchTasksList)
	}

	if taskSet.TaskType != nil {
		_ = d.Set("task_type", taskSet.TaskType)
	}

	if taskSet.Nodes != nil {
		_ = d.Set("nodes", taskSet.Nodes)
	}

	if taskSet.TaskId != nil {
		_ = d.Set("task_id", taskSet.TaskId)
	}

	if taskSet.Interval != nil {
		_ = d.Set("interval", taskSet.Interval)
	}

	if taskSet.Parameters != nil {
		_ = d.Set("parameters", taskSet.Parameters)
	}

	if taskSet.TaskCategory != nil {
		_ = d.Set("task_category", taskSet.TaskCategory)
	}

	if taskSet.Cron != nil {
		_ = d.Set("cron", taskSet.Cron)
	}

	if taskSet.Status != nil {
		_ = d.Set("status", taskSet.Status)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cat", "TaskId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCatTaskSetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cat_task_set.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cat.NewUpdateProbeTaskConfigurationListRequest()

	taskId := d.Id()

	request.TaskIds = []*string{helper.String(taskId)}

	if v, ok := d.GetOk("nodes"); ok {
		nodesSet := v.(*schema.Set).List()
		for i := range nodesSet {
			nodes := nodesSet[i].(string)
			request.Nodes = append(request.Nodes, &nodes)
		}
	}

	if v, ok := d.GetOk("interval"); ok {
		request.Interval = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("parameters"); ok {
		request.Parameters = helper.String(v.(string))
	}

	if d.HasChange("cron") {
		if v, ok := d.GetOk("cron"); ok {
			request.Cron = helper.String(v.(string))
		}
	}

	if d.HasChange("task_type") {
		return fmt.Errorf("`task_type` do not support change now.")
	}

	if d.HasChange("task_category") {
		return fmt.Errorf("`task_category` do not support change now.")
	}

	if d.HasChange("batch_tasks") {
		oldInterface, newInterface := d.GetChange("batch_tasks")
		oldMap := make(map[string]interface{})
		newMap := make(map[string]interface{})
		for _, item := range oldInterface.([]interface{}) {
			oldMap = item.(map[string]interface{})
		}
		for _, item := range newInterface.([]interface{}) {
			newMap = item.(map[string]interface{})
		}
		replace, _ := diffTags(oldMap, newMap)

		if _, ok := replace["target_address"]; ok {
			return fmt.Errorf("`target_address` do not support change now.")
		}

		if v, ok := replace["name"]; ok {
			requestTaskAttributes := cat.NewUpdateProbeTaskAttributesRequest()
			requestTaskAttributes.TaskId = &taskId
			requestTaskAttributes.Name = &v
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseCatClient().UpdateProbeTaskAttributes(requestTaskAttributes)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s Suspend cat task failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCatClient().UpdateProbeTaskConfigurationList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cat taskSet failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("operate") {
		if v, ok := d.GetOk("operate"); ok {
			operate := v.(string)
			if operate == "suspend" {
				requestSuspend := cat.NewSuspendProbeTaskRequest()
				requestSuspend.TaskIds = append(requestSuspend.TaskIds, &taskId)
				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseCatClient().SuspendProbeTask(requestSuspend)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s Suspend cat task failed, reason:%+v", logId, err)
					return err
				}
			} else if operate == "resume" {
				requestResume := cat.NewResumeProbeTaskRequest()
				requestResume.TaskIds = append(requestResume.TaskIds, &taskId)
				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseCatClient().ResumeProbeTask(requestResume)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s Resume cat task failed, reason:%+v", logId, err)
					return err
				}
			} else {
				return fmt.Errorf("`operate` only allows the input of suspend/resume.")
			}
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("cat", "TaskId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudCatTaskSetRead(d, meta)
}

func resourceTencentCloudCatTaskSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cat_task_set.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CatService{client: meta.(*TencentCloudClient).apiV3Conn}

	taskId := d.Id()

	if err := service.DeleteCatTaskSetById(ctx, taskId); err != nil {
		return err
	}

	err := resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeCatTaskSet(ctx, taskId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance == nil {
			return nil
		}
		if *instance.Status == 9 {
			return nil
		}
		if *instance.Status == 8 {
			return resource.NonRetryableError(fmt.Errorf("taskSet status is %v, operate failed.", *instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("taskSet status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}
	return nil
}
