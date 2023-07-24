/*
Provides a resource to create a vpc flow_log

~> **NOTE:** The cloud server instance specifications that support stream log collection include: M6ce, M6p, SA3se, S4m, DA3, ITA3, I6t, I6, S5se, SA2, SK1, S4, S5, SN3ne, S3ne, S2ne, SA2a, S3ne, SW3a, SW3b, SW3ne, ITA3, IT5c, IT5, IT5c, IT3, I3, D3, DA2, D2, M6, MA2, M4, C6, IT3a, IT3b, IT3c, C4ne, CN3ne, C3ne, GI1, PNV4, GNV4v, GNV4, GT4, GI3X, GN7, GN7vw.

~> **NOTE:** The following models no longer support the collection of new stream logs, and the stock stream logs will no longer be reported for data from July 25, 2022: Standard models: S3, SA1, S2, S1;Memory type: M3, M2, M1;Calculation type: C4, CN3, C3, C2;Batch calculation type: BC1, BS1;HPCC: HCCIC5, HCCG5v.

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

data "tencentcloud_images" "image" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones.zones.zones.0.name]
  }

  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_cls_logset" "logset" {
  logset_name = "delogsetmo"
  tags        = {
    "test" = "test"
  }
}

resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "topic"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-flow-log-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "vpc-flow-log-subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eni" "example" {
  name        = "vpc-flow-log-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc"
  ipv4_count  = 1
}

resource "tencentcloud_instance" "example" {
  instance_name            = "ci-test-eni-attach"
  availability_zone        = data.tencentcloud_availability_zones.zones.zones.0.name
  image_id                 = data.tencentcloud_images.image.images.0.image_id
  instance_type            = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type         = "CLOUD_PREMIUM"
  disable_security_service = true
  disable_monitor_service  = true
  vpc_id                   = tencentcloud_vpc.vpc.id
  subnet_id                = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_eni_attachment" "example" {
  eni_id      = tencentcloud_eni.example.id
  instance_id = tencentcloud_instance.example.id
}

resource "tencentcloud_vpc_flow_log" "example" {
  flow_log_name        = "tf-example-vpc-flow-log"
  resource_type        = "NETWORKINTERFACE"
  resource_id          = tencentcloud_eni_attachment.example.eni_id
  traffic_type         = "ACCEPT"
  vpc_id               = tencentcloud_vpc.vpc.id
  flow_log_description = "this is a testing flow log"
  cloud_log_id         = tencentcloud_cls_topic.topic.id
  storage_type         = "cls"
  tags                 = {
    "testKey" = "testValue"
  }
}
```

Import

vpc flow_log can be imported using the flow log Id combine vpc Id, e.g.

```
$ terraform import tencentcloud_vpc_flow_log.flow_log flow_log_id fl-xxxx1234#vpc-yyyy5678
```
*/
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
				Description: "Specify flow log storage id.",
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
