package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcFlowLogCreate,
		Read:   resourceTencentCloudVpcFlowLogRead,
		Update: resourceTencentCloudVpcFlowLogUpdate,
		Delete: resourceTencentCloudVpcFlowLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_log_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specify flow log name.",
			},
			"resource_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specify resource type. NOTE: Only support `NETWORKINTERFACE` for now. Values: `VPC`, `SUBNET`, `NETWORKINTERFACE`, `CCN`, `NAT`, `DCG`.",
			},
			"resource_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specify resource unique Id of `resource_type` configured.",
			},
			"traffic_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specify log traffic type, values: `ACCEPT`, `REJECT`, `ALL`.",
			},
			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specify vpc Id, ignore while `resource_type` is `CCN` (unsupported) but required while other types.",
			},
			"flow_log_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specify flow Log description.",
			},
			"cloud_log_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specify flow log storage id, just set cls topic id.",
			},
			"storage_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specify consumer type, values: `cls`, `ckafka`.",
			},
			"flow_log_storage": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Specify consumer detail, required while `storage_type` is `ckafka`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specify storage instance id, required while `storage_type` is `ckafka`.",
						},
						"storage_topic": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specify storage topic id, required while `storage_type` is `ckafka`.",
						},
					},
				},
			},
			"cloud_log_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify flow log storage region, default using current.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudVpcFlowLogCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = vpc.NewCreateFlowLogRequest()
		response  = vpc.NewCreateFlowLogResponse()
		flowLogId string
	)
	if v, ok := d.GetOk("flow_log_name"); ok {
		request.FlowLogName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("traffic_type"); ok {
		request.TrafficType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("flow_log_description"); ok {
		request.FlowLogDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cloud_log_id"); ok {
		request.CloudLogId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "flow_log_storage"); ok {
		flowLogStorage := vpc.FlowLogStorage{}
		if v, ok := dMap["storage_id"]; ok {
			flowLogStorage.StorageId = helper.String(v.(string))
		}
		if v, ok := dMap["storage_topic"]; ok {
			flowLogStorage.StorageTopic = helper.String(v.(string))
		}
		request.FlowLogStorage = &flowLogStorage
	}

	if v, ok := d.GetOk("cloud_log_region"); ok {
		request.CloudLogRegion = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateFlowLog(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc flowLog failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.FlowLog) == 0 {
		return fmt.Errorf("api %s returns nil response", request.GetAction())
	}

	flowLogId = *response.Response.FlowLog[0].FlowLogId

	resourceTencentCloudSetFlowLogId(d, flowLogId, d.Get("vpc_id").(string))

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		client := meta.(*TencentCloudClient).apiV3Conn
		tagService := TagService{client}
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		resourceName := BuildTagResourceName("vpc", "fl", client.Region, flowLogId)
		err := tagService.ModifyTags(ctx, resourceName, tags, nil)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudVpcFlowLogRead(d, meta)
}

func resourceTencentCloudVpcFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	flowLogId, vpcId, err := resourceTencentCloudGetFlowLogId(d)

	if err != nil {
		return err
	}

	flowLog, err := service.DescribeVpcFlowLogById(ctx, flowLogId, vpcId)
	if err != nil {
		return err
	}

	if flowLog == nil {
		d.SetId("")
		return fmt.Errorf("resource `vpc_flow_log` %s does not exist", flowLogId)
	}

	if flowLog.FlowLogName != nil {
		_ = d.Set("flow_log_name", flowLog.FlowLogName)
	}

	if flowLog.ResourceType != nil {
		_ = d.Set("resource_type", flowLog.ResourceType)
	}

	if flowLog.ResourceId != nil {
		_ = d.Set("resource_id", flowLog.ResourceId)
	}

	if flowLog.TrafficType != nil {
		_ = d.Set("traffic_type", flowLog.TrafficType)
	}

	if flowLog.VpcId != nil {
		_ = d.Set("vpc_id", flowLog.VpcId)
	}

	if flowLog.FlowLogDescription != nil {
		_ = d.Set("flow_log_description", flowLog.FlowLogDescription)
	}

	if flowLog.CloudLogId != nil {
		_ = d.Set("cloud_log_id", flowLog.CloudLogId)
	}

	if flowLog.StorageType != nil {
		_ = d.Set("storage_type", flowLog.StorageType)
	}

	if flowLog.FlowLogStorage != nil {
		flowLogStorageMap := map[string]interface{}{}

		if flowLog.FlowLogStorage.StorageId != nil {
			flowLogStorageMap["storage_id"] = flowLog.FlowLogStorage.StorageId
		}

		if flowLog.FlowLogStorage.StorageTopic != nil {
			flowLogStorageMap["storage_topic"] = flowLog.FlowLogStorage.StorageTopic
		}

		_ = d.Set("flow_log_storage", []interface{}{flowLogStorageMap})
	}

	if flowLog.CloudLogRegion != nil {
		_ = d.Set("cloud_log_region", flowLog.CloudLogRegion)
	}

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: client}
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "fl", client.Region, flowLogId)
	if err != nil {
		return err
	}

	if len(tags) > 0 {
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudVpcFlowLogUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := vpc.NewModifyFlowLogAttributeRequest()

	flowLogId, vpcId, err := resourceTencentCloudGetFlowLogId(d)

	if err != nil {
		return err
	}

	immutableArgs := []string{
		"resource_type",
		"resource_id",
		"vpc_id", // VPC now used as ID, means it cannot be modified for now
		"traffic_type",
		"cloud_log_id",
		"storage_type",
		"flow_log_storage",
		"cloud_log_region",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			o, _ := d.GetChange(v)
			_ = d.Set(v, o)
			return fmt.Errorf("argument %s cannot be changed", v)
		}
	}

	request.FlowLogId = &flowLogId
	request.VpcId = &vpcId

	if d.HasChange("flow_log_name") {
		if v, ok := d.GetOk("flow_log_name"); ok {
			request.FlowLogName = helper.String(v.(string))
		}
	}

	if d.HasChange("flow_log_description") {
		if v, ok := d.GetOk("flow_log_description"); ok {
			request.FlowLogDescription = helper.String(v.(string))
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyFlowLogAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc flowLog failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		client := meta.(*TencentCloudClient).apiV3Conn
		tagService := TagService{client}
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		resourceName := BuildTagResourceName("vpc", "fl", client.Region, flowLogId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudVpcFlowLogRead(d, meta)
}

func resourceTencentCloudVpcFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	flowLogId, vpcId, err := resourceTencentCloudGetFlowLogId(d)

	if err != nil {
		return err
	}

	if err := service.DeleteVpcFlowLogById(ctx, flowLogId, vpcId); err != nil {
		return nil
	}

	return nil
}

func resourceTencentCloudSetFlowLogId(d *schema.ResourceData, id, vpcId string) {
	d.SetId(id + FILED_SP + vpcId)
}

func resourceTencentCloudGetFlowLogId(d *schema.ResourceData) (id, vpcId string, err error) {
	rawId := d.Id()
	rawIdRE := regexp.MustCompile(`^fl-[0-9a-z]{8}#vpc-[0-9a-z]{8}$`)
	if !rawIdRE.MatchString(rawId) {
		err = fmt.Errorf("invalid id format %s, expect `fl-xxxxxxxx#vpc-xxxxxxxx`", rawId)
		return
	}
	ids := strings.Split(rawId, FILED_SP)
	id = ids[0]
	vpcId = ids[1]
	return
}
