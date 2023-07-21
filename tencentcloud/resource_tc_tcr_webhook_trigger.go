/*
Provides a resource to create a tcr webhook trigger

Example Usage

Create a tcr webhook trigger instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
	instance_id 	 = tencentcloud_tcr_instance.example.id
	name			 = "tf_example_ns_retention"
	is_public		 = true
	is_auto_scan	 = true
	is_prevent_vul = true
	severity		 = "medium"
	cve_whitelist_items	{
	  cve_id = "cve-xxxxx"
	}
  }

data "tencentcloud_tcr_namespaces" "example" {
	instance_id = tencentcloud_tcr_namespace.example.instance_id
  }

locals {
    ns_id = data.tencentcloud_tcr_namespaces.example.namespace_list.0.id
  }

resource "tencentcloud_tcr_webhook_trigger" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  namespace = tencentcloud_tcr_namespace.example.name
  trigger {
		name = "trigger-example"
		targets {
			address = "http://example.org/post"
			headers {
				key = "X-Custom-Header"
				values = ["a"]
			}
		}
		event_types = ["pushImage"]
		condition = ".*"
		enabled = true
		description = "example for trigger description"
		namespace_id = local.ns_id

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr webhook_trigger can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_webhook_trigger.example webhook_trigger_id
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
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcrWebhookTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrWebhookTriggerCreate,
		Read:   resourceTencentCloudTcrWebhookTriggerRead,
		Update: resourceTencentCloudTcrWebhookTriggerUpdate,
		Delete: resourceTencentCloudTcrWebhookTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance Id.",
			},

			"trigger": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "trigger parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "trigger name.",
						},
						"targets": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "trigger target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "target address.",
									},
									"headers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "custom Headers.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Header Key.",
												},
												"values": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Required:    true,
													Description: "Header Values.",
												},
											},
										},
									},
								},
							},
						},
						"event_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "trigger action.",
						},
						"condition": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "trigger rule.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "enable trigger.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "trigger Id.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "trigger description.",
						},
						"namespace_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "the namespace Id to which the trigger belongs.",
						},
					},
				},
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTcrWebhookTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_webhook_trigger.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tcr.NewCreateWebhookTriggerRequest()
		response      = tcr.NewCreateWebhookTriggerResponse()
		registryId    string
		triggerId     string
		namespaceName string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "trigger"); ok {
		webhookTrigger := tcr.WebhookTrigger{}
		if v, ok := dMap["name"]; ok {
			webhookTrigger.Name = helper.String(v.(string))
		}
		if v, ok := dMap["targets"]; ok {
			for _, item := range v.([]interface{}) {
				targetsMap := item.(map[string]interface{})
				webhookTarget := tcr.WebhookTarget{}
				if v, ok := targetsMap["address"]; ok {
					webhookTarget.Address = helper.String(v.(string))
				}
				if v, ok := targetsMap["headers"]; ok {
					for _, item := range v.([]interface{}) {
						headersMap := item.(map[string]interface{})
						header := tcr.Header{}
						if v, ok := headersMap["key"]; ok {
							header.Key = helper.String(v.(string))
						}
						if v, ok := headersMap["values"]; ok {
							valuesSet := v.(*schema.Set).List()
							for i := range valuesSet {
								values := valuesSet[i].(string)
								header.Values = append(header.Values, &values)
							}
						}
						webhookTarget.Headers = append(webhookTarget.Headers, &header)
					}
				}
				webhookTrigger.Targets = append(webhookTrigger.Targets, &webhookTarget)
			}
		}
		if v, ok := dMap["event_types"]; ok {
			eventTypesSet := v.(*schema.Set).List()
			for i := range eventTypesSet {
				eventTypes := eventTypesSet[i].(string)
				webhookTrigger.EventTypes = append(webhookTrigger.EventTypes, &eventTypes)
			}
		}
		if v, ok := dMap["condition"]; ok {
			webhookTrigger.Condition = helper.String(v.(string))
		}
		if v, ok := dMap["enabled"]; ok {
			webhookTrigger.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["description"]; ok {
			webhookTrigger.Description = helper.String(v.(string))
		}
		if v, ok := dMap["namespace_id"]; ok {
			webhookTrigger.NamespaceId = helper.IntInt64(v.(int))
		}
		request.Trigger = &webhookTrigger
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespaceName = v.(string)
		request.Namespace = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().CreateWebhookTrigger(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcr WebhookTrigger failed, reason:%+v", logId, err)
		return err
	}

	triggerId = helper.Int64ToStr(*response.Response.Trigger.Id)
	d.SetId(strings.Join([]string{registryId, namespaceName, triggerId}, FILED_SP))

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tcr:%s:uin/:repository/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrWebhookTriggerRead(d, meta)
}

func resourceTencentCloudTcrWebhookTriggerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_webhook_trigger.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	triggerId := helper.StrToInt64(idSplit[2])

	WebhookTrigger, err := service.DescribeTcrWebhookTriggerById(ctx, registryId, triggerId, namespaceName)
	if err != nil {
		return err
	}

	if WebhookTrigger == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrWebhookTrigger` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	triggerMap := map[string]interface{}{}

	if WebhookTrigger.Name != nil {
		triggerMap["name"] = WebhookTrigger.Name
	}

	if WebhookTrigger.Targets != nil {
		targetsList := []interface{}{}
		for _, targets := range WebhookTrigger.Targets {
			targetsMap := map[string]interface{}{}

			if targets.Address != nil {
				targetsMap["address"] = targets.Address
			}

			if targets.Headers != nil {
				headersList := []interface{}{}
				for _, headers := range targets.Headers {
					headersMap := map[string]interface{}{}

					if headers.Key != nil {
						headersMap["key"] = headers.Key
					}

					if headers.Values != nil {
						headersMap["values"] = headers.Values
					}

					headersList = append(headersList, headersMap)
				}
				targetsMap["headers"] = headersList
				// targetsMap["headers"] = []interface{}{headersList}
			}

			targetsList = append(targetsList, targetsMap)
		}
		triggerMap["targets"] = targetsList
		// triggerMap["targets"] = []interface{}{targetsList}
	}

	if WebhookTrigger.EventTypes != nil {
		triggerMap["event_types"] = WebhookTrigger.EventTypes
	}

	if WebhookTrigger.Condition != nil {
		triggerMap["condition"] = WebhookTrigger.Condition
	}

	if WebhookTrigger.Enabled != nil {
		triggerMap["enabled"] = WebhookTrigger.Enabled
	}

	if WebhookTrigger.Id != nil {
		triggerMap["id"] = WebhookTrigger.Id
	}

	if WebhookTrigger.Description != nil {
		triggerMap["description"] = WebhookTrigger.Description
	}

	if WebhookTrigger.NamespaceId != nil {
		triggerMap["namespace_id"] = WebhookTrigger.NamespaceId
	}

	// triggerList := []interface{}{}
	// triggerList = append(triggerList, triggerMap)
	_ = d.Set("trigger", []interface{}{triggerMap})

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tcr", "repository", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	_ = d.Set("registry_id", registryId)
	_ = d.Set("namespace", namespaceName)

	return nil
}

func resourceTencentCloudTcrWebhookTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_webhook_trigger.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcr.NewModifyWebhookTriggerRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	// triggerId := helper.StrToInt64Point(idSplit[2])

	request.RegistryId = &registryId
	request.Namespace = &namespaceName

	immutableArgs := []string{"registry_id", "namespace"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("trigger") {
		if dMap, ok := helper.InterfacesHeadMap(d, "trigger"); ok {
			webhookTrigger := tcr.WebhookTrigger{}
			if v, ok := dMap["name"]; ok {
				webhookTrigger.Name = helper.String(v.(string))
			}
			if v, ok := dMap["targets"]; ok {
				for _, item := range v.([]interface{}) {
					targetsMap := item.(map[string]interface{})
					webhookTarget := tcr.WebhookTarget{}
					if v, ok := targetsMap["address"]; ok {
						webhookTarget.Address = helper.String(v.(string))
					}
					if v, ok := targetsMap["headers"]; ok {
						for _, item := range v.([]interface{}) {
							headersMap := item.(map[string]interface{})
							header := tcr.Header{}
							if v, ok := headersMap["key"]; ok {
								header.Key = helper.String(v.(string))
							}
							if v, ok := headersMap["values"]; ok {
								valuesSet := v.(*schema.Set).List()
								for i := range valuesSet {
									values := valuesSet[i].(string)
									header.Values = append(header.Values, &values)
								}
							}
							webhookTarget.Headers = append(webhookTarget.Headers, &header)
						}
					}
					webhookTrigger.Targets = append(webhookTrigger.Targets, &webhookTarget)
				}
			}
			if v, ok := dMap["event_types"]; ok {
				eventTypesSet := v.(*schema.Set).List()
				for i := range eventTypesSet {
					eventTypes := eventTypesSet[i].(string)
					webhookTrigger.EventTypes = append(webhookTrigger.EventTypes, &eventTypes)
				}
			}
			if v, ok := dMap["condition"]; ok {
				webhookTrigger.Condition = helper.String(v.(string))
			}
			if v, ok := dMap["enabled"]; ok {
				webhookTrigger.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["id"]; ok {
				webhookTrigger.Id = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["description"]; ok {
				webhookTrigger.Description = helper.String(v.(string))
			}
			if v, ok := dMap["namespace_id"]; ok {
				webhookTrigger.NamespaceId = helper.IntInt64(v.(int))
			}
			request.Trigger = &webhookTrigger
		}
	}

	if d.HasChange("namespace") {
		if v, ok := d.GetOk("namespace"); ok {
			request.Namespace = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().ModifyWebhookTrigger(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcr WebhookTrigger failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tcr", "repository", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrWebhookTriggerRead(d, meta)
}

func resourceTencentCloudTcrWebhookTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_webhook_trigger.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	triggerId := idSplit[2]

	if err := service.DeleteTcrWebhookTriggerById(ctx, registryId, namespaceName, helper.StrToInt64(triggerId)); err != nil {
		return err
	}

	return nil
}
