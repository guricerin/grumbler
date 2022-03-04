store Stores.User {
  state rsrcUser : User = User.empty()

  fun getRsrcUser (userId : String) : Promise(Never, Void) {
    sequence {
      res =
        Api.getUser(userId)

      case (res) {
        Api.Status::Ok(user) => next { rsrcUser = user }
        => next { rsrcUser = User.empty() }
      }
    }
  }
}
