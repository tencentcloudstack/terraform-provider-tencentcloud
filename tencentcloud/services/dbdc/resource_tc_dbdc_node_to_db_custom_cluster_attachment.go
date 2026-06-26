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
		Delete: resourceTencentCloudDbdcNodeToDbCustomClusterAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Description: "OS image ID to reset the node to after it is added to the cluster.",
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

	return nil
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
