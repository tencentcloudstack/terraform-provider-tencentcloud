package cvm

import (
	"context"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudCvmChcConfigCreatePostHandleResponse3(ctx context.Context, resp *cvm.ConfigureChcAssistVpcResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		chcId string
	)
	if v, ok := d.GetOk("chc_id"); ok {
		chcId = v.(string)
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"READY"}, 20*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(chcId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func resourceTencentCloudCvmChcConfigCreatePostHandleResponse4(ctx context.Context, resp *cvm.ConfigureChcAssistVpcResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		chcId string
		vpcId string
	)
	if v, ok := d.GetOk("chc_id"); ok {
		chcId = v.(string)
	}
	if dMap, ok := helper.InterfacesHeadMap(d, "deploy_virtual_private_cloud"); ok {
		if v, ok := dMap["vpc_id"]; ok {
			vpcId = v.(string)
		}
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{vpcId}, 10*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(chcId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func resourceTencentCloudCvmChcConfigDeletePostHandleResponse0(ctx context.Context, resp *cvm.RemoveChcDeployVpcResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	chcId := d.Id()

	conf := tccommon.BuildStateChangeConf([]string{}, []string{""}, 5*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(chcId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}
	if len(chcHosts) > 0 && *chcHosts[0].InstanceState == "INIT" {
		return nil
	}

	return nil
}

func resourceTencentCloudCvmChcConfigDeletePostHandleResponse1(ctx context.Context, resp *cvm.RemoveChcAssistVpcResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"INIT"}, 10*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
