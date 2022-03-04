store Stores.PageUser {
  state rsrcUser : User = User.empty()
  state grumblesStatus : Api.Status(Grumbles) = Api.Status::Initial

  fun getUser (userId : String) : Promise(Never, Void) {
    sequence {
      res =
        Api.getUser(userId)

      case (res) {
        Api.Status::Ok(user) => next { rsrcUser = user }
        => next { rsrcUser = User.empty() }
      }
    }
  }

  fun getGrumbles (userId : String) : Promise(Never, Void) {
    sequence {
      status =
        Http.get("#{@ENDPOINT}/auth/user/#{userId}/grumbles")
        |> Api.send(Grumbles.decodes)

      case (status) {
        Api.Status::Ok(grumbles) =>
          sequence {
            next { grumblesStatus = Api.Status::Ok(grumbles) }
          }

        Api.Status::Initial => next { grumblesStatus = Api.Status::Initial }
        Api.Status::Error(err) => next { grumblesStatus = Api.Status::Error(err) }
      }
    }
  }
}
