record BookmarkReq {
  grumblePk : String,
  byUserId : String
}

record BookmarkRes {
  ok : Bool
}

module BookmarkRes {
  fun decodes (obj : Object) : Result(Object.Error, BookmarkRes) {
    decode obj as BookmarkRes
  }
}
