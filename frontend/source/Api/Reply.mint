record ReplyReq {
  content : String,
  dstGrumblePk : String
}

record ReplyRes {
  ok : Bool
}

module ReplyRes {
  fun decodes (obj : Object) : Result(Object.Error, ReplyRes) {
    decode obj as ReplyRes
  }
}
