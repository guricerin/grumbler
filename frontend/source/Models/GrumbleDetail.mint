record GrumbleDetail {
  root : Grumble,
  ancestors : Array(Grumble),
  replies : Array(Grumble)
}

module GrumbleDetail {
  fun empty : GrumbleDetail {
    {
      root = Grumble.empty(),
      ancestors = [],
      replies = []
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, GrumbleDetail) {
    decode obj as GrumbleDetail
  }
}
