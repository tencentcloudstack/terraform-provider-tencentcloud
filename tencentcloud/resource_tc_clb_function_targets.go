/*
Provides a resource to create a clb function_targets

Example Usage

```hcl
resource "tencentcloud_clb_function_targets" "function_targets" {
  load_balancer_id = &lt;nil&gt;
  listener_id = &lt;nil&gt;
  function_targets {
		function {
			function_namespace = &lt;nil&gt;
			function_name = &lt;nil&gt;
			function_qualifier = &lt;nil&gt;
			function_qualifier_type = &lt;nil&gt;
		}
		weight = &lt;nil&gt;

  }
  location_id = &lt;nil&gt;
  domain = &lt;nil&gt;
  url = &lt;nil&gt;
}
```

Import

clb function_targets can be imported using the id, e.g.

```
terraform import tencentcloud_clb_function_targets.function_targets function_targets_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudClbFunctionTargets() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbFunctionTargetsCreate,
		Read:   resourceTencentCloudClbFunctionTargetsRead,
		Delete: resourceTencentCloudClbFunctionTargetsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Load Balancer Instance ID.",
			},

			"listener_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Load Balancer Listener ID.",
			},

			"function_targets": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "List of cloud functions to be bound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Information about cloud functions.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"function_namespace": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The namespace of function.",
									},
									"function_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of function.",
									},
									"function_qualifier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The version name or alias of the function.",
									},
									"function_qualifier_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Identifies the type of FunctionQualifier parameter, possible values: VERSION, ALIAS.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Weight.",
						},
					},
				},
			},

			"location_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the target forwarding rule. When binding the cloud function to a layer-7 forwarding rule, this parameter or the Domain+Url parameter must be entered.",
			},

			"domain": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The domain name of the target forwarding rule. If the LocationId parameter has been entered, this parameter will not take effect.",
			},

			"url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The URL of the target forwarding rule. If the LocationId parameter has been entered, this parameter will not take effect.",
			},
		},
	}
}

func resourceTencentCloudClbFunctionTargetsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_function_targets.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = clb.NewRegisterFunctionTargetsRequest()
		response       = clb.NewRegisterFunctionTargetsResponse()
		loadBalancerId string
		listenerId     string
	)
	if v, ok := d.GetOk("load_balancer_id"); ok {
		loadBalancerId = v.(string)
		request.LoadBalancerId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("listener_id"); ok {
		listenerId = v.(string)
		request.ListenerId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_targets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			functionTarget := clb.FunctionTarget{}
			if functionMap, ok := helper.InterfaceToMap(dMap, "function"); ok {
				functionInfo := clb.FunctionInfo{}
				if v, ok := functionMap["function_namespace"]; ok {
					functionInfo.FunctionNamespace = helper.String(v.(string))
				}
				if v, ok := functionMap["function_name"]; ok {
					functionInfo.FunctionName = helper.String(v.(string))
				}
				if v, ok := functionMap["function_qualifier"]; ok {
					functionInfo.FunctionQualifier = helper.String(v.(string))
				}
				if v, ok := functionMap["function_qualifier_type"]; ok {
					functionInfo.FunctionQualifierType = helper.String(v.(string))
				}
				functionTarget.Function = &functionInfo
			}
			if v, ok := dMap["weight"]; ok {
				functionTarget.Weight = helper.IntUint64(v.(int))
			}
			request.FunctionTargets = append(request.FunctionTargets, &functionTarget)
		}
	}

	if v, ok := d.GetOk("location_id"); ok {
		request.LocationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().RegisterFunctionTargets(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb functionTargets failed, reason:%+v", logId, err)
		return err
	}

	loadBalancerId = *response.Response.LoadBalancerId
	d.SetId(strings.Join([]string{loadBalancerId, listenerId}, FILED_SP))

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"0"}, 60*readRetryTimeout, time.Second, service.ClbFunctionTargetsStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudClbFunctionTargetsRead(d, meta)
}

func resourceTencentCloudClbFunctionTargetsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_function_targets.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	loadBalancerId := idSplit[0]
	listenerId := idSplit[1]
	function := idSplit[2]

	functionTargets, err := service.DescribeClbFunctionTargetsById(ctx, loadBalancerId, listenerId, function)
	if err != nil {
		return err
	}

	if functionTargets == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbFunctionTargets` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if functionTargets.LoadBalancerId != nil {
		_ = d.Set("load_balancer_id", functionTargets.LoadBalancerId)
	}

	if functionTargets.ListenerId != nil {
		_ = d.Set("listener_id", functionTargets.ListenerId)
	}

	if functionTargets.FunctionTargets != nil {
		functionTargetsList := []interface{}{}
		for _, functionTargets := range functionTargets.FunctionTargets {
			functionTargetsMap := map[string]interface{}{}

			if functionTargets.FunctionTargets.Function != nil {
				functionMap := map[string]interface{}{}

				if functionTargets.FunctionTargets.Function.FunctionNamespace != nil {
					functionMap["function_namespace"] = functionTargets.FunctionTargets.Function.FunctionNamespace
				}

				if functionTargets.FunctionTargets.Function.FunctionName != nil {
					functionMap["function_name"] = functionTargets.FunctionTargets.Function.FunctionName
				}

				if functionTargets.FunctionTargets.Function.FunctionQualifier != nil {
					functionMap["function_qualifier"] = functionTargets.FunctionTargets.Function.FunctionQualifier
				}

				if functionTargets.FunctionTargets.Function.FunctionQualifierType != nil {
					functionMap["function_qualifier_type"] = functionTargets.FunctionTargets.Function.FunctionQualifierType
				}

				functionTargetsMap["function"] = []interface{}{functionMap}
			}

			if functionTargets.FunctionTargets.Weight != nil {
				functionTargetsMap["weight"] = functionTargets.FunctionTargets.Weight
			}

			functionTargetsList = append(functionTargetsList, functionTargetsMap)
		}

		_ = d.Set("function_targets", functionTargetsList)

	}

	if functionTargets.LocationId != nil {
		_ = d.Set("location_id", functionTargets.LocationId)
	}

	if functionTargets.Domain != nil {
		_ = d.Set("domain", functionTargets.Domain)
	}

	if functionTargets.Url != nil {
		_ = d.Set("url", functionTargets.Url)
	}

	return nil
}

func resourceTencentCloudClbFunctionTargetsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_function_targets.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	loadBalancerId := idSplit[0]
	listenerId := idSplit[1]
	function := idSplit[2]

	if err := service.DeleteClbFunctionTargetsById(ctx, loadBalancerId, listenerId, function); err != nil {
		return err
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"0"}, 60*readRetryTimeout, time.Second, service.ClbFunctionTargetsStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
