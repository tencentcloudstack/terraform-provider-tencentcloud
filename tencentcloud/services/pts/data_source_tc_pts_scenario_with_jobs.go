package pts

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPtsScenarioWithJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPtsScenarioWithJobsRead,
		Schema: map[string]*schema.Schema{
			"project_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Project ID list.",
			},

			"scenario_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Scenario ID list.",
			},

			"scenario_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scenario name.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The field column used for ordering.",
			},

			"ascend": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to use ascending order.",
			},

			"ignore_script": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to ignore the script content.",
			},

			"ignore_dataset": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to ignore the dataset.",
			},

			"scenario_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scenario type, e.g.: pts-http, pts-js, pts-trpc, pts-jmeter.",
			},

			"owner": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The job owner.",
			},

			"scenario_with_jobs_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The scenario configuration and its jobs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scenario": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The returned scenario.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scenario_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scenario ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scenario name.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scenario description.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scenario type, e.g.: pts-http, pts-js, pts-trpc, pts-jmeter.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Scenario status.",
									},
									"load": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Scenario is load test configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"load_spec": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Scenario is load specification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"concurrency": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The configuration for the concurrency mode.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"stages": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The configuration for the multi-stage load test.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"duration_seconds": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The execution time for the load test.",
																					},
																					"target_virtual_users": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The number of the target virtual users.",
																					},
																				},
																			},
																		},
																		"iteration_count": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The iteration count of the load test.",
																		},
																		"max_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The maximum RPS.",
																		},
																		"graceful_stop_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The waiting period for a graceful shutdown.",
																		},
																		"resources": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The resource count of the load test.",
																		},
																	},
																},
															},
															"requests_per_second": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The configuration of the RPS mode load test.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"max_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The maximum RPS.",
																		},
																		"duration_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The execution time of the load test.",
																		},
																		"target_virtual_users": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Deprecated.",
																		},
																		"resources": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The recource count of the load test.",
																		},
																		"start_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The starting minimum RPS.",
																		},
																		"target_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The target RPS.",
																		},
																		"graceful_stop_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The waiting period for a graceful shutdown.",
																		},
																	},
																},
															},
															"script_origin": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The script origin.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"machine_number": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The load test machine number.",
																		},
																		"machine_specification": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The load test machine specification.",
																		},
																		"duration_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The load test execution time.",
																		},
																	},
																},
															},
														},
													},
												},
												"vpc_load_distribution": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The distribution of the load source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"region_id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Region ID.",
															},
															"region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region.",
															},
															"vpc_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The VPC ID.",
															},
															"subnet_ids": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "The subnet ID list.",
															},
														},
													},
												},
												"geo_regions_load_distribution": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The geographical distribution of the load source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"region_id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Region ID.",
															},
															"region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region.",
															},
															"percentage": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Percentage.",
															},
														},
													},
												},
											},
										},
									},
									"encoded_scripts": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deprecated.",
									},
									"configs": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Deprecated.",
									},
									"extensions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Deprecated.",
									},
									"datasets": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The test data sets for the load test.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The file name of the test data sets.",
												},
												"split": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to split the test data.",
												},
												"header_in_file": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the first line contains the parameter names.",
												},
												"header_columns": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The parameter name list.",
												},
												"line_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The line count of the file.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The update time of the file.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The byte count of the file.",
												},
												"head_lines": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The header lines of the file.",
												},
												"tail_lines": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The tail lines of the file.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The file type.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The file ID.",
												},
											},
										},
									},
									"sla_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the SLA policy.",
									},
									"cron_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cron job ID.",
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The creation time of the scenario.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The updating time of the scenario.",
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Project ID.",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "AppId.",
									},
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Uin.",
									},
									"sub_account_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "SubAccountUin.",
									},
									"test_scripts": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The script of the load test.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"encoded_content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The base64 encoded content.",
												},
												"encoded_http_archive": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The base64 encoded HAR.",
												},
												"load_weight": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The weight of the script, ranging from 1 to 100.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"protocols": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The protocol file.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"request_files": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The files in the request.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"sla_policy": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The SLA policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sla_rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The SLA rules.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The load test metrics.",
															},
															"aggregation": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The aggregation method of the metrics.",
															},
															"condition": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The operator for checking the condition.",
															},
															"value": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold in the condition.",
															},
															"label_filter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The label filter.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"label_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Label name.",
																		},
																		"label_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Label value.",
																		},
																	},
																},
															},
															"abort_flag": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to abort the load test job.",
															},
															"for": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The duration for checking the condition.",
															},
														},
													},
												},
												"alert_channel": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The alert channel.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"notice_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The notice ID bound with this alert channel.",
															},
															"amp_consumer_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "AMP consumer ID.",
															},
														},
													},
												},
											},
										},
									},
									"plugins": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Plugins.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"domain_name_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The configuration for parsing domain names.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host_aliases": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The configuration for host aliases.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"host_names": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "Host names.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "IP.",
															},
														},
													},
												},
												"dns_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The DNS configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"nameservers": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "DNS IP list.",
															},
														},
													},
												},
											},
										},
									},
									"notification_hooks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The notification hooks.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The notification hook.",
												},
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The callback URL.",
												},
											},
										},
									},
									"owner": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The owner.",
									},
									"project_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Project name.",
									},
								},
							},
						},
						"jobs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Jobs related to the scenario.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"job_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job ID.",
									},
									"scenario_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scenario ID.",
									},
									"load": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The configuration of the load.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"load_spec": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The specification of the load configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"concurrency": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The configuration of the concurrency load test mode.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"stages": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The multi-stage configuration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"duration_seconds": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The execution time.",
																					},
																					"target_virtual_users": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The target count of the virtual users.",
																					},
																				},
																			},
																		},
																		"iteration_count": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The iteration count.",
																		},
																		"max_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The maximum RPS.",
																		},
																		"graceful_stop_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The waiting period for a graceful shutdown.",
																		},
																		"resources": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The count of the load test resource.",
																		},
																	},
																},
															},
															"requests_per_second": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The configuration of the RPS mode load test.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"max_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The maximum RPS.",
																		},
																		"duration_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The execution time.",
																		},
																		"target_virtual_users": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Deprecated.",
																		},
																		"resources": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The count of the load test resource.",
																		},
																		"start_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The starting minimum RPS.",
																		},
																		"target_requests_per_second": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The target RPS.",
																		},
																		"graceful_stop_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The waiting period for a gracefulshutdown.",
																		},
																	},
																},
															},
															"script_origin": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The script origin.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"machine_number": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Machine number.",
																		},
																		"machine_specification": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Machine specification.",
																		},
																		"duration_seconds": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The execution time.",
																		},
																	},
																},
															},
														},
													},
												},
												"vpc_load_distribution": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The distribution of the load source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"region_id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Region ID.",
															},
															"region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region.",
															},
															"vpc_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "VPC ID.",
															},
															"subnet_ids": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "The subnet ID list.",
															},
														},
													},
												},
												"geo_regions_load_distribution": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The geographical distribution of the load source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"region_id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Region ID.",
															},
															"region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region.",
															},
															"percentage": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Percentage.",
															},
														},
													},
												},
											},
										},
									},
									"configs": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Deprecated.",
									},
									"datasets": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The test data sets.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Test data set name.",
												},
												"split": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to split the test data.",
												},
												"header_in_file": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the first line contains the parameter names.",
												},
												"header_columns": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The parameter name list.",
												},
												"line_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The line count of the file.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"head_lines": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The header lines of the file.",
												},
												"tail_lines": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The tail lines of the file.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"extensions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Deprecated.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Job running status. JobUnknown: 0,JobCreated:1,JobPending:2, JobPreparing:3,JobSelectClustering:4,JobCreateTasking:5,JobSyncTasking:6 JobRunning:11,JobFinished:12,JobPrepareException:13,JobFinishException:14,JobAborting:15,JobAborted:16,JobAbortException:17,JobDeleted:18, JobSelectClusterException:19,JobCreateTaskException:20,JobSyncTaskException:21.",
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The job starting time.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The job ending time.",
									},
									"max_virtual_user_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum VU of the job.",
									},
									"note": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The note of the job.",
									},
									"error_rate": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Error rate.",
									},
									"job_owner": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job owner.",
									},
									"load_sources": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Deprecated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP of the load source.",
												},
												"pod_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The pod name of the load source.",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region.",
												},
											},
										},
									},
									"duration": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Job running duration.",
									},
									"max_requests_per_second": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum RPS.",
									},
									"request_total": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The total reqeust count.",
									},
									"requests_per_second": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "RPS.",
									},
									"response_time_average": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The average response time.",
									},
									"response_time_p99": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The 99 percentile of the response time.",
									},
									"response_time_p95": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The 95 percentile of the response time.",
									},
									"response_time_p90": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The 90 percentile of the response time.",
									},
									"scripts": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Deprecated.",
									},
									"response_time_max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The maximum response time.",
									},
									"response_time_min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The minimum response time.",
									},
									"load_source_infos": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The load source information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP of the load source.",
												},
												"pod_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The pod name of the load source.",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region.",
												},
											},
										},
									},
									"test_scripts": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Test scripts.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"encoded_content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The base64 encoded content.",
												},
												"encoded_http_archive": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The base64 encoded HAR.",
												},
												"load_weight": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The weight of the script, ranging from 1 to 100.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"protocols": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The protocol file.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"request_files": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The files in the request.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"plugins": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Plugins.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File name.",
												},
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "File size.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File type.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time of the most recent update.",
												},
												"file_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "File ID.",
												},
											},
										},
									},
									"cron_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cron job ID.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scenario type.",
									},
									"domain_name_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The configuration for parsing domain names.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host_aliases": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The configuration for host aliases.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"host_names": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "Host names.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "IP.",
															},
														},
													},
												},
												"dns_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The DNS configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"nameservers": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "DNS IP list.",
															},
														},
													},
												},
											},
										},
									},
									"debug": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to run the job in the debug mode. The default value is false.",
									},
									"abort_reason": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The reason for aborting the job.",
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The job creation time.",
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Project ID.",
									},
									"notification_hooks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Notification hooks.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Notification hook.",
												},
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The callback URL.",
												},
											},
										},
									},
									"network_receive_rate": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The rate of receiving bytes.",
									},
									"network_send_rate": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The rate of sending bytes.",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The message describing the job running status.",
									},
									"project_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Project name.",
									},
									"scenario_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scenario name.",
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

