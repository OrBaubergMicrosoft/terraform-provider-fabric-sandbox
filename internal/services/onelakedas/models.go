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
	WorkspaceID   customtypes.UUID                                               `tfsdk:"workspace_id"`
	ItemID        customtypes.UUID                                               `tfsdk:"item_id"`
	Name          types.String                                                   `tfsdk:"name"`
	Kind          types.String                                                   `tfsdk:"kind"`
	DecisionRules supertypes.ListNestedObjectValueOf[decisionRuleModel]          `tfsdk:"decision_rules"`
	Members       supertypes.SingleNestedObjectValueOf[membersModel]             `tfsdk:"members"`
}

func (to *baseOneLakeDataAccessSecurityModel) set(ctx context.Context, workspaceID, itemID string, from fabcore.DataAccessRoleBase) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	to.Name = types.StringPointerValue(from.Name)
	to.Kind = types.StringPointerValue((*string)(from.Kind))

	// DecisionRules
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
	WorkspaceID customtypes.UUID                                                        `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                                        `tfsdk:"item_id"`
	Values      supertypes.SetNestedObjectValueOf[baseOneLakeDataAccessSecurityModel]    `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                                         `tfsdk:"timeouts"`
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
		dr, diags := drModel.toSDK(ctx)
		if diags.HasError() {
			return diags
		}

		to.DecisionRules = append(to.DecisionRules, dr)
	}

	// Members
	membersPtr, diags := from.Members.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if membersPtr != nil {
		members, diags := membersPtr.toSDK(ctx)
		if diags.HasError() {
			return diags
		}

		to.Members = &members
	}

	return nil
}

/*
HELPER MODELS
*/

type decisionRuleModel struct {
	Effect      types.String                                                `tfsdk:"effect"`
	Permission  supertypes.ListNestedObjectValueOf[permissionScopeModel]    `tfsdk:"permission"`
	Constraints supertypes.SingleNestedObjectValueOf[constraintsModel]      `tfsdk:"constraints"`
}

