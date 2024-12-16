package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func ResourceTencentCloudClassicElasticPublicIpv6() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClassicElasticPublicIpv6Create,
		Read:   resourceTencentCloudClassicElasticPublicIpv6Read,
		Update: resourceTencentCloudClassicElasticPublicIpv6Update,
		Delete: resourceTencentCloudClassicElasticPublicIpv6Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ip6_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IPV6 addresses that require public network access.",
			},

			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Bandwidth in Mbps. Default is 1Mbps.",
			},

			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Network billing model. IPV6 currently supports `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. The default network charging mode is `TRAFFIC_POSTPAID_BY_HOUR`.",
			},

			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Bandwidth package id, move the account up, and you need to pass in the ipv6 address to apply for bandwidth package charging mode.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tags.",
			},
		},
	}
}

func resourceTencentCloudClassicElasticPublicIpv6Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_classic_elastic_public_ipv6.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		ipId string
	)
	var (
		request  = vpc.NewAllocateIp6AddressesBandwidthRequest()
		response = vpc.NewAllocateIp6AddressesBandwidthResponse()
	)

	if v, ok := d.GetOk("ip6_address"); ok {
		request.Ip6Addresses = []*string{helper.String(v.(string))}
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.BandwidthPackageId = helper.String(v.(string))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AllocateIp6AddressesBandwidthWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create classic elastic public ipv6 failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.AddressSet) < 1 {
		return fmt.Errorf("resource `tencentcloud_classic_elastic_public_ipv6` create failed.")
	}

	ipId = *response.Response.AddressSet[0]
	taskId := *response.Response.TaskId

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourceClassicElasticPublicIpv6CreateStateRefreshFunc_0_0(ctx, taskId),
		Target:     []string{"SUCCESS"},
		Timeout:    600 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}
	d.SetId(ipId)

	return resourceTencentCloudClassicElasticPublicIpv6Read(d, meta)
}

func resourceTencentCloudClassicElasticPublicIpv6Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_classic_elastic_public_ipv6.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ipId := d.Id()

	respData, err := service.DescribeClassicElasticPublicIpv6ById(ctx, ipId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `classic_elastic_public_ipv6` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if len(respData.AddressSet) > 0 {
		address := respData.AddressSet[0]
		_ = d.Set("ip6_address", address.AddressIp)
		_ = d.Set("internet_max_bandwidth_out", address.Bandwidth)
		_ = d.Set("internet_charge_type", address.InternetChargeType)
		_ = d.Set("bandwidth_package_id", address.BandwidthPackageId)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	region := tcClient.Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "eip", region, ipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eipv6 tags failed: %+v", logId, err)
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudClassicElasticPublicIpv6Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_classic_elastic_public_ipv6.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"ip6_address", "internet_charge_type", "bandwidth_package_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	ipId := d.Id()

	needChange := false
	mutableArgs := []string{"internet_max_bandwidth_out"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := vpc.NewModifyIp6AddressesBandwidthRequest()

		if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
			request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
		}

		request.Ip6AddressIds = []*string{helper.String(ipId)}

		var taskId *string
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyIp6AddressesBandwidthWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			taskId = result.Response.TaskId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update classic elastic public ipv6 failed, reason:%+v", logId, err)
			return err
		}
		if taskId == nil {
			return fmt.Errorf("taskId is nil")
		}
		if _, err := (&resource.StateChangeConf{
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
			Pending:    []string{},
			Refresh:    resourceElasticPublicIpv6CreateStateRefreshFunc_0_0(ctx, *taskId),
			Target:     []string{"SUCCESS"},
			Timeout:    600 * time.Second,
		}).WaitForStateContext(ctx); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		region := tcClient.Region
		resourceName := tccommon.BuildTagResourceName("vpc", "eip", region, ipId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			log.Printf("[CRITAL]%s update eipv6 tags failed: %+v", logId, err)
			return err
		}

	}

	return resourceTencentCloudClassicElasticPublicIpv6Read(d, meta)
}

func resourceTencentCloudClassicElasticPublicIpv6Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_classic_elastic_public_ipv6.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	ipId := d.Id()

	var (
		request = vpc.NewReleaseIp6AddressesBandwidthRequest()
		taskId  *string
	)

	request.Ip6AddressIds = []*string{helper.String(ipId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReleaseIp6AddressesBandwidthWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		taskId = result.Response.TaskId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete classic elastic public ipv6 failed, reason:%+v", logId, err)
		return err
	}

	if taskId == nil {
		return fmt.Errorf("taskId is nil")
	}
	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourceElasticPublicIpv6CreateStateRefreshFunc_0_0(ctx, *taskId),
		Target:     []string{"SUCCESS"},
		Timeout:    600 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}

func resourceClassicElasticPublicIpv6CreateStateRefreshFunc_0_0(ctx context.Context, taskId string) resource.StateRefreshFunc {
	var req *vpc.DescribeTaskResultRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}
		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}
			_ = d
			req = vpc.NewDescribeTaskResultRequest()
			req.TaskId = helper.StrToUint64Point(taskId)

		}
		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeTaskResultWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}
		state := fmt.Sprintf("%v", *resp.Response.Result)
		return resp.Response, state, nil
	}
}
