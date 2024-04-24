package teo

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0(ctx context.Context, resp *teo.AccelerationDomain) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		zoneId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	accelerationDomain := resp
	if accelerationDomain.Certificate != nil {
		certificate := accelerationDomain.Certificate
		zone, err := service.DescribeTeoZone(ctx, zoneId)
		if err != nil {
			return err
		}

		serverCertInfoList := []interface{}{}
		for _, serverCertInfo := range certificate.List {
			serverCertInfoMap := map[string]interface{}{}

			if serverCertInfo.CertId != nil {
				serverCertInfoMap["cert_id"] = serverCertInfo.CertId
			}

			if serverCertInfo.Alias != nil {
				serverCertInfoMap["alias"] = serverCertInfo.Alias
			}

			if serverCertInfo.Type != nil {
				serverCertInfoMap["type"] = serverCertInfo.Type
			}

			if serverCertInfo.ExpireTime != nil {
				serverCertInfoMap["expire_time"] = serverCertInfo.ExpireTime
			}

			if serverCertInfo.DeployTime != nil {
				serverCertInfoMap["deploy_time"] = serverCertInfo.DeployTime
			}

			if serverCertInfo.SignAlgo != nil {
				serverCertInfoMap["sign_algo"] = serverCertInfo.SignAlgo
			}

			if zone.ZoneName != nil {
				serverCertInfoMap["common_name"] = zone.ZoneName
			}

			serverCertInfoList = append(serverCertInfoList, serverCertInfoMap)
		}

		_ = d.Set("server_cert_info", serverCertInfoList)

		if certificate.Mode != nil {
			_ = d.Set("mode", certificate.Mode)
		}
	}

	return nil
}

func resourceTencentCloudTeoCertificateConfigUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		zoneId string
		host   string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}

	err := service.CheckAccelerationDomainStatus(ctx, zoneId, host, "")
	if err != nil {
		return err
	}

	return nil
}
