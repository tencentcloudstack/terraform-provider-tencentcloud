package dbdc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbdcNodeToDbCustomClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentCreate,
		Read:   resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentRead,
		Update: resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentUpdate,
		Delete: resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DB Custom cluster ID.",
			},

			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DB Custom node ID to add to the cluster.",
			},

			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "OS image ID to reset the node to after it is added to the cluster. Obtainable via the `DescribeDBCustomImages` API.",
			},

			"login_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Instance login settings. You can set the login method to password, key, or keep the original image login settings. Only one method can be set; for the key method, only a single key ID is supported.",
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

			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    20,
				Description: "Custom labels (Kubernetes labels) initialized after the node is added to the cluster. Up to 20 key-value pairs. Mutable via the `ModifyDBCustomClusterNodeConfig` API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label value.",
						},
					},
				},
			},

			"taints": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    5,
				Description: "Taints (Kubernetes taints) initialized after the node is added to the cluster. Up to 5 taints. Uniqueness key is (key, effect). Mutable via the `ModifyDBCustomClusterNodeConfig` API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Taint key.",
						},
						"effect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Taint effect. Valid values: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Taint value.",
						},
					},
				},
			},

			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Hostname of the node. Required when `host_name_type` is `1`; ignored otherwise. Uppercase letters and underscores (`_`) are not allowed; dots (`.`) and hyphens (`-`) cannot be the first/last character or used consecutively. Windows: 2-15 chars (letters, digits, `-`, no `.`); Linux/others: 2-60 chars (dot-separated segments). Supports pattern strings `{R:x}`, `{R:x,F:y}`, `{IP}`.",
			},

			"host_name_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Hostname source type. Valid values: `0` (reuse hostname set at node creation), `1` (re-specify HostName, must pass `host_name`), `2` (system auto-assign using NodeId).",
			},

			// computed
			"node_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node name.",
			},

			"lan_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet IP address of the node.",
			},

			"ssh_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSH endpoint to access the node, in the format `IP:Port`.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status of the node in the cluster.",
			},

			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability zone that the node belongs to.",
			},

			"node_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node spec.",
			},

			"network_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network mode. Valid values: `privatelink` (four-layer network connectivity), `cross_tenant_eni` (three-layer network connectivity, dual-NIC mode).",
			},

			"eni_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the network mode is `cross_tenant_eni`, this IP address is the user-accessible address.",
			},
		},
	}
}

func resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_node_to_db_custom_cluster_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request   = dbdcv20201029.NewAddNodesToDBCustomClusterRequest()
		response  = dbdcv20201029.NewAddNodesToDBCustomClusterResponse()
		clusterId string
		nodeId    string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(clusterId)
	}

	if v, ok := d.GetOk("node_id"); ok {
		nodeId = v.(string)
		request.NodeIds = []*string{helper.String(nodeId)}
	}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
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

	if v, ok := d.GetOk("labels"); ok {
		labelsList := v.([]interface{})
		for i := range labelsList {
			labelMap := labelsList[i].(map[string]interface{})
			label := dbdcv20201029.Label{}
			if key, ok := labelMap["key"]; ok {
				label.Key = helper.String(key.(string))
			}
			if value, ok := labelMap["value"]; ok && value.(string) != "" {
				label.Value = helper.String(value.(string))
			}
			request.Labels = append(request.Labels, &label)
		}
	}

	if v, ok := d.GetOk("taints"); ok {
		taintsList := v.([]interface{})
		for i := range taintsList {
			taintMap := taintsList[i].(map[string]interface{})
			taint := dbdcv20201029.Taint{}
			if key, ok := taintMap["key"]; ok {
				taint.Key = helper.String(key.(string))
			}
			if effect, ok := taintMap["effect"]; ok {
				taint.Effect = helper.String(effect.(string))
			}
			if value, ok := taintMap["value"]; ok && value.(string) != "" {
				taint.Value = helper.String(value.(string))
			}
			request.Taints = append(request.Taints, &taint)
		}
	}

	if v, ok := d.GetOk("host_name"); ok {
		request.HostName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("host_name_type"); ok {
		request.HostNameType = helper.Int64(int64(v.(int)))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().AddNodesToDBCustomClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Add nodes to dbdc db custom cluster failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s add nodes to dbdc db custom cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Create is async, wait for the task to succeed.
	if response.Response.TaskId != nil {
		if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	d.SetId(strings.Join([]string{clusterId, nodeId}, tccommon.FILED_SP))
	return resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentRead(d, meta)
}

func resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_node_to_db_custom_cluster_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	clusterId := idSplit[0]
	nodeId := idSplit[1]

	respData, err := service.DescribeDBCustomClusterNodeById(ctx, clusterId, nodeId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("node_id", nodeId)

	if respData.NodeName != nil {
		_ = d.Set("node_name", respData.NodeName)
	}

	if respData.LanIP != nil {
		_ = d.Set("lan_ip", respData.LanIP)
	}

	if respData.SSHEndpoint != nil {
		_ = d.Set("ssh_endpoint", respData.SSHEndpoint)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Zone != nil {
		_ = d.Set("zone", respData.Zone)
	}

	if respData.NodeType != nil {
		_ = d.Set("node_type", respData.NodeType)
	}

	if respData.NetworkMode != nil {
		_ = d.Set("network_mode", respData.NetworkMode)
	}

	if respData.EniIP != nil {
		_ = d.Set("eni_ip", respData.EniIP)
	}

	// Labels and taints are not part of DescribeDBCustomClusterNodes; query
	// them via the dedicated DescribeDBCustomClusterNodeConfig API.
	configRequest := dbdcv20201029.NewDescribeDBCustomClusterNodeConfigRequest()
	configRequest.ClusterId = helper.String(clusterId)
	configRequest.NodeIds = []*string{helper.String(nodeId)}
	configReqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().DescribeDBCustomClusterNodeConfigWithContext(ctx, configRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil || len(result.Response.NodeSet) == 0 {
			return nil
		}

		for _, nodeConfig := range result.Response.NodeSet {
			if nodeConfig == nil || nodeConfig.NodeId == nil || *nodeConfig.NodeId != nodeId {
				continue
			}

			labelsList := make([]interface{}, 0, len(nodeConfig.Labels))
			for _, label := range nodeConfig.Labels {
				if label == nil {
					continue
				}
				labelMap := map[string]interface{}{}
				if label.Key != nil {
					labelMap["key"] = *label.Key
				}
				if label.Value != nil {
					labelMap["value"] = *label.Value
				}
				labelsList = append(labelsList, labelMap)
			}
			_ = d.Set("labels", labelsList)

			taintsList := make([]interface{}, 0, len(nodeConfig.Taints))
			for _, taint := range nodeConfig.Taints {
				if taint == nil {
					continue
				}
				taintMap := map[string]interface{}{}
				if taint.Key != nil {
					taintMap["key"] = *taint.Key
				}
				if taint.Effect != nil {
					taintMap["effect"] = *taint.Effect
				}
				if taint.Value != nil {
					taintMap["value"] = *taint.Value
				}
				taintsList = append(taintsList, taintMap)
			}
			_ = d.Set("taints", taintsList)
			break
		}
		return nil
	})
	if configReqErr != nil {
		log.Printf("[CRITAL]%s describe dbdc db custom cluster node config failed, reason:%+v", logId, configReqErr)
		return configReqErr
	}

	return nil
}

func resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_node_to_db_custom_cluster_attachment.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	clusterId := idSplit[0]
	nodeId := idSplit[1]

	// Only labels and taints are mutable; other fields are ForceNew.
	if d.HasChange("labels") || d.HasChange("taints") {
		request := dbdcv20201029.NewModifyDBCustomClusterNodeConfigRequest()
		request.ClusterId = helper.String(clusterId)
		request.NodeIds = []*string{helper.String(nodeId)}

		// Reconcile labels to the desired full set: upsert all new labels,
		// delete old label keys absent from the new set.
		if d.HasChange("labels") {
			oldLabelsRaw, newLabelsRaw := d.GetChange("labels")
			newLabels := newLabelsRaw.([]interface{})
			oldLabels := oldLabelsRaw.([]interface{})

			newLabelKeys := make(map[string]bool)
			for _, l := range newLabels {
				lm := l.(map[string]interface{})
				label := dbdcv20201029.Label{}
				if k, ok := lm["key"]; ok && k.(string) != "" {
					label.Key = helper.String(k.(string))
					newLabelKeys[k.(string)] = true
				}
				if v, ok := lm["value"]; ok && v.(string) != "" {
					label.Value = helper.String(v.(string))
				}
				request.UpsertLabels = append(request.UpsertLabels, &label)
			}

			for _, l := range oldLabels {
				lm := l.(map[string]interface{})
				if k, ok := lm["key"]; ok && k.(string) != "" {
					if !newLabelKeys[k.(string)] {
						key := k.(string)
						request.DeleteLabelKeys = append(request.DeleteLabelKeys, &key)
					}
				}
			}
		}

		// Reconcile taints to the desired full set: upsert all new taints
		// (uniqueness key is (Key, Effect)), delete old taints absent from
		// the new set.
		if d.HasChange("taints") {
			oldTaintsRaw, newTaintsRaw := d.GetChange("taints")
			newTaints := newTaintsRaw.([]interface{})
			oldTaints := oldTaintsRaw.([]interface{})

			newTaintKeys := make(map[string]bool)
			for _, t := range newTaints {
				tm := t.(map[string]interface{})
				taint := dbdcv20201029.Taint{}
				key := ""
				effect := ""
				if k, ok := tm["key"]; ok && k.(string) != "" {
					taint.Key = helper.String(k.(string))
					key = k.(string)
				}
				if e, ok := tm["effect"]; ok && e.(string) != "" {
					taint.Effect = helper.String(e.(string))
					effect = e.(string)
				}
				if v, ok := tm["value"]; ok && v.(string) != "" {
					taint.Value = helper.String(v.(string))
				}
				request.UpsertTaints = append(request.UpsertTaints, &taint)
				newTaintKeys[key+"|"+effect] = true
			}

			for _, t := range oldTaints {
				tm := t.(map[string]interface{})
				key := ""
				effect := ""
				if k, ok := tm["key"]; ok {
					key = k.(string)
				}
				if e, ok := tm["effect"]; ok {
					effect = e.(string)
				}
				if key == "" {
					continue
				}
				if !newTaintKeys[key+"|"+effect] {
					taint := dbdcv20201029.Taint{
						Key:    helper.String(key),
						Effect: helper.String(effect),
					}
					request.DeleteTaints = append(request.DeleteTaints, &taint)
				}
			}
		}

		var taskId *uint64
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().ModifyDBCustomClusterNodeConfigWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify dbdc db custom cluster node config failed, Response is nil."))
			}

			taskId = result.Response.TaskId
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify dbdc db custom cluster node config failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// Modify is async, wait for the task to succeed.
		if taskId != nil {
			if err := waitDBCustomTaskSucceeded(ctx, &service, *taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentRead(d, meta)
}

func resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_node_to_db_custom_cluster_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = dbdcv20201029.NewRemoveNodesFromDBCustomClusterRequest()
		response = dbdcv20201029.NewRemoveNodesFromDBCustomClusterResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	clusterId := idSplit[0]
	nodeId := idSplit[1]

	request.ClusterId = helper.String(clusterId)
	request.NodeIds = []*string{helper.String(nodeId)}

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

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().RemoveNodesFromDBCustomClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Remove nodes from dbdc db custom cluster failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s remove nodes from dbdc db custom cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Delete is async, wait for the task to succeed.
	if response.Response.TaskId != nil {
		if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	return nil
}
