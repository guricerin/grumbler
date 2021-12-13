record SignUpUser {
  id : String,
  name : String,
  password : String
}

record User {
  id : String,
  name : String
}

module User {
  fun decode (obj : Object) : Result(Object.Error, User) {
    decode obj as User
  }
}
