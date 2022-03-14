record GrumbleReq {
  content : String
}

record GrumbleRes {
  ok : Bool
}

module GrumbleRes {
  fun decodes (obj : Object) : Result(Object.Error, GrumbleRes) {
    decode obj as GrumbleRes
  }
}

record DeleteGrumbleReq {
  grumblePk : String
}

record DeleteGrumbleRes {
  ok : Bool
}

module DeleteGrumbleRes {
  fun decodes (obj : Object) : Result(Object.Error, DeleteGrumbleRes) {
    decode obj as DeleteGrumbleRes
  }
}
