/*
Use this data source to query SCF functions.

Example Usage

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "scf-code-1234567890"
  cos_object_name   = "code.zip"
  cos_bucket_region = "ap-guangzhou"
}

data "tencentcloud_scf_functions" "foo" {
  name = "${tencentcloud_scf_function.foo.name}"
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfFunctions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SCF function to be queried.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace of the SCF function to be queried.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the SCF function to be queried.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the SCF function to be queried, can use up to 10 tags.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"functions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of functions. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SCF function.",
						},
						"handler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Handler of the SCF function.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the SCF function.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size of the SCF function runtime, unit is M.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timeout of the SCF function maximum execution time, unit is second.",
						},
						"environment": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Environment variable of the SCF function.",
						},
						"runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Runtime of the SCF function.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID of the SCF function.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID of the SCF function.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace of the SCF function.",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CAM role of the SCF function.",
						},
						"cls_logset_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS logset ID of the SCF function.",
						},
						"cls_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS topic ID of the SCF function.",
						},
						"l5_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable L5.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the SCF function.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the SCF function.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modify time of the SCF function.",
						},
						"code_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Code size of the SCF function.",
						},
						"code_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Code result of the SCF function.",
						},
						"code_error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Code error of the SCF function.",
						},
						"err_no": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Errno of the SCF function.",
						},
						"install_dependency": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to automatically install dependencies.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the SCF function.",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status description of the SCF function.",
						},
						"eip_fixed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether EIP is a fixed IP.",
						},
						"eips": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "EIP list of the SCF function.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host of the SCF function.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vip of the SCF function.",
						},
						"trigger_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Trigger details list the SCF function. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the SCF function trigger.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the SCF function trigger.",
									},
									"trigger_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TriggerDesc of the SCF function trigger.",
									},
									"enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable SCF function trigger.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time of the SCF function trigger.",
									},
									"modify_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Modify time of the SCF function trigger.",
									},
									"custom_argument": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "user-defined parameter of the SCF function trigger.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudScfFunctionsRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_functions.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		name      *string
		namespace *string
		desc      *string
	)

	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("namespace"); ok {
		namespace = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("description"); ok {
		desc = helper.String(raw.(string))
	}

	tags := helper.GetTags(d, "tags")
	if len(tags) > 10 {
		return errors.Errorf("can't set more than 10 tags")
	}

	respFunctions, err := service.DescribeFunctions(ctx, name, namespace, desc, tags)
	if err != nil {
		log.Printf("[CRITAL]%s get function list failed: %+v", logId, err)
		return err
	}

	functions := make([]map[string]interface{}, 0, len(respFunctions))
	ids := make([]string, 0, len(respFunctions))

	for _, fn := range respFunctions {
		ids = append(ids, fmt.Sprintf("%s+%s", *fn.Namespace, *fn.FunctionName))

		m := map[string]interface{}{
			"name":        fn.FunctionName,
			"description": fn.Description,
			"runtime":     fn.Runtime,
			"namespace":   fn.Namespace,
			"create_time": fn.AddTime,
			"modify_time": fn.ModTime,
			"status":      fn.Status,
			"status_desc": fn.StatusDesc,
		}

		rawResp, err := service.DescribeFunction(ctx, *fn.FunctionName, *fn.Namespace)
		if err != nil {
			log.Printf("[CRITAL]%s read function detail failed: %+v", logId, err)
			return err
		}
		resp := rawResp.Response

		m["handler"] = resp.Handler
		m["mem_size"] = resp.MemorySize
		m["timeout"] = resp.Timeout

		env := make(map[string]string, len(resp.Environment.Variables))
		for _, v := range resp.Environment.Variables {
			env[*v.Key] = *v.Value
		}
		m["environment"] = env

		m["vpc_id"] = resp.VpcConfig.VpcId
		m["subnet_id"] = resp.VpcConfig.SubnetId
		m["role"] = resp.Role
		m["cls_logset_id"] = resp.ClsLogsetId
		m["cls_topic_id"] = resp.ClsTopicId
		m["code_size"] = resp.CodeSize
		m["code_result"] = resp.CodeResult
		m["code_error"] = resp.CodeError
		m["err_no"] = resp.ErrNo
		m["install_dependency"] = *resp.InstallDependency == "TRUE"
		m["eip_fixed"] = *resp.EipConfig.EipFixed == "TRUE"
		m["eips"] = resp.EipConfig.Eips
		m["host"] = resp.AccessInfo.Host
		m["vip"] = resp.AccessInfo.Vip
		m["l5_enable"] = *resp.L5Enable == "TRUE"

		triggers := make([]map[string]interface{}, 0, len(resp.Triggers))
		for _, trigger := range resp.Triggers {
			switch *trigger.Type {
			case SCF_TRIGGER_TYPE_TIMER:
				data := struct {
					Cron string `json:"cron"`
				}{}
				if err := json.Unmarshal([]byte(*trigger.TriggerDesc), &data); err != nil {
					log.Printf("[WARN]%s unmarshal timer trigger trigger_desc failed: %+v", logId, errors.WithStack(err))
					continue
				}
				*trigger.TriggerDesc = data.Cron

			case SCF_TRIGGER_TYPE_COS:
				*trigger.TriggerName = strings.Replace(*trigger.TriggerName, SCF_TRIGGER_COS_NAME_SUFFIX, "", -1)
			}

			triggers = append(triggers, map[string]interface{}{
				"name":            trigger.TriggerName,
				"type":            trigger.Type,
				"trigger_desc":    trigger.TriggerDesc,
				"enable":          *trigger.Enable == 1,
				"create_time":     trigger.AddTime,
				"modify_time":     trigger.ModTime,
				"custom_argument": trigger.CustomArgument,
			})
		}
		m["trigger_info"] = triggers

		fnTags := make(map[string]string, len(resp.Tags))
		for _, tag := range resp.Tags {
			fnTags[*tag.Key] = *tag.Value
		}
		m["tags"] = fnTags

		functions = append(functions, m)
	}

	_ = d.Set("functions", functions)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), functions); err != nil {
			err = errors.WithStack(err)
			log.Printf("[CRITAL]%s output file[%s] fail, reason: %+v", logId, output.(string), err)
			return err
		}
	}

	return nil
}
