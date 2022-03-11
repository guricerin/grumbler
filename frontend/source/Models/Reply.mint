record Reply {
  toGrumblePk : String,
  toUserId : String
}

module Reply {
  fun empty : Reply {
    {
      toGrumblePk = "",
      toUserId = ""
    }
  }
}
