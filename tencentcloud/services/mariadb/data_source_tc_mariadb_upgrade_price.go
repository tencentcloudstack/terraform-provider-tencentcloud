package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbUpgradePrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbUpgradePriceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory size in GB, which can be obtained by querying the instance specification through the `DescribeDBInstanceSpecs` API.",
			},
			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Storage capacity in GB. The maximum and minimum storage space can be obtained by querying instance specification through the `DescribeDBInstanceSpecs` API.",
			},
			"node_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "New instance nodes, zero means not change.",
			},
			"amount_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Price unit. Valid values: `* pent` (cent), `* microPent` (microcent).",
			},
			"original_price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Original price * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).",
			},
			"price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The actual price may be different from the original price due to discounts. * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).",
			},
			"formula": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Price calculation formula.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMariadbUpgradePriceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_upgrade_price.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		price      *mariadb.DescribeUpgradePriceResponseParams
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("memory"); v != nil {
		paramMap["Memory"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("storage"); v != nil {
		paramMap["Storage"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("node_count"); v != nil {
		paramMap["NodeCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbUpgradePriceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		price = result
		return nil
	})

	if err != nil {
		return err
	}

	if price.OriginalPrice != nil {
		_ = d.Set("original_price", price.OriginalPrice)
	}

	if price.Price != nil {
		_ = d.Set("price", price.Price)
	}

	if price.Formula != nil {
		_ = d.Set("formula", price.Formula)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
