package cvm

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
	"log"
)

func resourceTencentCloudEipCreatePostFillRequest0(ctx context.Context, req *vpc.AllocateAddressesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			req.Tags = append(req.Tags, &tag)
		}
	}

	var internetChargeType string
	if v, ok := d.GetOk("internet_charge_type"); ok {
		internetChargeType = v.(string)
	}

	if internetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
		addressChargePrepaid := vpc.AddressChargePrepaid{}
		period := d.Get("prepaid_period")
		renewFlag := d.Get("auto_renew_flag")
		addressChargePrepaid.Period = helper.IntInt64(period.(int))
		addressChargePrepaid.AutoRenewFlag = helper.IntInt64(renewFlag.(int))
		req.AddressChargePrepaid = &addressChargePrepaid
	}

	return nil
}

func resourceTencentCloudEipCreatePostHandleResponse0(ctx context.Context, resp *vpc.AllocateAddressesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	logId := tccommon.GetLogId(tccommon.ContextNil)
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	tagService := svctag.NewTagService(client)
	region := client.Region
	eipId := *resp.Response.AddressSet[0]
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName(svcvpc.VPC_SERVICE_TYPE, svcvpc.EIP_RESOURCE_TYPE, region, eipId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			log.Printf("[CRITAL]%s set eip tags failed: %+v", logId, err)
			return err
		}
	}

	// wait for status
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		if eip != nil && *eip.AddressStatus == svcvpc.EIP_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("eip is still creating"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudEipReadPostHandleResponse0(ctx context.Context, resp *vpc.Address) error {
	d := tccommon.ResourceDataFromContext(ctx)
	logId := tccommon.GetLogId(tccommon.ContextNil)
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	tagService := svctag.NewTagService(client)
	region := client.Region
	eipId := d.Id()

	tags, err := tagService.DescribeResourceTags(ctx, svcvpc.VPC_SERVICE_TYPE, svcvpc.EIP_RESOURCE_TYPE, region, eipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
		return err
	}

	bgp, err := vpcService.DescribeVpcBandwidthPackageByEip(ctx, eipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
		return err
	}

	_ = d.Set("tags", tags)
	if bgp != nil {
		_ = d.Set("bandwidth_package_id", bgp.BandwidthPackageId)
	}

	return nil
}

func resourceTencentCloudEipUpdatePostFillRequest1(ctx context.Context, req *vpc.ModifyAddressInternetChargeTypeRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	period := d.Get("prepaid_period").(int)
	renewFlag := d.Get("auto_renew_flag").(int)

	if *req.InternetChargeType != "" && *req.InternetMaxBandwidthOut != 0 {
		if *req.InternetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
			addressChargePrepaid := vpc.AddressChargePrepaid{}
			addressChargePrepaid.AutoRenewFlag = helper.IntInt64(renewFlag)
			addressChargePrepaid.Period = helper.IntInt64(period)
			req.AddressChargePrepaid = &addressChargePrepaid
		}
	}

	return nil
}

func resourceTencentCloudEipDeletePostFillRequest0(ctx context.Context, req *vpc.ReleaseAddressesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := vpcService.UnattachEip(ctx, d.Id())
		if errRet != nil {
			return tccommon.RetryError(errRet, "DesOperation.MutexTaskRunning")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudEipDeletePostHandleResponse0(ctx context.Context, resp *vpc.ReleaseAddressesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	eipId := d.Id()
	var internetChargeType string
	if v, ok := d.GetOk("internet_charge_type"); ok {
		internetChargeType = v.(string)
	}

	if internetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
		// isolated
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			eip, errRet := vpcService.DescribeEipById(ctx, eipId)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			if !*eip.IsArrears {
				return resource.RetryableError(fmt.Errorf("eip is still isolate"))
			}
			return nil
		})
		if err != nil {
			return err
		}

		// release
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := vpcService.DeleteEip(ctx, eipId)
			if errRet != nil {
				return tccommon.RetryError(errRet, "DesOperation.MutexTaskRunning", "OperationDenied.MutexTaskRunning")
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		if eip != nil {
			return resource.RetryableError(fmt.Errorf("eip is still deleting"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudEipUpdateOnStart(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	eipId := d.Id()
	if d.HasChange("prepaid_period") || d.HasChange("auto_renew_flag") {
		period := d.Get("prepaid_period").(int)
		renewFlag := d.Get("auto_renew_flag").(int)
		err := vpcService.RenewAddress(ctx, eipId, period, renewFlag)
		if err != nil {
			return err
		}
	}
	return nil
}
