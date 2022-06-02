/*

 */
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsIndexCreate,
		Read:   resourceTencentCloudClsIndexRead,
		Delete: resourceTencentCloudClsIndexDelete,
		Update: resourceTencentCloudClsIndexUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log topic ID.",
			},
			"rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Index rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_text": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Full-Text index configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"case_sensitive": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Case sensitivity.",
									},
									"tokenizer": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Full-Text index delimiter. Each character in the string represents a delimiter.",
									},
									"contain_z_h": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether Chinese characters are contained.",
									},
								},
							},
						},
						"key_value": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Key-Value index configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"case_sensitive": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Case sensitivity.",
									},
									"key_values": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Key-Value pair information of the index to be created. Up to 100 key-value pairs can be configured.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
													Description: "When a key value or metafield index needs to be configured for a field, the metafield Key does not need to be prefixed with __TAG__. and is consistent " +
														"with the one when logs are uploaded. __TAG__. will be prefixed automatically for display in the console..",
												},
												"value": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Field index description information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Field type. Valid values: long, text, double.",
															},
															"tokenizer": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Field delimiter, which is meaningful only if the field type is text. Each character in the entered string represents a delimiter.",
															},
															"sql_flag": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Whether the analysis feature is enabled for the field.",
															},
															"contain_z_h": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Whether Chinese characters are contained.",
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
						"tag": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Metafield index configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"case_sensitive": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Case sensitivity.",
									},
									"key_values": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Key-Value pair information of the index to be created. Up to 100 key-value pairs can be configured.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
													Description: "When a key value or metafield index needs to be configured for a field, the metafield Key does not need to be prefixed with __TAG__. and is consistent " +
														"with the one when logs are uploaded. __TAG__. will be prefixed automatically for display in the console..",
												},
												"value": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Field index description information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Field type. Valid values: long, text, double.",
															},
															"tokenizer": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Field delimiter, which is meaningful only if the field type is text. Each character in the entered string represents a delimiter.",
															},
															"sql_flag": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Whether the analysis feature is enabled for the field.",
															},
															"contain_z_h": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Whether Chinese characters are contained.",
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
					},
				},
			},
			"status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to take effect. Default value: true.",
			},
			"include_internal_fields": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Internal field marker of full-text index. Default value: false. Valid value: false: excluding internal fields; true: including internal fields.",
			},
			"metadata_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Metadata flag. Default value: 0. Valid value: 0: full-text index (including the metadata field with key-value index enabled); 1: full-text index (including all metadata fields); 2: full-text index (excluding metadata fields)..",
			},
		},
	}
}

func resourceTencentCloudClsIndexCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_index.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateIndexRequest()
		response *cls.CreateIndexResponse
		indexId  string
	)

	if v, ok := d.GetOk("topic_id"); ok {
		indexId = v.(string)
		request.TopicId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		ruleInfo := cls.RuleInfo{}
		if fullText, ok := dMap["full_text"]; ok {
			fullTextMap := fullText.([]interface{})[0].(map[string]interface{})
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

		if tag, ok := dMap["tag"]; ok {
			tagMap := tag.([]interface{})[0].(map[string]interface{})
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
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls index failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(indexId)

	return resourceTencentCloudClsIndexRead(d, meta)
}

func resourceTencentCloudClsIndexRead(d *schema.ResourceData, meta interface{}) error {
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
			keyValueMap := map[string]interface{}{
				"case_sensitive": res.Rule.KeyValue.CaseSensitive,
			}
			if res.Rule.KeyValue.KeyValues != nil {
				KeyValuesMap := map[string]interface{}{}
				for _, item := range res.Rule.KeyValue.KeyValues {
					cls.KeyValueInfo
				}

			}
			ruleMap["full_text"] = []interface{}{keyValueMap}
		}

	}

	return nil
}

func resourceTencentCloudClsIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_shipper.update")()
	logId := getLogId(contextNil)
	request := cls.NewModifyShipperRequest()

	request.ShipperId = helper.String(d.Id())

	if d.HasChange("bucket") {
		if v, ok := d.GetOk("bucket"); ok {
			request.Bucket = helper.String(v.(string))
		}
	}

	if d.HasChange("prefix") {
		if v, ok := d.GetOk("prefix"); ok {
			request.Prefix = helper.String(v.(string))
		}
	}

	if d.HasChange("shipper_name") {
		if v, ok := d.GetOk("shipper_name"); ok {
			request.ShipperName = helper.String(v.(string))
		}
	}

	if d.HasChange("interval") {
		if v, ok := d.GetOk("interval"); ok {
			request.Interval = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_size") {
		if v, ok := d.GetOk("max_size"); ok {
			request.MaxSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("filter_rules") {
		if v, ok := d.GetOk("filter_rules"); ok {
			filterRules := make([]*cls.FilterRuleInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				filterRule := cls.FilterRuleInfo{}
				if v, ok := dMap["key"]; ok {
					filterRule.Key = helper.String(v.(string))
				}
				if v, ok := dMap["regex"]; ok {
					filterRule.Regex = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					filterRule.Value = helper.String(v.(string))
				}
				filterRules = append(filterRules, &filterRule)
			}
			request.FilterRules = filterRules
		}
	}

	if d.HasChange("partition") {
		if v, ok := d.GetOk("partition"); ok {
			request.Partition = helper.String(v.(string))
		}
	}

	if d.HasChange("compress") {
		if v, ok := d.GetOk("compress"); ok {
			compresses := make([]*cls.CompressInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one compress.")
			}
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				compress := cls.CompressInfo{}
				if v, ok := dMap["format"]; ok {
					compress.Format = helper.String(v.(string))
				}
				compresses = append(compresses, &compress)
			}
			request.Compress = compresses[0]
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			contents := make([]*cls.ContentInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one content.")
			}
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				content := cls.ContentInfo{}
				if v, ok := dMap["format"]; ok {
					content.Format = helper.String(v.(string))
				}
				if v, ok := dMap["csv"]; ok {
					if len(v.([]interface{})) == 1 {
						csv := v.([]interface{})[0].(map[string]interface{})
						csvInfo := cls.CsvInfo{}
						csvInfo.PrintKey = helper.Bool(csv["print_key"].(bool))
						keys := csv["keys"].(*schema.Set).List()
						for _, key := range keys {
							csvInfo.Keys = append(csvInfo.Keys, helper.String(key.(string)))
						}
						csvInfo.Delimiter = helper.String(csv["delimiter"].(string))
						csvInfo.EscapeChar = helper.String(csv["escape_char"].(string))
						csvInfo.NonExistingField = helper.String(csv["non_existing_field"].(string))
						content.Csv = &csvInfo
					}
				}
				if v, ok := dMap["json"]; ok {
					if len(v.([]interface{})) == 1 {

						json := v.([]interface{})[0].(map[string]interface{})
						jsonInfo := cls.JsonInfo{}
						jsonInfo.EnableTag = helper.Bool(json["enable_tag"].(bool))
						metaFields := json["meta_fields"].(*schema.Set).List()
						for _, metaField := range metaFields {
							jsonInfo.MetaFields = append(jsonInfo.MetaFields, helper.String(metaField.(string)))
						}
						content.Json = &jsonInfo
					}
				}
				contents = append(contents, &content)
			}
			request.Content = contents[0]
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyShipper(request)
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

	return resourceTencentCloudClsCosShipperRead(d, meta)
}

func resourceTencentCloudClsIndexDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_shipper.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteClsCosShipper(ctx, id); err != nil {
		return err
	}

	return nil
}
