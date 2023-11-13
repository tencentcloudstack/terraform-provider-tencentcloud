/*
Use this data source to query detailed information of tsf api_detail

Example Usage

```hcl
data "tencentcloud_tsf_api_detail" "api_detail" {
  microservice_id = "ms-yq3jo6jd"
  path = "/printRequest"
  method = "GET"
  pkg_version = "20210625192923"
  application_id = "application-a24x29xv"
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

func dataSourceTencentCloudTsfApiDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApiDetailRead,
		Schema: map[string]*schema.Schema{
			"microservice_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Microservice id.",
			},

			"path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Api path.",
			},

			"method": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Request method.",
			},

			"pkg_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Pkg version.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application id.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Api detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Api request description.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Param name.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type.",
									},
									"in": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Param position.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Param description.",
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Require or not.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value.",
									},
								},
							},
						},
						"response": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Api response.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Param description.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Param type.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Param description.",
									},
								},
							},
						},
						"definitions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Api data struct.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Object name.",
									},
									"properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Object property list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property name.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property type.",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property description.",
												},
											},
										},
									},
								},
							},
						},
						"request_content_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Api content type.",
						},
						"can_run": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Can debug or not.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "API status 0: offline 1: online, default 0. Note: This section may return null, indicating that no valid value can be obtained.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API description. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTsfApiDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_api_detail.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("microservice_id"); ok {
		paramMap["MicroserviceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		paramMap["Path"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("method"); ok {
		paramMap["Method"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pkg_version"); ok {
		paramMap["PkgVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tsf.ApiDetailResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApiDetailByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		apiDetailResponseMap := map[string]interface{}{}

		if result.Request != nil {
			requestList := []interface{}{}
			for _, request := range result.Request {
				requestMap := map[string]interface{}{}

				if request.Name != nil {
					requestMap["name"] = request.Name
				}

				if request.Type != nil {
					requestMap["type"] = request.Type
				}

				if request.In != nil {
					requestMap["in"] = request.In
				}

				if request.Description != nil {
					requestMap["description"] = request.Description
				}

				if request.Required != nil {
					requestMap["required"] = request.Required
				}

				if request.DefaultValue != nil {
					requestMap["default_value"] = request.DefaultValue
				}

				requestList = append(requestList, requestMap)
			}

			apiDetailResponseMap["request"] = []interface{}{requestList}
		}

		if result.Response != nil {
			responseList := []interface{}{}
			for _, response := range result.Response {
				responseMap := map[string]interface{}{}

				if response.Name != nil {
					responseMap["name"] = response.Name
				}

				if response.Type != nil {
					responseMap["type"] = response.Type
				}

				if response.Description != nil {
					responseMap["description"] = response.Description
				}

				responseList = append(responseList, responseMap)
			}

			apiDetailResponseMap["response"] = []interface{}{responseList}
		}

		if result.Definitions != nil {
			definitionsList := []interface{}{}
			for _, definitions := range result.Definitions {
				definitionsMap := map[string]interface{}{}

				if definitions.Name != nil {
					definitionsMap["name"] = definitions.Name
				}

				if definitions.Properties != nil {
					propertiesList := []interface{}{}
					for _, properties := range definitions.Properties {
						propertiesMap := map[string]interface{}{}

						if properties.Name != nil {
							propertiesMap["name"] = properties.Name
						}

						if properties.Type != nil {
							propertiesMap["type"] = properties.Type
						}

						if properties.Description != nil {
							propertiesMap["description"] = properties.Description
						}

						propertiesList = append(propertiesList, propertiesMap)
					}

					definitionsMap["properties"] = []interface{}{propertiesList}
				}

				definitionsList = append(definitionsList, definitionsMap)
			}

			apiDetailResponseMap["definitions"] = []interface{}{definitionsList}
		}

		if result.RequestContentType != nil {
			apiDetailResponseMap["request_content_type"] = result.RequestContentType
		}

		if result.CanRun != nil {
			apiDetailResponseMap["can_run"] = result.CanRun
		}

		if result.Status != nil {
			apiDetailResponseMap["status"] = result.Status
		}

		if result.Description != nil {
			apiDetailResponseMap["description"] = result.Description
		}

		ids = append(ids, *result.MicroserviceId)
		_ = d.Set("result", apiDetailResponseMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), apiDetailResponseMap); e != nil {
			return e
		}
	}
	return nil
}
