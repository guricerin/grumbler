record UserDetail {
  user : User,
  grumbles : Array(Grumble),
  follows : Array(Follow),
  followers : Array(Follow)
}

module UserDetail {
  fun empty : UserDetail {
    {
      user = User.empty(),
      grumbles = [],
      follows = [],
      followers = []
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, UserDetail) {
    decode obj as UserDetail
  }
}
