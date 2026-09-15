package bdrc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBdrcDisasterRecoveryVpcMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBdrcDisasterRecoveryVpcMappingCreate,
		Read:   resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead,
		Update: resourceTencentCloudBdrcDisasterRecoveryVpcMappingUpdate,
		Delete: resourceTencentCloudBdrcDisasterRecoveryVpcMappingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"site_pair_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site pair ID of the disaster recovery VPC mapping.",
			},

			"source_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source VPC ID of the disaster recovery VPC mapping.",
			},

			"source_subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source subnet ID of the disaster recovery VPC mapping.",
			},

			"target_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target VPC ID of the disaster recovery VPC mapping.",
			},

			"target_subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target subnet ID of the disaster recovery VPC mapping.",
			},

			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Mapping rule primary key ID.",
			},

			"source_vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source VPC ID returned by the cloud API.",
			},

			"source_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source subnet ID returned by the cloud API.",
			},

			"target_vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target VPC ID returned by the cloud API.",
			},

			"target_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target subnet ID returned by the cloud API.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mapping status.",
			},

			"life_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lifecycle state.",
			},
		},
	}
}

func resourceTencentCloudBdrcDisasterRecoveryVpcMappingCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_vpc_mapping.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request        = bdrcv20260330.NewCreateDisasterRecoveryVpcMappingRequest()
		sitePairId     string
		sourceVpcId    string
		sourceSubnetId string
		targetVpcId    string
		targetSubnetId string
	)

	if v, ok := d.GetOk("site_pair_id"); ok {
		request.SitePairId = helper.String(v.(string))
		sitePairId = v.(string)
	}

	if v, ok := d.GetOk("source_vpc_id"); ok {
		request.SourceVpcId = helper.String(v.(string))
		sourceVpcId = v.(string)
	}

	if v, ok := d.GetOk("source_subnet_id"); ok {
		request.SourceSubnetId = helper.String(v.(string))
		sourceSubnetId = v.(string)
	}

	if v, ok := d.GetOk("target_vpc_id"); ok {
		request.TargetVpcId = helper.String(v.(string))
		targetVpcId = v.(string)
	}

	if v, ok := d.GetOk("target_subnet_id"); ok {
		request.TargetSubnetId = helper.String(v.(string))
		targetSubnetId = v.(string)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().CreateDisasterRecoveryVpcMappingWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bdrc_disaster_recovery_vpc_mapping failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bdrc_disaster_recovery_vpc_mapping failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	vpcMappingIdStr := ""
	describeRequest := bdrcv20260330.NewDescribeVpcMappingsRequest()
	describeRequest.SitePairId = helper.String(sitePairId)
	describeRequest.Limit = helper.Int64(100)
	describeRequest.Filters = []*bdrcv20260330.FilterModel{
		{
			Name:   helper.String("source-vpc-id"),
			Values: []*string{&sourceVpcId},
		},
		{
			Name:   helper.String("source-subnet-id"),
			Values: []*string{&sourceSubnetId},
		},
	}

	describeErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().DescribeVpcMappingsWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bdrc_disaster_recovery_vpc_mapping failed, Response is nil."))
		}

		if len(result.Response.VpcMappingSet) == 0 {
			return resource.RetryableError(fmt.Errorf("bdrc_disaster_recovery_vpc_mapping not found yet, waiting for eventual consistency"))
		}

		var matched *bdrcv20260330.VpcMapping
		for _, item := range result.Response.VpcMappingSet {
			if item == nil {
				continue
			}
			if item.SourceVpc != nil && *item.SourceVpc == sourceVpcId &&
				item.SourceSubnet != nil && *item.SourceSubnet == sourceSubnetId &&
				item.TargetVpc != nil && *item.TargetVpc == targetVpcId &&
				item.TargetSubnet != nil && *item.TargetSubnet == targetSubnetId {
				if matched != nil {
					return resource.NonRetryableError(fmt.Errorf("multiple bdrc_disaster_recovery_vpc_mapping matched the same source/target, please check and use import instead."))
				}
				matched = item
			}
		}

		if matched == nil {
			return resource.RetryableError(fmt.Errorf("bdrc_disaster_recovery_vpc_mapping not found yet, waiting for eventual consistency"))
		}

		if matched.Id == nil {
			return resource.NonRetryableError(fmt.Errorf("bdrc_disaster_recovery_vpc_mapping Id is nil."))
		}

		vpcMappingIdStr = strconv.FormatUint(*matched.Id, 10)
		return nil
	})

	if describeErr != nil {
		log.Printf("[CRITAL]%s describe bdrc_disaster_recovery_vpc_mapping after create failed, reason:%+v", logId, describeErr)
		return describeErr
	}

	log.Printf("[CRITAL]%s create bdrc_disaster_recovery_vpc_mapping success, vpcMappingId=%s", logId, vpcMappingIdStr)
	d.SetId(strings.Join([]string{sitePairId, vpcMappingIdStr}, tccommon.FILED_SP))
	return resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead(d, meta)
}

func resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_vpc_mapping.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	sitePairId := idSplit[0]
	vpcMappingId := idSplit[1]

	request := bdrcv20260330.NewDescribeVpcMappingsRequest()
	request.SitePairId = helper.String(sitePairId)
	request.Limit = helper.Int64(100)

	var matched *bdrcv20260330.VpcMapping
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().DescribeVpcMappingsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bdrc_disaster_recovery_vpc_mapping failed, Response is nil."))
		}

		for _, item := range result.Response.VpcMappingSet {
			if item == nil || item.Id == nil {
				continue
			}
			if strconv.FormatUint(*item.Id, 10) == vpcMappingId {
				matched = item
				break
			}
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read bdrc_disaster_recovery_vpc_mapping failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if matched == nil {
		log.Printf("[CRUD] bdrc_disaster_recovery_vpc_mapping id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if matched.Id != nil {
		_ = d.Set("id", *matched.Id)
	}

	if matched.SitePairId != nil {
		_ = d.Set("site_pair_id", *matched.SitePairId)
	}

	if matched.SourceVpc != nil {
		_ = d.Set("source_vpc", *matched.SourceVpc)
	}

	if matched.SourceSubnet != nil {
		_ = d.Set("source_subnet", *matched.SourceSubnet)
	}

	if matched.TargetVpc != nil {
		_ = d.Set("target_vpc", *matched.TargetVpc)
	}

	if matched.TargetSubnet != nil {
		_ = d.Set("target_subnet", *matched.TargetSubnet)
	}

	if matched.Status != nil {
		_ = d.Set("status", *matched.Status)
	}

	if matched.LifeState != nil {
		_ = d.Set("life_state", *matched.LifeState)
	}

	return nil
}

func resourceTencentCloudBdrcDisasterRecoveryVpcMappingUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_vpc_mapping.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	needChange := false
	immutableArgs := []string{"site_pair_id", "source_vpc_id", "source_subnet_id", "target_vpc_id", "target_subnet_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		return fmt.Errorf("Update bdrc_disaster_recovery_vpc_mapping is not supported, all business arguments are immutable (CRD-only API), please recreate the resource.")
	}

	return resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead(d, meta)
}

func resourceTencentCloudBdrcDisasterRecoveryVpcMappingDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_vpc_mapping.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bdrcv20260330.NewDeleteDisasterRecoveryVpcMappingRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	vpcMappingId := idSplit[1]
	vpcMappingIdPoint := helper.StrToUint64Point(vpcMappingId)
	request.VpcMappingIds = []*uint64{vpcMappingIdPoint}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().DeleteDisasterRecoveryVpcMappingWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete bdrc_disaster_recovery_vpc_mapping failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bdrc_disaster_recovery_vpc_mapping failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
