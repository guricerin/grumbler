record GrumbleRes {
  ok : Bool
}

module GrumbleRes {
  fun decodes (obj : Object) : Result(Object.Error, GrumbleRes) {
    decode obj as GrumbleRes
  }
}
