record Grumble {
  pk : String,
  content : String,
  userId : String,
  userName : String,
  createdAt : String,
  bookmarkedCount : Number,
  isBookmarkedBySigninUser : Bool
}

module Grumble {
  fun empty : Grumble {
    {
      pk = "",
      content = "",
      userId = "",
      userName = "",
      createdAt = "",
      bookmarkedCount = 0,
      isBookmarkedBySigninUser = false
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, Grumble) {
    decode obj as Grumble
  }
}

record Grumbles {
  grumbles : Array(Grumble)
}

module Grumbles {
  fun empty : Grumbles {
    { grumbles = [] }
  }

  fun decodes (obj : Object) : Result(Object.Error, Grumbles) {
    decode obj as Grumbles
  }
}
