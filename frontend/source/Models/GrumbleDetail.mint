record GrumbleDetail {
  target : Grumble,
  ancestors : Array(Grumble),
  replies : Array(Grumble)
}

module GrumbleDetail {
  fun empty : GrumbleDetail {
    {
      target = Grumble.empty(),
      ancestors = [],
      replies = []
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, GrumbleDetail) {
    decode obj as GrumbleDetail
  }
}
