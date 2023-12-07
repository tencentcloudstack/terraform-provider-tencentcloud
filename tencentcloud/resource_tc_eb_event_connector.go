package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudEbEventConnector() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbEventConnectorCreate,
		Read:   resourceTencentCloudEbEventConnectorRead,
		Update: resourceTencentCloudEbEventConnectorUpdate,
		Delete: resourceTencentCloudEbEventConnectorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"connection_description": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Connector description.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource qcs six-segment style, more reference [resource six-segment style](https://cloud.tencent.com/document/product/598/10606).",
						},
						"api_gw_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "apigw parameter,Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "HTTPS.",
									},
									"method": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "POST.",
									},
								},
							},
						},
						"ckafka_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "ckafka parameter, note: this field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"offset": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "kafka offset.",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ckafka  topic.",
									},
								},
							},
						},
					},
				},
			},

			"event_bus_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "event bus Id.",
			},

			"connection_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "connector name.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "description.",
			},

			"enable": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "switch.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "type.",
			},
		},
	}
}

func resourceTencentCloudEbEventConnectorCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_connector.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = eb.NewCreateConnectionRequest()
		response     = eb.NewCreateConnectionResponse()
		connectionId string
		eventBusId   string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "connection_description"); ok {
		connectionDescription := eb.ConnectionDescription{}
		if v, ok := dMap["resource_description"]; ok {
			connectionDescription.ResourceDescription = helper.String(v.(string))
		}
		if aPIGWParamsMap, ok := helper.InterfaceToMap(dMap, "api_gw_params"); ok {
			aPIGWParams := eb.APIGWParams{}
			if v, ok := aPIGWParamsMap["protocol"]; ok {
				aPIGWParams.Protocol = helper.String(v.(string))
			}
			if v, ok := aPIGWParamsMap["method"]; ok {
				aPIGWParams.Method = helper.String(v.(string))
			}
			connectionDescription.APIGWParams = &aPIGWParams
		}
		if ckafkaParamsMap, ok := helper.InterfaceToMap(dMap, "ckafka_params"); ok {
			ckafkaParams := eb.CkafkaParams{}
			if v, ok := ckafkaParamsMap["offset"]; ok {
				ckafkaParams.Offset = helper.String(v.(string))
			}
			if v, ok := ckafkaParamsMap["topic_name"]; ok {
				ckafkaParams.TopicName = helper.String(v.(string))
			}
			connectionDescription.CkafkaParams = &ckafkaParams
		}
		request.ConnectionDescription = &connectionDescription
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("connection_name"); ok {
		request.ConnectionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().CreateConnection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eb eventConnector failed, reason:%+v", logId, err)
		return err
	}

	connectionId = *response.Response.ConnectionId
	d.SetId(eventBusId + FILED_SP + connectionId)

	return resourceTencentCloudEbEventConnectorRead(d, meta)
}

func resourceTencentCloudEbEventConnectorRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_connector.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	connectionId := idSplit[1]

	eventConnector, err := service.DescribeEbEventConnectorById(ctx, connectionId, eventBusId)
	if err != nil {
		return err
	}

	if eventConnector == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventConnector` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if eventConnector.ConnectionDescription != nil {
		connectionDescriptionMap := map[string]interface{}{}

		if eventConnector.ConnectionDescription.ResourceDescription != nil {
			description := *eventConnector.ConnectionDescription.ResourceDescription
			if eventConnector.Type != nil && *eventConnector.Type == "apigw" && strings.Contains(description, "serviceid") && strings.Contains(description, "API") {
				parts := strings.Split(description, "/")
				result := strings.Join(parts[:len(parts)-2], "/")
				connectionDescriptionMap["resource_description"] = &result
			} else {
				connectionDescriptionMap["resource_description"] = eventConnector.ConnectionDescription.ResourceDescription
			}
		}

		if eventConnector.ConnectionDescription.APIGWParams != nil {
			aPIGWParamsMap := map[string]interface{}{}

			if eventConnector.ConnectionDescription.APIGWParams.Protocol != nil {
				aPIGWParamsMap["protocol"] = eventConnector.ConnectionDescription.APIGWParams.Protocol
			}

			if eventConnector.ConnectionDescription.APIGWParams.Method != nil {
				aPIGWParamsMap["method"] = eventConnector.ConnectionDescription.APIGWParams.Method
			}

			connectionDescriptionMap["api_gw_params"] = []interface{}{aPIGWParamsMap}
		}

		if eventConnector.ConnectionDescription.CkafkaParams != nil {
			ckafkaParamsMap := map[string]interface{}{}

			if eventConnector.ConnectionDescription.CkafkaParams.Offset != nil {
				ckafkaParamsMap["offset"] = eventConnector.ConnectionDescription.CkafkaParams.Offset
			}

			if eventConnector.ConnectionDescription.CkafkaParams.TopicName != nil {
				ckafkaParamsMap["topic_name"] = eventConnector.ConnectionDescription.CkafkaParams.TopicName
			}

			connectionDescriptionMap["ckafka_params"] = []interface{}{ckafkaParamsMap}
		}

		_ = d.Set("connection_description", []interface{}{connectionDescriptionMap})
	}

	if eventConnector.EventBusId != nil {
		_ = d.Set("event_bus_id", eventConnector.EventBusId)
	}

	if eventConnector.ConnectionName != nil {
		_ = d.Set("connection_name", eventConnector.ConnectionName)
	}

	if eventConnector.Description != nil {
		_ = d.Set("description", eventConnector.Description)
	}

	if eventConnector.Enable != nil {
		_ = d.Set("enable", eventConnector.Enable)
	}

	if eventConnector.Type != nil {
		_ = d.Set("type", eventConnector.Type)
	}

	return nil
}

func resourceTencentCloudEbEventConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_connector.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := eb.NewUpdateConnectionRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	connectionId := idSplit[1]

	request.ConnectionId = &connectionId
	request.EventBusId = &eventBusId

	immutableArgs := []string{"connection_description", "type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("connection_name") {
		if v, ok := d.GetOk("connection_name"); ok {
			request.ConnectionName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("enable") {
		if v, ok := d.GetOkExists("enable"); ok {
			request.Enable = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().UpdateConnection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update eb eventConnector failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudEbEventConnectorRead(d, meta)
}

func resourceTencentCloudEbEventConnectorDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_connector.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	connectionId := idSplit[1]

	if err := service.DeleteEbEventConnectorById(ctx, connectionId, eventBusId); err != nil {
		return err
	}

	return nil
}
