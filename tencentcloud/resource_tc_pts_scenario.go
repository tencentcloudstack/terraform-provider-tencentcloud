/*
Provides a resource to create a pts scenario

Example Usage

```hcl
resource "tencentcloud_pts_scenario" "scenario" {
  name = "pts"
  type = &lt;nil&gt;
  project_id = &lt;nil&gt;
  description = &lt;nil&gt;
  load {
		load_spec {
			concurrency {
				stages {
					duration_seconds = &lt;nil&gt;
					target_virtual_users = &lt;nil&gt;
				}
				iteration_count = &lt;nil&gt;
				max_requests_per_second = &lt;nil&gt;
				graceful_stop_seconds = &lt;nil&gt;
			}
			requests_per_second {
				max_requests_per_second = &lt;nil&gt;
				duration_seconds = &lt;nil&gt;
				resources = &lt;nil&gt;
				start_requests_per_second = &lt;nil&gt;
				target_requests_per_second = &lt;nil&gt;
				graceful_stop_seconds = &lt;nil&gt;
			}
			script_origin {
				machine_number = &lt;nil&gt;
				machine_specification = &lt;nil&gt;
				duration_seconds = &lt;nil&gt;
			}
		}
		vpc_load_distribution {
			region_id = &lt;nil&gt;
			region = &lt;nil&gt;
			vpc_id = &lt;nil&gt;
			subnet_ids = &lt;nil&gt;
		}
		geo_regions_load_distribution {
			region_id = &lt;nil&gt;
			region = &lt;nil&gt;
			percentage = &lt;nil&gt;
		}

  }
  datasets {
		name = &lt;nil&gt;
		split = &lt;nil&gt;
		header_in_file = &lt;nil&gt;
		header_columns = &lt;nil&gt;
		line_count = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		size = &lt;nil&gt;
		head_lines = &lt;nil&gt;
		tail_lines = &lt;nil&gt;
		type = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  cron_id = &lt;nil&gt;
  test_scripts {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		encoded_content = &lt;nil&gt;
		encoded_http_archive = &lt;nil&gt;
		load_weight = &lt;nil&gt;

  }
  protocols {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  request_files {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  s_l_a_policy {
		s_l_a_rules {
			metric = &lt;nil&gt;
			aggregation = &lt;nil&gt;
			condition = &lt;nil&gt;
			value = &lt;nil&gt;
			label_filter {
				label_name = &lt;nil&gt;
				label_value = &lt;nil&gt;
			}
			abort_flag = &lt;nil&gt;
			for = &lt;nil&gt;
		}
		alert_channel {
			notice_id = &lt;nil&gt;
			a_m_p_consumer_id = &lt;nil&gt;
		}

  }
  plugins {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  domain_name_config {
		host_aliases {
			host_names = &lt;nil&gt;
			i_p = &lt;nil&gt;
		}
		d_n_s_config {
			nameservers = &lt;nil&gt;
		}

  }
              }
```

Import

pts scenario can be imported using the id, e.g.

```
terraform import tencentcloud_pts_scenario.scenario scenario_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudPtsScenario() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsScenarioCreate,
		Read:   resourceTencentCloudPtsScenarioRead,
		Update: resourceTencentCloudPtsScenarioUpdate,
		Delete: resourceTencentCloudPtsScenarioDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Pts Scenario name.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Pressure test engine type.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project id.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Pts Scenario Description.",
			},

			"load": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Pressure allocation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_spec": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Pressure allocation.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"concurrency": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Configuration of concurrent pressure mode.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"stages": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Multi-phase configuration array.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"duration_seconds": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Pressure time.",
															},
															"target_virtual_users": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Number of virtual users.",
															},
														},
													},
												},
												"iteration_count": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Number of runs.",
												},
												"max_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Maximum RPS.",
												},
												"graceful_stop_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Wait time for graceful termination of the task.",
												},
											},
										},
									},
									"requests_per_second": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Configuration of RPS pressure mode.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Maximum RPS.",
												},
												"duration_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Pressure time.",
												},
												"resources": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Number of resources.",
												},
												"start_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Initial RPS.",
												},
												"target_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Target RPS, invalid input parameter.",
												},
												"graceful_stop_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Elegant shutdown waiting time.",
												},
											},
										},
									},
									"script_origin": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Built-in stress mode in script.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"machine_number": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Number of machines.",
												},
												"machine_specification": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Machine specification.",
												},
												"duration_seconds": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Pressure testing time.",
												},
											},
										},
									},
								},
							},
						},
						"vpc_load_distribution": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Source of stress.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Regional ID.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPC ID.",
									},
									"subnet_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Subnet ID list.",
									},
								},
							},
						},
						"geo_regions_load_distribution": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Pressure distribution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Regional ID.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region.",
									},
									"percentage": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Percentage.",
									},
								},
							},
						},
					},
				},
			},

			"datasets": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Test data set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The file name where the test dataset is located.",
						},
						"split": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Test whether the dataset is fragmented.",
						},
						"header_in_file": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the first line is the parameter name.",
						},
						"header_columns": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Parameter name array.",
						},
						"line_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of file lines.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of file bytes.",
						},
						"head_lines": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Header data row.",
						},
						"tail_lines": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Trailing data row.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File ID.",
						},
					},
				},
			},

			"cron_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cron job ID.",
			},

			"test_scripts": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Test script file information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File size.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"encoded_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base64 encoded file content.",
						},
						"encoded_http_archive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base64 encoded har structure.",
						},
						"load_weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Script weight, range 1-100.",
						},
					},
				},
			},

			"protocols": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Protocol file path.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File name.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File ID.",
						},
					},
				},
			},

			"request_files": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Request file path.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File size.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File id.",
						},
					},
				},
			},

			"s_l_a_policy": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SLA strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s_l_a_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "SLA rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Pressure test index.",
									},
									"aggregation": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Aggregation method of pressure test index.",
									},
									"condition": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Pressure test index condition judgment symbol.",
									},
									"value": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold value.",
									},
									"label_filter": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Tag.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Label name.",
												},
												"label_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Label value.",
												},
											},
										},
									},
									"abort_flag": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to stop the stress test task.",
									},
									"for": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Duraion.",
									},
								},
							},
						},
						"alert_channel": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Alarm notification channel.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notice_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Notification template ID.",
									},
									"a_m_p_consumer_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AMP consumer ID.",
									},
								},
							},
						},
					},
				},
			},

			"plugins": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "SLA strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File size.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File id.",
						},
					},
				},
			},

			"domain_name_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Domain name resolution configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_aliases": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Domain name binding configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_names": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "List of domain names to be bound.",
									},
									"i_p": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The IP address to be bound.",
									},
								},
							},
						},
						"d_n_s_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "DNS configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nameservers": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "DNS IP List.",
									},
								},
							},
						},
					},
				},
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Scene statu.",
			},

			"created_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scene creation time.",
			},

			"updated_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scene modification time.",
			},

			"app_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "App ID.",
			},

			"uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "User ID.",
			},

			"sub_account_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sub-user ID.",
			},

			"notification_hooks": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Notification event callback.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"events": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Notification event.",
						},
						"u_r_l": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Callback URL.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPtsScenarioCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_scenario.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = pts.NewCreateScenarioRequest()
		response   = pts.NewCreateScenarioResponse()
		scenarioId string
		projectId  string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "load"); ok {
		load := pts.Load{}
		if loadSpecMap, ok := helper.InterfaceToMap(dMap, "load_spec"); ok {
			loadSpec := pts.LoadSpec{}
			if concurrencyMap, ok := helper.InterfaceToMap(loadSpecMap, "concurrency"); ok {
				concurrency := pts.Concurrency{}
				if v, ok := concurrencyMap["stages"]; ok {
					for _, item := range v.([]interface{}) {
						stagesMap := item.(map[string]interface{})
						stage := pts.Stage{}
						if v, ok := stagesMap["duration_seconds"]; ok {
							stage.DurationSeconds = helper.IntUint64(v.(int))
						}
						if v, ok := stagesMap["target_virtual_users"]; ok {
							stage.TargetVirtualUsers = helper.IntUint64(v.(int))
						}
						concurrency.Stages = append(concurrency.Stages, &stage)
					}
				}
				if v, ok := concurrencyMap["iteration_count"]; ok {
					concurrency.IterationCount = helper.IntUint64(v.(int))
				}
				if v, ok := concurrencyMap["max_requests_per_second"]; ok {
					concurrency.MaxRequestsPerSecond = helper.IntUint64(v.(int))
				}
				if v, ok := concurrencyMap["graceful_stop_seconds"]; ok {
					concurrency.GracefulStopSeconds = helper.IntUint64(v.(int))
				}
				loadSpec.Concurrency = &concurrency
			}
			if requestsPerSecondMap, ok := helper.InterfaceToMap(loadSpecMap, "requests_per_second"); ok {
				requestsPerSecond := pts.RequestsPerSecond{}
				if v, ok := requestsPerSecondMap["max_requests_per_second"]; ok {
					requestsPerSecond.MaxRequestsPerSecond = helper.IntUint64(v.(int))
				}
				if v, ok := requestsPerSecondMap["duration_seconds"]; ok {
					requestsPerSecond.DurationSeconds = helper.IntUint64(v.(int))
				}
				if v, ok := requestsPerSecondMap["resources"]; ok {
					requestsPerSecond.Resources = helper.IntUint64(v.(int))
				}
				if v, ok := requestsPerSecondMap["start_requests_per_second"]; ok {
					requestsPerSecond.StartRequestsPerSecond = helper.IntUint64(v.(int))
				}
				if v, ok := requestsPerSecondMap["target_requests_per_second"]; ok {
					requestsPerSecond.TargetRequestsPerSecond = helper.IntUint64(v.(int))
				}
				if v, ok := requestsPerSecondMap["graceful_stop_seconds"]; ok {
					requestsPerSecond.GracefulStopSeconds = helper.IntUint64(v.(int))
				}
				loadSpec.RequestsPerSecond = &requestsPerSecond
			}
			if scriptOriginMap, ok := helper.InterfaceToMap(loadSpecMap, "script_origin"); ok {
				scriptOrigin := pts.ScriptOrigin{}
				if v, ok := scriptOriginMap["machine_number"]; ok {
					scriptOrigin.MachineNumber = helper.IntUint64(v.(int))
				}
				if v, ok := scriptOriginMap["machine_specification"]; ok {
					scriptOrigin.MachineSpecification = helper.String(v.(string))
				}
				if v, ok := scriptOriginMap["duration_seconds"]; ok {
					scriptOrigin.DurationSeconds = helper.IntUint64(v.(int))
				}
				loadSpec.ScriptOrigin = &scriptOrigin
			}
			load.LoadSpec = &loadSpec
		}
		if vpcLoadDistributionMap, ok := helper.InterfaceToMap(dMap, "vpc_load_distribution"); ok {
			vpcLoadDistribution := pts.VpcLoadDistribution{}
			if v, ok := vpcLoadDistributionMap["region_id"]; ok {
				vpcLoadDistribution.RegionId = helper.IntUint64(v.(int))
			}
			if v, ok := vpcLoadDistributionMap["region"]; ok {
				vpcLoadDistribution.Region = helper.String(v.(string))
			}
			if v, ok := vpcLoadDistributionMap["vpc_id"]; ok {
				vpcLoadDistribution.VpcId = helper.String(v.(string))
			}
			if v, ok := vpcLoadDistributionMap["subnet_ids"]; ok {
				subnetIdsSet := v.(*schema.Set).List()
				for i := range subnetIdsSet {
					subnetIds := subnetIdsSet[i].(string)
					vpcLoadDistribution.SubnetIds = append(vpcLoadDistribution.SubnetIds, &subnetIds)
				}
			}
			load.VpcLoadDistribution = &vpcLoadDistribution
		}
		if v, ok := dMap["geo_regions_load_distribution"]; ok {
			for _, item := range v.([]interface{}) {
				geoRegionsLoadDistributionMap := item.(map[string]interface{})
				geoRegionsLoadItem := pts.GeoRegionsLoadItem{}
				if v, ok := geoRegionsLoadDistributionMap["region_id"]; ok {
					geoRegionsLoadItem.RegionId = helper.IntUint64(v.(int))
				}
				if v, ok := geoRegionsLoadDistributionMap["region"]; ok {
					geoRegionsLoadItem.Region = helper.String(v.(string))
				}
				if v, ok := geoRegionsLoadDistributionMap["percentage"]; ok {
					geoRegionsLoadItem.Percentage = helper.IntUint64(v.(int))
				}
				load.GeoRegionsLoadDistribution = append(load.GeoRegionsLoadDistribution, &geoRegionsLoadItem)
			}
		}
		request.Load = &load
	}

	if v, ok := d.GetOk("datasets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			testData := pts.TestData{}
			if v, ok := dMap["name"]; ok {
				testData.Name = helper.String(v.(string))
			}
			if v, ok := dMap["split"]; ok {
				testData.Split = helper.Bool(v.(bool))
			}
			if v, ok := dMap["header_in_file"]; ok {
				testData.HeaderInFile = helper.Bool(v.(bool))
			}
			if v, ok := dMap["header_columns"]; ok {
				headerColumnsSet := v.(*schema.Set).List()
				for i := range headerColumnsSet {
					headerColumns := headerColumnsSet[i].(string)
					testData.HeaderColumns = append(testData.HeaderColumns, &headerColumns)
				}
			}
			if v, ok := dMap["line_count"]; ok {
				testData.LineCount = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["updated_at"]; ok {
				testData.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				testData.Size = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["head_lines"]; ok {
				headLinesSet := v.(*schema.Set).List()
				for i := range headLinesSet {
					headLines := headLinesSet[i].(string)
					testData.HeadLines = append(testData.HeadLines, &headLines)
				}
			}
			if v, ok := dMap["tail_lines"]; ok {
				tailLinesSet := v.(*schema.Set).List()
				for i := range tailLinesSet {
					tailLines := tailLinesSet[i].(string)
					testData.TailLines = append(testData.TailLines, &tailLines)
				}
			}
			if v, ok := dMap["type"]; ok {
				testData.Type = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				testData.FileId = helper.String(v.(string))
			}
			request.Datasets = append(request.Datasets, &testData)
		}
	}

	if v, ok := d.GetOk("cron_id"); ok {
		request.CronId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("test_scripts"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			scriptInfo := pts.ScriptInfo{}
			if v, ok := dMap["name"]; ok {
				scriptInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				scriptInfo.Size = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["type"]; ok {
				scriptInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				scriptInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["encoded_content"]; ok {
				scriptInfo.EncodedContent = helper.String(v.(string))
			}
			if v, ok := dMap["encoded_http_archive"]; ok {
				scriptInfo.EncodedHttpArchive = helper.String(v.(string))
			}
			if v, ok := dMap["load_weight"]; ok {
				scriptInfo.LoadWeight = helper.IntUint64(v.(int))
			}
			request.TestScripts = append(request.TestScripts, &scriptInfo)
		}
	}

	if v, ok := d.GetOk("protocols"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			protocolInfo := pts.ProtocolInfo{}
			if v, ok := dMap["name"]; ok {
				protocolInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				protocolInfo.Size = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["type"]; ok {
				protocolInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				protocolInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				protocolInfo.FileId = helper.String(v.(string))
			}
			request.Protocols = append(request.Protocols, &protocolInfo)
		}
	}

	if v, ok := d.GetOk("request_files"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			fileInfo := pts.FileInfo{}
			if v, ok := dMap["name"]; ok {
				fileInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				fileInfo.Size = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["type"]; ok {
				fileInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				fileInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				fileInfo.FileId = helper.String(v.(string))
			}
			request.RequestFiles = append(request.RequestFiles, &fileInfo)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "s_l_a_policy"); ok {
		sLAPolicy := pts.SLAPolicy{}
		if v, ok := dMap["s_l_a_rules"]; ok {
			for _, item := range v.([]interface{}) {
				sLARulesMap := item.(map[string]interface{})
				sLARule := pts.SLARule{}
				if v, ok := sLARulesMap["metric"]; ok {
					sLARule.Metric = helper.String(v.(string))
				}
				if v, ok := sLARulesMap["aggregation"]; ok {
					sLARule.Aggregation = helper.String(v.(string))
				}
				if v, ok := sLARulesMap["condition"]; ok {
					sLARule.Condition = helper.String(v.(string))
				}
				if v, ok := sLARulesMap["value"]; ok {
					sLARule.Value = helper.IntInt64(v.(int))
				}
				if v, ok := sLARulesMap["label_filter"]; ok {
					for _, item := range v.([]interface{}) {
						labelFilterMap := item.(map[string]interface{})
						sLALabel := pts.SLALabel{}
						if v, ok := labelFilterMap["label_name"]; ok {
							sLALabel.LabelName = helper.String(v.(string))
						}
						if v, ok := labelFilterMap["label_value"]; ok {
							sLALabel.LabelValue = helper.String(v.(string))
						}
						sLARule.LabelFilter = append(sLARule.LabelFilter, &sLALabel)
					}
				}
				if v, ok := sLARulesMap["abort_flag"]; ok {
					sLARule.AbortFlag = helper.Bool(v.(bool))
				}
				if v, ok := sLARulesMap["for"]; ok {
					sLARule.For = helper.String(v.(string))
				}
				sLAPolicy.SLARules = append(sLAPolicy.SLARules, &sLARule)
			}
		}
		if alertChannelMap, ok := helper.InterfaceToMap(dMap, "alert_channel"); ok {
			alertChannel := pts.AlertChannel{}
			if v, ok := alertChannelMap["notice_id"]; ok {
				alertChannel.NoticeId = helper.String(v.(string))
			}
			if v, ok := alertChannelMap["a_m_p_consumer_id"]; ok {
				alertChannel.AMPConsumerId = helper.String(v.(string))
			}
			sLAPolicy.AlertChannel = &alertChannel
		}
		request.SLAPolicy = &sLAPolicy
	}

	if v, ok := d.GetOk("plugins"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			fileInfo := pts.FileInfo{}
			if v, ok := dMap["name"]; ok {
				fileInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				fileInfo.Size = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["type"]; ok {
				fileInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				fileInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				fileInfo.FileId = helper.String(v.(string))
			}
			request.Plugins = append(request.Plugins, &fileInfo)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "domain_name_config"); ok {
		domainNameConfig := pts.DomainNameConfig{}
		if v, ok := dMap["host_aliases"]; ok {
			for _, item := range v.([]interface{}) {
				hostAliasesMap := item.(map[string]interface{})
				hostAlias := pts.HostAlias{}
				if v, ok := hostAliasesMap["host_names"]; ok {
					hostNamesSet := v.(*schema.Set).List()
					for i := range hostNamesSet {
						hostNames := hostNamesSet[i].(string)
						hostAlias.HostNames = append(hostAlias.HostNames, &hostNames)
					}
				}
				if v, ok := hostAliasesMap["i_p"]; ok {
					hostAlias.IP = helper.String(v.(string))
				}
				domainNameConfig.HostAliases = append(domainNameConfig.HostAliases, &hostAlias)
			}
		}
		if dNSConfigMap, ok := helper.InterfaceToMap(dMap, "d_n_s_config"); ok {
			dNSConfig := pts.DNSConfig{}
			if v, ok := dNSConfigMap["nameservers"]; ok {
				nameserversSet := v.(*schema.Set).List()
				for i := range nameserversSet {
					nameservers := nameserversSet[i].(string)
					dNSConfig.Nameservers = append(dNSConfig.Nameservers, &nameservers)
				}
			}
			domainNameConfig.DNSConfig = &dNSConfig
		}
		request.DomainNameConfig = &domainNameConfig
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().CreateScenario(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create pts scenario failed, reason:%+v", logId, err)
		return err
	}

	scenarioId = *response.Response.ScenarioId
	d.SetId(strings.Join([]string{scenarioId, projectId}, FILED_SP))

	return resourceTencentCloudPtsScenarioRead(d, meta)
}

func resourceTencentCloudPtsScenarioRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_scenario.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	scenarioId := idSplit[0]
	projectId := idSplit[1]

	scenario, err := service.DescribePtsScenarioById(ctx, scenarioId, projectId)
	if err != nil {
		return err
	}

	if scenario == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PtsScenario` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if scenario.Name != nil {
		_ = d.Set("name", scenario.Name)
	}

	if scenario.Type != nil {
		_ = d.Set("type", scenario.Type)
	}

	if scenario.ProjectId != nil {
		_ = d.Set("project_id", scenario.ProjectId)
	}

	if scenario.Description != nil {
		_ = d.Set("description", scenario.Description)
	}

	if scenario.Load != nil {
		loadMap := map[string]interface{}{}

		if scenario.Load.LoadSpec != nil {
			loadSpecMap := map[string]interface{}{}

			if scenario.Load.LoadSpec.Concurrency != nil {
				concurrencyMap := map[string]interface{}{}

				if scenario.Load.LoadSpec.Concurrency.Stages != nil {
					stagesList := []interface{}{}
					for _, stages := range scenario.Load.LoadSpec.Concurrency.Stages {
						stagesMap := map[string]interface{}{}

						if stages.DurationSeconds != nil {
							stagesMap["duration_seconds"] = stages.DurationSeconds
						}

						if stages.TargetVirtualUsers != nil {
							stagesMap["target_virtual_users"] = stages.TargetVirtualUsers
						}

						stagesList = append(stagesList, stagesMap)
					}

					concurrencyMap["stages"] = []interface{}{stagesList}
				}

				if scenario.Load.LoadSpec.Concurrency.IterationCount != nil {
					concurrencyMap["iteration_count"] = scenario.Load.LoadSpec.Concurrency.IterationCount
				}

				if scenario.Load.LoadSpec.Concurrency.MaxRequestsPerSecond != nil {
					concurrencyMap["max_requests_per_second"] = scenario.Load.LoadSpec.Concurrency.MaxRequestsPerSecond
				}

				if scenario.Load.LoadSpec.Concurrency.GracefulStopSeconds != nil {
					concurrencyMap["graceful_stop_seconds"] = scenario.Load.LoadSpec.Concurrency.GracefulStopSeconds
				}

				loadSpecMap["concurrency"] = []interface{}{concurrencyMap}
			}

			if scenario.Load.LoadSpec.RequestsPerSecond != nil {
				requestsPerSecondMap := map[string]interface{}{}

				if scenario.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond != nil {
					requestsPerSecondMap["max_requests_per_second"] = scenario.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond
				}

				if scenario.Load.LoadSpec.RequestsPerSecond.DurationSeconds != nil {
					requestsPerSecondMap["duration_seconds"] = scenario.Load.LoadSpec.RequestsPerSecond.DurationSeconds
				}

				if scenario.Load.LoadSpec.RequestsPerSecond.Resources != nil {
					requestsPerSecondMap["resources"] = scenario.Load.LoadSpec.RequestsPerSecond.Resources
				}

				if scenario.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond != nil {
					requestsPerSecondMap["start_requests_per_second"] = scenario.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond
				}

				if scenario.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond != nil {
					requestsPerSecondMap["target_requests_per_second"] = scenario.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond
				}

				if scenario.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds != nil {
					requestsPerSecondMap["graceful_stop_seconds"] = scenario.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds
				}

				loadSpecMap["requests_per_second"] = []interface{}{requestsPerSecondMap}
			}

			if scenario.Load.LoadSpec.ScriptOrigin != nil {
				scriptOriginMap := map[string]interface{}{}

				if scenario.Load.LoadSpec.ScriptOrigin.MachineNumber != nil {
					scriptOriginMap["machine_number"] = scenario.Load.LoadSpec.ScriptOrigin.MachineNumber
				}

				if scenario.Load.LoadSpec.ScriptOrigin.MachineSpecification != nil {
					scriptOriginMap["machine_specification"] = scenario.Load.LoadSpec.ScriptOrigin.MachineSpecification
				}

				if scenario.Load.LoadSpec.ScriptOrigin.DurationSeconds != nil {
					scriptOriginMap["duration_seconds"] = scenario.Load.LoadSpec.ScriptOrigin.DurationSeconds
				}

				loadSpecMap["script_origin"] = []interface{}{scriptOriginMap}
			}

			loadMap["load_spec"] = []interface{}{loadSpecMap}
		}

		if scenario.Load.VpcLoadDistribution != nil {
			vpcLoadDistributionMap := map[string]interface{}{}

			if scenario.Load.VpcLoadDistribution.RegionId != nil {
				vpcLoadDistributionMap["region_id"] = scenario.Load.VpcLoadDistribution.RegionId
			}

			if scenario.Load.VpcLoadDistribution.Region != nil {
				vpcLoadDistributionMap["region"] = scenario.Load.VpcLoadDistribution.Region
			}

			if scenario.Load.VpcLoadDistribution.VpcId != nil {
				vpcLoadDistributionMap["vpc_id"] = scenario.Load.VpcLoadDistribution.VpcId
			}

			if scenario.Load.VpcLoadDistribution.SubnetIds != nil {
				vpcLoadDistributionMap["subnet_ids"] = scenario.Load.VpcLoadDistribution.SubnetIds
			}

			loadMap["vpc_load_distribution"] = []interface{}{vpcLoadDistributionMap}
		}

		if scenario.Load.GeoRegionsLoadDistribution != nil {
			geoRegionsLoadDistributionList := []interface{}{}
			for _, geoRegionsLoadDistribution := range scenario.Load.GeoRegionsLoadDistribution {
				geoRegionsLoadDistributionMap := map[string]interface{}{}

				if geoRegionsLoadDistribution.RegionId != nil {
					geoRegionsLoadDistributionMap["region_id"] = geoRegionsLoadDistribution.RegionId
				}

				if geoRegionsLoadDistribution.Region != nil {
					geoRegionsLoadDistributionMap["region"] = geoRegionsLoadDistribution.Region
				}

				if geoRegionsLoadDistribution.Percentage != nil {
					geoRegionsLoadDistributionMap["percentage"] = geoRegionsLoadDistribution.Percentage
				}

				geoRegionsLoadDistributionList = append(geoRegionsLoadDistributionList, geoRegionsLoadDistributionMap)
			}

			loadMap["geo_regions_load_distribution"] = []interface{}{geoRegionsLoadDistributionList}
		}

		_ = d.Set("load", []interface{}{loadMap})
	}

	if scenario.Datasets != nil {
		datasetsList := []interface{}{}
		for _, datasets := range scenario.Datasets {
			datasetsMap := map[string]interface{}{}

			if scenario.Datasets.Name != nil {
				datasetsMap["name"] = scenario.Datasets.Name
			}

			if scenario.Datasets.Split != nil {
				datasetsMap["split"] = scenario.Datasets.Split
			}

			if scenario.Datasets.HeaderInFile != nil {
				datasetsMap["header_in_file"] = scenario.Datasets.HeaderInFile
			}

			if scenario.Datasets.HeaderColumns != nil {
				datasetsMap["header_columns"] = scenario.Datasets.HeaderColumns
			}

			if scenario.Datasets.LineCount != nil {
				datasetsMap["line_count"] = scenario.Datasets.LineCount
			}

			if scenario.Datasets.UpdatedAt != nil {
				datasetsMap["updated_at"] = scenario.Datasets.UpdatedAt
			}

			if scenario.Datasets.Size != nil {
				datasetsMap["size"] = scenario.Datasets.Size
			}

			if scenario.Datasets.HeadLines != nil {
				datasetsMap["head_lines"] = scenario.Datasets.HeadLines
			}

			if scenario.Datasets.TailLines != nil {
				datasetsMap["tail_lines"] = scenario.Datasets.TailLines
			}

			if scenario.Datasets.Type != nil {
				datasetsMap["type"] = scenario.Datasets.Type
			}

			if scenario.Datasets.FileId != nil {
				datasetsMap["file_id"] = scenario.Datasets.FileId
			}

			datasetsList = append(datasetsList, datasetsMap)
		}

		_ = d.Set("datasets", datasetsList)

	}

	if scenario.CronId != nil {
		_ = d.Set("cron_id", scenario.CronId)
	}

	if scenario.TestScripts != nil {
		testScriptsList := []interface{}{}
		for _, testScripts := range scenario.TestScripts {
			testScriptsMap := map[string]interface{}{}

			if scenario.TestScripts.Name != nil {
				testScriptsMap["name"] = scenario.TestScripts.Name
			}

			if scenario.TestScripts.Size != nil {
				testScriptsMap["size"] = scenario.TestScripts.Size
			}

			if scenario.TestScripts.Type != nil {
				testScriptsMap["type"] = scenario.TestScripts.Type
			}

			if scenario.TestScripts.UpdatedAt != nil {
				testScriptsMap["updated_at"] = scenario.TestScripts.UpdatedAt
			}

			if scenario.TestScripts.EncodedContent != nil {
				testScriptsMap["encoded_content"] = scenario.TestScripts.EncodedContent
			}

			if scenario.TestScripts.EncodedHttpArchive != nil {
				testScriptsMap["encoded_http_archive"] = scenario.TestScripts.EncodedHttpArchive
			}

			if scenario.TestScripts.LoadWeight != nil {
				testScriptsMap["load_weight"] = scenario.TestScripts.LoadWeight
			}

			testScriptsList = append(testScriptsList, testScriptsMap)
		}

		_ = d.Set("test_scripts", testScriptsList)

	}

	if scenario.Protocols != nil {
		protocolsList := []interface{}{}
		for _, protocols := range scenario.Protocols {
			protocolsMap := map[string]interface{}{}

			if scenario.Protocols.Name != nil {
				protocolsMap["name"] = scenario.Protocols.Name
			}

			if scenario.Protocols.Size != nil {
				protocolsMap["size"] = scenario.Protocols.Size
			}

			if scenario.Protocols.Type != nil {
				protocolsMap["type"] = scenario.Protocols.Type
			}

			if scenario.Protocols.UpdatedAt != nil {
				protocolsMap["updated_at"] = scenario.Protocols.UpdatedAt
			}

			if scenario.Protocols.FileId != nil {
				protocolsMap["file_id"] = scenario.Protocols.FileId
			}

			protocolsList = append(protocolsList, protocolsMap)
		}

		_ = d.Set("protocols", protocolsList)

	}

	if scenario.RequestFiles != nil {
		requestFilesList := []interface{}{}
		for _, requestFiles := range scenario.RequestFiles {
			requestFilesMap := map[string]interface{}{}

			if scenario.RequestFiles.Name != nil {
				requestFilesMap["name"] = scenario.RequestFiles.Name
			}

			if scenario.RequestFiles.Size != nil {
				requestFilesMap["size"] = scenario.RequestFiles.Size
			}

			if scenario.RequestFiles.Type != nil {
				requestFilesMap["type"] = scenario.RequestFiles.Type
			}

			if scenario.RequestFiles.UpdatedAt != nil {
				requestFilesMap["updated_at"] = scenario.RequestFiles.UpdatedAt
			}

			if scenario.RequestFiles.FileId != nil {
				requestFilesMap["file_id"] = scenario.RequestFiles.FileId
			}

			requestFilesList = append(requestFilesList, requestFilesMap)
		}

		_ = d.Set("request_files", requestFilesList)

	}

	if scenario.SLAPolicy != nil {
		sLAPolicyMap := map[string]interface{}{}

		if scenario.SLAPolicy.SLARules != nil {
			sLARulesList := []interface{}{}
			for _, sLARules := range scenario.SLAPolicy.SLARules {
				sLARulesMap := map[string]interface{}{}

				if sLARules.Metric != nil {
					sLARulesMap["metric"] = sLARules.Metric
				}

				if sLARules.Aggregation != nil {
					sLARulesMap["aggregation"] = sLARules.Aggregation
				}

				if sLARules.Condition != nil {
					sLARulesMap["condition"] = sLARules.Condition
				}

				if sLARules.Value != nil {
					sLARulesMap["value"] = sLARules.Value
				}

				if sLARules.LabelFilter != nil {
					labelFilterList := []interface{}{}
					for _, labelFilter := range sLARules.LabelFilter {
						labelFilterMap := map[string]interface{}{}

						if labelFilter.LabelName != nil {
							labelFilterMap["label_name"] = labelFilter.LabelName
						}

						if labelFilter.LabelValue != nil {
							labelFilterMap["label_value"] = labelFilter.LabelValue
						}

						labelFilterList = append(labelFilterList, labelFilterMap)
					}

					sLARulesMap["label_filter"] = []interface{}{labelFilterList}
				}

				if sLARules.AbortFlag != nil {
					sLARulesMap["abort_flag"] = sLARules.AbortFlag
				}

				if sLARules.For != nil {
					sLARulesMap["for"] = sLARules.For
				}

				sLARulesList = append(sLARulesList, sLARulesMap)
			}

			sLAPolicyMap["s_l_a_rules"] = []interface{}{sLARulesList}
		}

		if scenario.SLAPolicy.AlertChannel != nil {
			alertChannelMap := map[string]interface{}{}

			if scenario.SLAPolicy.AlertChannel.NoticeId != nil {
				alertChannelMap["notice_id"] = scenario.SLAPolicy.AlertChannel.NoticeId
			}

			if scenario.SLAPolicy.AlertChannel.AMPConsumerId != nil {
				alertChannelMap["a_m_p_consumer_id"] = scenario.SLAPolicy.AlertChannel.AMPConsumerId
			}

			sLAPolicyMap["alert_channel"] = []interface{}{alertChannelMap}
		}

		_ = d.Set("s_l_a_policy", []interface{}{sLAPolicyMap})
	}

	if scenario.Plugins != nil {
		pluginsList := []interface{}{}
		for _, plugins := range scenario.Plugins {
			pluginsMap := map[string]interface{}{}

			if scenario.Plugins.Name != nil {
				pluginsMap["name"] = scenario.Plugins.Name
			}

			if scenario.Plugins.Size != nil {
				pluginsMap["size"] = scenario.Plugins.Size
			}

			if scenario.Plugins.Type != nil {
				pluginsMap["type"] = scenario.Plugins.Type
			}

			if scenario.Plugins.UpdatedAt != nil {
				pluginsMap["updated_at"] = scenario.Plugins.UpdatedAt
			}

			if scenario.Plugins.FileId != nil {
				pluginsMap["file_id"] = scenario.Plugins.FileId
			}

			pluginsList = append(pluginsList, pluginsMap)
		}

		_ = d.Set("plugins", pluginsList)

	}

	if scenario.DomainNameConfig != nil {
		domainNameConfigMap := map[string]interface{}{}

		if scenario.DomainNameConfig.HostAliases != nil {
			hostAliasesList := []interface{}{}
			for _, hostAliases := range scenario.DomainNameConfig.HostAliases {
				hostAliasesMap := map[string]interface{}{}

				if hostAliases.HostNames != nil {
					hostAliasesMap["host_names"] = hostAliases.HostNames
				}

				if hostAliases.IP != nil {
					hostAliasesMap["i_p"] = hostAliases.IP
				}

				hostAliasesList = append(hostAliasesList, hostAliasesMap)
			}

			domainNameConfigMap["host_aliases"] = []interface{}{hostAliasesList}
		}

		if scenario.DomainNameConfig.DNSConfig != nil {
			dNSConfigMap := map[string]interface{}{}

			if scenario.DomainNameConfig.DNSConfig.Nameservers != nil {
				dNSConfigMap["nameservers"] = scenario.DomainNameConfig.DNSConfig.Nameservers
			}

			domainNameConfigMap["d_n_s_config"] = []interface{}{dNSConfigMap}
		}

		_ = d.Set("domain_name_config", []interface{}{domainNameConfigMap})
	}

	if scenario.Status != nil {
		_ = d.Set("status", scenario.Status)
	}

	if scenario.CreatedAt != nil {
		_ = d.Set("created_at", scenario.CreatedAt)
	}

	if scenario.UpdatedAt != nil {
		_ = d.Set("updated_at", scenario.UpdatedAt)
	}

	if scenario.AppId != nil {
		_ = d.Set("app_id", scenario.AppId)
	}

	if scenario.Uin != nil {
		_ = d.Set("uin", scenario.Uin)
	}

	if scenario.SubAccountUin != nil {
		_ = d.Set("sub_account_uin", scenario.SubAccountUin)
	}

	if scenario.NotificationHooks != nil {
		notificationHooksList := []interface{}{}
		for _, notificationHooks := range scenario.NotificationHooks {
			notificationHooksMap := map[string]interface{}{}

			if scenario.NotificationHooks.Events != nil {
				notificationHooksMap["events"] = scenario.NotificationHooks.Events
			}

			if scenario.NotificationHooks.URL != nil {
				notificationHooksMap["u_r_l"] = scenario.NotificationHooks.URL
			}

			notificationHooksList = append(notificationHooksList, notificationHooksMap)
		}

		_ = d.Set("notification_hooks", notificationHooksList)

	}

	return nil
}

func resourceTencentCloudPtsScenarioUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_scenario.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := pts.NewUpdateScenarioRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	scenarioId := idSplit[0]
	projectId := idSplit[1]

	request.ScenarioId = &scenarioId
	request.ProjectId = &projectId

	immutableArgs := []string{"name", "type", "project_id", "description", "load", "datasets", "cron_id", "test_scripts", "protocols", "request_files", "s_l_a_policy", "plugins", "domain_name_config", "status", "created_at", "updated_at", "app_id", "uin", "sub_account_uin", "notification_hooks"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("load") {
		if dMap, ok := helper.InterfacesHeadMap(d, "load"); ok {
			load := pts.Load{}
			if loadSpecMap, ok := helper.InterfaceToMap(dMap, "load_spec"); ok {
				loadSpec := pts.LoadSpec{}
				if concurrencyMap, ok := helper.InterfaceToMap(loadSpecMap, "concurrency"); ok {
					concurrency := pts.Concurrency{}
					if v, ok := concurrencyMap["stages"]; ok {
						for _, item := range v.([]interface{}) {
							stagesMap := item.(map[string]interface{})
							stage := pts.Stage{}
							if v, ok := stagesMap["duration_seconds"]; ok {
								stage.DurationSeconds = helper.IntUint64(v.(int))
							}
							if v, ok := stagesMap["target_virtual_users"]; ok {
								stage.TargetVirtualUsers = helper.IntUint64(v.(int))
							}
							concurrency.Stages = append(concurrency.Stages, &stage)
						}
					}
					if v, ok := concurrencyMap["iteration_count"]; ok {
						concurrency.IterationCount = helper.IntUint64(v.(int))
					}
					if v, ok := concurrencyMap["max_requests_per_second"]; ok {
						concurrency.MaxRequestsPerSecond = helper.IntUint64(v.(int))
					}
					if v, ok := concurrencyMap["graceful_stop_seconds"]; ok {
						concurrency.GracefulStopSeconds = helper.IntUint64(v.(int))
					}
					loadSpec.Concurrency = &concurrency
				}
				if requestsPerSecondMap, ok := helper.InterfaceToMap(loadSpecMap, "requests_per_second"); ok {
					requestsPerSecond := pts.RequestsPerSecond{}
					if v, ok := requestsPerSecondMap["max_requests_per_second"]; ok {
						requestsPerSecond.MaxRequestsPerSecond = helper.IntUint64(v.(int))
					}
					if v, ok := requestsPerSecondMap["duration_seconds"]; ok {
						requestsPerSecond.DurationSeconds = helper.IntUint64(v.(int))
					}
					if v, ok := requestsPerSecondMap["resources"]; ok {
						requestsPerSecond.Resources = helper.IntUint64(v.(int))
					}
					if v, ok := requestsPerSecondMap["start_requests_per_second"]; ok {
						requestsPerSecond.StartRequestsPerSecond = helper.IntUint64(v.(int))
					}
					if v, ok := requestsPerSecondMap["target_requests_per_second"]; ok {
						requestsPerSecond.TargetRequestsPerSecond = helper.IntUint64(v.(int))
					}
					if v, ok := requestsPerSecondMap["graceful_stop_seconds"]; ok {
						requestsPerSecond.GracefulStopSeconds = helper.IntUint64(v.(int))
					}
					loadSpec.RequestsPerSecond = &requestsPerSecond
				}
				if scriptOriginMap, ok := helper.InterfaceToMap(loadSpecMap, "script_origin"); ok {
					scriptOrigin := pts.ScriptOrigin{}
					if v, ok := scriptOriginMap["machine_number"]; ok {
						scriptOrigin.MachineNumber = helper.IntUint64(v.(int))
					}
					if v, ok := scriptOriginMap["machine_specification"]; ok {
						scriptOrigin.MachineSpecification = helper.String(v.(string))
					}
					if v, ok := scriptOriginMap["duration_seconds"]; ok {
						scriptOrigin.DurationSeconds = helper.IntUint64(v.(int))
					}
					loadSpec.ScriptOrigin = &scriptOrigin
				}
				load.LoadSpec = &loadSpec
			}
			if vpcLoadDistributionMap, ok := helper.InterfaceToMap(dMap, "vpc_load_distribution"); ok {
				vpcLoadDistribution := pts.VpcLoadDistribution{}
				if v, ok := vpcLoadDistributionMap["region_id"]; ok {
					vpcLoadDistribution.RegionId = helper.IntUint64(v.(int))
				}
				if v, ok := vpcLoadDistributionMap["region"]; ok {
					vpcLoadDistribution.Region = helper.String(v.(string))
				}
				if v, ok := vpcLoadDistributionMap["vpc_id"]; ok {
					vpcLoadDistribution.VpcId = helper.String(v.(string))
				}
				if v, ok := vpcLoadDistributionMap["subnet_ids"]; ok {
					subnetIdsSet := v.(*schema.Set).List()
					for i := range subnetIdsSet {
						subnetIds := subnetIdsSet[i].(string)
						vpcLoadDistribution.SubnetIds = append(vpcLoadDistribution.SubnetIds, &subnetIds)
					}
				}
				load.VpcLoadDistribution = &vpcLoadDistribution
			}
			if v, ok := dMap["geo_regions_load_distribution"]; ok {
				for _, item := range v.([]interface{}) {
					geoRegionsLoadDistributionMap := item.(map[string]interface{})
					geoRegionsLoadItem := pts.GeoRegionsLoadItem{}
					if v, ok := geoRegionsLoadDistributionMap["region_id"]; ok {
						geoRegionsLoadItem.RegionId = helper.IntUint64(v.(int))
					}
					if v, ok := geoRegionsLoadDistributionMap["region"]; ok {
						geoRegionsLoadItem.Region = helper.String(v.(string))
					}
					if v, ok := geoRegionsLoadDistributionMap["percentage"]; ok {
						geoRegionsLoadItem.Percentage = helper.IntUint64(v.(int))
					}
					load.GeoRegionsLoadDistribution = append(load.GeoRegionsLoadDistribution, &geoRegionsLoadItem)
				}
			}
			request.Load = &load
		}
	}

	if d.HasChange("datasets") {
		if v, ok := d.GetOk("datasets"); ok {
			for _, item := range v.([]interface{}) {
				testData := pts.TestData{}
				if v, ok := dMap["name"]; ok {
					testData.Name = helper.String(v.(string))
				}
				if v, ok := dMap["split"]; ok {
					testData.Split = helper.Bool(v.(bool))
				}
				if v, ok := dMap["header_in_file"]; ok {
					testData.HeaderInFile = helper.Bool(v.(bool))
				}
				if v, ok := dMap["header_columns"]; ok {
					headerColumnsSet := v.(*schema.Set).List()
					for i := range headerColumnsSet {
						headerColumns := headerColumnsSet[i].(string)
						testData.HeaderColumns = append(testData.HeaderColumns, &headerColumns)
					}
				}
				if v, ok := dMap["line_count"]; ok {
					testData.LineCount = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["updated_at"]; ok {
					testData.UpdatedAt = helper.String(v.(string))
				}
				if v, ok := dMap["size"]; ok {
					testData.Size = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["head_lines"]; ok {
					headLinesSet := v.(*schema.Set).List()
					for i := range headLinesSet {
						headLines := headLinesSet[i].(string)
						testData.HeadLines = append(testData.HeadLines, &headLines)
					}
				}
				if v, ok := dMap["tail_lines"]; ok {
					tailLinesSet := v.(*schema.Set).List()
					for i := range tailLinesSet {
						tailLines := tailLinesSet[i].(string)
						testData.TailLines = append(testData.TailLines, &tailLines)
					}
				}
				if v, ok := dMap["type"]; ok {
					testData.Type = helper.String(v.(string))
				}
				if v, ok := dMap["file_id"]; ok {
					testData.FileId = helper.String(v.(string))
				}
				request.Datasets = append(request.Datasets, &testData)
			}
		}
	}

	if d.HasChange("cron_id") {
		if v, ok := d.GetOk("cron_id"); ok {
			request.CronId = helper.String(v.(string))
		}
	}

	if d.HasChange("test_scripts") {
		if v, ok := d.GetOk("test_scripts"); ok {
			for _, item := range v.([]interface{}) {
				scriptInfo := pts.ScriptInfo{}
				if v, ok := dMap["name"]; ok {
					scriptInfo.Name = helper.String(v.(string))
				}
				if v, ok := dMap["size"]; ok {
					scriptInfo.Size = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["type"]; ok {
					scriptInfo.Type = helper.String(v.(string))
				}
				if v, ok := dMap["updated_at"]; ok {
					scriptInfo.UpdatedAt = helper.String(v.(string))
				}
				if v, ok := dMap["encoded_content"]; ok {
					scriptInfo.EncodedContent = helper.String(v.(string))
				}
				if v, ok := dMap["encoded_http_archive"]; ok {
					scriptInfo.EncodedHttpArchive = helper.String(v.(string))
				}
				if v, ok := dMap["load_weight"]; ok {
					scriptInfo.LoadWeight = helper.IntUint64(v.(int))
				}
				request.TestScripts = append(request.TestScripts, &scriptInfo)
			}
		}
	}

	if d.HasChange("protocols") {
		if v, ok := d.GetOk("protocols"); ok {
			for _, item := range v.([]interface{}) {
				protocolInfo := pts.ProtocolInfo{}
				if v, ok := dMap["name"]; ok {
					protocolInfo.Name = helper.String(v.(string))
				}
				if v, ok := dMap["size"]; ok {
					protocolInfo.Size = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["type"]; ok {
					protocolInfo.Type = helper.String(v.(string))
				}
				if v, ok := dMap["updated_at"]; ok {
					protocolInfo.UpdatedAt = helper.String(v.(string))
				}
				if v, ok := dMap["file_id"]; ok {
					protocolInfo.FileId = helper.String(v.(string))
				}
				request.Protocols = append(request.Protocols, &protocolInfo)
			}
		}
	}

	if d.HasChange("request_files") {
		if v, ok := d.GetOk("request_files"); ok {
			for _, item := range v.([]interface{}) {
				fileInfo := pts.FileInfo{}
				if v, ok := dMap["name"]; ok {
					fileInfo.Name = helper.String(v.(string))
				}
				if v, ok := dMap["size"]; ok {
					fileInfo.Size = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["type"]; ok {
					fileInfo.Type = helper.String(v.(string))
				}
				if v, ok := dMap["updated_at"]; ok {
					fileInfo.UpdatedAt = helper.String(v.(string))
				}
				if v, ok := dMap["file_id"]; ok {
					fileInfo.FileId = helper.String(v.(string))
				}
				request.RequestFiles = append(request.RequestFiles, &fileInfo)
			}
		}
	}

	if d.HasChange("s_l_a_policy") {
		if dMap, ok := helper.InterfacesHeadMap(d, "s_l_a_policy"); ok {
			sLAPolicy := pts.SLAPolicy{}
			if v, ok := dMap["s_l_a_rules"]; ok {
				for _, item := range v.([]interface{}) {
					sLARulesMap := item.(map[string]interface{})
					sLARule := pts.SLARule{}
					if v, ok := sLARulesMap["metric"]; ok {
						sLARule.Metric = helper.String(v.(string))
					}
					if v, ok := sLARulesMap["aggregation"]; ok {
						sLARule.Aggregation = helper.String(v.(string))
					}
					if v, ok := sLARulesMap["condition"]; ok {
						sLARule.Condition = helper.String(v.(string))
					}
					if v, ok := sLARulesMap["value"]; ok {
						sLARule.Value = helper.IntInt64(v.(int))
					}
					if v, ok := sLARulesMap["label_filter"]; ok {
						for _, item := range v.([]interface{}) {
							labelFilterMap := item.(map[string]interface{})
							sLALabel := pts.SLALabel{}
							if v, ok := labelFilterMap["label_name"]; ok {
								sLALabel.LabelName = helper.String(v.(string))
							}
							if v, ok := labelFilterMap["label_value"]; ok {
								sLALabel.LabelValue = helper.String(v.(string))
							}
							sLARule.LabelFilter = append(sLARule.LabelFilter, &sLALabel)
						}
					}
					if v, ok := sLARulesMap["abort_flag"]; ok {
						sLARule.AbortFlag = helper.Bool(v.(bool))
					}
					if v, ok := sLARulesMap["for"]; ok {
						sLARule.For = helper.String(v.(string))
					}
					sLAPolicy.SLARules = append(sLAPolicy.SLARules, &sLARule)
				}
			}
			if alertChannelMap, ok := helper.InterfaceToMap(dMap, "alert_channel"); ok {
				alertChannel := pts.AlertChannel{}
				if v, ok := alertChannelMap["notice_id"]; ok {
					alertChannel.NoticeId = helper.String(v.(string))
				}
				if v, ok := alertChannelMap["a_m_p_consumer_id"]; ok {
					alertChannel.AMPConsumerId = helper.String(v.(string))
				}
				sLAPolicy.AlertChannel = &alertChannel
			}
			request.SLAPolicy = &sLAPolicy
		}
	}

	if d.HasChange("plugins") {
		if v, ok := d.GetOk("plugins"); ok {
			for _, item := range v.([]interface{}) {
				fileInfo := pts.FileInfo{}
				if v, ok := dMap["name"]; ok {
					fileInfo.Name = helper.String(v.(string))
				}
				if v, ok := dMap["size"]; ok {
					fileInfo.Size = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["type"]; ok {
					fileInfo.Type = helper.String(v.(string))
				}
				if v, ok := dMap["updated_at"]; ok {
					fileInfo.UpdatedAt = helper.String(v.(string))
				}
				if v, ok := dMap["file_id"]; ok {
					fileInfo.FileId = helper.String(v.(string))
				}
				request.Plugins = append(request.Plugins, &fileInfo)
			}
		}
	}

	if d.HasChange("domain_name_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "domain_name_config"); ok {
			domainNameConfig := pts.DomainNameConfig{}
			if v, ok := dMap["host_aliases"]; ok {
				for _, item := range v.([]interface{}) {
					hostAliasesMap := item.(map[string]interface{})
					hostAlias := pts.HostAlias{}
					if v, ok := hostAliasesMap["host_names"]; ok {
						hostNamesSet := v.(*schema.Set).List()
						for i := range hostNamesSet {
							hostNames := hostNamesSet[i].(string)
							hostAlias.HostNames = append(hostAlias.HostNames, &hostNames)
						}
					}
					if v, ok := hostAliasesMap["i_p"]; ok {
						hostAlias.IP = helper.String(v.(string))
					}
					domainNameConfig.HostAliases = append(domainNameConfig.HostAliases, &hostAlias)
				}
			}
			if dNSConfigMap, ok := helper.InterfaceToMap(dMap, "d_n_s_config"); ok {
				dNSConfig := pts.DNSConfig{}
				if v, ok := dNSConfigMap["nameservers"]; ok {
					nameserversSet := v.(*schema.Set).List()
					for i := range nameserversSet {
						nameservers := nameserversSet[i].(string)
						dNSConfig.Nameservers = append(dNSConfig.Nameservers, &nameservers)
					}
				}
				domainNameConfig.DNSConfig = &dNSConfig
			}
			request.DomainNameConfig = &domainNameConfig
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().UpdateScenario(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update pts scenario failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsScenarioRead(d, meta)
}

func resourceTencentCloudPtsScenarioDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_scenario.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	scenarioId := idSplit[0]
	projectId := idSplit[1]

	if err := service.DeletePtsScenarioById(ctx, scenarioId, projectId); err != nil {
		return err
	}

	return nil
}
