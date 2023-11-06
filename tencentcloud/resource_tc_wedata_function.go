/*
Provides a resource to create a wedata function

Example Usage

```hcl
resource "tencentcloud_wedata_function" "example" {
  type               = "HIVE"
  kind               = "ANALYSIS"
  name               = "tf_example"
  cluster_identifier = "emr-m6u3qgk0"
  db_name            = "tf_db_example"
  project_id         = "1612982498218618880"
  class_name         = "tf_class_example"
  resource_list {
    path = "/wedata-demo-1314991481/untitled3-1.0-SNAPSHOT.jar"
    name = "untitled3-1.0-SNAPSHOT.jar"
    id   = "5b28bcdf-a0e6-4022-927d-927d399c4593"
    type = "cos"
  }
  description = "description."
  usage       = "usage info."
  param_desc  = "param info."
  return_desc = "return value info."
  example     = "example info."
  comment     = "V1"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWedataFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataFunctionCreate,
		Read:   resourceTencentCloudWedataFunctionRead,
		Update: resourceTencentCloudWedataFunctionUpdate,
		Delete: resourceTencentCloudWedataFunctionDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function Type, Enum: HIVE, SPARK, DLC.",
			},
			"kind": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function Kind, Enum: ANALYSIS, ENCRYPTION, AGGREGATE, LOGIC, DATE_AND_TIME, MATH, CONVERSION, STRING, IP_AND_DOMAIN, WINDOW, OTHER.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function Name.",
			},
			"cluster_identifier": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},
			"class_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Class name of function entry.",
			},
			"resource_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Resource of the function, stored in WeData COS(.jar,...).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource Path.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource Name.",
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource ID.",
						},
						"md5": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource MD5 Value.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource Type.",
						},
					},
				},
			},
			"description": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Description of the function.",
			},
			"usage": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Usage of the function.",
			},
			"param_desc": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Description of the Parameter.",
			},
			"return_desc": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Description of the Return value.",
			},
			"example": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Example of the function.",
			},
			"comment": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Comment.",
			},
			"function_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Function ID.",
			},
		},
	}
}

func resourceTencentCloudWedataFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_function.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                        = getLogId(contextNil)
		createCustomFunctionRequest  = wedata.NewCreateCustomFunctionRequest()
		createCustomFunctionResponse = wedata.NewCreateCustomFunctionResponse()
		saveCustomFunctionRequest    = wedata.NewSaveCustomFunctionRequest()
		submitCustomFunctionRequest  = wedata.NewSubmitCustomFunctionRequest()
		functionId                   string
		funcType                     string
		funcName                     string
		projectId                    string
		clusterIdentifier            string
	)

	if v, ok := d.GetOk("type"); ok {
		createCustomFunctionRequest.Type = helper.String(v.(string))
		funcType = v.(string)
	}

	if v, ok := d.GetOk("kind"); ok {
		createCustomFunctionRequest.Kind = helper.String(v.(string))
		saveCustomFunctionRequest.Kind = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		createCustomFunctionRequest.Name = helper.String(v.(string))
		funcName = v.(string)
	}

	if v, ok := d.GetOk("cluster_identifier"); ok {
		createCustomFunctionRequest.ClusterIdentifier = helper.String(v.(string))
		saveCustomFunctionRequest.ClusterIdentifier = helper.String(v.(string))
		submitCustomFunctionRequest.ClusterIdentifier = helper.String(v.(string))
		clusterIdentifier = v.(string)
	}

	if v, ok := d.GetOk("db_name"); ok {
		createCustomFunctionRequest.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		createCustomFunctionRequest.ProjectId = helper.String(v.(string))
		submitCustomFunctionRequest.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateCustomFunction(createCustomFunctionRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, createCustomFunctionRequest.GetAction(), createCustomFunctionRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.FunctionId == nil {
			e = fmt.Errorf("wedata function not exists")
			if result.Response.ErrorMessage != nil {
				e = fmt.Errorf(*result.Response.ErrorMessage)
			}

			return resource.NonRetryableError(e)
		}

		createCustomFunctionResponse = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata function failed, reason:%+v", logId, err)
		return err
	}

	functionId = *createCustomFunctionResponse.Response.FunctionId
	d.SetId(strings.Join([]string{functionId, funcType, funcName, projectId, clusterIdentifier}, FILED_SP))

	saveCustomFunctionRequest.FunctionId = &functionId
	submitCustomFunctionRequest.FunctionId = &functionId

	if v, ok := d.GetOk("class_name"); ok {
		saveCustomFunctionRequest.ClassName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			functionResource := wedata.FunctionResource{}
			if v, ok := dMap["path"]; ok {
				functionResource.Path = helper.String(v.(string))
			}

			if v, ok := dMap["name"]; ok {
				functionResource.Name = helper.String(v.(string))
			}

			if v, ok := dMap["id"]; ok {
				functionResource.Id = helper.String(v.(string))
			}

			if v, ok := dMap["md5"]; ok {
				functionResource.Md5 = helper.String(v.(string))
			}

			if v, ok := dMap["type"]; ok {
				functionResource.Type = helper.String(v.(string))
			}

			saveCustomFunctionRequest.ResourceList = append(saveCustomFunctionRequest.ResourceList, &functionResource)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		saveCustomFunctionRequest.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("usage"); ok {
		saveCustomFunctionRequest.Usage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_desc"); ok {
		saveCustomFunctionRequest.ParamDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("return_desc"); ok {
		saveCustomFunctionRequest.ReturnDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("example"); ok {
		saveCustomFunctionRequest.Example = helper.String(v.(string))
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().SaveCustomFunction(saveCustomFunctionRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, saveCustomFunctionRequest.GetAction(), saveCustomFunctionRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata function not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s save wedata function failed, reason:%+v", logId, err)
		return err
	}

	if v, ok := d.GetOk("comment"); ok {
		submitCustomFunctionRequest.Comment = helper.String(v.(string))
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().SubmitCustomFunction(submitCustomFunctionRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, submitCustomFunctionRequest.GetAction(), submitCustomFunctionRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata function not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s submit wedata function failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataFunctionRead(d, meta)
}

func resourceTencentCloudWedataFunctionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_function.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionId := idSplit[0]
	funcType := idSplit[1]
	funcName := idSplit[2]
	projectId := idSplit[3]

	function, err := service.DescribeWedataFunctionById(ctx, functionId, funcType, funcName, projectId)
	if err != nil {
		return err
	}

	if function == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataFunction` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("function_id", functionId)
	_ = d.Set("project_id", projectId)

	if function.Type != nil {
		_ = d.Set("type", function.Type)
	}

	if function.Kind != nil {
		_ = d.Set("kind", function.Kind)
	}

	if function.Name != nil {
		_ = d.Set("name", function.Name)
	}

	if function.ClusterIdentifier != nil {
		_ = d.Set("cluster_identifier", function.ClusterIdentifier)
	}

	if function.DbName != nil {
		_ = d.Set("db_name", function.DbName)
	}

	if function.ClassName != nil {
		_ = d.Set("class_name", function.ClassName)
	}

	if function.Description != nil {
		_ = d.Set("description", function.Description)
	}

	if function.Usage != nil {
		_ = d.Set("usage", function.Usage)
	}

	if function.ParamDesc != nil {
		_ = d.Set("param_desc", function.ParamDesc)
	}

	if function.ReturnDesc != nil {
		_ = d.Set("return_desc", function.ReturnDesc)
	}

	if function.Example != nil {
		_ = d.Set("example", function.Example)
	}

	return nil
}

func resourceTencentCloudWedataFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_function.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                       = getLogId(contextNil)
		saveCustomFunctionRequest   = wedata.NewSaveCustomFunctionRequest()
		submitCustomFunctionRequest = wedata.NewSubmitCustomFunctionRequest()
	)

	immutableArgs := []string{"type", "name", "cluster_identifier", "db_name", "project_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionId := idSplit[0]
	projectId := idSplit[3]
	clusterIdentifier := idSplit[3]

	saveCustomFunctionRequest.FunctionId = &functionId
	saveCustomFunctionRequest.ClusterIdentifier = &clusterIdentifier

	if v, ok := d.GetOk("kind"); ok {
		saveCustomFunctionRequest.Kind = helper.String(v.(string))
	}

	if v, ok := d.GetOk("class_name"); ok {
		saveCustomFunctionRequest.ClassName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			functionResource := wedata.FunctionResource{}
			if v, ok := dMap["path"]; ok {
				functionResource.Path = helper.String(v.(string))
			}

			if v, ok := dMap["name"]; ok {
				functionResource.Name = helper.String(v.(string))
			}

			if v, ok := dMap["id"]; ok {
				functionResource.Id = helper.String(v.(string))
			}

			if v, ok := dMap["md5"]; ok {
				functionResource.Md5 = helper.String(v.(string))
			}

			if v, ok := dMap["type"]; ok {
				functionResource.Type = helper.String(v.(string))
			}

			saveCustomFunctionRequest.ResourceList = append(saveCustomFunctionRequest.ResourceList, &functionResource)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		saveCustomFunctionRequest.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("usage"); ok {
		saveCustomFunctionRequest.Usage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_desc"); ok {
		saveCustomFunctionRequest.ParamDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("return_desc"); ok {
		saveCustomFunctionRequest.ReturnDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("example"); ok {
		saveCustomFunctionRequest.Example = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().SaveCustomFunction(saveCustomFunctionRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, saveCustomFunctionRequest.GetAction(), saveCustomFunctionRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata function not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update save wedata function failed, reason:%+v", logId, err)
		return err
	}

	submitCustomFunctionRequest.FunctionId = &functionId
	submitCustomFunctionRequest.ClusterIdentifier = &clusterIdentifier
	submitCustomFunctionRequest.ProjectId = &projectId
	if v, ok := d.GetOk("comment"); ok {
		submitCustomFunctionRequest.Comment = helper.String(v.(string))
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().SubmitCustomFunction(submitCustomFunctionRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, submitCustomFunctionRequest.GetAction(), submitCustomFunctionRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata function not exists")
			return resource.NonRetryableError(e)
		}

		return nil

	})
	if err != nil {
		log.Printf("[CRITAL]%s update submit wedata function failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataFunctionRead(d, meta)
}

func resourceTencentCloudWedataFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_function.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionId := idSplit[0]
	projectId := idSplit[3]
	clusterIdentifier := idSplit[4]

	if err := service.DeleteWedataFunctionById(ctx, functionId, projectId, clusterIdentifier); err != nil {
		return err
	}

	return nil
}
