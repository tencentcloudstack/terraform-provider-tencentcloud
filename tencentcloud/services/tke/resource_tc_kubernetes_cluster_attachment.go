package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
)

func ResourceTencentCloudKubernetesClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterAttachmentCreate,
		Read:   resourceTencentCloudKubernetesClusterAttachmentRead,
		Delete: resourceTencentCloudKubernetesClusterAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the cluster.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CVM instance, this cvm will reinstall the system.",
			},

			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of Node image.",
			},

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				Description:  "Password to access, should be set if `key_ids` not set.",
				ValidateFunc: tccommon.ValidateAsConfigPassword,
			},

			"key_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The key pair to use for the instance, it looks like skey-16jig7tx, it should be set if `password` not set.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The host name of the attached instance. Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).",
			},

			"worker_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Deploy the machine configuration information of the 'WORKER', commonly used to attach existing instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Mount target. Default is not mounting.",
						},
						"docker_graph_path": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "/var/lib/docker",
							Description: "Docker graph path. Default is `/var/lib/docker`.",
						},
						"data_disk": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    11,
							Description: "Configurations of data disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Default:      "CLOUD_PREMIUM",
										Description:  "Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Default:     0,
										Description: "Volume of disk in GB. Default is `0`.",
									},
									"file_system": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "File system, e.g. `ext3/ext4/xfs`.",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate whether to auto format and mount or not. Default is `false`.",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "Mount target.",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The name of the device or partition to mount. NOTE: this argument doesn't support setting in node pool, or will leads to mount error.",
									},
								},
							},
						},
						"extra_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Custom parameter information related to the node. This is a white-list parameter.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Base64-encoded User Data text, the length limit is 16KB.",
						},
						"pre_start_user_script": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Base64-encoded user script, executed before initializing the node, currently only effective for adding existing nodes.",
						},
						"taints": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Node taint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Key of the taint.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Value of the taint.",
									},
									"effect": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Effect of the taint.",
									},
								},
							},
						},
						"is_schedule": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Deprecated:  "This argument was deprecated, use `unschedulable` instead.",
							Description: "Indicate to schedule the adding node or not. Default is true.",
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicate to set desired pod number in node. valid when the cluster is podCIDR.",
						},
						"gpu_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "GPU driver parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mig_enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether to enable MIG.",
									},
									"driver": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "GPU driver version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cuda": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "CUDA  version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cudnn": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "cuDNN version. Format like: `{ version: String, name: String, doc_name: String, dev_name: String }`. `version`: cuDNN version; `name`: cuDNN name; `doc_name`: Doc name of cuDNN; `dev_name`: Dev name of cuDNN.",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"custom_driver": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Custom GPU driver. Format like: `{address: String}`. `address`: URL of custom GPU driver address.",
									},
								},
							},
						},
					},
				},
			},

			"worker_config_overrides": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Override variable worker_config, commonly used to attach existing instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Mount target. Default is not mounting.",
						},
						"docker_graph_path": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "/var/lib/docker",
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Docker graph path. Default is `/var/lib/docker`.",
						},
						"data_disk": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    11,
							Description: "Configurations of data disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Default:      "CLOUD_PREMIUM",
										Description:  "Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Default:     0,
										Description: "Volume of disk in GB. Default is `0`.",
									},
									"file_system": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "File system, e.g. `ext3/ext4/xfs`.",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate whether to auto format and mount or not. Default is `false`.",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "Mount target.",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The name of the device or partition to mount. NOTE: this argument doesn't support setting in node pool, or will leads to mount error.",
									},
								},
							},
						},
						"extra_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Custom parameter information related to the node. This is a white-list parameter.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Base64-encoded User Data text, the length limit is 16KB.",
						},
						"pre_start_user_script": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Base64-encoded user script, executed before initializing the node, currently only effective for adding existing nodes.",
						},
						"is_schedule": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Deprecated:  "This argument was deprecated, use `unschedulable` instead.",
							Description: "Indicate to schedule the adding node or not. Default is true.",
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicate to set desired pod number in node. valid when the cluster is podCIDR.",
						},
						"gpu_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "GPU driver parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mig_enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether to enable MIG.",
									},
									"driver": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "GPU driver version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cuda": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "CUDA  version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cudnn": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "cuDNN version. Format like: `{ version: String, name: String, doc_name: String, dev_name: String }`. `version`: cuDNN version; `name`: cuDNN name; `doc_name`: Doc name of cuDNN; `dev_name`: Dev name of cuDNN.",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"custom_driver": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Custom GPU driver. Format like: `{address: String}`. `address`: URL of custom GPU driver address.",
									},
								},
							},
						},
					},
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Labels of tke attachment exits CVM.",
			},

			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Sets whether the joining node participates in the schedule. Default is `0`, which means it participates in scheduling. Non-zero(eg: `1`) number means it does not participate in scheduling.",
			},

			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "A list of security group IDs after attach to cluster.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the node.",
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
		clusterId  string
	)
	var (
		request  = tke.NewAddExistedInstancesRequest()
		response = tke.NewAddExistedInstancesResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	request.InstanceIds = []*string{helper.String(instanceId)}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	loginSettings := tke.LoginSettings{}
	if v, ok := d.GetOk("password"); ok {
		loginSettings.Password = helper.String(v.(string))
	}
	request.LoginSettings = &loginSettings

	if instanceAdvancedSettingsMap, ok := helper.InterfacesHeadMap(d, "worker_config"); ok {
		instanceAdvancedSettings := tke.InstanceAdvancedSettings{}
		if v, ok := instanceAdvancedSettingsMap["mount_target"]; ok {
			instanceAdvancedSettings.MountTarget = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["data_disk"]; ok {
			for _, item := range v.([]interface{}) {
				dataDisksMap := item.(map[string]interface{})
				dataDisk := tke.DataDisk{}
				if v, ok := dataDisksMap["disk_type"]; ok {
					dataDisk.DiskType = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["file_system"]; ok {
					dataDisk.FileSystem = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["auto_format_and_mount"]; ok {
					dataDisk.AutoFormatAndMount = helper.Bool(v.(bool))
				}
				if v, ok := dataDisksMap["mount_target"]; ok {
					dataDisk.MountTarget = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["disk_partition"]; ok {
					dataDisk.DiskPartition = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["disk_size"]; ok {
					dataDisk.DiskSize = helper.IntInt64(v.(int))
				}
				instanceAdvancedSettings.DataDisks = append(instanceAdvancedSettings.DataDisks, &dataDisk)
			}
		}
		if v, ok := instanceAdvancedSettingsMap["user_data"]; ok {
			instanceAdvancedSettings.UserScript = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["pre_start_user_script"]; ok {
			instanceAdvancedSettings.PreStartUserScript = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["taints"]; ok {
			for _, item := range v.([]interface{}) {
				taintsMap := item.(map[string]interface{})
				taint := tke.Taint{}
				if v, ok := taintsMap["key"]; ok {
					taint.Key = helper.String(v.(string))
				}
				if v, ok := taintsMap["value"]; ok {
					taint.Value = helper.String(v.(string))
				}
				if v, ok := taintsMap["effect"]; ok {
					taint.Effect = helper.String(v.(string))
				}
				instanceAdvancedSettings.Taints = append(instanceAdvancedSettings.Taints, &taint)
			}
		}
		if v, ok := instanceAdvancedSettingsMap["docker_graph_path"]; ok {
			instanceAdvancedSettings.DockerGraphPath = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["desired_pod_num"]; ok {
			instanceAdvancedSettings.DesiredPodNumber = helper.IntInt64(v.(int))
		}
		if gPUArgsMap, ok := helper.ConvertInterfacesHeadToMap(instanceAdvancedSettingsMap["gpu_args"]); ok {
			gPUArgs := tke.GPUArgs{}
			if v, ok := gPUArgsMap["mig_enable"]; ok {
				gPUArgs.MIGEnable = helper.Bool(v.(bool))
			}
			instanceAdvancedSettings.GPUArgs = &gPUArgs
		}
		request.InstanceAdvancedSettings = &instanceAdvancedSettings
	}

	if v, ok := d.GetOk("hostname"); ok {
		request.HostName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_groups"); ok {
		securityGroupIdsSet := v.([]interface{})
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroupIds))
		}
	}

	if v, ok := d.GetOk("worker_config_overrides"); ok {
		for _, item := range v.([]interface{}) {
			instanceAdvancedSettingsOverridesMap := item.(map[string]interface{})
			instanceAdvancedSettings := tke.InstanceAdvancedSettings{}
			if v, ok := instanceAdvancedSettingsOverridesMap["mount_target"]; ok {
				instanceAdvancedSettings.MountTarget = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["data_disk"]; ok {
				for _, item := range v.([]interface{}) {
					dataDisksMap := item.(map[string]interface{})
					dataDisk := tke.DataDisk{}
					if v, ok := dataDisksMap["disk_type"]; ok {
						dataDisk.DiskType = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["file_system"]; ok {
						dataDisk.FileSystem = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["auto_format_and_mount"]; ok {
						dataDisk.AutoFormatAndMount = helper.Bool(v.(bool))
					}
					if v, ok := dataDisksMap["mount_target"]; ok {
						dataDisk.MountTarget = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["disk_partition"]; ok {
						dataDisk.DiskPartition = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["disk_size"]; ok {
						dataDisk.DiskSize = helper.IntInt64(v.(int))
					}
					instanceAdvancedSettings.DataDisks = append(instanceAdvancedSettings.DataDisks, &dataDisk)
				}
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["user_data"]; ok {
				instanceAdvancedSettings.UserScript = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["pre_start_user_script"]; ok {
				instanceAdvancedSettings.PreStartUserScript = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["docker_graph_path"]; ok {
				instanceAdvancedSettings.DockerGraphPath = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["desired_pod_num"]; ok {
				instanceAdvancedSettings.DesiredPodNumber = helper.IntInt64(v.(int))
			}
			if gPUArgsMap, ok := helper.ConvertInterfacesHeadToMap(instanceAdvancedSettingsOverridesMap["gpu_args"]); ok {
				gPUArgs2 := tke.GPUArgs{}
				if v, ok := gPUArgsMap["mig_enable"]; ok {
					gPUArgs2.MIGEnable = helper.Bool(v.(bool))
				}
				instanceAdvancedSettings.GPUArgs = &gPUArgs2
			}
			request.InstanceAdvancedSettingsOverrides = append(request.InstanceAdvancedSettingsOverrides, &instanceAdvancedSettings)
		}
	}

	if err := resourceTencentCloudKubernetesClusterAttachmentCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().AddExistedInstancesWithContext(ctx, request)
		if e != nil {
			return resourceTencentCloudKubernetesClusterAttachmentCreateRequestOnError0(ctx, request, e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes cluster attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	if err := resourceTencentCloudKubernetesClusterAttachmentCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{instanceId, clusterId}, "_"))

	return resourceTencentCloudKubernetesClusterAttachmentRead(d, meta)
}

func resourceTencentCloudKubernetesClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), "_")
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterId := idSplit[1]

	_ = d.Set("instance_id", instanceId)

	_ = d.Set("cluster_id", clusterId)

	respData, err := service.DescribeKubernetesClusterAttachmentById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	respData1, err := service.DescribeKubernetesClusterAttachmentById1(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData1 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData1.LoginSettings != nil {
		if respData1.LoginSettings.KeyIds != nil {
			_ = d.Set("key_ids", respData1.LoginSettings.KeyIds)
		}

	}

	if respData1.SecurityGroupIds != nil {
		_ = d.Set("security_groups", respData1.SecurityGroupIds)
	}

	if respData1.ImageId != nil {
		_ = d.Set("image_id", respData1.ImageId)
	}

	var respData2 *tke.Instance
	reqErr2 := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterAttachmentById2(ctx, instanceId, clusterId)
		if e != nil {
			return resourceTencentCloudKubernetesClusterAttachmentReadRequestOnError2(ctx, result, e)
		}
		if err := resourceTencentCloudKubernetesClusterAttachmentReadRequestOnSuccess2(ctx, result); err != nil {
			return err
		}
		respData2 = result
		return nil
	})
	if reqErr2 != nil {
		log.Printf("[CRITAL]%s read kubernetes cluster attachment failed, reason:%+v", logId, reqErr2)
		return reqErr2
	}

	if respData2 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData2.InstanceAdvancedSettings != nil {
		if respData2.InstanceAdvancedSettings.Unschedulable != nil {
			_ = d.Set("unschedulable", respData2.InstanceAdvancedSettings.Unschedulable)
		}

	}

	if respData2.InstanceState != nil {
		_ = d.Set("state", respData2.InstanceState)
	}

	return nil
}

func resourceTencentCloudKubernetesClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), "_")
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterId := idSplit[1]

	var (
		request  = tke.NewDeleteClusterInstancesRequest()
		response = tke.NewDeleteClusterInstancesResponse()
	)

	request.ClusterId = helper.String(clusterId)

	request.InstanceIds = []*string{helper.String(instanceId)}

	instanceDeleteMode := "retain"
	request.InstanceDeleteMode = &instanceDeleteMode

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DeleteClusterInstancesWithContext(ctx, request)
		if e != nil {
			return resourceTencentCloudKubernetesClusterAttachmentDeleteRequestOnError0(ctx, e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes cluster attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
