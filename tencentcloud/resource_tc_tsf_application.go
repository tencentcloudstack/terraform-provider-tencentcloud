/*
Provides a resource to create a tsf application

Example Usage

```hcl
resource "tencentcloud_tsf_application" "application" {
	application_name = "terraform-test"
	application_type = "V"
	microservice_type = "N"
	application_desc = "application desc"
	application_runtime_type = ""
	# application_remark_name = "remark-name"
	program_id = ""
  	service_config_list {
		name = "service-name"
		ports {
			target_port = 9090
			protocol = "TCP"
		}
		health_check {
			path = "/"
		}

  	}
	ignore_create_image_repository = true
	# program_id_list =

	tags = {
		"createdBy" = "terraform"
	}
}
```

Import

tsf application can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application.application application_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationCreate,
		Read:   resourceTencentCloudTsfApplicationRead,
		Update: resourceTencentCloudTsfApplicationUpdate,
		Delete: resourceTencentCloudTsfApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "application ID.",
			},

			"application_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application name.",
			},

			"application_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application type, `V`: virtual machine application; `C`: container application; `S`: serverless application.",
			},

			"microservice_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application microservice type, `M`: service mesh application; `N`: common application; `G`: gateway application.",
			},

			"application_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Application description.",
			},

			"application_runtime_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "application runtime type.",
			},

			"program_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The dataset ID to be bound.",
			},

			"service_config_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of service configuration information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Service Name.",
						},
						"ports": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Port information list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "service port.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "port protocol.",
									},
								},
							},
						},
						"health_check": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Health check configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "health check path.",
									},
								},
							},
						},
					},
				},
			},

			"ignore_create_image_repository": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Ignore creating a mirror repository.",
			},

			"prog_lang": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Programming language.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"apigateway_service_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "ServiceId of the Apigateway.",
			},

			"application_remark_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application note name.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tsf.NewCreateApplicationRequest()
		response      = tsf.NewCreateApplicationResponse()
		applicationId string
	)
	if v, ok := d.GetOk("application_name"); ok {
		request.ApplicationName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_type"); ok {
		request.ApplicationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_type"); ok {
		request.MicroserviceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_desc"); ok {
		request.ApplicationDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_runtime_type"); ok {
		request.ApplicationRuntimeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id"); ok {
		request.ProgramId = helper.String(v.(string))
	}

	if _, ok := d.GetOk("application_remark_name"); ok {
		return fmt.Errorf("There is no need to remark name when creating an application")
	}

	if v, ok := d.GetOk("service_config_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			serviceConfig := tsf.ServiceConfig{}
			if v, ok := dMap["name"]; ok {
				serviceConfig.Name = helper.String(v.(string))
			}
			if v, ok := dMap["ports"]; ok {
				for _, item := range v.([]interface{}) {
					portsMap := item.(map[string]interface{})
					ports := tsf.Ports{}
					if v, ok := portsMap["target_port"]; ok {
						ports.TargetPort = helper.IntUint64(v.(int))
					}
					if v, ok := portsMap["protocol"]; ok {
						ports.Protocol = helper.String(v.(string))
					}
					serviceConfig.Ports = append(serviceConfig.Ports, &ports)
				}
			}
			if healthCheckMap, ok := helper.InterfaceToMap(dMap, "health_check"); ok {
				healthCheckConfig := tsf.HealthCheckConfig{}
				if v, ok := healthCheckMap["path"]; ok {
					healthCheckConfig.Path = helper.String(v.(string))
				}
				serviceConfig.HealthCheck = &healthCheckConfig
			}
			request.ServiceConfigList = append(request.ServiceConfigList, &serviceConfig)
		}
	}

	if v, _ := d.GetOk("ignore_create_image_repository"); v != nil {
		request.IgnoreCreateImageRepository = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateApplication(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf application failed, reason:%+v", logId, err)
		return err
	}

	applicationId = *response.Response.Result
	d.SetId(applicationId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:application/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfApplicationRead(d, meta)
}

func resourceTencentCloudTsfApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationId := d.Id()

	application, err := service.DescribeTsfApplicationById(ctx, applicationId)
	if err != nil {
		return err
	}

	if application == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplication` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if application.ApplicationId != nil {
		_ = d.Set("application_id", application.ApplicationId)
	}

	if application.ApplicationName != nil {
		_ = d.Set("application_name", application.ApplicationName)
	}

	if application.ApplicationType != nil {
		_ = d.Set("application_type", application.ApplicationType)
	}

	if application.MicroserviceType != nil {
		_ = d.Set("microservice_type", application.MicroserviceType)
	}

	if application.ApplicationDesc != nil {
		_ = d.Set("application_desc", application.ApplicationDesc)
	}

	if application.ApplicationRuntimeType != nil {
		_ = d.Set("application_runtime_type", application.ApplicationRuntimeType)
	}

	// if application.ProgramId != nil {
	// 	_ = d.Set("program_id", application.ProgramId)
	// }

	if application.ServiceConfigList != nil {
		serviceConfigListList := []interface{}{}
		for _, serviceConfigList := range application.ServiceConfigList {
			serviceConfigListMap := map[string]interface{}{}

			if serviceConfigList.Name != nil {
				serviceConfigListMap["name"] = serviceConfigList.Name
			}

			if serviceConfigList.Ports != nil {
				portsList := []interface{}{}
				for _, ports := range serviceConfigList.Ports {
					portsMap := map[string]interface{}{}

					if ports.TargetPort != nil {
						portsMap["target_port"] = ports.TargetPort
					}

					if ports.Protocol != nil {
						portsMap["protocol"] = ports.Protocol
					}

					portsList = append(portsList, portsMap)
				}

				serviceConfigListMap["ports"] = []interface{}{portsList}
			}

			if serviceConfigList.HealthCheck != nil {
				healthCheckMap := map[string]interface{}{}

				if serviceConfigList.HealthCheck.Path != nil {
					healthCheckMap["path"] = serviceConfigList.HealthCheck.Path
				}

				serviceConfigListMap["health_check"] = []interface{}{healthCheckMap}
			}

			serviceConfigListList = append(serviceConfigListList, serviceConfigListMap)
		}

		_ = d.Set("service_config_list", serviceConfigListList)

	}

	// if application.IgnoreCreateImageRepository != nil {
	// 	_ = d.Set("ignore_create_image_repository", application.IgnoreCreateImageRepository)
	// }

	if application.ProgLang != nil {
		_ = d.Set("prog_lang", application.ProgLang)
	}

	if application.CreateTime != nil {
		_ = d.Set("create_time", application.CreateTime)
	}

	if application.UpdateTime != nil {
		_ = d.Set("update_time", application.UpdateTime)
	}

	if application.ApigatewayServiceId != nil {
		_ = d.Set("apigateway_service_id", application.ApigatewayServiceId)
	}

	if application.ApplicationRemarkName != nil {
		_ = d.Set("application_remark_name", application.ApplicationRemarkName)
	}

	// if application.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", application.ProgramIdList)
	// }

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "application", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyApplicationRequest()

	applicationId := d.Id()

	request.ApplicationId = &applicationId

	immutableArgs := []string{"application_id", "application_name", "application_type", "microservice_type", "application_runtime_type", "program_id", "ignore_create_image_repository", "prog_lang", "create_time", "update_time", "apigateway_service_id", "program_id_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("application_desc") {
		if v, ok := d.GetOk("application_desc"); ok {
			request.ApplicationDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("application_remark_name") {
		if v, ok := d.GetOk("application_remark_name"); ok {
			request.ApplicationRemarkName = helper.String(v.(string))
		}
	}

	if d.HasChange("service_config_list") {
		if v, ok := d.GetOk("service_config_list"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				serviceConfig := tsf.ServiceConfig{}
				if v, ok := dMap["name"]; ok {
					serviceConfig.Name = helper.String(v.(string))
				}
				if v, ok := dMap["ports"]; ok {
					for _, item := range v.([]interface{}) {
						portsMap := item.(map[string]interface{})
						ports := tsf.Ports{}
						if v, ok := portsMap["target_port"]; ok {
							ports.TargetPort = helper.IntUint64(v.(int))
						}
						if v, ok := portsMap["protocol"]; ok {
							ports.Protocol = helper.String(v.(string))
						}
						serviceConfig.Ports = append(serviceConfig.Ports, &ports)
					}
				}
				if healthCheckMap, ok := helper.InterfaceToMap(dMap, "health_check"); ok {
					healthCheckConfig := tsf.HealthCheckConfig{}
					if v, ok := healthCheckMap["path"]; ok {
						healthCheckConfig.Path = helper.String(v.(string))
					}
					serviceConfig.HealthCheck = &healthCheckConfig
				}
				request.ServiceConfigList = append(request.ServiceConfigList, &serviceConfig)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyApplication(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf application failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "application", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfApplicationRead(d, meta)
}

func resourceTencentCloudTsfApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	applicationId := d.Id()

	if err := service.DeleteTsfApplicationById(ctx, applicationId); err != nil {
		return err
	}

	return nil
}
