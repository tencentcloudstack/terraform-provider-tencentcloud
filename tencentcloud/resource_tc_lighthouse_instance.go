/*
Provides a resource to create a cls index.

Example Usage

```hcl
resource "tencentcloud_cls_index" "index" {
  topic_id = "0937e56f-4008-49d2-ad2d-69c52a9f11cc"

  rule {
    full_text {
      case_sensitive = true
      tokenizer      = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
      contain_z_h    = true
    }

    key_value {
      case_sensitive = true
      key_values {
        key = "hello"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
          type        = "text"
        }
      }

      key_values {
        key = "world"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
          type        = "text"
        }
      }
    }

    tag {
      case_sensitive = true
      key_values {
        key = "terraform"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
          type        = "text"
        }
      }
    }
  }
  status                  = true
  include_internal_fields = true
  metadata_flag           = 1
}
```

Import

cls cos index can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_index.index 0937e56f-4008-49d2-ad2d-69c52a9f11cc
```
*/
package tencentcloud

import (
	"context"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseInstanceCreate,
		Read:   resourceTencentCloudLighthouseInstanceRead,
		Delete: resourceTencentCloudLighthouseInstanceDelete,
		Update: resourceTencentCloudLighthouseInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bundle_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Lighthouse package.",
			},
			"blueprint_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Lighthouse image.",
			},
			"instance_charge_prepaid": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Monthly subscription information for the instance, including the purchase period, setting of auto-renewal, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Subscription period in months. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.",
						},
						"RenewFlag": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Auto-Renewal flag. Valid values: NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically; NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically. You need to manually renew DISABLE_NOTIFY_AND_AUTO_RENEW: neither notify upon expiration nor renew automatically. " +
								"Default value: NOTIFY_AND_MANUAL_RENEW.",
						},
					},
				},
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the Lighthouse instance.",
			},
			"zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of availability zones. A random AZ is selected by default.",
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether the request is a dry run only." +
					"true: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available. If the dry run fails, the corresponding error code will be returned.If the dry run succeeds, the RequestId will be returned." +
					"false (default value): send a normal request and create instance(s) if all the requirements are met.",
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.",
			},
			"login_configuration": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Login password of the instance. It’s only available for Windows instances. If it’s not specified, it means that the user choose to set the login password after the instance creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_generate_password": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "whether auto generate password. if false, need set password.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login password.",
						},
					},
				},
			},
			"containers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the containers to create.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Container image address.",
						},
						"container_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Container name.",
						},
						"envs": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of environment variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Environment variable key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Environment variable value.",
									},
								},
							},
						},
						"publish_ports": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of mappings of container ports and host ports.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Host port.",
									},
									"container_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Container port.",
									},
									"ip": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "External IP. It defaults to 0.0.0.0.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The protocol defaults to tcp. Valid values: tcp, udp and sctp.",
									},
								},
							},
						},
						"volumes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of container mount volumes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Container path.",
									},
									"host_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Host path.",
									},
								},
							},
						},
						"command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The command to run.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudLighthouseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = lighthouse.CreateInstancesRequest{}
		indexId string
	)

	if v, ok := d.GetOk("bundle_id"); ok {
		request.BundleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("blueprint_id"); ok {
		request.BlueprintId = helper.String(v.(string))
	}

	if instanceChargePrepaidMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		ruleInfo := lighthouse.InstanceChargePrepaid{}
		if instanceChargePrepaidMap, ok := helper.InterfaceToMap(instanceChargePrepaidMap, )
			fullTextInfo := cls.FullTextInfo{}
			if v, ok := fullTextMap["case_sensitive"]; ok {
				fullTextInfo.CaseSensitive = helper.Bool(v.(bool))
			}
			if v, ok := fullTextMap["tokenizer"]; ok {
				fullTextInfo.Tokenizer = helper.String(v.(string))
			}
			if v, ok := fullTextMap["contain_z_h"]; ok {
				fullTextInfo.ContainZH = helper.Bool(v.(bool))
			}
			ruleInfo.FullText = &fullTextInfo
		}

		if keyValue, ok := dMap["key_value"]; ok {
			ruleKeyValueMap := keyValue.([]interface{})[0].(map[string]interface{})
			ruleKeyValueInfo := cls.RuleKeyValueInfo{}
			if v, ok := ruleKeyValueMap["case_sensitive"]; ok {
				ruleKeyValueInfo.CaseSensitive = helper.Bool(v.(bool))
			}
			if v, ok := ruleKeyValueMap["key_values"]; ok {
				for _, keyValue := range v.([]interface{}) {
					keyValueInfo := cls.KeyValueInfo{}
					keyValueMap := keyValue.(map[string]interface{})
					if v, ok := keyValueMap["key"]; ok {
						keyValueInfo.Key = helper.String(v.(string))
					}
					if v, ok := keyValueMap["value"]; ok {
						valueMap := v.([]interface{})[0].(map[string]interface{})
						valueInfo := cls.ValueInfo{}
						if v, ok := valueMap["type"]; ok {
							valueInfo.Type = helper.String(v.(string))
						}
						if v, ok := valueMap["tokenizer"]; ok {
							valueInfo.Tokenizer = helper.String(v.(string))
						}
						if v, ok := valueMap["sql_flag"]; ok {
							valueInfo.SqlFlag = helper.Bool(v.(bool))
						}
						if v, ok := valueMap["contain_z_h"]; ok {
							valueInfo.ContainZH = helper.Bool(v.(bool))
						}
						keyValueInfo.Value = &valueInfo
					}
					ruleKeyValueInfo.KeyValues = append(ruleKeyValueInfo.KeyValues, &keyValueInfo)
				}
			}

			ruleInfo.KeyValue = &ruleKeyValueInfo
		}

		if tagMap, ok := helper.InterfaceToMap(dMap, "tag"); ok {
			ruleTagInfo := cls.RuleTagInfo{}
			if v, ok := tagMap["case_sensitive"]; ok {
				ruleTagInfo.CaseSensitive = helper.Bool(v.(bool))
			}
			if v, ok := tagMap["key_values"]; ok {
				for _, keyValue := range v.([]interface{}) {
					keyValueInfo := cls.KeyValueInfo{}
					keyValueMap := keyValue.(map[string]interface{})
					if v, ok := keyValueMap["key"]; ok {
						keyValueInfo.Key = helper.String(v.(string))
					}
					if v, ok := keyValueMap["value"]; ok {
						valueMap := v.([]interface{})[0].(map[string]interface{})
						valueInfo := cls.ValueInfo{}
						if v, ok := valueMap["type"]; ok {
							valueInfo.Type = helper.String(v.(string))
						}
						if v, ok := valueMap["tokenizer"]; ok {
							valueInfo.Tokenizer = helper.String(v.(string))
						}
						if v, ok := valueMap["sql_flag"]; ok {
							valueInfo.SqlFlag = helper.Bool(v.(bool))
						}
						if v, ok := valueMap["contain_z_h"]; ok {
							valueInfo.ContainZH = helper.Bool(v.(bool))
						}
						keyValueInfo.Value = &valueInfo
					}
					ruleTagInfo.KeyValues = append(ruleTagInfo.KeyValues, &keyValueInfo)
				}
			}
			ruleInfo.Tag = &ruleTagInfo
		}
		request.Rule = &ruleInfo
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("include_internal_fields"); ok {
		request.IncludeInternalFields = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("metadata_flag"); ok {
		request.MetadataFlag = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateIndex(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls index failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(indexId)

	return resourceTencentCloudLighthouseInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_index.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = cls.NewDescribeIndexRequest()
		result  *cls.DescribeIndexResponse
	)
	id := d.Id()

	request.TopicId = &id

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().DescribeIndex(request)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cls index failed, reason:%s\n", logId, err.Error())
		return err
	}

	res := result.Response

	if res.TopicId != nil {
		_ = d.Set("topic_id", res.TopicId)
	}

	if res.Rule != nil {
		ruleMap := map[string]interface{}{}

		if res.Rule.FullText != nil {
			fullTextMap := map[string]interface{}{
				"case_sensitive": res.Rule.FullText.CaseSensitive,
				"tokenizer":      res.Rule.FullText.Tokenizer,
				"contain_z_h":    res.Rule.FullText.ContainZH,
			}
			ruleMap["full_text"] = []interface{}{fullTextMap}
		}

		if res.Rule.KeyValue != nil {
			ruleKeyValueMap := map[string]interface{}{
				"case_sensitive": res.Rule.KeyValue.CaseSensitive,
			}
			if res.Rule.KeyValue.KeyValues != nil {
				keyValuesList := []interface{}{}
				for _, keyValueInfo := range res.Rule.KeyValue.KeyValues {
					keyValueInfoMap := map[string]interface{}{
						"key": keyValueInfo.Key,
					}
					if keyValueInfo.Value != nil {
						valueInfoMap := map[string]interface{}{
							"type":        keyValueInfo.Value.Type,
							"tokenizer":   keyValueInfo.Value.Tokenizer,
							"sql_flag":    keyValueInfo.Value.SqlFlag,
							"contain_z_h": keyValueInfo.Value.ContainZH,
						}
						keyValueInfoMap["value"] = []interface{}{valueInfoMap}
					}
					keyValuesList = append(keyValuesList, keyValueInfoMap)
				}
				ruleKeyValueMap["key_values"] = keyValuesList
			}
			ruleMap["key_value"] = []interface{}{ruleKeyValueMap}
		}

		if res.Rule.Tag != nil {
			ruleTagMap := map[string]interface{}{
				"case_sensitive": res.Rule.Tag.CaseSensitive,
			}
			if res.Rule.Tag.KeyValues != nil {
				keyValuesList := []interface{}{}
				for _, keyValueInfo := range res.Rule.Tag.KeyValues {
					keyValueInfoMap := map[string]interface{}{
						"key": keyValueInfo.Key,
					}
					if keyValueInfo.Value != nil {
						valueInfoMap := map[string]interface{}{
							"type":        keyValueInfo.Value.Type,
							"tokenizer":   keyValueInfo.Value.Tokenizer,
							"sql_flag":    keyValueInfo.Value.SqlFlag,
							"contain_z_h": keyValueInfo.Value.ContainZH,
						}
						keyValueInfoMap["value"] = []interface{}{valueInfoMap}
					}
					keyValuesList = append(keyValuesList, keyValueInfoMap)
				}
				ruleTagMap["key_values"] = keyValuesList
			}
			ruleMap["tag"] = []interface{}{ruleTagMap}
		}

		_ = d.Set("rule", []interface{}{ruleMap})
	}

	if res.Status != nil {
		_ = d.Set("status", res.Status)
	}

	if res.IncludeInternalFields != nil {
		_ = d.Set("include_internal_fields", res.IncludeInternalFields)
	}

	if res.MetadataFlag != nil {
		_ = d.Set("metadata_flag", res.MetadataFlag)
	}

	return nil
}

func resourceTencentCloudLighthouseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_index.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = cls.NewModifyIndexRequest()
	)
	id := d.Id()

	request.TopicId = &id

	if d.HasChange("rule") || d.HasChange("status") || d.HasChange("include_internal_fields") || d.HasChange("metadata_flag") {
		if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
			ruleInfo := cls.RuleInfo{}
			if fullTextMap, ok := helper.InterfaceToMap(dMap, "full_text"); ok {
				fullTextInfo := cls.FullTextInfo{}
				if v, ok := fullTextMap["case_sensitive"]; ok {
					fullTextInfo.CaseSensitive = helper.Bool(v.(bool))
				}
				if v, ok := fullTextMap["tokenizer"]; ok {
					fullTextInfo.Tokenizer = helper.String(v.(string))
				}
				if v, ok := fullTextMap["contain_z_h"]; ok {
					fullTextInfo.ContainZH = helper.Bool(v.(bool))
				}
				ruleInfo.FullText = &fullTextInfo
			}

			if ruleKeyValueMap, ok := helper.InterfaceToMap(dMap, "key_value"); ok {
				ruleKeyValueInfo := cls.RuleKeyValueInfo{}
				if v, ok := ruleKeyValueMap["case_sensitive"]; ok {
					ruleKeyValueInfo.CaseSensitive = helper.Bool(v.(bool))
				}
				if v, ok := ruleKeyValueMap["key_values"]; ok {
					for _, keyValue := range v.([]interface{}) {
						keyValueInfo := cls.KeyValueInfo{}
						keyValueMap := keyValue.(map[string]interface{})
						if v, ok := keyValueMap["key"]; ok {
							keyValueInfo.Key = helper.String(v.(string))
						}
						if v, ok := keyValueMap["value"]; ok {
							valueMap := v.([]interface{})[0].(map[string]interface{})
							valueInfo := cls.ValueInfo{}
							if v, ok := valueMap["type"]; ok {
								valueInfo.Type = helper.String(v.(string))
							}
							if v, ok := valueMap["tokenizer"]; ok {
								valueInfo.Tokenizer = helper.String(v.(string))
							}
							if v, ok := valueMap["sql_flag"]; ok {
								valueInfo.SqlFlag = helper.Bool(v.(bool))
							}
							if v, ok := valueMap["contain_z_h"]; ok {
								valueInfo.ContainZH = helper.Bool(v.(bool))
							}
							keyValueInfo.Value = &valueInfo
						}
						ruleKeyValueInfo.KeyValues = append(ruleKeyValueInfo.KeyValues, &keyValueInfo)
					}
				}

				ruleInfo.KeyValue = &ruleKeyValueInfo
			}

			if tagMap, ok := helper.InterfaceToMap(dMap, "tag"); ok {
				ruleTagInfo := cls.RuleTagInfo{}
				if v, ok := tagMap["case_sensitive"]; ok {
					ruleTagInfo.CaseSensitive = helper.Bool(v.(bool))
				}
				if v, ok := tagMap["key_values"]; ok {
					for _, keyValue := range v.([]interface{}) {
						keyValueInfo := cls.KeyValueInfo{}
						keyValueMap := keyValue.(map[string]interface{})
						if v, ok := keyValueMap["key"]; ok {
							keyValueInfo.Key = helper.String(v.(string))
						}
						if v, ok := keyValueMap["value"]; ok {
							valueMap := v.([]interface{})[0].(map[string]interface{})
							valueInfo := cls.ValueInfo{}
							if v, ok := valueMap["type"]; ok {
								valueInfo.Type = helper.String(v.(string))
							}
							if v, ok := valueMap["tokenizer"]; ok {
								valueInfo.Tokenizer = helper.String(v.(string))
							}
							if v, ok := valueMap["sql_flag"]; ok {
								valueInfo.SqlFlag = helper.Bool(v.(bool))
							}
							if v, ok := valueMap["contain_z_h"]; ok {
								valueInfo.ContainZH = helper.Bool(v.(bool))
							}
							keyValueInfo.Value = &valueInfo
						}
						ruleTagInfo.KeyValues = append(ruleTagInfo.KeyValues, &keyValueInfo)
					}
				}
				ruleInfo.Tag = &ruleTagInfo
			}
			request.Rule = &ruleInfo
		}

		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("include_internal_fields"); ok {
			request.IncludeInternalFields = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("metadata_flag"); ok {
			request.MetadataFlag = helper.IntUint64(v.(int))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyIndex(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return resourceTencentCloudLighthouseInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_shipper.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteClsIndex(ctx, id); err != nil {
		return err
	}
	return nil
}
