/*
Provides a resource to create a tcm cluster_attachment

Example Usage

```hcl
resource "tencentcloud_tcm_cluster_attachment" "cluster_attachment" {
  mesh_id = "mesh-xxxxxxxx"
  cluster_list {
		cluster_id = "cls-xxxxxxxx"
		region = "ap-shanghai"
		role = "REMOTE"
		vpc_id = "vpc-xxxxxxxx"
		subnet_id = "subnet-xxxxxxx"
		type = "TKE or EKS"

  }
}
```

Import

tcm cluster_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tcm_cluster_attachment.cluster_attachment cluster_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTcmClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmClusterAttachmentCreate,
		Read:   resourceTencentCloudTcmClusterAttachmentRead,
		Delete: resourceTencentCloudTcmClusterAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Mesh ID.",
			},

			"cluster_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
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
							Description: "Cluster&amp;#39;s VpcId.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Subnet id, only needed if it&amp;#39;s standalone mesh.",
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

	logId := getLogId(contextNil)

	var (
		request  = tcm.NewLinkClusterListRequest()
		response = tcm.NewLinkClusterListResponse()
		meshId   string
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcm clusterAttachment failed, reason:%+v", logId, err)
		return err
	}

	meshId = *response.Response.MeshId
	d.SetId(meshId)

	return resourceTencentCloudTcmClusterAttachmentRead(d, meta)
}

func resourceTencentCloudTcmClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_cluster_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterAttachmentId := d.Id()

	clusterAttachment, err := service.DescribeTcmClusterAttachmentById(ctx, meshId)
	if err != nil {
		return err
	}

	if clusterAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcmClusterAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clusterAttachment.MeshId != nil {
		_ = d.Set("mesh_id", clusterAttachment.MeshId)
	}

	if clusterAttachment.ClusterList != nil {
		clusterListList := []interface{}{}
		for _, clusterList := range clusterAttachment.ClusterList {
			clusterListMap := map[string]interface{}{}

			if clusterAttachment.ClusterList.ClusterId != nil {
				clusterListMap["cluster_id"] = clusterAttachment.ClusterList.ClusterId
			}

			if clusterAttachment.ClusterList.Region != nil {
				clusterListMap["region"] = clusterAttachment.ClusterList.Region
			}

			if clusterAttachment.ClusterList.Role != nil {
				clusterListMap["role"] = clusterAttachment.ClusterList.Role
			}

			if clusterAttachment.ClusterList.VpcId != nil {
				clusterListMap["vpc_id"] = clusterAttachment.ClusterList.VpcId
			}

			if clusterAttachment.ClusterList.SubnetId != nil {
				clusterListMap["subnet_id"] = clusterAttachment.ClusterList.SubnetId
			}

			if clusterAttachment.ClusterList.Type != nil {
				clusterListMap["type"] = clusterAttachment.ClusterList.Type
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
	clusterAttachmentId := d.Id()

	if err := service.DeleteTcmClusterAttachmentById(ctx, meshId); err != nil {
		return err
	}

	return nil
}
