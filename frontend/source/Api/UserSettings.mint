record UserSettingsReq {
  name : String,
  profile : String
}

record UserSettingsRes {
  ok : Bool
}

module UserSettingsRes {
  fun decodes (obj : Object) : Result(Object.Error, UserSettingsRes) {
    decode obj as UserSettingsRes
  }
}
