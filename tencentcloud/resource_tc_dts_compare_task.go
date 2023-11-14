/*
Provides a resource to create a dts compare_task

Example Usage

```hcl
resource "tencentcloud_dts_compare_task" "compare_task" {
  job_id = &lt;nil&gt;
  task_name = &lt;nil&gt;
  object_mode = &lt;nil&gt;
  objects {
		object_mode = &lt;nil&gt;
		object_items {
			db_name = &lt;nil&gt;
			db_mode = &lt;nil&gt;
			schema_name = &lt;nil&gt;
			table_mode = &lt;nil&gt;
			tables {
				table_name = &lt;nil&gt;
			}
			view_mode = &lt;nil&gt;
			views {
				view_name = &lt;nil&gt;
			}
		}

  }
  }
```

Import

dts compare_task can be imported using the id, e.g.

```
terraform import tencentcloud_dts_compare_task.compare_task compare_task_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDtsCompareTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsCompareTaskCreate,
		Read:   resourceTencentCloudDtsCompareTaskRead,
		Update: resourceTencentCloudDtsCompareTaskUpdate,
		Delete: resourceTencentCloudDtsCompareTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"task_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Task name.",
			},

			"object_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Object mode, optional value is sameAsMigrate(migrate all) or custom.",
			},

			"objects": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object mode, optional value is all or partial.",
						},
						"object_items": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Object items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database name.",
									},
									"db_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database mode.",
									},
									"schema_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Schema name.",
									},
									"table_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table mode.",
									},
									"tables": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Table list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Table name.",
												},
											},
										},
									},
									"view_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "View mode.",
									},
									"views": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "View list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"view_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "View name.",
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
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Compare task id.",
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
		response      = dts.NewCreateCompareTaskResponse()
		jobId         string
		compareTaskId string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
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
				objectItemsMap := item.(map[string]interface{})
				compareObjectItem := dts.CompareObjectItem{}
				if v, ok := objectItemsMap["db_name"]; ok {
					compareObjectItem.DbName = helper.String(v.(string))
				}
				if v, ok := objectItemsMap["db_mode"]; ok {
					compareObjectItem.DbMode = helper.String(v.(string))
				}
				if v, ok := objectItemsMap["schema_name"]; ok {
					compareObjectItem.SchemaName = helper.String(v.(string))
				}
				if v, ok := objectItemsMap["table_mode"]; ok {
					compareObjectItem.TableMode = helper.String(v.(string))
				}
				if v, ok := objectItemsMap["tables"]; ok {
					for _, item := range v.([]interface{}) {
						tablesMap := item.(map[string]interface{})
						compareTableItem := dts.CompareTableItem{}
						if v, ok := tablesMap["table_name"]; ok {
							compareTableItem.TableName = helper.String(v.(string))
						}
						compareObjectItem.Tables = append(compareObjectItem.Tables, &compareTableItem)
					}
				}
				if v, ok := objectItemsMap["view_mode"]; ok {
					compareObjectItem.ViewMode = helper.String(v.(string))
				}
				if v, ok := objectItemsMap["views"]; ok {
					for _, item := range v.([]interface{}) {
						viewsMap := item.(map[string]interface{})
						compareViewItem := dts.CompareViewItem{}
						if v, ok := viewsMap["view_name"]; ok {
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateCompareTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts compareTask failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(strings.Join([]string{jobId, compareTaskId}, FILED_SP))

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

	compareTask, err := service.DescribeDtsCompareTaskById(ctx, jobId, compareTaskId)
	if err != nil {
		return err
	}

	if compareTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsCompareTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if compareTask.JobId != nil {
		_ = d.Set("job_id", compareTask.JobId)
	}

	if compareTask.TaskName != nil {
		_ = d.Set("task_name", compareTask.TaskName)
	}

	if compareTask.ObjectMode != nil {
		_ = d.Set("object_mode", compareTask.ObjectMode)
	}

	if compareTask.Objects != nil {
		objectsMap := map[string]interface{}{}

		if compareTask.Objects.ObjectMode != nil {
			objectsMap["object_mode"] = compareTask.Objects.ObjectMode
		}

		if compareTask.Objects.ObjectItems != nil {
			objectItemsList := []interface{}{}
			for _, objectItems := range compareTask.Objects.ObjectItems {
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

					objectItemsMap["tables"] = []interface{}{tablesList}
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

					objectItemsMap["views"] = []interface{}{viewsList}
				}

				objectItemsList = append(objectItemsList, objectItemsMap)
			}

			objectsMap["object_items"] = []interface{}{objectItemsList}
		}

		_ = d.Set("objects", []interface{}{objectsMap})
	}

	if compareTask.CompareTaskId != nil {
		_ = d.Set("compare_task_id", compareTask.CompareTaskId)
	}

	return nil
}

func resourceTencentCloudDtsCompareTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dts.NewModifyCompareTaskRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	compareTaskId := idSplit[1]

	request.JobId = &jobId
	request.CompareTaskId = &compareTaskId

	immutableArgs := []string{"job_id", "task_name", "object_mode", "objects", "compare_task_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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
					objectItemsMap := item.(map[string]interface{})
					compareObjectItem := dts.CompareObjectItem{}
					if v, ok := objectItemsMap["db_name"]; ok {
						compareObjectItem.DbName = helper.String(v.(string))
					}
					if v, ok := objectItemsMap["db_mode"]; ok {
						compareObjectItem.DbMode = helper.String(v.(string))
					}
					if v, ok := objectItemsMap["schema_name"]; ok {
						compareObjectItem.SchemaName = helper.String(v.(string))
					}
					if v, ok := objectItemsMap["table_mode"]; ok {
						compareObjectItem.TableMode = helper.String(v.(string))
					}
					if v, ok := objectItemsMap["tables"]; ok {
						for _, item := range v.([]interface{}) {
							tablesMap := item.(map[string]interface{})
							compareTableItem := dts.CompareTableItem{}
							if v, ok := tablesMap["table_name"]; ok {
								compareTableItem.TableName = helper.String(v.(string))
							}
							compareObjectItem.Tables = append(compareObjectItem.Tables, &compareTableItem)
						}
					}
					if v, ok := objectItemsMap["view_mode"]; ok {
						compareObjectItem.ViewMode = helper.String(v.(string))
					}
					if v, ok := objectItemsMap["views"]; ok {
						for _, item := range v.([]interface{}) {
							viewsMap := item.(map[string]interface{})
							compareViewItem := dts.CompareViewItem{}
							if v, ok := viewsMap["view_name"]; ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts compareTask failed, reason:%+v", logId, err)
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
