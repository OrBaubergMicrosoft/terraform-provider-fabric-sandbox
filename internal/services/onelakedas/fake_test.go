// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedas_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func newRandomDataAccessRoleBase() fabcore.DataAccessRoleBase {
	return fabcore.DataAccessRoleBase{
		Name: new("TestRole"),
		Kind: azto.Ptr(fabcore.DataAccessRoleKindPolicy),
		DecisionRules: []fabcore.DecisionRule{
			{
				Effect: azto.Ptr(fabcore.EffectPermit),
				Permission: []fabcore.PermissionScope{
					{
						AttributeName:            azto.Ptr(fabcore.AttributeNamePath),
						AttributeValueIncludedIn: []string{"Tables/TestTable"},
					},
				},
			},
		},
		Members: &fabcore.Members{
			MicrosoftEntraMembers: []fabcore.MicrosoftEntraMember{
				{
					ObjectID:   new(testhelp.RandomUUID()),
					TenantID:   new(testhelp.RandomUUID()),
					ObjectType: azto.Ptr(fabcore.ObjectTypeUser),
				},
			},
		},
	}
}

func newRandomDataAccessRoleListItem(base fabcore.DataAccessRoleBase) fabcore.DataAccessRoleListItem {
	return fabcore.DataAccessRoleListItem{
		Name:          base.Name,
		Kind:          base.Kind,
		DecisionRules: base.DecisionRules,
		Members:       base.Members,
	}
}

func fakeCreateOrUpdateSingleDataAccessRole() func(ctx context.Context, workspaceID, itemID string, createOrUpdateSingleDataAccessRoleRequest fabcore.DataAccessRoleBase, options *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ fabcore.DataAccessRoleBase, _ *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetDataAccessRole(
	entity fabcore.DataAccessRoleBase,
) func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse{
			DataAccessRoleBase: entity,
		}, nil)

		return resp, errResp
	}
}

func fakeDeleteDataAccessRole() func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse{}, nil)

		return resp, errResp
	}
}

func fakeListDataAccessRoles(
	entities []fabcore.DataAccessRoleListItem,
) func(ctx context.Context, workspaceID, itemID string, options *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse{
			DataAccessRoles: fabcore.DataAccessRoles{
				Value: entities,
			},
		}, nil)

		return resp, errResp
	}
}
