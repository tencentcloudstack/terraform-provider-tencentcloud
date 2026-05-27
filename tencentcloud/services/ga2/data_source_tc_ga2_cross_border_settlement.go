package ga2

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGa2CrossBorderSettlement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGa2CrossBorderSettlementRead,
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Global accelerator instance ID.",
			},
			"accelerate_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Acceleration region.",
			},
			"endpoint_group_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint group region.",
			},
			"settlement_month": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Billing year-month time.",
			},
			"traffic": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Traffic usage in GB with 6 decimal places precision.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudGa2CrossBorderSettlementRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ga2_cross_border_settlement.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(nil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client              = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		globalAcceleratorId = d.Get("global_accelerator_id").(string)
		accelerateRegion    = d.Get("accelerate_region").(string)
		endpointGroupRegion = d.Get("endpoint_group_region").(string)
		settlementMonth     = d.Get("settlement_month").(int)
	)

	request := ga2v20250115.NewDescribeCrossBorderSettlementRequest()
	request.GlobalAcceleratorId = helper.String(globalAcceleratorId)
	request.AccelerateRegion = helper.String(accelerateRegion)
	request.EndpointGroupRegion = helper.String(endpointGroupRegion)
	request.SettlementMonth = helper.Uint64(uint64(settlementMonth))

	var response *ga2v20250115.DescribeCrossBorderSettlementResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := client.UseGa2V20250115Client().DescribeCrossBorderSettlementWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeCrossBorderSettlement response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe ga2 cross border settlement failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Traffic != nil {
		_ = d.Set("traffic", *response.Response.Traffic)
	}

	d.SetId(globalAcceleratorId + tccommon.FILED_SP + accelerateRegion + tccommon.FILED_SP + endpointGroupRegion + tccommon.FILED_SP + strconv.Itoa(settlementMonth))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
