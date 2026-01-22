package tke

import (
	"context"
	"fmt"
	"time"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudKubernetesEncryptionProtectionCreatePostHandleResponse0(ctx context.Context, resp *tke.EnableEncryptionProtectionResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf(
		[]string{},
		[]string{"Opened"},
		5*tccommon.ReadRetryTimeout,
		time.Second,
		service.TkeEncryptionProtectionStateRefreshFunc(clusterId, []string{}),
	)

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}

func resourceTencentCloudKubernetesEncryptionProtectionDeletePostHandleResponse0(ctx context.Context, resp *tke.DisableEncryptionProtectionResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf(
		[]string{},
		[]string{"Closed"},
		5*tccommon.ReadRetryTimeout,
		time.Second,
		service.TkeEncryptionProtectionStateRefreshFunc(d.Id(), []string{}),
	)

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}
