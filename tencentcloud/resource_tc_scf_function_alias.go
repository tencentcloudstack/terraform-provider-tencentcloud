/*
Provides a resource to create a scf function_alias

Example Usage

```hcl
// by weight
resource "tencentcloud_scf_function_alias" "function_alias" {
  description      = "weight test"
  function_name    = "keep-1676351130"
  function_version = "$LATEST"
  name             = "weight"
  namespace        = "default"

  routing_config {
    additional_version_weights {
      version = "2"
      weight  = 0.4
    }
  }
}

// by route
resource "tencentcloud_scf_function_alias" "function_alias" {
  description      = "matchs for test 12312312"
  function_name    = "keep-1676351130"
  function_version = "3"
  name             = "matchs"
  namespace        = "default"

  routing_config {
    additional_version_matches {
      expression = "testuser"
      key        = "invoke.headers.User"
      method     = "exact"
      version    = "2"
    }
  }
}
```

Import

scf function_alias can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_alias.function_alias namespace#functionName#name
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
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudScfFunctionAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfFunctionAliasCreate,
		Read:   resourceTencentCloudScfFunctionAliasRead,
		Update: resourceTencentCloudScfFunctionAliasUpdate,
		Delete: resourceTencentCloudScfFunctionAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Alias name, which must be unique in the function, can contain 1 to 64 letters, digits, _, and -, and must begin with a letter.",
			},

			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"function_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Master version pointed to by the alias.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},

			"routing_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Request routing configuration of alias.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_version_weights": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional version with random weight-based routing.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Function version name.",
									},
									"weight": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Version weight.",
									},
								},
							},
						},
						"additional_version_matches": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional version with rule-based routing.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Function version name.",
									},
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching rule key. When the API is called, pass in the key to route the request to the specified version based on the matching ruleHeader method:Enter invoke.headers.User for key and pass in RoutingKey:{User:value} when invoking a function through invoke for invocation based on rule matching.",
									},
									"method": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match method. Valid values:range: Range matchexact: exact string match.",
									},
									"expression": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Rule requirements for range match:It should be described in an open or closed range, i.e., (a,b) or [a,b], where both a and b are integersRule requirements for exact match:Exact string match.",
									},
								},
							},
						},
					},
				},
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Alias description information.",
			},
		},
	}
}

func resourceTencentCloudScfFunctionAliasCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_alias.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = scf.NewCreateAliasRequest()
		namespace    string
		functionName string
		name         string
	)
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_version"); ok {
		request.FunctionVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
		request.Namespace = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "routing_config"); ok {
		routingConfig := scf.RoutingConfig{}
		if v, ok := dMap["additional_version_weights"]; ok {
			for _, item := range v.([]interface{}) {
				additionalVersionWeightsMap := item.(map[string]interface{})
				versionWeight := scf.VersionWeight{}
				if v, ok := additionalVersionWeightsMap["version"]; ok {
					versionWeight.Version = helper.String(v.(string))
				}
				if v, ok := additionalVersionWeightsMap["weight"]; ok {
					versionWeight.Weight = helper.Float64(v.(float64))
				}
				routingConfig.AdditionalVersionWeights = append(routingConfig.AdditionalVersionWeights, &versionWeight)
			}
		}
		if v, ok := dMap["additional_version_matches"]; ok {
			for _, item := range v.([]interface{}) {
				addtionVersionMatchsMap := item.(map[string]interface{})
				versionMatch := scf.VersionMatch{}
				if v, ok := addtionVersionMatchsMap["version"]; ok {
					versionMatch.Version = helper.String(v.(string))
				}
				if v, ok := addtionVersionMatchsMap["key"]; ok {
					versionMatch.Key = helper.String(v.(string))
				}
				if v, ok := addtionVersionMatchsMap["method"]; ok {
					versionMatch.Method = helper.String(v.(string))
				}
				if v, ok := addtionVersionMatchsMap["expression"]; ok {
					versionMatch.Expression = helper.String(v.(string))
				}
				routingConfig.AddtionVersionMatchs = append(routingConfig.AddtionVersionMatchs, &versionMatch)
			}
		}
		request.RoutingConfig = &routingConfig
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().CreateAlias(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf FunctionAlias failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(namespace + FILED_SP + functionName + FILED_SP + name)

	return resourceTencentCloudScfFunctionAliasRead(d, meta)
}

func resourceTencentCloudScfFunctionAliasRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_alias.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]
	name := idSplit[2]

	functionAlias, err := service.DescribeScfFunctionAliasById(ctx, namespace, functionName, name)
	if err != nil {
		return err
	}

	if functionAlias == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfFunctionAlias` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("name", name)
	_ = d.Set("function_name", functionName)
	_ = d.Set("namespace", namespace)

	if functionAlias.Response.FunctionVersion != nil {
		_ = d.Set("function_version", functionAlias.Response.FunctionVersion)
	}

	if functionAlias.Response.RoutingConfig != nil {
		routingConfigMap := map[string]interface{}{}

		routingConfig := functionAlias.Response.RoutingConfig

		if routingConfig.AdditionalVersionWeights != nil {
			additionalVersionWeightsList := []interface{}{}
			for _, additionalVersionWeights := range routingConfig.AdditionalVersionWeights {
				additionalVersionWeightsMap := map[string]interface{}{}

				if additionalVersionWeights.Version != nil {
					additionalVersionWeightsMap["version"] = additionalVersionWeights.Version
				}

				if additionalVersionWeights.Weight != nil {
					additionalVersionWeightsMap["weight"] = additionalVersionWeights.Weight
				}

				additionalVersionWeightsList = append(additionalVersionWeightsList, additionalVersionWeightsMap)
			}

			routingConfigMap["additional_version_weights"] = additionalVersionWeightsList
		}

		if routingConfig.AddtionVersionMatchs != nil {
			addtionVersionMatchsList := []interface{}{}
			for _, addtionVersionMatchs := range routingConfig.AddtionVersionMatchs {
				addtionVersionMatchsMap := map[string]interface{}{}

				if addtionVersionMatchs.Version != nil {
					addtionVersionMatchsMap["version"] = addtionVersionMatchs.Version
				}

				if addtionVersionMatchs.Key != nil {
					addtionVersionMatchsMap["key"] = addtionVersionMatchs.Key
				}

				if addtionVersionMatchs.Method != nil {
					addtionVersionMatchsMap["method"] = addtionVersionMatchs.Method
				}

				if addtionVersionMatchs.Expression != nil {
					addtionVersionMatchsMap["expression"] = addtionVersionMatchs.Expression
				}

				addtionVersionMatchsList = append(addtionVersionMatchsList, addtionVersionMatchsMap)
			}

			routingConfigMap["additional_version_matches"] = addtionVersionMatchsList
		}

		_ = d.Set("routing_config", []interface{}{routingConfigMap})
	}

	if functionAlias.Response.Description != nil {
		_ = d.Set("description", functionAlias.Response.Description)
	}

	return nil
}

func resourceTencentCloudScfFunctionAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_alias.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := scf.NewUpdateAliasRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]
	name := idSplit[2]

	request.Namespace = &namespace
	request.FunctionName = &functionName
	request.Name = &name
	request.FunctionVersion = helper.String(d.Get("function_version").(string))

	mutableArgs := []string{"routing_config", "description"}

	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if dMap, ok := helper.InterfacesHeadMap(d, "routing_config"); ok {
			routingConfig := scf.RoutingConfig{}
			if v, ok := dMap["additional_version_weights"]; ok {
				for _, item := range v.([]interface{}) {
					additionalVersionWeightsMap := item.(map[string]interface{})
					versionWeight := scf.VersionWeight{}
					if v, ok := additionalVersionWeightsMap["version"]; ok {
						versionWeight.Version = helper.String(v.(string))
					}
					if v, ok := additionalVersionWeightsMap["weight"]; ok {
						versionWeight.Weight = helper.Float64(v.(float64))
					}
					routingConfig.AdditionalVersionWeights = append(routingConfig.AdditionalVersionWeights, &versionWeight)
				}
			}
			if v, ok := dMap["additional_version_matches"]; ok {
				for _, item := range v.([]interface{}) {
					addtionVersionMatchsMap := item.(map[string]interface{})
					versionMatch := scf.VersionMatch{}
					if v, ok := addtionVersionMatchsMap["version"]; ok {
						versionMatch.Version = helper.String(v.(string))
					}
					if v, ok := addtionVersionMatchsMap["key"]; ok {
						versionMatch.Key = helper.String(v.(string))
					}
					if v, ok := addtionVersionMatchsMap["method"]; ok {
						versionMatch.Method = helper.String(v.(string))
					}
					if v, ok := addtionVersionMatchsMap["expression"]; ok {
						versionMatch.Expression = helper.String(v.(string))
					}
					routingConfig.AddtionVersionMatchs = append(routingConfig.AddtionVersionMatchs, &versionMatch)
				}
			}
			request.RoutingConfig = &routingConfig
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().UpdateAlias(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update scf FunctionAlias failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudScfFunctionAliasRead(d, meta)
}

func resourceTencentCloudScfFunctionAliasDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_alias.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]
	name := idSplit[2]

	if err := service.DeleteScfFunctionAliasById(ctx, namespace, functionName, name); err != nil {
		return err
	}

	return nil
}
