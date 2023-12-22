package tsf

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfApiDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApiDetailRead,
		Schema: map[string]*schema.Schema{
			"microservice_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "microservice id.",
			},

			"path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "api path.",
			},

			"method": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "request method.",
			},

			"pkg_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "pkg version.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "application id.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "api detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "api request description.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "param name.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type.",
									},
									"in": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "param position.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "param description.",
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "require or not.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "default value.",
									},
								},
							},
						},
						"response": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "api response.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "param description.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "param type.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "param description.",
									},
								},
							},
						},
						"definitions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "api data struct.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "object name.",
									},
									"properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "object property list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "property name.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "property type.",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "property description.",
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
							Description: "api content type.",
						},
						"can_run": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "can debug or not.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_api_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		microserviceId string
		path           string
		method         string
		pkgVersion     string
		applicationId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("microservice_id"); ok {
		microserviceId = v.(string)
		paramMap["MicroserviceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
		paramMap["Path"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("method"); ok {
		method = v.(string)
		paramMap["Method"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pkg_version"); ok {
		pkgVersion = v.(string)
		paramMap["PkgVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		applicationId = v.(string)
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var apiDetail *tsf.ApiDetailResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApiDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		apiDetail = result
		return nil
	})
	if err != nil {
		return err
	}

	apiDetailResponseMap := map[string]interface{}{}
	if apiDetail != nil {
		if apiDetail.Request != nil {
			requestList := []interface{}{}
			for _, request := range apiDetail.Request {
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

			apiDetailResponseMap["request"] = requestList
		}

		if apiDetail.Response != nil {
			responseList := []interface{}{}
			for _, response := range apiDetail.Response {
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

			apiDetailResponseMap["response"] = responseList
		}

		if apiDetail.Definitions != nil {
			definitionsList := []interface{}{}
			for _, definitions := range apiDetail.Definitions {
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

					definitionsMap["properties"] = propertiesList
				}

				definitionsList = append(definitionsList, definitionsMap)
			}

			apiDetailResponseMap["definitions"] = definitionsList
		}

		if apiDetail.RequestContentType != nil {
			apiDetailResponseMap["request_content_type"] = apiDetail.RequestContentType
		}

		if apiDetail.CanRun != nil {
			apiDetailResponseMap["can_run"] = apiDetail.CanRun
		}

		if apiDetail.Status != nil {
			apiDetailResponseMap["status"] = apiDetail.Status
		}

		if apiDetail.Description != nil {
			apiDetailResponseMap["description"] = apiDetail.Description
		}

		_ = d.Set("result", []interface{}{apiDetailResponseMap})
	}

	d.SetId(strings.Join([]string{microserviceId, path, method, pkgVersion, applicationId}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), apiDetailResponseMap); e != nil {
			return e
		}
	}
	return nil
}
