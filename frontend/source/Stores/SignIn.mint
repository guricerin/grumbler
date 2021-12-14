store Stores.SignIn {
  state userId : String = ""
  state password : String = ""
  state apiStatus : Api.Status(User) = Api.Status::Initial

  fun setUserId (v : String) : Promise(Never, Void) {
    next { userId = v }
  }

  fun setPassword (v : String) : Promise(Never, Void) {
    next { password = v }
  }

  fun setApiStatus (v : Api.Status(User)) : Promise(Never, Void) {
    next { apiStatus = v }
  }
}
