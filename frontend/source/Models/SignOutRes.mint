record SignOutRes {
  ok : Bool
}

module SignOutRes {
  fun decodes (obj : Object) : Result(Object.Error, SignOutRes) {
    decode obj as SignOutRes
  }
}
