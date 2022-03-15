record FollowReq {
  dstUserId : String
}

record FollowRes {
  ok : Bool
}

module FollowRes {
  fun decodes (obj : Object) : Result(Object.Error, FollowRes) {
    decode obj as FollowRes
  }
}
