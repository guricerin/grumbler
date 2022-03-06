record Follow {
  srcUserId : String,
  dstUserId : String
}

module Follow {
  fun empty : Follow {
    {
      srcUserId = "",
      dstUserId = ""
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, Follow) {
    decode obj as Follow
  }
}