func dataSourceTencentCloudPtsScenarioWithJobsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_pts_scenario_with_jobs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		paramMap["ProjectIds"] = helper.InterfacesStringsPoint(projectIdsSet)
	}

	if v, ok := d.GetOk("scenario_ids"); ok {
		scenarioIdsSet := v.(*schema.Set).List()
		paramMap["ScenarioIds"] = helper.InterfacesStringsPoint(scenarioIdsSet)
	}

	if v, ok := d.GetOk("scenario_name"); ok {
		paramMap["ScenarioName"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("scenario_status"); v != nil {
		paramMap["ScenarioStatus"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ascend"); v != nil {
		paramMap["Ascend"] = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("ignore_script"); v != nil {
		paramMap["IgnoreScript"] = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("ignore_dataset"); v != nil {
		paramMap["IgnoreDataset"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("scenario_type"); ok {
		paramMap["ScenarioType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner"); ok {
		paramMap["Owner"] = helper.String(v.(string))
	}

	service := PtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var scenarioWithJobsSet []*pts.ScenarioWithJobs

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePtsScenarioWithJobsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		scenarioWithJobsSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(scenarioWithJobsSet))
	tmpList := make([]map[string]interface{}, 0, len(scenarioWithJobsSet))

	if scenarioWithJobsSet != nil {
		for _, scenarioWithJobs := range scenarioWithJobsSet {
			scenarioWithJobsMap := map[string]interface{}{}

			if scenarioWithJobs.Scenario != nil {
				scenario := scenarioWithJobs.Scenario
				scenarioMap := map[string]interface{}{}

				if scenario.ScenarioId != nil {
					scenarioMap["scenario_id"] = scenario.ScenarioId
					ids = append(ids, *scenario.ScenarioId)
				}

				if scenario.Name != nil {
					scenarioMap["name"] = scenario.Name
				}

				if scenario.Description != nil {
					scenarioMap["description"] = scenario.Description
				}

				if scenario.Type != nil {
					scenarioMap["type"] = scenario.Type
				}

				if scenario.Status != nil {
					scenarioMap["status"] = scenario.Status
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

								concurrencyMap["stages"] = stagesList
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

							if scenario.Load.LoadSpec.Concurrency.Resources != nil {
								concurrencyMap["resources"] = scenario.Load.LoadSpec.Concurrency.Resources
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

							if scenario.Load.LoadSpec.RequestsPerSecond.TargetVirtualUsers != nil {
								requestsPerSecondMap["target_virtual_users"] = scenario.Load.LoadSpec.RequestsPerSecond.TargetVirtualUsers
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

						loadMap["geo_regions_load_distribution"] = geoRegionsLoadDistributionList
					}

					scenarioMap["load"] = []interface{}{loadMap}
				}

				if scenario.EncodedScripts != nil {
					scenarioMap["encoded_scripts"] = scenario.EncodedScripts
				}

				if scenario.Configs != nil {
					scenarioMap["configs"] = scenario.Configs
				}

				if scenario.Extensions != nil {
					scenarioMap["extensions"] = scenario.Extensions
				}

				if scenario.Datasets != nil {
					datasetsList := []interface{}{}
					for _, datasets := range scenario.Datasets {
						datasetsMap := map[string]interface{}{}

						if datasets.Name != nil {
							datasetsMap["name"] = datasets.Name
						}

						if datasets.Split != nil {
							datasetsMap["split"] = datasets.Split
						}

						if datasets.HeaderInFile != nil {
							datasetsMap["header_in_file"] = datasets.HeaderInFile
						}

						if datasets.HeaderColumns != nil {
							datasetsMap["header_columns"] = datasets.HeaderColumns
						}

						if datasets.LineCount != nil {
							datasetsMap["line_count"] = datasets.LineCount
						}

						if datasets.UpdatedAt != nil {
							datasetsMap["updated_at"] = datasets.UpdatedAt
						}

						if datasets.Size != nil {
							datasetsMap["size"] = datasets.Size
						}

						if datasets.HeadLines != nil {
							datasetsMap["head_lines"] = datasets.HeadLines
						}

						if datasets.TailLines != nil {
							datasetsMap["tail_lines"] = datasets.TailLines
						}

						if datasets.Type != nil {
							datasetsMap["type"] = datasets.Type
						}

						if datasets.FileId != nil {
							datasetsMap["file_id"] = datasets.FileId
						}

						datasetsList = append(datasetsList, datasetsMap)
					}

					scenarioMap["datasets"] = datasetsList
				}

				if scenario.SLAId != nil {
					scenarioMap["sla_id"] = scenario.SLAId
				}

				if scenario.CronId != nil {
					scenarioMap["cron_id"] = scenario.CronId
				}

				if scenario.CreatedAt != nil {
					scenarioMap["created_at"] = scenario.CreatedAt
				}

				if scenario.UpdatedAt != nil {
					scenarioMap["updated_at"] = scenario.UpdatedAt
				}

				if scenario.ProjectId != nil {
					scenarioMap["project_id"] = scenario.ProjectId
				}

				if scenario.AppId != nil {
					scenarioMap["app_id"] = scenario.AppId
				}

				if scenario.Uin != nil {
					scenarioMap["uin"] = scenario.Uin
				}

				if scenario.SubAccountUin != nil {
					scenarioMap["sub_account_uin"] = scenario.SubAccountUin
				}

				if scenario.TestScripts != nil {
					testScriptsList := []interface{}{}
					for _, testScripts := range scenario.TestScripts {
						testScriptsMap := map[string]interface{}{}

						if testScripts.Name != nil {
							testScriptsMap["name"] = testScripts.Name
						}

						if testScripts.Size != nil {
							testScriptsMap["size"] = testScripts.Size
						}

						if testScripts.Type != nil {
							testScriptsMap["type"] = testScripts.Type
						}

						if testScripts.UpdatedAt != nil {
							testScriptsMap["updated_at"] = testScripts.UpdatedAt
						}

						if testScripts.EncodedContent != nil {
							testScriptsMap["encoded_content"] = testScripts.EncodedContent
						}

						if testScripts.EncodedHttpArchive != nil {
							testScriptsMap["encoded_http_archive"] = testScripts.EncodedHttpArchive
						}

						if testScripts.LoadWeight != nil {
							testScriptsMap["load_weight"] = testScripts.LoadWeight
						}

						if testScripts.FileId != nil {
							testScriptsMap["file_id"] = testScripts.FileId
						}

						testScriptsList = append(testScriptsList, testScriptsMap)
					}

					scenarioMap["test_scripts"] = testScriptsList
				}

				if scenario.Protocols != nil {
					protocolsList := []interface{}{}
					for _, protocols := range scenario.Protocols {
						protocolsMap := map[string]interface{}{}

						if protocols.Name != nil {
							protocolsMap["name"] = protocols.Name
						}

						if protocols.Size != nil {
							protocolsMap["size"] = protocols.Size
						}

						if protocols.Type != nil {
							protocolsMap["type"] = protocols.Type
						}

						if protocols.UpdatedAt != nil {
							protocolsMap["updated_at"] = protocols.UpdatedAt
						}

						if protocols.FileId != nil {
							protocolsMap["file_id"] = protocols.FileId
						}

						protocolsList = append(protocolsList, protocolsMap)
					}

					scenarioMap["protocols"] = protocolsList
				}

				if scenario.RequestFiles != nil {
					requestFilesList := []interface{}{}
					for _, requestFiles := range scenario.RequestFiles {
						requestFilesMap := map[string]interface{}{}

						if requestFiles.Name != nil {
							requestFilesMap["name"] = requestFiles.Name
						}

						if requestFiles.Size != nil {
							requestFilesMap["size"] = requestFiles.Size
						}

						if requestFiles.Type != nil {
							requestFilesMap["type"] = requestFiles.Type
						}

						if requestFiles.UpdatedAt != nil {
							requestFilesMap["updated_at"] = requestFiles.UpdatedAt
						}

						if requestFiles.FileId != nil {
							requestFilesMap["file_id"] = requestFiles.FileId
						}

						requestFilesList = append(requestFilesList, requestFilesMap)
					}

					scenarioMap["request_files"] = requestFilesList
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

								sLARulesMap["label_filter"] = labelFilterList
							}

							if sLARules.AbortFlag != nil {
								sLARulesMap["abort_flag"] = sLARules.AbortFlag
							}

							if sLARules.For != nil {
								sLARulesMap["for"] = sLARules.For
							}

							sLARulesList = append(sLARulesList, sLARulesMap)
						}

						sLAPolicyMap["sla_rules"] = sLARulesList
					}

					if scenario.SLAPolicy.AlertChannel != nil {
						alertChannelMap := map[string]interface{}{}

						if scenario.SLAPolicy.AlertChannel.NoticeId != nil {
							alertChannelMap["notice_id"] = scenario.SLAPolicy.AlertChannel.NoticeId
						}

						if scenario.SLAPolicy.AlertChannel.AMPConsumerId != nil {
							alertChannelMap["amp_consumer_id"] = scenario.SLAPolicy.AlertChannel.AMPConsumerId
						}

						sLAPolicyMap["alert_channel"] = []interface{}{alertChannelMap}
					}

					scenarioMap["sla_policy"] = []interface{}{sLAPolicyMap}
				}

				if scenario.Plugins != nil {
					pluginsList := []interface{}{}
					for _, plugins := range scenario.Plugins {
						pluginsMap := map[string]interface{}{}

						if plugins.Name != nil {
							pluginsMap["name"] = plugins.Name
						}

						if plugins.Size != nil {
							pluginsMap["size"] = plugins.Size
						}

						if plugins.Type != nil {
							pluginsMap["type"] = plugins.Type
						}

						if plugins.UpdatedAt != nil {
							pluginsMap["updated_at"] = plugins.UpdatedAt
						}

						if plugins.FileId != nil {
							pluginsMap["file_id"] = plugins.FileId
						}

						pluginsList = append(pluginsList, pluginsMap)
					}

					scenarioMap["plugins"] = pluginsList
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
								hostAliasesMap["ip"] = hostAliases.IP
							}

							hostAliasesList = append(hostAliasesList, hostAliasesMap)
						}

						domainNameConfigMap["host_aliases"] = hostAliasesList
					}

					if scenario.DomainNameConfig.DNSConfig != nil {
						dNSConfigMap := map[string]interface{}{}

						if scenario.DomainNameConfig.DNSConfig.Nameservers != nil {
							dNSConfigMap["nameservers"] = scenario.DomainNameConfig.DNSConfig.Nameservers
						}

						domainNameConfigMap["dns_config"] = []interface{}{dNSConfigMap}
					}

					scenarioMap["domain_name_config"] = []interface{}{domainNameConfigMap}
				}

				if scenario.NotificationHooks != nil {
					notificationHooksList := []interface{}{}
					for _, notificationHooks := range scenario.NotificationHooks {
						notificationHooksMap := map[string]interface{}{}

						if notificationHooks.Events != nil {
							notificationHooksMap["events"] = notificationHooks.Events
						}

						if notificationHooks.URL != nil {
							notificationHooksMap["url"] = notificationHooks.URL
						}

						notificationHooksList = append(notificationHooksList, notificationHooksMap)
					}

					scenarioMap["notification_hooks"] = notificationHooksList
				}

				if scenario.Owner != nil {
					scenarioMap["owner"] = scenario.Owner
				}

				if scenario.ProjectName != nil {
					scenarioMap["project_name"] = scenario.ProjectName
				}

				scenarioWithJobsMap["scenario"] = []interface{}{scenarioMap}
			}

			if scenarioWithJobs.Jobs != nil {
				jobsList := []interface{}{}
				for _, jobs := range scenarioWithJobs.Jobs {
					jobsMap := map[string]interface{}{}

					if jobs.JobId != nil {
						jobsMap["job_id"] = jobs.JobId
						ids = append(ids, *jobs.JobId)
					}

					if jobs.ScenarioId != nil {
						jobsMap["scenario_id"] = jobs.ScenarioId
					}

					if jobs.Load != nil {
						loadMap := map[string]interface{}{}

						if jobs.Load.LoadSpec != nil {
							loadSpecMap := map[string]interface{}{}

							if jobs.Load.LoadSpec.Concurrency != nil {
								concurrencyMap := map[string]interface{}{}

								if jobs.Load.LoadSpec.Concurrency.Stages != nil {
									stagesList := []interface{}{}
									for _, stages := range jobs.Load.LoadSpec.Concurrency.Stages {
										stagesMap := map[string]interface{}{}

										if stages.DurationSeconds != nil {
											stagesMap["duration_seconds"] = stages.DurationSeconds
										}

										if stages.TargetVirtualUsers != nil {
											stagesMap["target_virtual_users"] = stages.TargetVirtualUsers
										}

										stagesList = append(stagesList, stagesMap)
									}

									concurrencyMap["stages"] = stagesList
								}

								if jobs.Load.LoadSpec.Concurrency.IterationCount != nil {
									concurrencyMap["iteration_count"] = jobs.Load.LoadSpec.Concurrency.IterationCount
								}

								if jobs.Load.LoadSpec.Concurrency.MaxRequestsPerSecond != nil {
									concurrencyMap["max_requests_per_second"] = jobs.Load.LoadSpec.Concurrency.MaxRequestsPerSecond
								}

								if jobs.Load.LoadSpec.Concurrency.GracefulStopSeconds != nil {
									concurrencyMap["graceful_stop_seconds"] = jobs.Load.LoadSpec.Concurrency.GracefulStopSeconds
								}

								if jobs.Load.LoadSpec.Concurrency.Resources != nil {
									concurrencyMap["resources"] = jobs.Load.LoadSpec.Concurrency.Resources
								}

								loadSpecMap["concurrency"] = []interface{}{concurrencyMap}
							}

							if jobs.Load.LoadSpec.RequestsPerSecond != nil {
								requestsPerSecondMap := map[string]interface{}{}

								if jobs.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond != nil {
									requestsPerSecondMap["max_requests_per_second"] = jobs.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond
								}

								if jobs.Load.LoadSpec.RequestsPerSecond.DurationSeconds != nil {
									requestsPerSecondMap["duration_seconds"] = jobs.Load.LoadSpec.RequestsPerSecond.DurationSeconds
								}

								if jobs.Load.LoadSpec.RequestsPerSecond.TargetVirtualUsers != nil {
									requestsPerSecondMap["target_virtual_users"] = jobs.Load.LoadSpec.RequestsPerSecond.TargetVirtualUsers
								}

								if jobs.Load.LoadSpec.RequestsPerSecond.Resources != nil {
									requestsPerSecondMap["resources"] = jobs.Load.LoadSpec.RequestsPerSecond.Resources
								}

								if jobs.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond != nil {
									requestsPerSecondMap["start_requests_per_second"] = jobs.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond
								}

								if jobs.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond != nil {
									requestsPerSecondMap["target_requests_per_second"] = jobs.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond
								}

								if jobs.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds != nil {
									requestsPerSecondMap["graceful_stop_seconds"] = jobs.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds
								}

								loadSpecMap["requests_per_second"] = []interface{}{requestsPerSecondMap}
							}

							if jobs.Load.LoadSpec.ScriptOrigin != nil {
								scriptOriginMap := map[string]interface{}{}

								if jobs.Load.LoadSpec.ScriptOrigin.MachineNumber != nil {
									scriptOriginMap["machine_number"] = jobs.Load.LoadSpec.ScriptOrigin.MachineNumber
								}

								if jobs.Load.LoadSpec.ScriptOrigin.MachineSpecification != nil {
									scriptOriginMap["machine_specification"] = jobs.Load.LoadSpec.ScriptOrigin.MachineSpecification
								}

								if jobs.Load.LoadSpec.ScriptOrigin.DurationSeconds != nil {
									scriptOriginMap["duration_seconds"] = jobs.Load.LoadSpec.ScriptOrigin.DurationSeconds
								}

								loadSpecMap["script_origin"] = []interface{}{scriptOriginMap}
							}

							loadMap["load_spec"] = []interface{}{loadSpecMap}
						}

						if jobs.Load.VpcLoadDistribution != nil {
							vpcLoadDistributionMap := map[string]interface{}{}

							if jobs.Load.VpcLoadDistribution.RegionId != nil {
								vpcLoadDistributionMap["region_id"] = jobs.Load.VpcLoadDistribution.RegionId
							}

							if jobs.Load.VpcLoadDistribution.Region != nil {
								vpcLoadDistributionMap["region"] = jobs.Load.VpcLoadDistribution.Region
							}

							if jobs.Load.VpcLoadDistribution.VpcId != nil {
								vpcLoadDistributionMap["vpc_id"] = jobs.Load.VpcLoadDistribution.VpcId
							}

							if jobs.Load.VpcLoadDistribution.SubnetIds != nil {
								vpcLoadDistributionMap["subnet_ids"] = jobs.Load.VpcLoadDistribution.SubnetIds
							}

							loadMap["vpc_load_distribution"] = []interface{}{vpcLoadDistributionMap}
						}

						if jobs.Load.GeoRegionsLoadDistribution != nil {
							geoRegionsLoadDistributionList := []interface{}{}
							for _, geoRegionsLoadDistribution := range jobs.Load.GeoRegionsLoadDistribution {
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

							loadMap["geo_regions_load_distribution"] = geoRegionsLoadDistributionList
						}

						jobsMap["load"] = []interface{}{loadMap}
					}

					if jobs.Configs != nil {
						jobsMap["configs"] = jobs.Configs
					}

					if jobs.Datasets != nil {
						datasetsList := []interface{}{}
						for _, datasets := range jobs.Datasets {
							datasetsMap := map[string]interface{}{}

							if datasets.Name != nil {
								datasetsMap["name"] = datasets.Name
							}

							if datasets.Split != nil {
								datasetsMap["split"] = datasets.Split
							}

							if datasets.HeaderInFile != nil {
								datasetsMap["header_in_file"] = datasets.HeaderInFile
							}

							if datasets.HeaderColumns != nil {
								datasetsMap["header_columns"] = datasets.HeaderColumns
							}

							if datasets.LineCount != nil {
								datasetsMap["line_count"] = datasets.LineCount
							}

							if datasets.UpdatedAt != nil {
								datasetsMap["updated_at"] = datasets.UpdatedAt
							}

							if datasets.Size != nil {
								datasetsMap["size"] = datasets.Size
							}

							if datasets.HeadLines != nil {
								datasetsMap["head_lines"] = datasets.HeadLines
							}

							if datasets.TailLines != nil {
								datasetsMap["tail_lines"] = datasets.TailLines
							}

							if datasets.Type != nil {
								datasetsMap["type"] = datasets.Type
							}

							if datasets.FileId != nil {
								datasetsMap["file_id"] = datasets.FileId
							}

							datasetsList = append(datasetsList, datasetsMap)
						}

						jobsMap["datasets"] = datasetsList
					}

					if jobs.Extensions != nil {
						jobsMap["extensions"] = jobs.Extensions
					}

					if jobs.Status != nil {
						jobsMap["status"] = jobs.Status
					}

					if jobs.StartTime != nil {
						jobsMap["start_time"] = jobs.StartTime
					}

					if jobs.EndTime != nil {
						jobsMap["end_time"] = jobs.EndTime
					}

					if jobs.MaxVirtualUserCount != nil {
						jobsMap["max_virtual_user_count"] = jobs.MaxVirtualUserCount
					}

					if jobs.Note != nil {
						jobsMap["note"] = jobs.Note
					}

					if jobs.ErrorRate != nil {
						jobsMap["error_rate"] = jobs.ErrorRate
					}

					if jobs.JobOwner != nil {
						jobsMap["job_owner"] = jobs.JobOwner
					}

					if jobs.LoadSources != nil {
						loadSourcesMap := map[string]interface{}{}

						if jobs.LoadSources.IP != nil {
							loadSourcesMap["ip"] = jobs.LoadSources.IP
						}

						if jobs.LoadSources.PodName != nil {
							loadSourcesMap["pod_name"] = jobs.LoadSources.PodName
						}

						if jobs.LoadSources.Region != nil {
							loadSourcesMap["region"] = jobs.LoadSources.Region
						}

						jobsMap["load_sources"] = []interface{}{loadSourcesMap}
					}

					if jobs.Duration != nil {
						jobsMap["duration"] = jobs.Duration
					}

					if jobs.MaxRequestsPerSecond != nil {
						jobsMap["max_requests_per_second"] = jobs.MaxRequestsPerSecond
					}

					if jobs.RequestTotal != nil {
						jobsMap["request_total"] = jobs.RequestTotal
					}

					if jobs.RequestsPerSecond != nil {
						jobsMap["requests_per_second"] = jobs.RequestsPerSecond
					}

					if jobs.ResponseTimeAverage != nil {
						jobsMap["response_time_average"] = jobs.ResponseTimeAverage
					}

					if jobs.ResponseTimeP99 != nil {
						jobsMap["response_time_p99"] = jobs.ResponseTimeP99
					}

					if jobs.ResponseTimeP95 != nil {
						jobsMap["response_time_p95"] = jobs.ResponseTimeP95
					}

					if jobs.ResponseTimeP90 != nil {
						jobsMap["response_time_p90"] = jobs.ResponseTimeP90
					}

					if jobs.Scripts != nil {
						jobsMap["scripts"] = jobs.Scripts
					}

					if jobs.ResponseTimeMax != nil {
						jobsMap["response_time_max"] = jobs.ResponseTimeMax
					}

					if jobs.ResponseTimeMin != nil {
						jobsMap["response_time_min"] = jobs.ResponseTimeMin
					}

					if jobs.LoadSourceInfos != nil {
						loadSourceInfosList := []interface{}{}
						for _, loadSourceInfos := range jobs.LoadSourceInfos {
							loadSourceInfosMap := map[string]interface{}{}

							if loadSourceInfos.IP != nil {
								loadSourceInfosMap["ip"] = loadSourceInfos.IP
							}

							if loadSourceInfos.PodName != nil {
								loadSourceInfosMap["pod_name"] = loadSourceInfos.PodName
							}

							if loadSourceInfos.Region != nil {
								loadSourceInfosMap["region"] = loadSourceInfos.Region
							}

							loadSourceInfosList = append(loadSourceInfosList, loadSourceInfosMap)
						}

						jobsMap["load_source_infos"] = loadSourceInfosList
					}

					if jobs.TestScripts != nil {
						testScriptsList := []interface{}{}
						for _, testScripts := range jobs.TestScripts {
							testScriptsMap := map[string]interface{}{}

							if testScripts.Name != nil {
								testScriptsMap["name"] = testScripts.Name
							}

							if testScripts.Size != nil {
								testScriptsMap["size"] = testScripts.Size
							}

							if testScripts.Type != nil {
								testScriptsMap["type"] = testScripts.Type
							}

							if testScripts.UpdatedAt != nil {
								testScriptsMap["updated_at"] = testScripts.UpdatedAt
							}

							if testScripts.EncodedContent != nil {
								testScriptsMap["encoded_content"] = testScripts.EncodedContent
							}

							if testScripts.EncodedHttpArchive != nil {
								testScriptsMap["encoded_http_archive"] = testScripts.EncodedHttpArchive
							}

							if testScripts.LoadWeight != nil {
								testScriptsMap["load_weight"] = testScripts.LoadWeight
							}

							if testScripts.FileId != nil {
								testScriptsMap["file_id"] = testScripts.FileId
							}

							testScriptsList = append(testScriptsList, testScriptsMap)
						}

						jobsMap["test_scripts"] = testScriptsList
					}

					if jobs.Protocols != nil {
						protocolsList := []interface{}{}
						for _, protocols := range jobs.Protocols {
							protocolsMap := map[string]interface{}{}

							if protocols.Name != nil {
								protocolsMap["name"] = protocols.Name
							}

							if protocols.Size != nil {
								protocolsMap["size"] = protocols.Size
							}

							if protocols.Type != nil {
								protocolsMap["type"] = protocols.Type
							}

							if protocols.UpdatedAt != nil {
								protocolsMap["updated_at"] = protocols.UpdatedAt
							}

							if protocols.FileId != nil {
								protocolsMap["file_id"] = protocols.FileId
							}

							protocolsList = append(protocolsList, protocolsMap)
						}

						jobsMap["protocols"] = protocolsList
					}

					if jobs.RequestFiles != nil {
						requestFilesList := []interface{}{}
						for _, requestFiles := range jobs.RequestFiles {
							requestFilesMap := map[string]interface{}{}

							if requestFiles.Name != nil {
								requestFilesMap["name"] = requestFiles.Name
							}

							if requestFiles.Size != nil {
								requestFilesMap["size"] = requestFiles.Size
							}

							if requestFiles.Type != nil {
								requestFilesMap["type"] = requestFiles.Type
							}

							if requestFiles.UpdatedAt != nil {
								requestFilesMap["updated_at"] = requestFiles.UpdatedAt
							}

							if requestFiles.FileId != nil {
								requestFilesMap["file_id"] = requestFiles.FileId
							}

							requestFilesList = append(requestFilesList, requestFilesMap)
						}

						jobsMap["request_files"] = requestFilesList
					}

					if jobs.Plugins != nil {
						pluginsList := []interface{}{}
						for _, plugins := range jobs.Plugins {
							pluginsMap := map[string]interface{}{}

							if plugins.Name != nil {
								pluginsMap["name"] = plugins.Name
							}

							if plugins.Size != nil {
								pluginsMap["size"] = plugins.Size
							}

							if plugins.Type != nil {
								pluginsMap["type"] = plugins.Type
							}

							if plugins.UpdatedAt != nil {
								pluginsMap["updated_at"] = plugins.UpdatedAt
							}

							if plugins.FileId != nil {
								pluginsMap["file_id"] = plugins.FileId
							}

							pluginsList = append(pluginsList, pluginsMap)
						}

						jobsMap["plugins"] = pluginsList
					}

					if jobs.CronId != nil {
						jobsMap["cron_id"] = jobs.CronId
					}

					if jobs.Type != nil {
						jobsMap["type"] = jobs.Type
					}

					if jobs.DomainNameConfig != nil {
						domainNameConfigMap := map[string]interface{}{}

						if jobs.DomainNameConfig.HostAliases != nil {
							hostAliasesList := []interface{}{}
							for _, hostAliases := range jobs.DomainNameConfig.HostAliases {
								hostAliasesMap := map[string]interface{}{}

								if hostAliases.HostNames != nil {
									hostAliasesMap["host_names"] = hostAliases.HostNames
								}

								if hostAliases.IP != nil {
									hostAliasesMap["ip"] = hostAliases.IP
								}

								hostAliasesList = append(hostAliasesList, hostAliasesMap)
							}

							domainNameConfigMap["host_aliases"] = hostAliasesList
						}

						if jobs.DomainNameConfig.DNSConfig != nil {
							dNSConfigMap := map[string]interface{}{}

							if jobs.DomainNameConfig.DNSConfig.Nameservers != nil {
								dNSConfigMap["nameservers"] = jobs.DomainNameConfig.DNSConfig.Nameservers
							}

							domainNameConfigMap["dns_config"] = []interface{}{dNSConfigMap}
						}

						jobsMap["domain_name_config"] = []interface{}{domainNameConfigMap}
					}

					if jobs.Debug != nil {
						jobsMap["debug"] = jobs.Debug
					}

					if jobs.AbortReason != nil {
						jobsMap["abort_reason"] = jobs.AbortReason
					}

					if jobs.CreatedAt != nil {
						jobsMap["created_at"] = jobs.CreatedAt
					}

					if jobs.ProjectId != nil {
						jobsMap["project_id"] = jobs.ProjectId
					}

					if jobs.NotificationHooks != nil {
						notificationHooksList := []interface{}{}
						for _, notificationHooks := range jobs.NotificationHooks {
							notificationHooksMap := map[string]interface{}{}

							if notificationHooks.Events != nil {
								notificationHooksMap["events"] = notificationHooks.Events
							}

							if notificationHooks.URL != nil {
								notificationHooksMap["url"] = notificationHooks.URL
							}

							notificationHooksList = append(notificationHooksList, notificationHooksMap)
						}

						jobsMap["notification_hooks"] = notificationHooksList
					}

					if jobs.NetworkReceiveRate != nil {
						jobsMap["network_receive_rate"] = jobs.NetworkReceiveRate
					}

					if jobs.NetworkSendRate != nil {
						jobsMap["network_send_rate"] = jobs.NetworkSendRate
					}

					if jobs.Message != nil {
						jobsMap["message"] = jobs.Message
					}

					if jobs.ProjectName != nil {
						jobsMap["project_name"] = jobs.ProjectName
					}

					if jobs.ScenarioName != nil {
						jobsMap["scenario_name"] = jobs.ScenarioName
					}

					jobsList = append(jobsList, jobsMap)
				}

				scenarioWithJobsMap["jobs"] = jobsList
			}

			tmpList = append(tmpList, scenarioWithJobsMap)
		}

		_ = d.Set("scenario_with_jobs_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
