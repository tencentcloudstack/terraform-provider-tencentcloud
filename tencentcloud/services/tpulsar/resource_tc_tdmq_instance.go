package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func ResourceTencentCloudTdmqInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqCreate,
		Read:   resourceTencentCloudTdmqRead,
		Update: resourceTencentCloudTdmqUpdate,
		Delete: resourceTencentCloudTdmqDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of tdmq cluster to be created.",
			},
			"bind_cluster_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The Dedicated Cluster Id.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the tdmq cluster.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTdmqCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.create")()

	return fmt.Errorf("tencentcloud_tdmq_instance creation has been deprecated, please use tencentcloud_tdmq_professional_cluster instead")
}

func resourceTencentCloudTdmqRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tdmqService = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id          = d.Id()
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqInstanceById(ctx, id)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("cluster_name", info.ClusterName)
		_ = d.Set("remark", info.Remark)
		return nil
	})

	if err != nil {
		return err
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "tdmq", "cluster", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTdmqUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.update")()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id          = d.Id()
		clusterName string
		remark      string
	)

	immutableArgs := []string{"bind_cluster_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	//internal version: replace var begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace var end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	old, now := d.GetChange("cluster_name")
	if d.HasChange("cluster_name") {
		clusterName = now.(string)
	} else {
		clusterName = old.(string)
	}

	old, now = d.GetChange("remark")
	if d.HasChange("remark") {
		remark = now.(string)
	} else {
		remark = old.(string)
	}

	if err := service.ModifyTdmqInstanceAttribute(ctx, id, clusterName, remark); err != nil {
		return err
	}

	if d.HasChange("tags") {
		//internal version: replace setTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("tdmq", "cluster", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		//internal version: replace setTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	}

	return resourceTencentCloudTdmqRead(d, meta)
}

func resourceTencentCloudTdmqDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.delete")()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		clusterId = d.Id()
	)

	if err := service.DeleteTdmqInstance(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
