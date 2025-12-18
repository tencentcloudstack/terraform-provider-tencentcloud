package antiddos

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosBgpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosBgpInstanceCreate,
		Read:   resourceTencentCloudAntiddosBgpInstanceRead,
		Update: resourceTencentCloudAntiddosBgpInstanceUpdate,
		Delete: resourceTencentCloudAntiddosBgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_charge_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Payment Type: Payment Mode: PREPAID (Prepaid) / POSTPAID_BY_MONTH (Postpaid).",
			},

			"package_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "High-defense package types: Enterprise, Standard, StandardPlus (Standard Edition 2.0).",
			},

			"instance_charge_prepaid": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Prepaid configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Purchase period in months.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "OTIFY_AND_MANUAL_RENEW: Notify the user of the expiration date and do not automatically renew. NOTIFY_AND_AUTO_RENEW: Notify the user of the expiration date and automatically renew. DISABLE_NOTIFY_AND_MANUAL_RENEW: Do not notify the user of the expiration date and do not automatically renew. The default is: Notify the user of the expiration date and do not automatically renew.",
						},
					},
				},
			},

			"enterprise_package_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enterprise package configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The region where the high-defense package was purchased.",
						},
						"protect_ip_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of protected IPs.",
						},
						"basic_protect_bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Guaranteed protection bandwidth.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Service bandwidth scale.",
						},
						"elastic_protect_bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Elastic bandwidth (Gbps), selectable elastic bandwidth [0, 400, 500, 600, 800, 1000], default is 0.",
						},
						"elastic_bandwidth_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable elastic service bandwidth. The default value is false.",
						},
					},
				},
			},

			"standard_package_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Standard package configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The region where the high-defense package was purchased.",
						},
						"protect_ip_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of protected IPs.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Protected service bandwidth 50Mbps.",
						},
						"elastic_bandwidth_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable elastic service bandwidth. The default value is false.",
						},
					},
				},
			},

			"standard_plus_package_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Standard Plus package configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The region where the high-defense package was purchased.",
						},
						"protect_count": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protection Count: TWO_TIMES: Two full-power protections; UNLIMITED: Infinite protections.",
						},
						"protect_ip_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of protected IPs.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "50Mbps protected bandwidth.",
						},
						"elastic_bandwidth_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable elastic service bandwidth. The default value is false.",
						},
					},
				},
			},

			"tag_info_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Prepaid configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAntiddosBgpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_bgp_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = antiddos.NewCreateBgpInstanceRequest()
		response   = antiddos.NewCreateBgpInstanceResponse()
		resourceId string
		region     string
	)

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		request.PackageType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		icpMap := antiddos.InstanceChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			icpMap.Period = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["renew_flag"]; ok {
			icpMap.RenewFlag = helper.String(v.(string))
		}

		request.InstanceChargePrepaid = &icpMap
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "enterprise_package_config"); ok {
		epcMap := antiddos.EnterprisePackageConfig{}
		if v, ok := dMap["region"]; ok {
			epcMap.Region = helper.String(v.(string))
			region = v.(string)
		}

		if v, ok := dMap["protect_ip_count"]; ok {
			epcMap.ProtectIpCount = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["basic_protect_bandwidth"]; ok {
			epcMap.BasicProtectBandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["bandwidth"]; ok {
			epcMap.Bandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["elastic_protect_bandwidth"]; ok {
			epcMap.ElasticProtectBandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["elastic_bandwidth_flag"]; ok {
			epcMap.ElasticBandwidthFlag = helper.Bool(v.(bool))
		}

		request.EnterprisePackageConfig = &epcMap
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "standard_package_config"); ok {
		spcMap := antiddos.StandardPackageConfig{}
		if v, ok := dMap["region"]; ok {
			spcMap.Region = helper.String(v.(string))
			region = v.(string)
		}

		if v, ok := dMap["protect_ip_count"]; ok {
			spcMap.ProtectIpCount = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["bandwidth"]; ok {
			spcMap.Bandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["elastic_bandwidth_flag"]; ok {
			spcMap.ElasticBandwidthFlag = helper.Bool(v.(bool))
		}

		request.StandardPackageConfig = &spcMap
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "standard_plus_package_config"); ok {
		sppcMap := antiddos.StandardPlusPackageConfig{}
		if v, ok := dMap["region"]; ok {
			sppcMap.Region = helper.String(v.(string))
			region = v.(string)
		}

		if v, ok := dMap["protect_count"]; ok {
			sppcMap.ProtectCount = helper.String(v.(string))
		}

		if v, ok := dMap["protect_ip_count"]; ok {
			sppcMap.ProtectIpCount = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["bandwidth"]; ok {
			sppcMap.Bandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["elastic_bandwidth_flag"]; ok {
			sppcMap.ElasticBandwidthFlag = helper.Bool(v.(bool))
		}

		request.StandardPlusPackageConfig = &sppcMap
	}

	request.InstanceCount = helper.IntUint64(1)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateBgpInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ResourceIds == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bgp instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bgp instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.ResourceIds) == 0 {
		return fmt.Errorf("ResourceIds is nil.")
	}

	resourceId = *response.Response.ResourceIds[0]
	d.SetId(strings.Join([]string{resourceId, region}, tccommon.FILED_SP))
	return resourceTencentCloudAntiddosBgpInstanceRead(d, meta)
}

func resourceTencentCloudAntiddosBgpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_bgp_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	resourceId := idSplit[0]
	region := idSplit[1]

	respData, err := service.DescribeAntiddosBgpInstancesById(ctx, resourceId, region)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_antiddos_bgp_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceChargePrepaid != nil {
		_ = d.Set("instance_charge_type", *respData.InstanceChargePrepaid)
	}

	if respData.PackageType != nil {
		_ = d.Set("package_type", *respData.PackageType)
	}

	if respData.InstanceChargePrepaid != nil {
		icpMap := map[string]interface{}{}
		if respData.InstanceChargePrepaid.Period != nil {
			icpMap["period"] = *respData.InstanceChargePrepaid.Period
		}

		if respData.InstanceChargePrepaid.RenewFlag != nil {
			icpMap["renew_flag"] = *respData.InstanceChargePrepaid.RenewFlag
		}

		_ = d.Set("instance_charge_prepaid", []interface{}{icpMap})
	}

	if respData.EnterprisePackageConfig != nil {
		epcMap := map[string]interface{}{}
		if respData.EnterprisePackageConfig.Region != nil {
			epcMap["region"] = *respData.EnterprisePackageConfig.Region
		}

		if respData.EnterprisePackageConfig.ProtectIpCount != nil {
			epcMap["protect_ip_count"] = *respData.EnterprisePackageConfig.ProtectIpCount
		}

		if respData.EnterprisePackageConfig.BasicProtectBandwidth != nil {
			epcMap["basic_protect_bandwidth"] = *respData.EnterprisePackageConfig.BasicProtectBandwidth
		}

		if respData.EnterprisePackageConfig.Bandwidth != nil {
			epcMap["bandwidth"] = *respData.EnterprisePackageConfig.Bandwidth
		}

		if respData.EnterprisePackageConfig.ElasticProtectBandwidth != nil {
			epcMap["elastic_protect_bandwidth"] = *respData.EnterprisePackageConfig.ElasticProtectBandwidth
		}

		if respData.EnterprisePackageConfig.ElasticBandwidthFlag != nil {
			epcMap["elastic_bandwidth_flag"] = *respData.EnterprisePackageConfig.ElasticBandwidthFlag
		}

		_ = d.Set("enterprise_package_config", []interface{}{epcMap})
	}

	return nil
}

func resourceTencentCloudAntiddosBgpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_bgp_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudAntiddosBgpInstanceRead(d, meta)
}

func resourceTencentCloudAntiddosBgpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_bgp_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
