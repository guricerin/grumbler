record Grumble {
  pk : String,
  content : String,
  userId : String,
  userName : String,
  createdAt : String
}

module Grumble {
  fun empty : Grumble {
    {
      pk = "",
      content = "",
      userId = "",
      userName = "",
      createdAt = ""
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
