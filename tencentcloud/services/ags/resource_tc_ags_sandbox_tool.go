package ags

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAgsSandboxTool() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAgsSandboxToolCreate,
		Read:   resourceTencentCloudAgsSandboxToolRead,
		Update: resourceTencentCloudAgsSandboxToolUpdate,
		Delete: resourceTencentCloudAgsSandboxToolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tool_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sandbox tool name, length 1-50 characters, supports Chinese, English, numbers, underscores. Must be unique under the same AppId.",
			},
			"tool_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sandbox tool type. Currently supports: browser, code-interpreter, custom, mobile, osworld, swebench, aio, android-world, waa.",
			},
			"network_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Network configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Network mode. Valid values: PUBLIC, VPC.",
						},
						"vpc_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "VPC network configuration. Required when network_mode is VPC.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_ids": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "VPC subnet IDs.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"security_group_ids": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Security group IDs.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sandbox tool description, max length 200 characters.",
			},
			"default_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Default timeout, supports formats: 5m, 300s, 1h etc. If not specified, system default (5 minutes) is used. Max 24 hours.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags for the sandbox tool, up to 50 tags.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Idempotency token, max length 64 characters.",
			},
			"role_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Role ARN bound to the sandbox tool.",
			},
			"storage_mounts": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Sandbox tool storage mount configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage mount configuration name.",
						},
						"storage_source": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Storage source configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cos": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "COS storage configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"endpoint": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "COS endpoint.",
												},
												"bucket_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "COS bucket name.",
												},
												"bucket_path": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "COS bucket path.",
												},
											},
										},
									},
									"image": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Image storage configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"reference": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Image reference address.",
												},
												"image_registry_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Image registry type: enterprise, personal, custom.",
												},
												"sub_path": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sub path in the image.",
												},
												"digest": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Image digest.",
												},
											},
										},
									},
									"cfs": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "CFS storage configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"file_system_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "CFS file system ID.",
												},
												"path": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "CFS mount path.",
												},
											},
										},
									},
									"agent_bucket": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "AgentBucket storage configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"library_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Library ID.",
												},
												"space_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Space ID.",
												},
												"access_domain": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Access domain.",
												},
											},
										},
									},
								},
							},
						},
						"mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Local mount path on the sandbox instance.",
						},
						"read_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Read-only mount permission, defaults to false.",
						},
						"volume": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Volume reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Volume name.",
									},
									"volume_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Volume ID.",
									},
								},
							},
						},
						"volume_template": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Volume template reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_template_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Volume template name.",
									},
									"reuse_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Reuse key for the volume template derivation.",
									},
								},
							},
						},
					},
				},
			},
			"custom_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Sandbox tool custom configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Image address.",
						},
						"image_registry_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Image registry type: enterprise, personal, custom.",
						},
						"command": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Startup command.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"args": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Startup arguments.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"env": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Environment variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Environment variable name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Environment variable value.",
									},
								},
							},
						},
						"ports": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Port configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Port name.",
									},
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Port number.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Port protocol.",
									},
								},
							},
						},
						"resources": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Resource configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "CPU resource specification.",
									},
									"memory": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Memory resource specification.",
									},
									"storage": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Storage resource specification.",
									},
								},
							},
						},
						"probe": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Probe configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_get": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "HTTP GET probe configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "HTTP path.",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "HTTP port.",
												},
												"scheme": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "HTTP scheme.",
												},
											},
										},
									},
									"ready_timeout_ms": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Ready timeout in milliseconds.",
									},
									"probe_timeout_ms": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Probe timeout in milliseconds.",
									},
									"probe_period_ms": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Probe period in milliseconds.",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Success threshold.",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Failure threshold.",
									},
								},
							},
						},
						"dns_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "DNS configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"servers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "DNS server addresses.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"searches": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "DNS search domains.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"options": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "DNS options.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"security_context": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Security context.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"privileged": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to enable privileged mode. Default: true.",
									},
									"capabilities": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "POSIX capabilities for the container.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"add": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Capabilities to add.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"drop": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Capabilities to drop.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"computer_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Computer environment sandbox configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"waa_configuration": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "WAA configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "WAA image ID.",
									},
								},
							},
						},
						"os_world_configuration": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "OSWorld configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "OSWorld version.",
									},
								},
							},
						},
					},
				},
			},
			"log_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Sandbox tool log push configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cls_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "CLS log configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CLS topic ID.",
									},
								},
							},
						},
						"log_sources": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Log source configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"files": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Log file paths to collect, must be in /logs/ directory, max 10 files.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"persistent": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Persistent sandbox flag.",
			},

			// computed fields
			"tool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sandbox tool ID.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sandbox tool status. Values: CREATING, ACTIVE, DELETING, FAILED.",
			},
			"default_timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Default timeout in seconds.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sandbox tool creation time, ISO8601 format.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sandbox tool update time, ISO8601 format.",
			},
			"status_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason for the current status.",
			},
			"image_digest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image digest of the custom configuration.",
			},
		},
	}
}

func resourceTencentCloudAgsSandboxToolCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ags_sandbox_tool.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = agsv20250920.NewCreateSandboxToolRequest()
		response = agsv20250920.NewCreateSandboxToolResponse()
	)

	if v, ok := d.GetOk("tool_name"); ok {
		request.ToolName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tool_type"); ok {
		request.ToolType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_configuration"); ok {
		request.NetworkConfiguration = expandNetworkConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("default_timeout"); ok {
		request.DefaultTimeout = helper.String(v.(string))
	}

	if _, ok := d.GetOk("tags"); ok {
		tags := helper.GetTags(d, "tags")
		for k, v := range tags {
			tag := agsv20250920.Tag{
				Key:   helper.String(k),
				Value: helper.String(v),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_arn"); ok {
		request.RoleArn = helper.String(v.(string))
	}

	if v, ok := d.GetOk("storage_mounts"); ok {
		request.StorageMounts = expandStorageMounts(v.([]interface{}))
	}

	if v, ok := d.GetOk("custom_configuration"); ok {
		request.CustomConfiguration = expandCustomConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("computer_configuration"); ok {
		request.ComputerConfiguration = expandComputerConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("log_configuration"); ok {
		request.LogConfiguration = expandLogConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("persistent"); ok {
		request.Persistent = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAgsV20250920Client().CreateSandboxToolWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create sandbox_tool failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create sandbox_tool failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create sandbox_tool, d.Id() = %s, logId = %s", logId, d.Id(), logId)

	if response.Response.ToolId == nil || *response.Response.ToolId == "" {
		return fmt.Errorf("ToolId is nil or empty.")
	}

	d.SetId(*response.Response.ToolId)
	return resourceTencentCloudAgsSandboxToolRead(d, meta)
}

func resourceTencentCloudAgsSandboxToolRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ags_sandbox_tool.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = AgsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		toolId  = d.Id()
	)

	respData, err := service.DescribeSandboxToolById(ctx, toolId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ags_sandbox_tool` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ToolId != nil {
		_ = d.Set("tool_id", respData.ToolId)
	}

	if respData.ToolName != nil {
		_ = d.Set("tool_name", respData.ToolName)
	}

	if respData.ToolType != nil {
		_ = d.Set("tool_type", respData.ToolType)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Persistent != nil {
		_ = d.Set("persistent", respData.Persistent)
	}

	if respData.DefaultTimeoutSeconds != nil {
		_ = d.Set("default_timeout_seconds", respData.DefaultTimeoutSeconds)
	}

	if respData.NetworkConfiguration != nil {
		_ = d.Set("network_configuration", flattenNetworkConfiguration(respData.NetworkConfiguration))
	}

	if respData.Tags != nil {
		tagsMap := make(map[string]interface{})
		for _, tag := range respData.Tags {
			if tag.Key != nil && tag.Value != nil {
				tagsMap[*tag.Key] = *tag.Value
			}
		}
		_ = d.Set("tags", tagsMap)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.RoleArn != nil {
		_ = d.Set("role_arn", respData.RoleArn)
	}

	if respData.StorageMounts != nil {
		_ = d.Set("storage_mounts", flattenStorageMounts(respData.StorageMounts))
	}

	if respData.CustomConfiguration != nil {
		_ = d.Set("custom_configuration", flattenCustomConfigurationDetail(respData.CustomConfiguration))
	}

	if respData.LogConfiguration != nil {
		_ = d.Set("log_configuration", flattenLogConfiguration(respData.LogConfiguration))
	}

	if respData.ComputerConfiguration != nil {
		_ = d.Set("computer_configuration", flattenComputerConfiguration(respData.ComputerConfiguration))
	}

	if respData.StatusReason != nil {
		_ = d.Set("status_reason", respData.StatusReason)
	}

	return nil
}

func resourceTencentCloudAgsSandboxToolUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ags_sandbox_tool.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		toolId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"description", "network_configuration", "tags", "custom_configuration", "computer_configuration"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := agsv20250920.NewUpdateSandboxToolRequest()
		request.ToolId = &toolId

		if d.HasChange("description") {
			if v, ok := d.GetOk("description"); ok {
				request.Description = helper.String(v.(string))
			}
		}

		if d.HasChange("network_configuration") {
			if v, ok := d.GetOk("network_configuration"); ok {
				request.NetworkConfiguration = expandNetworkConfiguration(v.([]interface{}))
			}
		}

		if d.HasChange("tags") {
			if _, ok := d.GetOk("tags"); ok {
				tags := helper.GetTags(d, "tags")
				for k, v := range tags {
					tag := agsv20250920.Tag{
						Key:   helper.String(k),
						Value: helper.String(v),
					}
					request.Tags = append(request.Tags, &tag)
				}
			}
		}

		if d.HasChange("custom_configuration") {
			if v, ok := d.GetOk("custom_configuration"); ok {
				request.CustomConfiguration = expandCustomConfiguration(v.([]interface{}))
			}
		}

		if d.HasChange("computer_configuration") {
			if v, ok := d.GetOk("computer_configuration"); ok {
				request.ComputerConfiguration = expandComputerConfiguration(v.([]interface{}))
			}
		}

		reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAgsV20250920Client().UpdateSandboxToolWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update sandbox_tool failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudAgsSandboxToolRead(d, meta)
}

func resourceTencentCloudAgsSandboxToolDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ags_sandbox_tool.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = agsv20250920.NewDeleteSandboxToolRequest()
		toolId  = d.Id()
	)

	request.ToolId = &toolId
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAgsV20250920Client().DeleteSandboxToolWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete sandbox_tool failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// expand functions
func expandNetworkConfiguration(v []interface{}) *agsv20250920.NetworkConfiguration {
	if len(v) == 0 {
		return nil
	}
	config := agsv20250920.NetworkConfiguration{}
	data := v[0].(map[string]interface{})
	if val, ok := data["network_mode"].(string); ok && val != "" {
		config.NetworkMode = helper.String(val)
	}
	if val, ok := data["vpc_config"]; ok {
		vpcList := val.([]interface{})
		if len(vpcList) > 0 {
			vpcConfig := agsv20250920.VPCConfig{}
			vpcData := vpcList[0].(map[string]interface{})
			if subnets, ok := vpcData["subnet_ids"].(*schema.Set); ok {
				for _, subnet := range subnets.List() {
					vpcConfig.SubnetIds = append(vpcConfig.SubnetIds, helper.String(subnet.(string)))
				}
			}
			if sgs, ok := vpcData["security_group_ids"].(*schema.Set); ok {
				for _, sg := range sgs.List() {
					vpcConfig.SecurityGroupIds = append(vpcConfig.SecurityGroupIds, helper.String(sg.(string)))
				}
			}
			config.VpcConfig = &vpcConfig
		}
	}
	return &config
}

func expandStorageMounts(v []interface{}) []*agsv20250920.StorageMount {
	if len(v) == 0 {
		return nil
	}
	result := make([]*agsv20250920.StorageMount, 0, len(v))
	for _, item := range v {
		data := item.(map[string]interface{})
		mount := agsv20250920.StorageMount{}
		if val, ok := data["name"].(string); ok && val != "" {
			mount.Name = helper.String(val)
		}
		if val, ok := data["storage_source"]; ok {
			storageSourceList := val.([]interface{})
			if len(storageSourceList) > 0 {
				mount.StorageSource = expandStorageSource(storageSourceList[0].(map[string]interface{}))
			}
		}
		if val, ok := data["mount_path"].(string); ok && val != "" {
			mount.MountPath = helper.String(val)
		}
		if val, ok := data["read_only"].(bool); ok {
			mount.ReadOnly = helper.Bool(val)
		}
		if val, ok := data["volume"]; ok {
			volumeList := val.([]interface{})
			if len(volumeList) > 0 {
				volumeData := volumeList[0].(map[string]interface{})
				vol := agsv20250920.VolumeRef{}
				if vn, ok := volumeData["volume_name"].(string); ok && vn != "" {
					vol.VolumeName = helper.String(vn)
				}
				if vi, ok := volumeData["volume_id"].(string); ok && vi != "" {
					vol.VolumeId = helper.String(vi)
				}
				mount.Volume = &vol
			}
		}
		if val, ok := data["volume_template"]; ok {
			vtList := val.([]interface{})
			if len(vtList) > 0 {
				vtData := vtList[0].(map[string]interface{})
				vt := agsv20250920.VolumeTemplateRef{}
				if vtn, ok := vtData["volume_template_name"].(string); ok && vtn != "" {
					vt.VolumeTemplateName = helper.String(vtn)
				}
				if rk, ok := vtData["reuse_key"].(string); ok && rk != "" {
					vt.ReuseKey = helper.String(rk)
				}
				mount.VolumeTemplate = &vt
			}
		}
		result = append(result, &mount)
	}
	return result
}

func expandStorageSource(data map[string]interface{}) *agsv20250920.StorageSource {
	source := agsv20250920.StorageSource{}
	if val, ok := data["cos"]; ok {
		cosList := val.([]interface{})
		if len(cosList) > 0 {
			cosData := cosList[0].(map[string]interface{})
			cos := agsv20250920.CosStorageSource{}
			if endpoint, ok := cosData["endpoint"].(string); ok && endpoint != "" {
				cos.Endpoint = helper.String(endpoint)
			}
			if bucket, ok := cosData["bucket_name"].(string); ok && bucket != "" {
				cos.BucketName = helper.String(bucket)
			}
			if path, ok := cosData["bucket_path"].(string); ok && path != "" {
				cos.BucketPath = helper.String(path)
			}
			source.Cos = &cos
		}
	}
	if val, ok := data["image"]; ok {
		imageList := val.([]interface{})
		if len(imageList) > 0 {
			imageData := imageList[0].(map[string]interface{})
			image := agsv20250920.ImageStorageSource{}
			if ref, ok := imageData["reference"].(string); ok && ref != "" {
				image.Reference = helper.String(ref)
			}
			if irt, ok := imageData["image_registry_type"].(string); ok && irt != "" {
				image.ImageRegistryType = helper.String(irt)
			}
			if sp, ok := imageData["sub_path"].(string); ok && sp != "" {
				image.SubPath = helper.String(sp)
			}
			source.Image = &image
		}
	}
	if val, ok := data["cfs"]; ok {
		cfsList := val.([]interface{})
		if len(cfsList) > 0 {
			cfsData := cfsList[0].(map[string]interface{})
			cfs := agsv20250920.CfsStorageSource{}
			if fsId, ok := cfsData["file_system_id"].(string); ok && fsId != "" {
				cfs.FileSystemId = helper.String(fsId)
			}
			if path, ok := cfsData["path"].(string); ok && path != "" {
				cfs.Path = helper.String(path)
			}
			source.Cfs = &cfs
		}
	}
	if val, ok := data["agent_bucket"]; ok {
		abList := val.([]interface{})
		if len(abList) > 0 {
			abData := abList[0].(map[string]interface{})
			ab := agsv20250920.AgentBucketStorageSource{}
			if libId, ok := abData["library_id"].(string); ok && libId != "" {
				ab.LibraryId = helper.String(libId)
			}
			if spaceId, ok := abData["space_id"].(string); ok && spaceId != "" {
				ab.SpaceId = helper.String(spaceId)
			}
			if domain, ok := abData["access_domain"].(string); ok && domain != "" {
				ab.AccessDomain = helper.String(domain)
			}
			source.AgentBucket = &ab
		}
	}
	return &source
}

func expandCustomConfiguration(v []interface{}) *agsv20250920.CustomConfiguration {
	if len(v) == 0 {
		return nil
	}
	config := agsv20250920.CustomConfiguration{}
	data := v[0].(map[string]interface{})
	if val, ok := data["image"].(string); ok && val != "" {
		config.Image = helper.String(val)
	}
	if val, ok := data["image_registry_type"].(string); ok && val != "" {
		config.ImageRegistryType = helper.String(val)
	}
	if val, ok := data["command"].([]interface{}); ok {
		for _, cmd := range val {
			config.Command = append(config.Command, helper.String(cmd.(string)))
		}
	}
	if val, ok := data["args"].([]interface{}); ok {
		for _, arg := range val {
			config.Args = append(config.Args, helper.String(arg.(string)))
		}
	}
	if val, ok := data["env"].(*schema.Set); ok {
		for _, envItem := range val.List() {
			envData := envItem.(map[string]interface{})
			env := agsv20250920.EnvVar{}
			if name, ok := envData["name"].(string); ok && name != "" {
				env.Name = helper.String(name)
			}
			if value, ok := envData["value"].(string); ok && value != "" {
				env.Value = helper.String(value)
			}
			config.Env = append(config.Env, &env)
		}
	}
	if val, ok := data["ports"]; ok {
		portsList := val.([]interface{})
		for _, portItem := range portsList {
			portData := portItem.(map[string]interface{})
			port := agsv20250920.PortConfiguration{}
			if name, ok := portData["name"].(string); ok && name != "" {
				port.Name = helper.String(name)
			}
			if p, ok := portData["port"].(int); ok {
				port.Port = helper.IntInt64(p)
			}
			if proto, ok := portData["protocol"].(string); ok && proto != "" {
				port.Protocol = helper.String(proto)
			}
			config.Ports = append(config.Ports, &port)
		}
	}
	if val, ok := data["resources"]; ok {
		resList := val.([]interface{})
		if len(resList) > 0 {
			resData := resList[0].(map[string]interface{})
			res := agsv20250920.ResourceConfiguration{}
			if cpu, ok := resData["cpu"].(string); ok && cpu != "" {
				res.CPU = helper.String(cpu)
			}
			if mem, ok := resData["memory"].(string); ok && mem != "" {
				res.Memory = helper.String(mem)
			}
			if storage, ok := resData["storage"].(string); ok && storage != "" {
				res.Storage = helper.String(storage)
			}
			config.Resources = &res
		}
	}
	if val, ok := data["probe"]; ok {
		probeList := val.([]interface{})
		if len(probeList) > 0 {
			config.Probe = expandProbeConfiguration(probeList[0].(map[string]interface{}))
		}
	}
	if val, ok := data["dns_config"]; ok {
		dnsList := val.([]interface{})
		if len(dnsList) > 0 {
			config.DNSConfig = expandDNSConfig(dnsList[0].(map[string]interface{}))
		}
	}
	if val, ok := data["security_context"]; ok {
		scList := val.([]interface{})
		if len(scList) > 0 {
			config.SecurityContext = expandSecurityContext(scList[0].(map[string]interface{}))
		}
	}
	return &config
}

func expandProbeConfiguration(data map[string]interface{}) *agsv20250920.ProbeConfiguration {
	probe := agsv20250920.ProbeConfiguration{}
	if val, ok := data["http_get"]; ok {
		hgList := val.([]interface{})
		if len(hgList) > 0 {
			hgData := hgList[0].(map[string]interface{})
			hg := agsv20250920.HttpGetAction{}
			if path, ok := hgData["path"].(string); ok && path != "" {
				hg.Path = helper.String(path)
			}
			if port, ok := hgData["port"].(int); ok {
				hg.Port = helper.IntInt64(port)
			}
			if scheme, ok := hgData["scheme"].(string); ok && scheme != "" {
				hg.Scheme = helper.String(scheme)
			}
			probe.HttpGet = &hg
		}
	}
	if val, ok := data["ready_timeout_ms"].(int); ok && val != 0 {
		probe.ReadyTimeoutMs = helper.IntInt64(val)
	}
	if val, ok := data["probe_timeout_ms"].(int); ok && val != 0 {
		probe.ProbeTimeoutMs = helper.IntInt64(val)
	}
	if val, ok := data["probe_period_ms"].(int); ok && val != 0 {
		probe.ProbePeriodMs = helper.IntInt64(val)
	}
	if val, ok := data["success_threshold"].(int); ok && val != 0 {
		probe.SuccessThreshold = helper.IntInt64(val)
	}
	if val, ok := data["failure_threshold"].(int); ok && val != 0 {
		probe.FailureThreshold = helper.IntInt64(val)
	}
	return &probe
}

func expandDNSConfig(data map[string]interface{}) *agsv20250920.DNSConfig {
	dns := agsv20250920.DNSConfig{}
	if val, ok := data["servers"].([]interface{}); ok {
		for _, s := range val {
			dns.Servers = append(dns.Servers, helper.String(s.(string)))
		}
	}
	if val, ok := data["searches"].([]interface{}); ok {
		for _, s := range val {
			dns.Searches = append(dns.Searches, helper.String(s.(string)))
		}
	}
	if val, ok := data["options"].([]interface{}); ok {
		for _, s := range val {
			dns.Options = append(dns.Options, helper.String(s.(string)))
		}
	}
	return &dns
}

func expandSecurityContext(data map[string]interface{}) *agsv20250920.SecurityContext {
	sc := agsv20250920.SecurityContext{}
	if val, ok := data["privileged"].(bool); ok {
		sc.Privileged = helper.Bool(val)
	}
	if val, ok := data["capabilities"]; ok {
		capList := val.([]interface{})
		if len(capList) > 0 {
			capData := capList[0].(map[string]interface{})
			caps := agsv20250920.Capabilities{}
			if add, ok := capData["add"].([]interface{}); ok {
				for _, a := range add {
					caps.Add = append(caps.Add, helper.String(a.(string)))
				}
			}
			if drop, ok := capData["drop"].([]interface{}); ok {
				for _, d := range drop {
					caps.Drop = append(caps.Drop, helper.String(d.(string)))
				}
			}
			sc.Capabilities = &caps
		}
	}
	return &sc
}

func expandComputerConfiguration(v []interface{}) *agsv20250920.ComputerConfiguration {
	if len(v) == 0 {
		return nil
	}
	config := agsv20250920.ComputerConfiguration{}
	data := v[0].(map[string]interface{})
	if val, ok := data["waa_configuration"]; ok {
		waaList := val.([]interface{})
		if len(waaList) > 0 {
			waaData := waaList[0].(map[string]interface{})
			waa := agsv20250920.WAAConfiguration{}
			if imageId, ok := waaData["image_id"].(string); ok && imageId != "" {
				waa.ImageId = helper.String(imageId)
			}
			config.WAAConfiguration = &waa
		}
	}
	if val, ok := data["os_world_configuration"]; ok {
		oswList := val.([]interface{})
		if len(oswList) > 0 {
			oswData := oswList[0].(map[string]interface{})
			osw := agsv20250920.OSWorldConfiguration{}
			if version, ok := oswData["version"].(string); ok && version != "" {
				osw.Version = helper.String(version)
			}
			config.OSWorldConfiguration = &osw
		}
	}
	return &config
}

func expandLogConfiguration(v []interface{}) *agsv20250920.LogConfiguration {
	if len(v) == 0 {
		return nil
	}
	config := agsv20250920.LogConfiguration{}
	data := v[0].(map[string]interface{})
	if val, ok := data["cls_config"]; ok {
		clsList := val.([]interface{})
		if len(clsList) > 0 {
			clsData := clsList[0].(map[string]interface{})
			cls := agsv20250920.CLSConfig{}
			if topicId, ok := clsData["topic_id"].(string); ok && topicId != "" {
				cls.TopicId = helper.String(topicId)
			}
			config.CLSConfig = &cls
		}
	}
	if val, ok := data["log_sources"]; ok {
		lsList := val.([]interface{})
		if len(lsList) > 0 {
			lsData := lsList[0].(map[string]interface{})
			ls := agsv20250920.LogSources{}
			if files, ok := lsData["files"].([]interface{}); ok {
				for _, f := range files {
					ls.Files = append(ls.Files, helper.String(f.(string)))
				}
			}
			config.LogSources = &ls
		}
	}
	return &config
}

// flatten functions
func flattenNetworkConfiguration(config *agsv20250920.NetworkConfiguration) []interface{} {
	if config == nil {
		return nil
	}
	result := map[string]interface{}{}
	if config.NetworkMode != nil {
		result["network_mode"] = *config.NetworkMode
	}
	if config.VpcConfig != nil {
		vpcList := make([]map[string]interface{}, 0)
		vpcMap := map[string]interface{}{}
		if config.VpcConfig.SubnetIds != nil {
			vpcMap["subnet_ids"] = helper.StringsInterfaces(config.VpcConfig.SubnetIds)
		}
		if config.VpcConfig.SecurityGroupIds != nil {
			vpcMap["security_group_ids"] = helper.StringsInterfaces(config.VpcConfig.SecurityGroupIds)
		}
		vpcList = append(vpcList, vpcMap)
		result["vpc_config"] = vpcList
	}
	return []interface{}{result}
}

func flattenStorageMounts(mounts []*agsv20250920.StorageMount) []interface{} {
	result := make([]interface{}, 0, len(mounts))
	for _, mount := range mounts {
		m := map[string]interface{}{}
		if mount.Name != nil {
			m["name"] = *mount.Name
		}
		if mount.StorageSource != nil {
			m["storage_source"] = flattenStorageSource(mount.StorageSource)
		}
		if mount.MountPath != nil {
			m["mount_path"] = *mount.MountPath
		}
		if mount.ReadOnly != nil {
			m["read_only"] = *mount.ReadOnly
		}
		if mount.Volume != nil {
			volList := make([]map[string]interface{}, 0)
			volMap := map[string]interface{}{}
			if mount.Volume.VolumeName != nil {
				volMap["volume_name"] = *mount.Volume.VolumeName
			}
			if mount.Volume.VolumeId != nil {
				volMap["volume_id"] = *mount.Volume.VolumeId
			}
			volList = append(volList, volMap)
			m["volume"] = volList
		}
		if mount.VolumeTemplate != nil {
			vtList := make([]map[string]interface{}, 0)
			vtMap := map[string]interface{}{}
			if mount.VolumeTemplate.VolumeTemplateName != nil {
				vtMap["volume_template_name"] = *mount.VolumeTemplate.VolumeTemplateName
			}
			if mount.VolumeTemplate.ReuseKey != nil {
				vtMap["reuse_key"] = *mount.VolumeTemplate.ReuseKey
			}
			vtList = append(vtList, vtMap)
			m["volume_template"] = vtList
		}
		result = append(result, m)
	}
	return result
}

func flattenStorageSource(source *agsv20250920.StorageSource) []interface{} {
	if source == nil {
		return nil
	}
	result := map[string]interface{}{}
	if source.Cos != nil {
		cosList := make([]map[string]interface{}, 0)
		cosMap := map[string]interface{}{}
		if source.Cos.Endpoint != nil {
			cosMap["endpoint"] = *source.Cos.Endpoint
		}
		if source.Cos.BucketName != nil {
			cosMap["bucket_name"] = *source.Cos.BucketName
		}
		if source.Cos.BucketPath != nil {
			cosMap["bucket_path"] = *source.Cos.BucketPath
		}
		cosList = append(cosList, cosMap)
		result["cos"] = cosList
	}
	if source.Image != nil {
		imageList := make([]map[string]interface{}, 0)
		imageMap := map[string]interface{}{}
		if source.Image.Reference != nil {
			imageMap["reference"] = *source.Image.Reference
		}
		if source.Image.ImageRegistryType != nil {
			imageMap["image_registry_type"] = *source.Image.ImageRegistryType
		}
		if source.Image.SubPath != nil {
			imageMap["sub_path"] = *source.Image.SubPath
		}
		if source.Image.Digest != nil {
			imageMap["digest"] = *source.Image.Digest
		}
		imageList = append(imageList, imageMap)
		result["image"] = imageList
	}
	if source.Cfs != nil {
		cfsList := make([]map[string]interface{}, 0)
		cfsMap := map[string]interface{}{}
		if source.Cfs.FileSystemId != nil {
			cfsMap["file_system_id"] = *source.Cfs.FileSystemId
		}
		if source.Cfs.Path != nil {
			cfsMap["path"] = *source.Cfs.Path
		}
		cfsList = append(cfsList, cfsMap)
		result["cfs"] = cfsList
	}
	if source.AgentBucket != nil {
		abList := make([]map[string]interface{}, 0)
		abMap := map[string]interface{}{}
		if source.AgentBucket.LibraryId != nil {
			abMap["library_id"] = *source.AgentBucket.LibraryId
		}
		if source.AgentBucket.SpaceId != nil {
			abMap["space_id"] = *source.AgentBucket.SpaceId
		}
		if source.AgentBucket.AccessDomain != nil {
			abMap["access_domain"] = *source.AgentBucket.AccessDomain
		}
		abList = append(abList, abMap)
		result["agent_bucket"] = abList
	}
	return []interface{}{result}
}

func flattenCustomConfigurationDetail(config *agsv20250920.CustomConfigurationDetail) []interface{} {
	if config == nil {
		return nil
	}
	result := map[string]interface{}{}
	if config.Image != nil {
		result["image"] = *config.Image
	}
	if config.ImageRegistryType != nil {
		result["image_registry_type"] = *config.ImageRegistryType
	}
	if config.ImageDigest != nil {
		result["image_digest"] = *config.ImageDigest
	}
	if config.Command != nil {
		result["command"] = helper.StringsInterfaces(config.Command)
	}
	if config.Args != nil {
		result["args"] = helper.StringsInterfaces(config.Args)
	}
	if config.Env != nil {
		envList := make([]map[string]interface{}, 0, len(config.Env))
		for _, env := range config.Env {
			envMap := map[string]interface{}{}
			if env.Name != nil {
				envMap["name"] = *env.Name
			}
			if env.Value != nil {
				envMap["value"] = *env.Value
			}
			envList = append(envList, envMap)
		}
		result["env"] = envList
	}
	if config.Ports != nil {
		portsList := make([]map[string]interface{}, 0, len(config.Ports))
		for _, port := range config.Ports {
			portMap := map[string]interface{}{}
			if port.Name != nil {
				portMap["name"] = *port.Name
			}
			if port.Port != nil {
				portMap["port"] = *port.Port
			}
			if port.Protocol != nil {
				portMap["protocol"] = *port.Protocol
			}
			portsList = append(portsList, portMap)
		}
		result["ports"] = portsList
	}
	if config.Resources != nil {
		resList := make([]map[string]interface{}, 0)
		resMap := map[string]interface{}{}
		if config.Resources.CPU != nil {
			resMap["cpu"] = *config.Resources.CPU
		}
		if config.Resources.Memory != nil {
			resMap["memory"] = *config.Resources.Memory
		}
		if config.Resources.Storage != nil {
			resMap["storage"] = *config.Resources.Storage
		}
		resList = append(resList, resMap)
		result["resources"] = resList
	}
	if config.Probe != nil {
		result["probe"] = flattenProbeConfiguration(config.Probe)
	}
	if config.DNSConfig != nil {
		result["dns_config"] = flattenDNSConfig(config.DNSConfig)
	}
	if config.SecurityContext != nil {
		result["security_context"] = flattenSecurityContext(config.SecurityContext)
	}
	return []interface{}{result}
}

func flattenProbeConfiguration(probe *agsv20250920.ProbeConfiguration) []interface{} {
	if probe == nil {
		return nil
	}
	result := map[string]interface{}{}
	if probe.HttpGet != nil {
		hgList := make([]map[string]interface{}, 0)
		hgMap := map[string]interface{}{}
		if probe.HttpGet.Path != nil {
			hgMap["path"] = *probe.HttpGet.Path
		}
		if probe.HttpGet.Port != nil {
			hgMap["port"] = *probe.HttpGet.Port
		}
		if probe.HttpGet.Scheme != nil {
			hgMap["scheme"] = *probe.HttpGet.Scheme
		}
		hgList = append(hgList, hgMap)
		result["http_get"] = hgList
	}
	if probe.ReadyTimeoutMs != nil {
		result["ready_timeout_ms"] = *probe.ReadyTimeoutMs
	}
	if probe.ProbeTimeoutMs != nil {
		result["probe_timeout_ms"] = *probe.ProbeTimeoutMs
	}
	if probe.ProbePeriodMs != nil {
		result["probe_period_ms"] = *probe.ProbePeriodMs
	}
	if probe.SuccessThreshold != nil {
		result["success_threshold"] = *probe.SuccessThreshold
	}
	if probe.FailureThreshold != nil {
		result["failure_threshold"] = *probe.FailureThreshold
	}
	return []interface{}{result}
}

func flattenDNSConfig(dns *agsv20250920.DNSConfig) []interface{} {
	if dns == nil {
		return nil
	}
	result := map[string]interface{}{}
	if dns.Servers != nil {
		result["servers"] = helper.StringsInterfaces(dns.Servers)
	}
	if dns.Searches != nil {
		result["searches"] = helper.StringsInterfaces(dns.Searches)
	}
	if dns.Options != nil {
		result["options"] = helper.StringsInterfaces(dns.Options)
	}
	return []interface{}{result}
}

func flattenSecurityContext(sc *agsv20250920.SecurityContext) []interface{} {
	if sc == nil {
		return nil
	}
	result := map[string]interface{}{}
	if sc.Privileged != nil {
		result["privileged"] = *sc.Privileged
	}
	if sc.Capabilities != nil {
		capList := make([]map[string]interface{}, 0)
		capMap := map[string]interface{}{}
		if sc.Capabilities.Add != nil {
			capMap["add"] = helper.StringsInterfaces(sc.Capabilities.Add)
		}
		if sc.Capabilities.Drop != nil {
			capMap["drop"] = helper.StringsInterfaces(sc.Capabilities.Drop)
		}
		capList = append(capList, capMap)
		result["capabilities"] = capList
	}
	return []interface{}{result}
}

func flattenComputerConfiguration(config *agsv20250920.ComputerConfiguration) []interface{} {
	if config == nil {
		return nil
	}
	result := map[string]interface{}{}
	if config.WAAConfiguration != nil {
		waaList := make([]map[string]interface{}, 0)
		waaMap := map[string]interface{}{}
		if config.WAAConfiguration.ImageId != nil {
			waaMap["image_id"] = *config.WAAConfiguration.ImageId
		}
		waaList = append(waaList, waaMap)
		result["waa_configuration"] = waaList
	}
	if config.OSWorldConfiguration != nil {
		oswList := make([]map[string]interface{}, 0)
		oswMap := map[string]interface{}{}
		if config.OSWorldConfiguration.Version != nil {
			oswMap["version"] = *config.OSWorldConfiguration.Version
		}
		oswList = append(oswList, oswMap)
		result["os_world_configuration"] = oswList
	}
	return []interface{}{result}
}

func flattenLogConfiguration(config *agsv20250920.LogConfiguration) []interface{} {
	if config == nil {
		return nil
	}
	result := map[string]interface{}{}
	if config.CLSConfig != nil {
		clsList := make([]map[string]interface{}, 0)
		clsMap := map[string]interface{}{}
		if config.CLSConfig.TopicId != nil {
			clsMap["topic_id"] = *config.CLSConfig.TopicId
		}
		clsList = append(clsList, clsMap)
		result["cls_config"] = clsList
	}
	if config.LogSources != nil {
		lsList := make([]map[string]interface{}, 0)
		lsMap := map[string]interface{}{}
		if config.LogSources.Files != nil {
			lsMap["files"] = helper.StringsInterfaces(config.LogSources.Files)
		}
		lsList = append(lsList, lsMap)
		result["log_sources"] = lsList
	}
	return []interface{}{result}
}
