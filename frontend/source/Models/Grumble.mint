record Grumble {
  pk : String,
  content : String,
  userId : String,
  userName : String,
  createdAt : String,
  reply : Reply,
  regrumble : Regrumble,
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
      reply = Reply.empty(),
      regrumble = Regrumble.empty(),
      bookmarkedCount = 0,
      isBookmarkedBySigninUser = false
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, Grumble) {
    decode obj as Grumble
  }

  fun isReply (g : Grumble) : Bool {
    if (g.reply.dstUserId == "") {
      false
    } else {
      true
    }
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
