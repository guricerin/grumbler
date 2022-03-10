record GrumbleDetail {
  root : Grumble,
  replies : Array(Grumble)
}

module GrumbleDetail {
  fun empty : GrumbleDetail {
    {
      root = Grumble.empty(),
      replies = []
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, GrumbleDetail) {
    decode obj as GrumbleDetail
  }
}
