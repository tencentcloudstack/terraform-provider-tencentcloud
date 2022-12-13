/*
Provides a resource to create a tem application_service

Example Usage

```hcl
resource "tencentcloud_tem_application_service" "application_service" {
  environment_id = "en-dpxyydl5"
  application_id = "app-jrl3346j"
  service {
		type = "CLUSTER"
		service_name = "test0-1"
		port_mapping_item_list {
			port = 80
			target_port = 80
			protocol = "tcp"
		}
  }
}
```

Import

tem application_service can be imported using the environmentId#applicationId#serviceName, e.g.

```
terraform import tencentcloud_tem_application_service.application_service application_service_id
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
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTemApplicationService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemApplicationServiceCreate,
		Read:   resourceTencentCloudTemApplicationServiceRead,
		Update: resourceTencentCloudTemApplicationServiceUpdate,
		Delete: resourceTencentCloudTemApplicationServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "environment ID.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "application ID.",
			},

			"service": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "service detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "application service type: EXTERNAL | VPC | CLUSTER.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "application service name.",
						},
						"port_mapping_item_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "port mapping item list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "container port.",
									},
									"target_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "application listen port.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "udp or tcp.",
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

func resourceTencentCloudTemApplicationServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewCreateApplicationServiceRequest()
		environmentId string
		applicationId string
		serviceName   string
	)
	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		applicationId = v.(string)
		request.ApplicationId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service"); ok {
		servicePortMapping := tem.ServicePortMapping{}
		if v, ok := dMap["type"]; ok {
			servicePortMapping.Type = helper.String(v.(string))
		}
		if v, ok := dMap["service_name"]; ok {
			serviceName = v.(string)
			servicePortMapping.ServiceName = helper.String(v.(string))
		}
		if v, ok := dMap["port_mapping_item_list"]; ok {
			for _, item := range v.([]interface{}) {
				portMappingItemListMap := item.(map[string]interface{})
				servicePortMappingItem := tem.ServicePortMappingItem{}
				if v, ok := portMappingItemListMap["port"]; ok {
					servicePortMappingItem.Port = helper.IntInt64(v.(int))
				}
				if v, ok := portMappingItemListMap["target_port"]; ok {
					servicePortMappingItem.TargetPort = helper.IntInt64(v.(int))
				}
				if v, ok := portMappingItemListMap["protocol"]; ok {
					servicePortMappingItem.Protocol = helper.String(v.(string))
				}
				servicePortMapping.PortMappingItemList = append(servicePortMapping.PortMappingItemList, &servicePortMappingItem)
			}
		}
		request.Service = &servicePortMapping
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateApplicationService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem applicationService failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(environmentId + FILED_SP + applicationId + FILED_SP + serviceName)

	return resourceTencentCloudTemApplicationServiceRead(d, meta)
}

func resourceTencentCloudTemApplicationServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	serviceName := idSplit[2]

	res, err := service.DescribeTemApplicationServiceById(ctx, environmentId, applicationId)
	if err != nil {
		return err
	}

	if res == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("environment_id", environmentId)
	_ = d.Set("application_id", applicationId)

	var applicationService *tem.ServicePortMapping
	for _, v := range res.Result.ServicePortMappingList {
		if *v.ServiceName == serviceName {
			applicationService = v
		}
	}

	if applicationService != nil {
		serviceMap := map[string]interface{}{}

		if applicationService.Type != nil {
			serviceMap["type"] = applicationService.Type
		}

		if applicationService.ServiceName != nil {
			serviceMap["service_name"] = applicationService.ServiceName
		}

		if applicationService.PortMappingItemList != nil {
			portMappingItemListList := []interface{}{}
			for _, portMappingItemList := range applicationService.PortMappingItemList {
				portMappingItemListMap := map[string]interface{}{}

				if portMappingItemList.Port != nil {
					portMappingItemListMap["port"] = portMappingItemList.Port
				}

				if portMappingItemList.TargetPort != nil {
					portMappingItemListMap["target_port"] = portMappingItemList.TargetPort
				}

				if portMappingItemList.Protocol != nil {
					portMappingItemListMap["protocol"] = portMappingItemList.Protocol
				}

				portMappingItemListList = append(portMappingItemListList, portMappingItemListMap)
			}

			serviceMap["port_mapping_item_list"] = portMappingItemListList
		}

		err = d.Set("service", []interface{}{serviceMap})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudTemApplicationServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	unsupportedUpdateFields := []string{
		"service",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_tem_application_service update on %s is not support yet", field)
		}
	}

	request := tem.NewModifyApplicationServiceRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	serviceName := idSplit[2]

	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	if d.HasChange("service") {
		if dMap, ok := helper.InterfacesHeadMap(d, "service"); ok {
			servicePortMapping := tem.ServicePortMapping{}
			if v, ok := dMap["type"]; ok {
				servicePortMapping.Type = helper.String(v.(string))
			}

			servicePortMapping.ServiceName = &serviceName
			if v, ok := dMap["port_mapping_item_list"]; ok {
				for _, item := range v.([]interface{}) {
					portMappingItemListMap := item.(map[string]interface{})
					servicePortMappingItem := tem.ServicePortMappingItem{}
					if v, ok := portMappingItemListMap["port"]; ok {
						servicePortMappingItem.Port = helper.IntInt64(v.(int))
					}
					if v, ok := portMappingItemListMap["target_port"]; ok {
						servicePortMappingItem.TargetPort = helper.IntInt64(v.(int))
					}
					if v, ok := portMappingItemListMap["protocol"]; ok {
						servicePortMappingItem.Protocol = helper.String(v.(string))
					}
					servicePortMapping.PortMappingItemList = append(servicePortMapping.PortMappingItemList, &servicePortMappingItem)
				}
			}
			request.Data = &servicePortMapping
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyApplicationService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem applicationService failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTemApplicationServiceRead(d, meta)
}

func resourceTencentCloudTemApplicationServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	serviceName := idSplit[2]

	if err := service.DeleteTemApplicationServiceById(ctx, environmentId, applicationId, serviceName); err != nil {
		return err
	}

	return nil
}
