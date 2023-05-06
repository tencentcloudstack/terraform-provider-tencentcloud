/*
Use this data source to query detailed information of tsf application

Example Usage

```hcl
data "tencentcloud_tsf_application" "application" {
  application_type = "V"
  microservice_type = "N"
  # application_resource_type_list = [""]
  application_id_list = ["application-a24x29xv"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApplicationRead,
		Schema: map[string]*schema.Schema{
			"application_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The application type. V OR C, V means VM, C means container.",
			},

			"microservice_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The microservice type of the application.",
			},

			"application_resource_type_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "An array of application resource types.",
			},

			"application_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Id list.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The application paging list information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of applications.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of application information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the application.",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the application.",
									},
									"application_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the application.",
									},
									"application_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the application.",
									},
									"microservice_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The microservice type of the application.",
									},
									"prog_lang": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Programming language.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "create time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "update time.",
									},
									"application_resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "application resource type.",
									},
									"application_runtime_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "application runtime type.",
									},
									"apigateway_service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "gateway service id.",
									},
									"application_remark_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "remark name.",
									},
									"service_config_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "service config list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "serviceName.",
												},
												"ports": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "port list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"target_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "service port.",
															},
															"protocol": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "protocol.",
															},
														},
													},
												},
												"health_check": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "health check setting.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "health check path.",
															},
														},
													},
												},
											},
										},
									},
									"ignore_create_image_repository": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether ignore create image repository.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTsfApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_application.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("application_type"); ok {
		paramMap["ApplicationType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_type"); ok {
		paramMap["MicroserviceType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_resource_type_list"); ok {
		applicationResourceTypeListSet := v.(*schema.Set).List()
		paramMap["ApplicationResourceTypeList"] = helper.InterfacesStringsPoint(applicationResourceTypeListSet)
	}

	if v, ok := d.GetOk("application_id_list"); ok {
		applicationIdListSet := v.(*schema.Set).List()
		paramMap["ApplicationIdList"] = helper.InterfacesStringsPoint(applicationIdListSet)
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var application *tsf.TsfPageApplication
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApplicationByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		application = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(application.Content))
	tsfPageApplicationMap := map[string]interface{}{}
	if application != nil {
		if application.TotalCount != nil {
			tsfPageApplicationMap["total_count"] = application.TotalCount
		}

		if application.Content != nil {
			contentList := []interface{}{}
			for _, content := range application.Content {
				contentMap := map[string]interface{}{}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.ApplicationDesc != nil {
					contentMap["application_desc"] = content.ApplicationDesc
				}

				if content.ApplicationType != nil {
					contentMap["application_type"] = content.ApplicationType
				}

				if content.MicroserviceType != nil {
					contentMap["microservice_type"] = content.MicroserviceType
				}

				if content.ProgLang != nil {
					contentMap["prog_lang"] = content.ProgLang
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.ApplicationResourceType != nil {
					contentMap["application_resource_type"] = content.ApplicationResourceType
				}

				if content.ApplicationRuntimeType != nil {
					contentMap["application_runtime_type"] = content.ApplicationRuntimeType
				}

				if content.ApigatewayServiceId != nil {
					contentMap["apigateway_service_id"] = content.ApigatewayServiceId
				}

				if content.ApplicationRemarkName != nil {
					contentMap["application_remark_name"] = content.ApplicationRemarkName
				}

				if content.ServiceConfigList != nil {
					serviceConfigListList := []interface{}{}
					for _, serviceConfigList := range content.ServiceConfigList {
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

							serviceConfigListMap["ports"] = portsList
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

					contentMap["service_config_list"] = serviceConfigListList
				}

				if content.IgnoreCreateImageRepository != nil {
					contentMap["ignore_create_image_repository"] = content.IgnoreCreateImageRepository
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ApplicationId)
			}

			tsfPageApplicationMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageApplicationMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageApplicationMap); e != nil {
			return e
		}
	}
	return nil
}
