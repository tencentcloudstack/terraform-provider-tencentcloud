/*
Provides a resource to create a tcr webhook_trigger

Example Usage

```hcl
resource "tencentcloud_tcr_webhook_trigger" "webhook_trigger" {
  registry_id = "tcr-xxx"
  trigger {
		name = "trigger"
		targets {
			address = "http://example.org/post"
			headers {
				key = "X-Custom-Header"
				values =
			}
		}
		event_types =
		condition = ".*"
		enabled = true
		id = 20
		description = "this is trigger description"
		namespace_id = 10

  }
  namespace = "trigger"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr webhook_trigger can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_webhook_trigger.webhook_trigger webhook_trigger_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
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
				Description: "Instance Id.",
			},

			"trigger": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Trigger parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger name.",
						},
						"targets": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Trigger target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Target address.",
									},
									"headers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Custom Headers.",
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
							Description: "Trigger action.",
						},
						"condition": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger rule.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable trigger.",
						},
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "TriggerId.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trigger description.",
						},
						"namespace_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The namespace Id to which the trigger belongs.",
						},
					},
				},
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
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
		request    = tcr.NewCreateWebhookTriggerRequest()
		response   = tcr.NewCreateWebhookTriggerResponse()
		registryId string
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

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().CreateWebhookTrigger(request)
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

	registryId = *response.Response.RegistryId
	d.SetId(registryId)

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

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	webhookTriggerId := d.Id()

	WebhookTrigger, err := service.DescribeTcrWebhookTriggerById(ctx, registryId)
	if err != nil {
		return err
	}

	if WebhookTrigger == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrWebhookTrigger` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if WebhookTrigger.RegistryId != nil {
		_ = d.Set("registry_id", WebhookTrigger.RegistryId)
	}

	if WebhookTrigger.Trigger != nil {
		triggerMap := map[string]interface{}{}

		if WebhookTrigger.Trigger.Name != nil {
			triggerMap["name"] = WebhookTrigger.Trigger.Name
		}

		if WebhookTrigger.Trigger.Targets != nil {
			targetsList := []interface{}{}
			for _, targets := range WebhookTrigger.Trigger.Targets {
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

					targetsMap["headers"] = []interface{}{headersList}
				}

				targetsList = append(targetsList, targetsMap)
			}

			triggerMap["targets"] = []interface{}{targetsList}
		}

		if WebhookTrigger.Trigger.EventTypes != nil {
			triggerMap["event_types"] = WebhookTrigger.Trigger.EventTypes
		}

		if WebhookTrigger.Trigger.Condition != nil {
			triggerMap["condition"] = WebhookTrigger.Trigger.Condition
		}

		if WebhookTrigger.Trigger.Enabled != nil {
			triggerMap["enabled"] = WebhookTrigger.Trigger.Enabled
		}

		if WebhookTrigger.Trigger.Id != nil {
			triggerMap["id"] = WebhookTrigger.Trigger.Id
		}

		if WebhookTrigger.Trigger.Description != nil {
			triggerMap["description"] = WebhookTrigger.Trigger.Description
		}

		if WebhookTrigger.Trigger.NamespaceId != nil {
			triggerMap["namespace_id"] = WebhookTrigger.Trigger.NamespaceId
		}

		_ = d.Set("trigger", []interface{}{triggerMap})
	}

	if WebhookTrigger.Namespace != nil {
		_ = d.Set("namespace", WebhookTrigger.Namespace)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tcr", "repository", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTcrWebhookTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_webhook_trigger.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcr.NewModifyWebhookTriggerRequest()

	webhookTriggerId := d.Id()

	request.RegistryId = &registryId

	immutableArgs := []string{"registry_id", "trigger", "namespace"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("registry_id") {
		if v, ok := d.GetOk("registry_id"); ok {
			request.RegistryId = helper.String(v.(string))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().ModifyWebhookTrigger(request)
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

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}
	webhookTriggerId := d.Id()

	if err := service.DeleteTcrWebhookTriggerById(ctx, registryId); err != nil {
		return err
	}

	return nil
}
