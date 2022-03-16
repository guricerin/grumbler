record RegrumbleReq {
  grumblePk : String
}

record RegrumbleRes {
  ok : Bool
}

module RegrumbleRes {
  fun decodes (obj : Object) : Result(Object.Error, RegrumbleRes) {
    decode obj as RegrumbleRes
  }
}
