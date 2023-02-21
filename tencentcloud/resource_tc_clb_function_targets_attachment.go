/*
Provides a resource to create a clb function_targets_attachment

Example Usage

```hcl
resource "tencentcloud_clb_function_targets_attachment" "function_targets" {
  domain           = "xxx.com"
  listener_id      = "lbl-nonkgvc2"
  load_balancer_id = "lb-5dnrkgry"
  url              = "/"

  function_targets {
    weight = 10

    function {
      function_name           = "keep-tf-test-1675954233"
      function_namespace      = "default"
      function_qualifier      = "$LATEST"
      function_qualifier_type = "VERSION"
    }
  }
}
```

Import

clb function_targets_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_clb_function_targets_attachment.function_targets loadBalancerId#listenerId#locationId or loadBalancerId#listenerId#domain#rule
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClbFunctionTargetsAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbFunctionTargetsAttachmentCreate,
		Read:   resourceTencentCloudClbFunctionTargetsAttachmentRead,
		Update: resourceTencentCloudClbFunctionTargetsAttachmentUpdate,
		Delete: resourceTencentCloudClbFunctionTargetsAttachmentDelete,
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
				MaxItems:    1,
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
							Default:     10,
							Optional:    true,
							Description: "Weight. The default is `10`.",
						},
					},
				},
			},

			"location_id": {
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"domain", "url"},
				Description:   "The ID of the target forwarding rule. When binding the cloud function to a layer-7 forwarding rule, this parameter or the Domain+Url parameter must be entered.",
			},

			"domain": {
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Type:         schema.TypeString,
				RequiredWith: []string{"url"},
				Description:  "The domain name of the target forwarding rule. If the LocationId parameter has been entered, this parameter will not take effect.",
			},

			"url": {
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Type:         schema.TypeString,
				RequiredWith: []string{"domain"},
				Description:  "The URL of the target forwarding rule. If the LocationId parameter has been entered, this parameter will not take effect.",
			},
		},
	}
}

func resourceTencentCloudClbFunctionTargetsAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_function_targets_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = clb.NewRegisterFunctionTargetsRequest()
		loadBalancerId string
		listenerId     string
		locationId     string
		domain         string
		url            string
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
		locationId = v.(string)
		request.LocationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
		url = v.(string)
		request.Url = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().RegisterFunctionTargets(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb functionTargetsAttachment failed, reason:%+v", logId, err)
		return err
	}

	if locationId != "" {
		d.SetId(loadBalancerId + FILED_SP + listenerId + FILED_SP + locationId)
	} else {
		d.SetId(loadBalancerId + FILED_SP + listenerId + FILED_SP + domain + FILED_SP + url)
	}

	return resourceTencentCloudClbFunctionTargetsAttachmentRead(d, meta)
}

func resourceTencentCloudClbFunctionTargetsAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_function_targets_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)

	var (
		loadBalancerId string
		listenerId     string
		locationId     string
		domain         string
		url            string
	)
	if len(idSplit) == 3 {
		loadBalancerId = idSplit[0]
		listenerId = idSplit[1]
		locationId = idSplit[2]
	} else if len(idSplit) == 4 {
		loadBalancerId = idSplit[0]
		listenerId = idSplit[1]
		domain = idSplit[2]
		url = idSplit[3]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	functionTargets, locationId, domain, url, err := service.DescribeClbFunctionTargetsAttachmentById(ctx, loadBalancerId, listenerId, locationId, domain, url)
	if err != nil {
		return err
	}

	if functionTargets == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbFunctionTargetsAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("load_balancer_id", loadBalancerId)

	_ = d.Set("listener_id", listenerId)

	if functionTargets != nil {
		functionTargetsList := []interface{}{}
		for _, functionTarget := range functionTargets {
			functionTargetsMap := map[string]interface{}{}

			if functionTarget.Function != nil {
				functionMap := map[string]interface{}{}

				if functionTarget.Function.FunctionNamespace != nil {
					functionMap["function_namespace"] = functionTarget.Function.FunctionNamespace
				}

				if functionTarget.Function.FunctionName != nil {
					functionMap["function_name"] = functionTarget.Function.FunctionName
				}

				if functionTarget.Function.FunctionQualifier != nil {
					functionMap["function_qualifier"] = functionTarget.Function.FunctionQualifier
				}

				if functionTarget.Function.FunctionQualifierType != nil {
					functionMap["function_qualifier_type"] = functionTarget.Function.FunctionQualifierType
				}

				functionTargetsMap["function"] = []interface{}{functionMap}
			}

			if functionTarget.Weight != nil {
				functionTargetsMap["weight"] = functionTarget.Weight
			}

			functionTargetsList = append(functionTargetsList, functionTargetsMap)
		}

		_ = d.Set("function_targets", functionTargetsList)

	}

	if locationId != "" {
		_ = d.Set("location_id", locationId)
	}

	if domain != "" {
		_ = d.Set("domain", domain)
	}

	if url != "" {
		_ = d.Set("url", url)
	}

	return nil
}

func resourceTencentCloudClbFunctionTargetsAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_function_targets_attachment.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	idSplit := strings.Split(d.Id(), FILED_SP)
	var (
		request        = clb.NewModifyFunctionTargetsRequest()
		loadBalancerId string
		listenerId     string
		locationId     string
		domain         string
		url            string
	)

	if len(idSplit) == 3 {
		loadBalancerId = idSplit[0]
		listenerId = idSplit[1]
		locationId = idSplit[2]
	} else if len(idSplit) == 4 {
		loadBalancerId = idSplit[0]
		listenerId = idSplit[1]
		domain = idSplit[2]
		url = idSplit[3]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerId = helper.String(listenerId)

	if locationId != "" {
		request.LocationId = helper.String(locationId)
	} else {
		request.Url = helper.String(url)
		request.Domain = helper.String(domain)
	}

	if d.HasChange("function_targets") {

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
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyFunctionTargets(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update clb functionTargetsAttachment failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClbFunctionTargetsAttachmentRead(d, meta)
}

func resourceTencentCloudClbFunctionTargetsAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_function_targets_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	idSplit := strings.Split(d.Id(), FILED_SP)
	var (
		request        = clb.NewDeregisterFunctionTargetsRequest()
		loadBalancerId string
		listenerId     string
		locationId     string
		domain         string
		url            string
	)

	if len(idSplit) == 3 {
		loadBalancerId = idSplit[0]
		listenerId = idSplit[1]
		locationId = idSplit[2]
	} else if len(idSplit) == 4 {
		loadBalancerId = idSplit[0]
		listenerId = idSplit[1]
		domain = idSplit[2]
		url = idSplit[3]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerId = helper.String(listenerId)

	if locationId != "" {
		request.LocationId = helper.String(locationId)
	} else {
		request.Url = helper.String(url)
		request.Domain = helper.String(domain)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().DeregisterFunctionTargets(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb functionTargetsAttachment failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
