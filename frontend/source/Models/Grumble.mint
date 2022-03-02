record Grumble {
  pk : String,
  content : String,
  userId : String,
  createdAt : Time
}

module Grumble {
  fun empty : Grumble {
    {
      pk = "",
      content = "",
      userId = "",
      createdAt = Time.now()
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
