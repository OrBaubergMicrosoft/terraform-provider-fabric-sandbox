// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedas

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseOneLakeDataAccessSecurityModel struct {
	WorkspaceID   customtypes.UUID                                      `tfsdk:"workspace_id"`
	ItemID        customtypes.UUID                                      `tfsdk:"item_id"`
	Name          types.String                                          `tfsdk:"name"`
	Kind          types.String                                          `tfsdk:"kind"`
	DecisionRules supertypes.ListNestedObjectValueOf[decisionRuleModel] `tfsdk:"decision_rules"`
	Members       supertypes.SingleNestedObjectValueOf[membersModel]    `tfsdk:"members"`
}

func (to *baseOneLakeDataAccessSecurityModel) set(ctx context.Context, workspaceID, itemID string, from fabcore.DataAccessRoleBase) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	to.Name = types.StringPointerValue(from.Name)
	to.Kind = types.StringPointerValue((*string)(from.Kind))

	// DecisionRules
	to.DecisionRules.SetNull(ctx)

	if from.DecisionRules != nil {
		drSlice := make([]*decisionRuleModel, 0, len(from.DecisionRules))

		for _, dr := range from.DecisionRules {
			drModel := &decisionRuleModel{}

			if diags := drModel.set(ctx, dr); diags.HasError() {
				return diags
			}

			drSlice = append(drSlice, drModel)
		}

		if diags := to.DecisionRules.Set(ctx, drSlice); diags.HasError() {
			return diags
		}
	}

	// Members
	to.Members.SetNull(ctx)

	if from.Members != nil {
		membersModel := &membersModel{}

		if diags := membersModel.set(ctx, *from.Members); diags.HasError() {
			return diags
		}

		if diags := to.Members.Set(ctx, membersModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceOneLakeDataAccessSecurityModel struct {
	baseOneLakeDataAccessSecurityModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceOneLakeDataAccessSecuritiesModel struct {
	WorkspaceID customtypes.UUID                                                      `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                                      `tfsdk:"item_id"`
	Values      supertypes.SetNestedObjectValueOf[baseOneLakeDataAccessSecurityModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                                       `tfsdk:"timeouts"`
}

func (to *dataSourceOneLakeDataAccessSecuritiesModel) setValues(ctx context.Context, workspaceID, itemID string, from []fabcore.DataAccessRoleListItem) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	slice := make([]*baseOneLakeDataAccessSecurityModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseOneLakeDataAccessSecurityModel

		// Convert DataAccessRoleListItem to DataAccessRoleBase for reuse
		base := fabcore.DataAccessRoleBase{
			DecisionRules: entity.DecisionRules,
			Name:          entity.Name,
			Kind:          entity.Kind,
			Members:       entity.Members,
		}

		if diags := entityModel.set(ctx, workspaceID, itemID, base); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceOneLakeDataAccessSecurityModel struct {
	baseOneLakeDataAccessSecurityModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateOrUpdateOneLakeDataAccessSecurity struct {
	fabcore.DataAccessRoleBase
}

func (to *requestCreateOrUpdateOneLakeDataAccessSecurity) set(ctx context.Context, from resourceOneLakeDataAccessSecurityModel) diag.Diagnostics {
	to.Name = from.Name.ValueStringPointer()
	to.Kind = (*fabcore.DataAccessRoleKind)(from.Kind.ValueStringPointer())

	// DecisionRules
	decisionRules, diags := from.DecisionRules.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.DecisionRules = make([]fabcore.DecisionRule, 0, len(decisionRules))

	for _, drModel := range decisionRules {
		dr := fabcore.DecisionRule{
			Effect: (*fabcore.Effect)(drModel.Effect.ValueStringPointer()),
		}

		// Permission
		permModels, d := drModel.Permission.Get(ctx)
		if d.HasError() {
			return d
		}

		dr.Permission = make([]fabcore.PermissionScope, 0, len(permModels))

		for _, pm := range permModels {
			ps := fabcore.PermissionScope{
				AttributeName: (*fabcore.AttributeName)(pm.AttributeName.ValueStringPointer()),
			}

			elements, d := pm.AttributeValueIncludedIn.Get(ctx)
			if d.HasError() {
				return d
			}

			for _, e := range elements {
				ps.AttributeValueIncludedIn = append(ps.AttributeValueIncludedIn, e.ValueString())
			}

			dr.Permission = append(dr.Permission, ps)
		}

		// Constraints
		cPtr, d := drModel.Constraints.Get(ctx)
		if d.HasError() {
			return d
		}

		if cPtr != nil {
			constraints := fabcore.DecisionRuleConstraints{}

			colModels, d := cPtr.Columns.Get(ctx)
			if d.HasError() {
				return d
			}

			if colModels != nil {
				constraints.Columns = make([]fabcore.ColumnConstraint, 0, len(colModels))

				for _, cm := range colModels {
					cc := fabcore.ColumnConstraint{
						ColumnEffect: (*fabcore.ColumnEffect)(cm.ColumnEffect.ValueStringPointer()),
						TablePath:    cm.TablePath.ValueStringPointer(),
					}

					actions, d := cm.ColumnAction.Get(ctx)
					if d.HasError() {
						return d
					}

					for _, a := range actions {
						cc.ColumnAction = append(cc.ColumnAction, fabcore.ColumnAction(a.ValueString()))
					}

					names, d := cm.ColumnNames.Get(ctx)
					if d.HasError() {
						return d
					}

					for _, n := range names {
						cc.ColumnNames = append(cc.ColumnNames, n.ValueString())
					}

					constraints.Columns = append(constraints.Columns, cc)
				}
			}

			rowModels, d := cPtr.Rows.Get(ctx)
			if d.HasError() {
				return d
			}

			if rowModels != nil {
				constraints.Rows = make([]fabcore.RowConstraint, 0, len(rowModels))

				for _, rm := range rowModels {
					constraints.Rows = append(constraints.Rows, fabcore.RowConstraint{
						TablePath: rm.TablePath.ValueStringPointer(),
						Value:     rm.Value.ValueStringPointer(),
					})
				}
			}

			dr.Constraints = &constraints
		}

		to.DecisionRules = append(to.DecisionRules, dr)
	}

	// Members
	membersPtr, diags := from.Members.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if membersPtr != nil {
		members := fabcore.Members{}

		fimModels, d := membersPtr.FabricItemMembers.Get(ctx)
		if d.HasError() {
			return d
		}

		if fimModels != nil {
			members.FabricItemMembers = make([]fabcore.FabricItemMember, 0, len(fimModels))

			for _, fimModel := range fimModels {
				fim := fabcore.FabricItemMember{
					SourcePath: fimModel.SourcePath.ValueStringPointer(),
				}

				accesses, d := fimModel.ItemAccess.Get(ctx)
				if d.HasError() {
					return d
				}

				for _, a := range accesses {
					fim.ItemAccess = append(fim.ItemAccess, fabcore.ItemAccess(a.ValueString()))
				}

				members.FabricItemMembers = append(members.FabricItemMembers, fim)
			}
		}

		memModels, d := membersPtr.MicrosoftEntraMembers.Get(ctx)
		if d.HasError() {
			return d
		}

		if memModels != nil {
			members.MicrosoftEntraMembers = make([]fabcore.MicrosoftEntraMember, 0, len(memModels))

			for _, memModel := range memModels {
				members.MicrosoftEntraMembers = append(members.MicrosoftEntraMembers, fabcore.MicrosoftEntraMember{
					ObjectID:   memModel.ObjectID.ValueStringPointer(),
					TenantID:   memModel.TenantID.ValueStringPointer(),
					ObjectType: (*fabcore.ObjectType)(memModel.ObjectType.ValueStringPointer()),
				})
			}
		}

		to.Members = &members
	}

	return nil
}

/*
HELPER MODELS
*/

type decisionRuleModel struct {
	Effect      types.String                                             `tfsdk:"effect"`
	Permission  supertypes.ListNestedObjectValueOf[permissionScopeModel] `tfsdk:"permission"`
	Constraints supertypes.SingleNestedObjectValueOf[constraintsModel]   `tfsdk:"constraints"`
}

func (to *decisionRuleModel) set(ctx context.Context, from fabcore.DecisionRule) diag.Diagnostics {
	to.Effect = types.StringPointerValue((*string)(from.Effect))
	to.Permission.SetNull(ctx)
	to.Constraints.SetNull(ctx)

	// Permission
	if from.Permission != nil {
		permSlice := make([]*permissionScopeModel, 0, len(from.Permission))

		for _, p := range from.Permission {
			pModel := &permissionScopeModel{}

			if diags := pModel.set(ctx, p); diags.HasError() {
				return diags
			}

			permSlice = append(permSlice, pModel)
		}

		if diags := to.Permission.Set(ctx, permSlice); diags.HasError() {
			return diags
		}
	}

	// Constraints
	if from.Constraints != nil {
		cModel := &constraintsModel{}

		if diags := cModel.set(ctx, *from.Constraints); diags.HasError() {
			return diags
		}

		if diags := to.Constraints.Set(ctx, cModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

type permissionScopeModel struct {
	AttributeName            types.String                         `tfsdk:"attribute_name"`
	AttributeValueIncludedIn supertypes.ListValueOf[types.String] `tfsdk:"attribute_value_included_in"`
}

func (to *permissionScopeModel) set(ctx context.Context, from fabcore.PermissionScope) diag.Diagnostics {
	to.AttributeName = types.StringPointerValue((*string)(from.AttributeName))
	to.AttributeValueIncludedIn.SetNull(ctx)

	if from.AttributeValueIncludedIn != nil {
		values := make([]types.String, 0, len(from.AttributeValueIncludedIn))
		for _, v := range from.AttributeValueIncludedIn {
			values = append(values, types.StringValue(v))
		}

		if diags := to.AttributeValueIncludedIn.Set(ctx, values); diags.HasError() {
			return diags
		}
	}

	return nil
}

type constraintsModel struct {
	Columns supertypes.ListNestedObjectValueOf[columnConstraintModel] `tfsdk:"columns"`
	Rows    supertypes.ListNestedObjectValueOf[rowConstraintModel]    `tfsdk:"rows"`
}

func (to *constraintsModel) set(ctx context.Context, from fabcore.DecisionRuleConstraints) diag.Diagnostics {
	to.Columns.SetNull(ctx)
	to.Rows.SetNull(ctx)

	if from.Columns != nil {
		colSlice := make([]*columnConstraintModel, 0, len(from.Columns))

		for _, c := range from.Columns {
			cModel := &columnConstraintModel{}

			if diags := cModel.set(ctx, c); diags.HasError() {
				return diags
			}

			colSlice = append(colSlice, cModel)
		}

		if diags := to.Columns.Set(ctx, colSlice); diags.HasError() {
			return diags
		}
	}

	if from.Rows != nil {
		rowSlice := make([]*rowConstraintModel, 0, len(from.Rows))

		for _, r := range from.Rows {
			rModel := &rowConstraintModel{}
			rModel.set(r)
			rowSlice = append(rowSlice, rModel)
		}

		if diags := to.Rows.Set(ctx, rowSlice); diags.HasError() {
			return diags
		}
	}

	return nil
}

type columnConstraintModel struct {
	ColumnAction supertypes.ListValueOf[types.String] `tfsdk:"column_action"`
	ColumnEffect types.String                         `tfsdk:"column_effect"`
	ColumnNames  supertypes.ListValueOf[types.String] `tfsdk:"column_names"`
	TablePath    types.String                         `tfsdk:"table_path"`
}

func (to *columnConstraintModel) set(ctx context.Context, from fabcore.ColumnConstraint) diag.Diagnostics {
	to.ColumnAction.SetNull(ctx)
	to.ColumnNames.SetNull(ctx)
	to.ColumnEffect = types.StringPointerValue((*string)(from.ColumnEffect))
	to.TablePath = types.StringPointerValue(from.TablePath)

	if from.ColumnAction != nil {
		actions := make([]types.String, 0, len(from.ColumnAction))
		for _, a := range from.ColumnAction {
			actions = append(actions, types.StringValue(string(a)))
		}

		if diags := to.ColumnAction.Set(ctx, actions); diags.HasError() {
			return diags
		}
	}

	if from.ColumnNames != nil {
		names := make([]types.String, 0, len(from.ColumnNames))
		for _, n := range from.ColumnNames {
			names = append(names, types.StringValue(n))
		}

		if diags := to.ColumnNames.Set(ctx, names); diags.HasError() {
			return diags
		}
	}

	return nil
}

type rowConstraintModel struct {
	TablePath types.String `tfsdk:"table_path"`
	Value     types.String `tfsdk:"value"`
}

func (to *rowConstraintModel) set(from fabcore.RowConstraint) {
	to.TablePath = types.StringPointerValue(from.TablePath)
	to.Value = types.StringPointerValue(from.Value)
}

type membersModel struct {
	FabricItemMembers     supertypes.ListNestedObjectValueOf[fabricItemMemberModel]     `tfsdk:"fabric_item_members"`
	MicrosoftEntraMembers supertypes.ListNestedObjectValueOf[microsoftEntraMemberModel] `tfsdk:"microsoft_entra_members"`
}

func (to *membersModel) set(ctx context.Context, from fabcore.Members) diag.Diagnostics {
	to.FabricItemMembers.SetNull(ctx)
	to.MicrosoftEntraMembers.SetNull(ctx)

	if from.FabricItemMembers != nil {
		fimSlice := make([]*fabricItemMemberModel, 0, len(from.FabricItemMembers))

		for _, fim := range from.FabricItemMembers {
			fimModel := &fabricItemMemberModel{}

			if diags := fimModel.set(ctx, fim); diags.HasError() {
				return diags
			}

			fimSlice = append(fimSlice, fimModel)
		}

		if diags := to.FabricItemMembers.Set(ctx, fimSlice); diags.HasError() {
			return diags
		}
	}

	if from.MicrosoftEntraMembers != nil {
		memSlice := make([]*microsoftEntraMemberModel, 0, len(from.MicrosoftEntraMembers))

		for _, mem := range from.MicrosoftEntraMembers {
			memModel := &microsoftEntraMemberModel{}
			memModel.set(mem)
			memSlice = append(memSlice, memModel)
		}

		if diags := to.MicrosoftEntraMembers.Set(ctx, memSlice); diags.HasError() {
			return diags
		}
	}

	return nil
}

type fabricItemMemberModel struct {
	ItemAccess supertypes.ListValueOf[types.String] `tfsdk:"item_access"`
	SourcePath types.String                         `tfsdk:"source_path"`
}

func (to *fabricItemMemberModel) set(ctx context.Context, from fabcore.FabricItemMember) diag.Diagnostics {
	to.ItemAccess.SetNull(ctx)
	to.SourcePath = types.StringPointerValue(from.SourcePath)

	if from.ItemAccess != nil {
		accesses := make([]types.String, 0, len(from.ItemAccess))
		for _, a := range from.ItemAccess {
			accesses = append(accesses, types.StringValue(string(a)))
		}

		if diags := to.ItemAccess.Set(ctx, accesses); diags.HasError() {
			return diags
		}
	}

	return nil
}

type microsoftEntraMemberModel struct {
	ObjectID   customtypes.UUID `tfsdk:"object_id"`
	TenantID   customtypes.UUID `tfsdk:"tenant_id"`
	ObjectType types.String     `tfsdk:"object_type"`
}

func (to *microsoftEntraMemberModel) set(from fabcore.MicrosoftEntraMember) {
	to.ObjectID = customtypes.NewUUIDPointerValue(from.ObjectID)
	to.TenantID = customtypes.NewUUIDPointerValue(from.TenantID)
	to.ObjectType = types.StringPointerValue((*string)(from.ObjectType))
}
