record Reply {
  dstGrumblePk : String,
  dstUserId : String,
  repliedCount : Number
}

module Reply {
  fun empty : Reply {
    {
      dstGrumblePk = "",
      dstUserId = "",
      repliedCount = 0
    }
  }
}
