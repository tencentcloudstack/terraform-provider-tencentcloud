package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsCompareTask() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDtsCompareTaskRead,
		Create: resourceTencentCloudDtsCompareTaskCreate,
		Update: resourceTencentCloudDtsCompareTaskUpdate,
		Delete: resourceTencentCloudDtsCompareTaskDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "job id.",
			},

			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "task name.",
			},

			"object_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "object mode.",
			},

			"objects": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "object mode.",
						},
						"object_items": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "object items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "database name.",
									},
									"db_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "database mode.",
									},
									"schema_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "schema name.",
									},
									"table_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "table mode.",
									},
									"tables": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "table list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "table name.",
												},
											},
										},
									},
									"view_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "view mode.",
									},
									"views": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "view list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"view_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "view name.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"compare_task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "compare task id.",
			},
		},
	}
}

func resourceTencentCloudDtsCompareTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = dts.NewCreateCompareTaskRequest()
		response      *dts.CreateCompareTaskResponse
		startRequest  = dts.NewStartCompareRequest()
		service       = DtsService{client: meta.(*TencentCloudClient).apiV3Conn}
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		jobId         string
		compareTaskId string
	)

	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
		startRequest.JobId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("object_mode"); ok {
		request.ObjectMode = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "objects"); ok {
		compareObject := dts.CompareObject{}

		if v, ok := dMap["object_mode"]; ok {
			compareObject.ObjectMode = helper.String(v.(string))
		}

		if v, ok := dMap["object_items"]; ok {
			for _, item := range v.([]interface{}) {
				ObjectItemsMap := item.(map[string]interface{})
				compareObjectItem := dts.CompareObjectItem{}
				if v, ok := ObjectItemsMap["db_name"]; ok {
					compareObjectItem.DbName = helper.String(v.(string))
				}
				if v, ok := ObjectItemsMap["db_mode"]; ok {
					compareObjectItem.DbMode = helper.String(v.(string))
				}
				if v, ok := ObjectItemsMap["schema_name"]; ok {
					compareObjectItem.SchemaName = helper.String(v.(string))
				}
				if v, ok := ObjectItemsMap["table_mode"]; ok {
					compareObjectItem.TableMode = helper.String(v.(string))
				}
				if v, ok := ObjectItemsMap["tables"]; ok {
					for _, item := range v.([]interface{}) {
						TablesMap := item.(map[string]interface{})
						compareTableItem := dts.CompareTableItem{}
						if v, ok := TablesMap["table_name"]; ok {
							compareTableItem.TableName = helper.String(v.(string))
						}
						compareObjectItem.Tables = append(compareObjectItem.Tables, &compareTableItem)
					}
				}
				if v, ok := ObjectItemsMap["view_mode"]; ok {
					compareObjectItem.ViewMode = helper.String(v.(string))
				}
				if v, ok := ObjectItemsMap["views"]; ok {
					for _, item := range v.([]interface{}) {
						ViewsMap := item.(map[string]interface{})
						compareViewItem := dts.CompareViewItem{}
						if v, ok := ViewsMap["view_name"]; ok {
							compareViewItem.ViewName = helper.String(v.(string))
						}
						compareObjectItem.Views = append(compareObjectItem.Views, &compareViewItem)
					}
				}
				compareObject.ObjectItems = append(compareObject.ObjectItems, &compareObjectItem)
			}
		}

		request.Objects = &compareObject
	}

	// create compareTask
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateCompareTask(request)
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
		log.Printf("[CRITAL]%s create dts compareTask failed, reason:%+v", logId, err)
		return err
	}

	// wait created
	if err = service.PollingCompareTaskStatusUntil(ctx, jobId, compareTaskId, "created"); err != nil {
		return err
	}

	// start compareTask
	compareTaskId = *response.Response.CompareTaskId
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		startRequest.CompareTaskId = helper.String(compareTaskId)
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().StartCompare(startRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, startRequest.GetAction(), startRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s start dts compareTask failed, reason:%+v", logId, err)
		return err
	}

	// wait running
	if err = service.PollingCompareTaskStatusUntil(ctx, jobId, compareTaskId, "running"); err != nil {
		return err
	}

	d.SetId(jobId + FILED_SP + compareTaskId)
	return resourceTencentCloudDtsCompareTaskRead(d, meta)
}

func resourceTencentCloudDtsCompareTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	compareTaskId := idSplit[1]

	ret, err := service.DescribeDtsCompareTask(ctx, helper.String(jobId), helper.String(compareTaskId))

	if err != nil {
		return err
	}

	if ret == nil {
		d.SetId("")
		return fmt.Errorf("resource `compareTask` %s does not exist", compareTaskId)
	}

	if len(ret) > 0 {
		compareTask := ret[0]

		if compareTask.JobId != nil {
			_ = d.Set("job_id", compareTask.JobId)
		}

		if compareTask.TaskName != nil {
			_ = d.Set("task_name", compareTask.TaskName)
		}

		if compareTask.Config != nil {
			objects := compareTask.Config
			// SDK do not support this field ObjectMode
			// if objects.ObjectMode != nil {
			// 	_ = d.Set("object_mode", objects.ObjectMode)
			// }

			//objects
			objectsMap := map[string]interface{}{}
			if objects.ObjectMode != nil {
				objectsMap["object_mode"] = objects.ObjectMode
			}

			if objects.ObjectItems != nil {
				objectItemsList := []interface{}{}
				// object_items
				for _, objectItems := range objects.ObjectItems {
					objectItemsMap := map[string]interface{}{}
					if objectItems.DbName != nil {
						objectItemsMap["db_name"] = objectItems.DbName
					}
					if objectItems.DbMode != nil {
						objectItemsMap["db_mode"] = objectItems.DbMode
					}
					if objectItems.SchemaName != nil {
						objectItemsMap["schema_name"] = objectItems.SchemaName
					}
					if objectItems.TableMode != nil {
						objectItemsMap["table_mode"] = objectItems.TableMode
					}
					if objectItems.Tables != nil {
						tablesList := []interface{}{}
						for _, tables := range objectItems.Tables {
							tablesMap := map[string]interface{}{}
							if tables.TableName != nil {
								tablesMap["table_name"] = tables.TableName
							}

							tablesList = append(tablesList, tablesMap)
						}
						objectItemsMap["tables"] = tablesList
					}
					if objectItems.ViewMode != nil {
						objectItemsMap["view_mode"] = objectItems.ViewMode
					}
					if objectItems.Views != nil {
						viewsList := []interface{}{}
						for _, views := range objectItems.Views {
							viewsMap := map[string]interface{}{}
							if views.ViewName != nil {
								viewsMap["view_name"] = views.ViewName
							}

							viewsList = append(viewsList, viewsMap)
						}
						objectItemsMap["views"] = viewsList
					}
					objectItemsList = append(objectItemsList, objectItemsMap)
				}
				objectsMap["object_items"] = objectItemsList
			}

			_ = d.Set("objects", []interface{}{objectsMap})
		}

		if compareTask.CompareTaskId != nil {
			_ = d.Set("compare_task_id", compareTask.CompareTaskId)
		}

	}
	return nil
}

func resourceTencentCloudDtsCompareTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := dts.NewModifyCompareTaskRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	compareTaskId := idSplit[1]

	request.JobId = &jobId
	request.CompareTaskId = &compareTaskId

	if d.HasChange("job_id") {
		return fmt.Errorf("`job_id` do not support change now.")
	}

	if d.HasChange("task_name") {
		if v, ok := d.GetOk("task_name"); ok {
			request.TaskName = helper.String(v.(string))
		}
	}

	if d.HasChange("object_mode") {
		if v, ok := d.GetOk("object_mode"); ok {
			request.ObjectMode = helper.String(v.(string))
		}
	}

	if d.HasChange("objects") {
		if dMap, ok := helper.InterfacesHeadMap(d, "objects"); ok {
			compareObject := dts.CompareObject{}
			if v, ok := dMap["object_mode"]; ok {
				compareObject.ObjectMode = helper.String(v.(string))
			}
			if v, ok := dMap["object_items"]; ok {
				for _, item := range v.([]interface{}) {
					ObjectItemsMap := item.(map[string]interface{})
					compareObjectItem := dts.CompareObjectItem{}
					if v, ok := ObjectItemsMap["db_name"]; ok {
						compareObjectItem.DbName = helper.String(v.(string))
					}
					if v, ok := ObjectItemsMap["db_mode"]; ok {
						compareObjectItem.DbMode = helper.String(v.(string))
					}
					if v, ok := ObjectItemsMap["schema_name"]; ok {
						compareObjectItem.SchemaName = helper.String(v.(string))
					}
					if v, ok := ObjectItemsMap["table_mode"]; ok {
						compareObjectItem.TableMode = helper.String(v.(string))
					}
					if v, ok := ObjectItemsMap["tables"]; ok {
						for _, item := range v.([]interface{}) {
							TablesMap := item.(map[string]interface{})
							compareTableItem := dts.CompareTableItem{}
							if v, ok := TablesMap["table_name"]; ok {
								compareTableItem.TableName = helper.String(v.(string))
							}
							compareObjectItem.Tables = append(compareObjectItem.Tables, &compareTableItem)
						}
					}
					if v, ok := ObjectItemsMap["view_mode"]; ok {
						compareObjectItem.ViewMode = helper.String(v.(string))
					}
					if v, ok := ObjectItemsMap["views"]; ok {
						for _, item := range v.([]interface{}) {
							ViewsMap := item.(map[string]interface{})
							compareViewItem := dts.CompareViewItem{}
							if v, ok := ViewsMap["view_name"]; ok {
								compareViewItem.ViewName = helper.String(v.(string))
							}
							compareObjectItem.Views = append(compareObjectItem.Views, &compareViewItem)
						}
					}
					compareObject.ObjectItems = append(compareObject.ObjectItems, &compareObjectItem)
				}
			}

			request.Objects = &compareObject
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ModifyCompareTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dts compareTask failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDtsCompareTaskRead(d, meta)
}

func resourceTencentCloudDtsCompareTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	compareTaskId := idSplit[1]

	if err := service.DeleteDtsCompareTaskById(ctx, jobId, compareTaskId); err != nil {
		return err
	}

	return nil
}
