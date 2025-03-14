package machine_key

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/authn"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveMachineKey(helper.CtxWithOrgID(ctx, d), &management.RemoveMachineKeyRequest{
		UserId: d.Get(UserIDVar).(string),
		KeyId:  d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete machine key: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	keyType := d.Get(keyTypeVar).(string)
	req := &management.AddMachineKeyRequest{
		UserId: d.Get(UserIDVar).(string),
		Type:   authn.KeyType(authn.KeyType_value[keyType]),
	}

	if publicKey, ok := d.GetOk(PublicKeyVar); ok {
		req.PublicKey = []byte(publicKey.(string))
	}

	if expiration, ok := d.GetOk(ExpirationDateVar); ok {
		t, err := time.Parse(time.RFC3339, expiration.(string))
		if err != nil {
			return diag.Errorf("failed to parse time: %v", err)
		}
		req.ExpirationDate = timestamppb.New(t)
	}

	resp, err := client.AddMachineKey(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetKeyId())
	if keyDetails := resp.GetKeyDetails(); keyDetails != nil {
		if err := d.Set(KeyDetailsVar, string(keyDetails)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(UserIDVar).(string)
	resp, err := client.GetMachineKeyByIDs(helper.CtxWithOrgID(ctx, d), &management.GetMachineKeyByIDsRequest{
		UserId: userID,
		KeyId:  d.Id(),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get machine key")
	}

	d.SetId(resp.GetKey().GetId())
	set := map[string]interface{}{
		ExpirationDateVar: resp.GetKey().GetExpirationDate().AsTime().Format(time.RFC3339),
		UserIDVar:         userID,
		helper.OrgIDVar:   orgID,
		keyTypeVar:        resp.GetKey().GetType().String(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of machine key: %v", k, err)
		}
	}
	return nil
}
