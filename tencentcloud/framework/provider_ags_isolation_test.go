package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

type isolationMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *isolationMeta) GetAPIV3Conn() *connectivity.TencentCloudClient { return m.client }

func TestAGSFrameworkConfigurationKeepsOwnAliasClient(t *testing.T) {
	t.Cleanup(sharedmeta.ResetSharedMetaForTest)
	ctx := context.Background()
	makePair := func(region string) (*sdk.Provider, *Provider, *connectivity.TencentCloudClient) {
		client := &connectivity.TencentCloudClient{Region: region}
		primary := &sdk.Provider{ConfigureFunc: func(*sdk.ResourceData) (any, error) {
			sharedmeta.SetSharedMeta(client)
			return &isolationMeta{client}, nil
		}}
		return primary, NewProvider(primary).(*Provider), client
	}
	a, fa, ca := makePair("ap-guangzhou")
	b, fb, cb := makePair("ap-shanghai")
	if _, err := a.ConfigureFunc(nil); err != nil {
		t.Fatal(err)
	}
	if _, err := b.ConfigureFunc(nil); err != nil {
		t.Fatal(err)
	}
	// Deliberately interleave SDKv2 A, SDKv2 B, Framework A, Framework B.
	for _, item := range []struct {
		p    *Provider
		want *connectivity.TencentCloudClient
	}{{fa, ca}, {fb, cb}} {
		resp := provider.ConfigureResponse{}
		item.p.Configure(ctx, provider.ConfigureRequest{}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatal(resp.Diagnostics)
		}
		meta, ok := resp.ResourceData.(*sharedmeta.ProviderMeta)
		if !ok || meta.Client != item.want {
			t.Fatalf("alias client crossed configuration boundary: got=%v want_region=%s", meta, item.want.Region)
		}
	}
}
