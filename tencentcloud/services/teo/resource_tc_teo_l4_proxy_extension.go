package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoL4ProxyCreatePostHandleResponse0(ctx context.Context, resp *teo.CreateL4ProxyResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		zoneId  string
		proxyId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	proxyId = *resp.Response.ProxyId

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"online"}, 10*tccommon.ReadRetryTimeout, time.Second, teoL4proxyStateRefreshFunc(meta, zoneId, proxyId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func resourceTencentCloudTeoL4ProxyDeletePostFillRequest0(ctx context.Context, req *teo.DeleteL4ProxyRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	proxy, err := service.DescribeTeoL4ProxyById(ctx, zoneId, proxyId)
	if err != nil {
		return err
	}

	if *proxy.Status == "online" {
		logId := tccommon.GetLogId(tccommon.ContextNil)

		request := teo.NewModifyL4ProxyStatusRequest()

		request.ZoneId = &zoneId
		request.ProxyId = &proxyId
		request.Status = helper.String("offline")

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyL4ProxyStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo l4 proxy failed, reason:%+v", logId, err)
			return err
		}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"offline"}, 10*tccommon.ReadRetryTimeout, time.Second, teoL4proxyStateRefreshFunc(meta, zoneId, proxyId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	return nil
}

func resourceTencentCloudTeoL4ProxyUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	status := "online"
	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{status}, 10*tccommon.ReadRetryTimeout, time.Second, teoL4proxyStateRefreshFunc(meta, zoneId, proxyId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func teoL4proxyStateRefreshFunc(meta interface{}, zoneId, proxyId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		object, err := service.DescribeTeoL4ProxyById(ctx, zoneId, proxyId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}
