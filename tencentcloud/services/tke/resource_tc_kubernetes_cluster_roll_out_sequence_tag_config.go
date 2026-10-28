package tke

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesClusterRollOutSequenceTagConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigCreate,
		Read:   resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigRead,
		Update: resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigUpdate,
		Delete: resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},

			"tags": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Cluster roll-out sequence tags. Supported tags: key `Env` with values [`Test`, `Pre-Production`, `Production`]; key `Protection-Level` with values [`Low`, `Medium`, `High`].",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},

						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	clusterId := d.Get("cluster_id").(string)
	d.SetId(clusterId)

	return resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigUpdate(d, meta)
}

func resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	clusterId := d.Id()
	_ = d.Set("cluster_id", clusterId)

	clusterTag, err := service.DescribeKubernetesClusterRollOutSequenceTagConfigById(ctx, clusterId)
	if err != nil {
		return err
	}

	if clusterTag == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	tagsList := make([]map[string]interface{}, 0, len(clusterTag))
	for _, item := range clusterTag {
		if item != nil && item.Tags != nil && len(item.Tags) > 0 {
			for _, tag := range item.Tags {
				tagMap := map[string]interface{}{}
				if tag.Key != nil {
					tagMap["key"] = tag.Key
				}

				if tag.Value != nil {
					tagMap["value"] = tag.Value
				}

				tagsList = append(tagsList, tagMap)
			}
		}
	}

	_ = d.Set("tags", tagsList)

	return nil
}

func resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20180525.NewModifyClusterRollOutSequenceTagsRequest()
	)

	clusterId := d.Id()
	request.ClusterID = helper.String(clusterId)

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagMap := item.(map[string]interface{})
			tag := tkev20180525.Tag{}
			if v, ok := tagMap["key"].(string); ok && v != "" {
				tag.Key = helper.String(v)
			}

			if v, ok := tagMap["value"].(string); ok && v != "" {
				tag.Value = helper.String(v)
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterRollOutSequenceTagsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update kubernetes cluster roll out sequence tag config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigRead(d, meta)
}

func resourceTencentCloudKubernetesClusterRollOutSequenceTagConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20180525.NewModifyClusterRollOutSequenceTagsRequest()
	)

	clusterId := d.Id()
	request.ClusterID = helper.String(clusterId)
	// An empty Tags list means removing all roll-out sequence tags from the cluster.
	request.Tags = []*tkev20180525.Tag{}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterRollOutSequenceTagsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete kubernetes cluster roll out sequence tag config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
