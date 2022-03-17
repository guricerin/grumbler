record Regrumble {
  createdAt : String,
  isRegrumble : Bool,
  byUserId : String,
  regrumbledCount : Number,
  isRegrumbledBySigninUser : Bool
}

module Regrumble {
  fun empty : Regrumble {
    {
      createdAt = "",
      isRegrumble = false,
      byUserId = "",
      regrumbledCount = 0,
      isRegrumbledBySigninUser = false
    }
  }
}