func (to *decisionRuleModel) set(ctx context.Context, from fabcore.DecisionRule) diag.Diagnostics {
	to.Effect = types.StringPointerValue((*string)(from.Effect))

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

func (from *decisionRuleModel) toSDK(ctx context.Context) (fabcore.DecisionRule, diag.Diagnostics) {
	result := fabcore.DecisionRule{
		Effect: (*fabcore.Effect)(from.Effect.ValueStringPointer()),
	}

	// Permission
	permModels, diags := from.Permission.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	result.Permission = make([]fabcore.PermissionScope, 0, len(permModels))

	for _, pm := range permModels {
		ps, d := pm.toSDK(ctx)
		if d.HasError() {
			return result, d
		}

		result.Permission = append(result.Permission, ps)
	}

	// Constraints
	cPtr, diags := from.Constraints.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	if cPtr != nil {
		c, diags := cPtr.toSDK(ctx)
		if diags.HasError() {
			return result, diags
		}

		result.Constraints = &c
	}

	return result, nil
}

type permissionScopeModel struct {
	AttributeName            types.String                       `tfsdk:"attribute_name"`
	AttributeValueIncludedIn supertypes.ListValueOf[types.String] `tfsdk:"attribute_value_included_in"`
}

func (to *permissionScopeModel) set(ctx context.Context, from fabcore.PermissionScope) diag.Diagnostics {
	to.AttributeName = types.StringPointerValue((*string)(from.AttributeName))

	values := make([]types.String, 0, len(from.AttributeValueIncludedIn))
	for _, v := range from.AttributeValueIncludedIn {
		values = append(values, types.StringValue(v))
	}

	return to.AttributeValueIncludedIn.Set(ctx, values)
}

func (from *permissionScopeModel) toSDK(ctx context.Context) (fabcore.PermissionScope, diag.Diagnostics) {
	result := fabcore.PermissionScope{
		AttributeName: (*fabcore.AttributeName)(from.AttributeName.ValueStringPointer()),
	}

	elements, diags := from.AttributeValueIncludedIn.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	for _, e := range elements {
		result.AttributeValueIncludedIn = append(result.AttributeValueIncludedIn, e.ValueString())
	}

	return result, nil
}

type constraintsModel struct {
	Columns supertypes.ListNestedObjectValueOf[columnConstraintModel] `tfsdk:"columns"`
	Rows    supertypes.ListNestedObjectValueOf[rowConstraintModel]    `tfsdk:"rows"`
}

func (to *constraintsModel) set(ctx context.Context, from fabcore.DecisionRuleConstraints) diag.Diagnostics {
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

func (from *constraintsModel) toSDK(ctx context.Context) (fabcore.DecisionRuleConstraints, diag.Diagnostics) {
	result := fabcore.DecisionRuleConstraints{}

	colModels, diags := from.Columns.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	if colModels != nil {
		result.Columns = make([]fabcore.ColumnConstraint, 0, len(colModels))

		for _, cm := range colModels {
			cc, d := cm.toSDK(ctx)
			if d.HasError() {
				return result, d
			}

			result.Columns = append(result.Columns, cc)
		}
	}

	rowModels, diags := from.Rows.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	if rowModels != nil {
		result.Rows = make([]fabcore.RowConstraint, 0, len(rowModels))

		for _, rm := range rowModels {
			result.Rows = append(result.Rows, rm.toSDK())
		}
	}

	return result, nil
}

type columnConstraintModel struct {
	ColumnAction supertypes.ListValueOf[types.String] `tfsdk:"column_action"`
	ColumnEffect types.String                         `tfsdk:"column_effect"`
	ColumnNames  supertypes.ListValueOf[types.String] `tfsdk:"column_names"`
	TablePath    types.String                         `tfsdk:"table_path"`
}

func (to *columnConstraintModel) set(ctx context.Context, from fabcore.ColumnConstraint) diag.Diagnostics {
	actions := make([]types.String, 0, len(from.ColumnAction))
	for _, a := range from.ColumnAction {
		actions = append(actions, types.StringValue(string(a)))
	}

	if diags := to.ColumnAction.Set(ctx, actions); diags.HasError() {
		return diags
	}

	to.ColumnEffect = types.StringPointerValue((*string)(from.ColumnEffect))

	names := make([]types.String, 0, len(from.ColumnNames))
	for _, n := range from.ColumnNames {
		names = append(names, types.StringValue(n))
	}

	if diags := to.ColumnNames.Set(ctx, names); diags.HasError() {
		return diags
	}

	to.TablePath = types.StringPointerValue(from.TablePath)

	return nil
}

func (from *columnConstraintModel) toSDK(ctx context.Context) (fabcore.ColumnConstraint, diag.Diagnostics) {
	result := fabcore.ColumnConstraint{
		ColumnEffect: (*fabcore.ColumnEffect)(from.ColumnEffect.ValueStringPointer()),
		TablePath:    from.TablePath.ValueStringPointer(),
	}

	actions, diags := from.ColumnAction.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	for _, a := range actions {
		result.ColumnAction = append(result.ColumnAction, fabcore.ColumnAction(a.ValueString()))
	}

	names, diags := from.ColumnNames.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	for _, n := range names {
		result.ColumnNames = append(result.ColumnNames, n.ValueString())
	}

	return result, nil
}

type rowConstraintModel struct {
	TablePath types.String `tfsdk:"table_path"`
	Value     types.String `tfsdk:"value"`
}

func (to *rowConstraintModel) set(from fabcore.RowConstraint) {
	to.TablePath = types.StringPointerValue(from.TablePath)
	to.Value = types.StringPointerValue(from.Value)
}

func (from *rowConstraintModel) toSDK() fabcore.RowConstraint {
	return fabcore.RowConstraint{
		TablePath: from.TablePath.ValueStringPointer(),
		Value:     from.Value.ValueStringPointer(),
	}
}

type membersModel struct {
	FabricItemMembers     supertypes.ListNestedObjectValueOf[fabricItemMemberModel]     `tfsdk:"fabric_item_members"`
	MicrosoftEntraMembers supertypes.ListNestedObjectValueOf[microsoftEntraMemberModel] `tfsdk:"microsoft_entra_members"`
}

func (to *membersModel) set(ctx context.Context, from fabcore.Members) diag.Diagnostics {
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

func (from *membersModel) toSDK(ctx context.Context) (fabcore.Members, diag.Diagnostics) {
	result := fabcore.Members{}

	fimModels, diags := from.FabricItemMembers.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	if fimModels != nil {
		result.FabricItemMembers = make([]fabcore.FabricItemMember, 0, len(fimModels))

		for _, fimModel := range fimModels {
			fim, d := fimModel.toSDK(ctx)
			if d.HasError() {
				return result, d
			}

			result.FabricItemMembers = append(result.FabricItemMembers, fim)
		}
	}

	memModels, diags := from.MicrosoftEntraMembers.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	if memModels != nil {
		result.MicrosoftEntraMembers = make([]fabcore.MicrosoftEntraMember, 0, len(memModels))

		for _, memModel := range memModels {
			result.MicrosoftEntraMembers = append(result.MicrosoftEntraMembers, memModel.toSDK())
		}
	}

	return result, nil
}

type fabricItemMemberModel struct {
	ItemAccess supertypes.ListValueOf[types.String] `tfsdk:"item_access"`
	SourcePath types.String                         `tfsdk:"source_path"`
}

func (to *fabricItemMemberModel) set(ctx context.Context, from fabcore.FabricItemMember) diag.Diagnostics {
	accesses := make([]types.String, 0, len(from.ItemAccess))
	for _, a := range from.ItemAccess {
		accesses = append(accesses, types.StringValue(string(a)))
	}

	if diags := to.ItemAccess.Set(ctx, accesses); diags.HasError() {
		return diags
	}

	to.SourcePath = types.StringPointerValue(from.SourcePath)

	return nil
}

func (from *fabricItemMemberModel) toSDK(ctx context.Context) (fabcore.FabricItemMember, diag.Diagnostics) {
	result := fabcore.FabricItemMember{
		SourcePath: from.SourcePath.ValueStringPointer(),
	}

	accesses, diags := from.ItemAccess.Get(ctx)
	if diags.HasError() {
		return result, diags
	}

	for _, a := range accesses {
		result.ItemAccess = append(result.ItemAccess, fabcore.ItemAccess(a.ValueString()))
	}

	return result, nil
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

func (from *microsoftEntraMemberModel) toSDK() fabcore.MicrosoftEntraMember {
	return fabcore.MicrosoftEntraMember{
		ObjectID:   from.ObjectID.ValueStringPointer(),
		TenantID:   from.TenantID.ValueStringPointer(),
		ObjectType: (*fabcore.ObjectType)(from.ObjectType.ValueStringPointer()),
	}
}
