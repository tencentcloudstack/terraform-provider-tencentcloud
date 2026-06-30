package dbdc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbdcDbCustomNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbdcDbCustomNodeCreate,
		Read:   resourceTencentCloudDbdcDbCustomNodeRead,
		Update: resourceTencentCloudDbdcDbCustomNodeUpdate,
		Delete: resourceTencentCloudDbdcDbCustomNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Availability zone supported by the product, e.g. `ap-shanghai-5`, `ap-shanghai-8`, `ap-nanjing-3`.",
			},

			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image ID, format `img-xxxxxxx`. Must be an image owned by the DB Custom product under the current account.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID used to establish the SSH connection for the node. Must be owned by the current account and cannot be cross-region.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet ID used to establish the SSH connection for the node. Must belong to the VPC and match the availability zone.",
			},

			"node_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Node spec, e.g. `DB.AT5.8XLARGE128`, `DB.AT5.16XLARGE256`, `DB.AT5.32XLARGE512`, `DB.AT5.64XLARGE1152`, `DB.AT5.128XLARGE2304`.",
			},

			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Description: "Purchase duration in months. Valid values: 1/2/3/4/5/6/7/8/9/10/11/12/24/36. Default value is `1`.",
			},

			"node_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Node name. Up to 128 characters.",
			},

			"login_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Instance login settings. You can set the login method to password, key, or keep the original image login settings. Only one method can be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: "Instance login password. Password complexity limits vary by operating system type.",
						},
						"key_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Key pair ID list. Only a single ID is supported currently. Password and key cannot be specified at the same time.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"keep_image_login": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Whether to keep the original login settings of the image. Valid values: `true`, `false`. Cannot be specified together with Password or KeyIds.",
						},
					},
				},
			},

			"auto_voucher": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to use voucher to deduct automatically. Valid values: `1` (use), `0` (not use). Default value is `0`.",
			},

			"voucher_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Voucher ID list. Must be undeducted voucher IDs owned by the current account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"auto_renew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Auto-renew flag. Valid values: `1` (auto-renew), `2` (not auto-renew). Mutable via the renew API.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Node tags.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// computed
			"node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node ID.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID that the node belongs to.",
			},

			"ssh_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSH endpoint to access this node, in the format `IP:Port`.",
			},

			"lan_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet communication IP address of the node.",
			},

			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Node CPU size, unit: core.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Node memory, unit: GiB.",
			},

			"os_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating system name of the node.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node status. Valid values: `Creating`, `Running`, `Isolating`, `Isolated`, `Activating`, `Destroying`.",
			},

			"charge_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Charge type. Valid values: `PREPAID`.",
			},

			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node expiration time.",
			},

			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node creation time.",
			},

			"isolated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node isolation time.",
			},

			"system_disk": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System disk information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk type.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size, unit: GiB.",
						},
					},
				},
			},

			"data_disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data disk information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk type.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size, unit: GiB.",
						},
						"disk_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk name.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDbdcDbCustomNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_node.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = dbdcv20201029.NewCreateDBCustomNodesRequest()
		response = dbdcv20201029.NewCreateDBCustomNodesResponse()
		nodeId   string
	)

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("node_type"); ok {
		request.NodeType = helper.String(v.(string))
	}

	request.NodeCount = helper.IntInt64(1)

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("node_name"); ok {
		request.NodeName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "login_settings"); ok {
		loginSettings := dbdcv20201029.LoginSettings{}
		if v, ok := dMap["password"]; ok && v.(string) != "" {
			loginSettings.Password = helper.String(v.(string))
		}

		if v, ok := dMap["key_ids"]; ok {
			keyIdsList := v.([]interface{})
			for i := range keyIdsList {
				keyId := keyIdsList[i].(string)
				loginSettings.KeyIds = append(loginSettings.KeyIds, &keyId)
			}
		}

		if v, ok := dMap["keep_image_login"]; ok && v.(string) != "" {
			loginSettings.KeepImageLogin = helper.String(v.(string))
		}

		request.LoginSettings = &loginSettings
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsList := v.([]interface{})
		for i := range voucherIdsList {
			voucherId := voucherIdsList[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherId)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		for tagKey, tagValue := range v.(map[string]interface{}) {
			tag := dbdcv20201029.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue.(string)),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().CreateDBCustomNodesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dbdc db custom node failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dbdc db custom node failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.NodeIds) == 0 || response.Response.NodeIds[0] == nil {
		return fmt.Errorf("NodeIds is empty.")
	}

	nodeId = *response.Response.NodeIds[0]
	d.SetId(nodeId)

	// Create is async, wait for the task to succeed.
	if response.Response.TaskId != nil {
		if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	return resourceTencentCloudDbdcDbCustomNodeRead(d, meta)
}

func resourceTencentCloudDbdcDbCustomNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_node.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		nodeId  = d.Id()
	)

	respData, err := service.DescribeDBCustomNodeById(ctx, nodeId)
	if err != nil {
		return err
	}

	if respData == nil || respData.NodeId == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dbdc_db_custom_node` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Zone != nil {
		_ = d.Set("zone", respData.Zone)
	}

	if respData.ImageId != nil {
		_ = d.Set("image_id", respData.ImageId)
	}

	if respData.VpcId != nil {
		_ = d.Set("vpc_id", respData.VpcId)
	}

	if respData.SubnetId != nil {
		_ = d.Set("subnet_id", respData.SubnetId)
	}

	if respData.NodeType != nil {
		_ = d.Set("node_type", respData.NodeType)
	}

	if respData.NodeName != nil {
		_ = d.Set("node_name", respData.NodeName)
	}

	if respData.AutoRenew != nil {
		_ = d.Set("auto_renew", respData.AutoRenew)
	}

	if respData.Tags != nil {
		tags := make(map[string]interface{}, len(respData.Tags))
		for _, tag := range respData.Tags {
			if tag == nil || tag.Key == nil {
				continue
			}

			if tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			} else {
				tags[*tag.Key] = ""
			}
		}

		_ = d.Set("tags", tags)
	}

	if respData.NodeId != nil {
		_ = d.Set("node_id", respData.NodeId)
	}

	if respData.ClusterId != nil {
		_ = d.Set("cluster_id", respData.ClusterId)
	}

	if respData.SSHEndpoint != nil {
		_ = d.Set("ssh_endpoint", respData.SSHEndpoint)
	}

	if respData.LanIP != nil {
		_ = d.Set("lan_ip", respData.LanIP)
	}

	if respData.CPU != nil {
		_ = d.Set("cpu", respData.CPU)
	}

	if respData.Memory != nil {
		_ = d.Set("memory", respData.Memory)
	}

	if respData.OsName != nil {
		_ = d.Set("os_name", respData.OsName)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.ChargeType != nil {
		_ = d.Set("charge_type", respData.ChargeType)
	}

	if respData.ExpireTime != nil {
		_ = d.Set("expire_time", respData.ExpireTime)
	}

	if respData.CreatedTime != nil {
		_ = d.Set("created_time", respData.CreatedTime)
	}

	if respData.IsolatedTime != nil {
		_ = d.Set("isolated_time", respData.IsolatedTime)
	}

	if respData.SystemDisk != nil {
		systemDiskMap := map[string]interface{}{}
		if respData.SystemDisk.DiskType != nil {
			systemDiskMap["disk_type"] = respData.SystemDisk.DiskType
		}

		if respData.SystemDisk.DiskSize != nil {
			systemDiskMap["disk_size"] = respData.SystemDisk.DiskSize
		}

		_ = d.Set("system_disk", []interface{}{systemDiskMap})
	}

	if respData.DataDisks != nil {
		dataDisksList := make([]interface{}, 0, len(respData.DataDisks))
		for _, dataDisk := range respData.DataDisks {
			if dataDisk == nil {
				continue
			}

			dataDiskMap := map[string]interface{}{}
			if dataDisk.DiskType != nil {
				dataDiskMap["disk_type"] = dataDisk.DiskType
			}

			if dataDisk.DiskSize != nil {
				dataDiskMap["disk_size"] = dataDisk.DiskSize
			}

			if dataDisk.DiskName != nil {
				dataDiskMap["disk_name"] = dataDisk.DiskName
			}

			dataDisksList = append(dataDisksList, dataDiskMap)
		}

		_ = d.Set("data_disks", dataDisksList)
	}

	return nil
}

func resourceTencentCloudDbdcDbCustomNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_node.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		nodeId = d.Id()
	)

	if d.HasChange("tags") {
		oldRaw, newRaw := d.GetChange("tags")
		oldTags := oldRaw.(map[string]interface{})
		newTags := newRaw.(map[string]interface{})

		request := dbdcv20201029.NewModifyDBCustomNodeTagsRequest()
		request.NodeId = helper.String(nodeId)

		for tagKey, tagValue := range newTags {
			tag := dbdcv20201029.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue.(string)),
			}

			request.AddTags = append(request.AddTags, &tag)
		}

		for tagKey := range oldTags {
			if _, ok := newTags[tagKey]; !ok {
				request.DeleteTagKeys = append(request.DeleteTagKeys, helper.String(tagKey))
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().ModifyDBCustomNodeTagsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify dbdc db custom node tags failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dbdc db custom node tags failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("auto_renew") {
		request := dbdcv20201029.NewRenewDBCustomNodeRequest()
		request.NodeId = helper.String(nodeId)

		if v, ok := d.GetOkExists("period"); ok {
			request.Period = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("auto_renew"); ok {
			request.AutoRenew = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("auto_voucher"); ok {
			request.AutoVoucher = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("voucher_ids"); ok {
			voucherIdsList := v.([]interface{})
			for i := range voucherIdsList {
				voucherId := voucherIdsList[i].(string)
				request.VoucherIds = append(request.VoucherIds, &voucherId)
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().RenewDBCustomNodeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Renew dbdc db custom node failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s renew dbdc db custom node failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudDbdcDbCustomNodeRead(d, meta)
}

func resourceTencentCloudDbdcDbCustomNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_node.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		nodeId   = d.Id()
		response = dbdcv20201029.NewDestroyDBCustomNodeResponse()
	)

	// Stage 1: isolate the node (this API does not return a TaskId).
	isolateRequest := dbdcv20201029.NewIsolateDBCustomNodeRequest()
	isolateRequest.NodeId = helper.String(nodeId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().IsolateDBCustomNodeWithContext(ctx, isolateRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, isolateRequest.GetAction(), isolateRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Isolate dbdc db custom node failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s isolate dbdc db custom node failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Wait until the node reaches the `Isolated` status before destroying it.
	if err := waitDBCustomNodeStatus(ctx, &service, nodeId, "Isolated", d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	// Stage 2: destroy the node (async, returns a TaskId).
	destroyRequest := dbdcv20201029.NewDestroyDBCustomNodeRequest()
	destroyRequest.NodeId = helper.String(nodeId)
	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().DestroyDBCustomNodeWithContext(ctx, destroyRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, destroyRequest.GetAction(), destroyRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Destroy dbdc db custom node failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s destroy dbdc db custom node failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TaskId != nil {
		if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	return nil
}

// waitDBCustomNodeStatus polls DescribeDBCustomNodes until the node reaches the
// target status. If the node can no longer be found it is treated as gone. The
// poll loop honors the given timeout.
func waitDBCustomNodeStatus(ctx context.Context, service *DbdcService, nodeId, target string, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		node, e := service.DescribeDBCustomNodeById(ctx, nodeId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if node == nil || node.NodeId == nil {
			return nil
		}

		if node.Status != nil && *node.Status == target {
			return nil
		}

		status := ""
		if node.Status != nil {
			status = *node.Status
		}

		return resource.RetryableError(fmt.Errorf("dbdc db custom node [%s] status is `%s`, waiting for `%s`.", nodeId, status, target))
	})
}
