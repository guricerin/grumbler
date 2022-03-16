record Regrumble {
  isRegrumble : Bool,
  byUserId : String,
  regrumbledCount : Number,
  isRegrumbledBySigninUser : Bool
}

module Regrumble {
  fun empty : Regrumble {
    {
      isRegrumble = false,
      byUserId = "",
      regrumbledCount = 0,
      isRegrumbledBySigninUser = false
    }
  }
}
