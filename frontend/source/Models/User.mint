record SignUpUser {
  id : String,
  name : String,
  password : String
}

record SignInUser {
  id : String,
  password : String
}

record GetUserReq {
  id : String
}

record User {
  id : String,
  name : String,
  profile : String
}

module User {
  fun empty : User {
    {
      id = "",
      name = "",
      profile = ""
    }
  }

  fun decodes (obj : Object) : Result(Object.Error, User) {
    decode obj as User
  }
}

record Users {
  users : Array(User)
}

module Users {
  fun empty : Users {
    { users = [] }
  }

  fun decodes (obj : Object) : Result(Object.Error, Users) {
    decode obj as Users
  }
}
