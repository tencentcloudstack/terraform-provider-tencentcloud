package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcmClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTcmClusterAttachmentRead,
		Create: resourceTencentCloudTcmClusterAttachmentCreate,
		Delete: resourceTencentCloudTcmClusterAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Mesh ID.",
			},

			"cluster_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Cluster list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "TKE Cluster id.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "TKE cluster region.",
						},
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cluster role in mesh, REMOTE or MASTER.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cluster&#39;s VpcId.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Subnet id, only needed if it&#39;s standalone mesh.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cluster type.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcmClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_cluster_attachment.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		request   = tcm.NewLinkClusterListRequest()
		meshId    string
		clusterId string
	)

	if v, ok := d.GetOk("mesh_id"); ok {
		meshId = v.(string)
		request.MeshId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			cluster := tcm.Cluster{}
			if v, ok := dMap["cluster_id"]; ok {
				clusterId = v.(string)
				cluster.ClusterId = helper.String(v.(string))
			}
			if v, ok := dMap["region"]; ok {
				cluster.Region = helper.String(v.(string))
			}
			if v, ok := dMap["role"]; ok {
				cluster.Role = helper.String(v.(string))
			}
			if v, ok := dMap["vpc_id"]; ok {
				cluster.VpcId = helper.String(v.(string))
			}
			if v, ok := dMap["subnet_id"]; ok {
				cluster.SubnetId = helper.String(v.(string))
			}
			if v, ok := dMap["type"]; ok {
				cluster.Type = helper.String(v.(string))
			}

			request.ClusterList = append(request.ClusterList, &cluster)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().LinkClusterList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcm clusterAttachment failed, reason:%+v", logId, err)
		return err
	}

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		mesh, errRet := service.DescribeTcmMesh(ctx, meshId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		clusterList := mesh.Mesh.ClusterList
		if len(clusterList) < 1 {
			return resource.RetryableError(fmt.Errorf("link is being created, retry..."))
		}
		var linkState string
		for _, v := range clusterList {
			if *v.ClusterId != clusterId {
				continue
			}
			linkState = *v.Status.LinkState
			if linkState == "LINKED" {
				return nil
			}
			if linkState == "LINK_FAILED" {
				return resource.NonRetryableError(fmt.Errorf("link status is %v, operate failed.", linkState))
			}
		}
		return resource.RetryableError(fmt.Errorf("link status is %v, retry...", linkState))
	})
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{meshId, clusterId}, FILED_SP))
	return resourceTencentCloudTcmClusterAttachmentRead(d, meta)
}

func resourceTencentCloudTcmClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_cluster_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	meshId := ids[0]
	clusterId := ids[1]

	mesh, err := service.DescribeTcmMesh(ctx, meshId)

	if err != nil {
		return err
	}

	if mesh == nil || mesh.Mesh == nil || len(mesh.Mesh.ClusterList) < 1 {
		d.SetId("")
		return fmt.Errorf("resource `clusterAttachment` %s does not exist", meshId)
	}

	_ = d.Set("mesh_id", meshId)

	if len(mesh.Mesh.ClusterList) > 0 {
		clusterAttachment := mesh.Mesh
		clusterListList := []interface{}{}
		for _, clusterList := range clusterAttachment.ClusterList {
			if *clusterList.ClusterId != clusterId {
				continue
			}
			clusterListMap := map[string]interface{}{}
			if clusterList.ClusterId != nil {
				clusterListMap["cluster_id"] = clusterList.ClusterId
			}
			if clusterList.Region != nil {
				clusterListMap["region"] = clusterList.Region
			}
			if clusterList.Role != nil {
				clusterListMap["role"] = clusterList.Role
			}
			if clusterList.VpcId != nil {
				clusterListMap["vpc_id"] = clusterList.VpcId
			}
			if clusterList.SubnetId != nil {
				clusterListMap["subnet_id"] = clusterList.SubnetId
			}
			if clusterList.Type != nil {
				clusterListMap["type"] = clusterList.Type
			}

			clusterListList = append(clusterListList, clusterListMap)
		}
		_ = d.Set("cluster_list", clusterListList)
	}

	return nil
}

func resourceTencentCloudTcmClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_cluster_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	meshId := ids[0]
	clusterId := ids[1]

	if err := service.DeleteTcmClusterAttachmentById(ctx, meshId, clusterId); err != nil {
		return err
	}

	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		mesh, errRet := service.DescribeTcmMesh(ctx, meshId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		clusterList := mesh.Mesh.ClusterList
		if len(clusterList) < 1 {
			return nil
		}
		var linkState string
		for _, v := range clusterList {
			if *v.ClusterId != clusterId {
				continue
			}
			linkState = *v.Status.LinkState
			if linkState == "UNLINK_FAILED" {
				return resource.NonRetryableError(fmt.Errorf("link status is %v, operate failed.", linkState))
			}
		}
		return resource.RetryableError(fmt.Errorf("link status is %v, retry...", linkState))
	})
	if err != nil {
		return err
	}

	return nil
}
