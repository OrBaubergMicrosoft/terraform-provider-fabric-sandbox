resource "fabric_onelake_data_access_security" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyRole"

  decision_rules {
    effect = "Permit"

    permission {
      attribute_name              = "Path"
      attribute_value_included_in = ["Tables/MyTable"]
    }
  }

  members {
    microsoft_entra_members {
      object_id   = "22222222-2222-2222-2222-222222222222"
      tenant_id   = "33333333-3333-3333-3333-333333333333"
      object_type = "User"
    }
  }
}
