package clb

import (
	"context"
	"fmt"
	"log"
	"sync"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var clsActionMu = &sync.Mutex{}

func ResourceTencentCloudClbLogTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceTopicCreate,
		Read:   resourceTencentCloudClbInstanceTopicRead,
		Delete: resourceTencentCloudClbInstanceTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"log_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic of CLB instance.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic of CLB instance.",
			},
			//"partition_count": {
			//	Type:         schema.TypeInt,
			//	Optional:     true,
			//	ValidateFunc: tccommon.ValidateIntegerInRange(1, 10),
			//	Description:  "Topic partition count of CLB instance.(Default 1).",
			//},
			//compute
			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The status of log topic.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Log topic creation time.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceTopicCreate(d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clsService := ClsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if v, ok := d.GetOk("log_set_id"); ok {
		info, err := clsService.DescribeClsLogset(ctx, v.(string))
		if err != nil {
			return err
		}
		if info == nil {
			return fmt.Errorf("resource `log_set` %s does not exist", v.(string))
		}
	}

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	params := make(map[string]interface{})
	if topicName, ok := d.GetOk("topic_name"); ok {
		params["topic_name"] = topicName
	}
	if partitionCount, ok := d.GetOk("partition_count"); ok {
		params["partition_count"] = partitionCount
	}
	resp, err := clbService.CreateTopic(ctx, params)
	if err != nil {
		log.Printf("[CRITAL]%s create clb topic failed, reason:%+v", logId, err)
		return err
	}
	d.SetId(*resp.Response.TopicId)
	return resourceTencentCloudClbInstanceTopicRead(d, meta)
}

func resourceTencentCloudClbInstanceTopicRead(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	clsService := ClsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	res, err := clsService.DescribeClsTopicById(ctx, id)
	if err != nil {
		return err
	}
	if res == nil {
		d.SetId("")
		return fmt.Errorf("resource `logTopic` %s does not exist", id)
	}
	_ = d.Set("log_set_id", res.LogsetId)
	_ = d.Set("topic_name", res.TopicName)
	_ = d.Set("create_time", res.CreateTime)
	_ = d.Set("status", res.Status)
	return nil

}

func resourceTencentCloudClbInstanceTopicDelete(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	clsService := ClsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := clsService.DeleteClsTopic(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
