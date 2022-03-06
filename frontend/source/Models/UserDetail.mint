record UserDetail {
  user : User,
  grumbles : Array(Grumble),
  follows : Array(Follow),
  followers : Array(Follow),
  isFollow : Bool,
  isFollower : Bool
}

module UserDetail {
  fun empty : UserDetail {
    {
      user = User.empty(),
      grumbles = [],
      follows = [],
      followers = [],
      isFollow = false,
      isFollower = false
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, UserDetail) {
    decode obj as UserDetail
  }
}
