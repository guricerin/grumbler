store Stores.PageUser {
  state rsrcUser : UserDetail = UserDetail.empty()
  state grumblesStatus : Api.Status(Grumbles) = Api.Status::Initial

  fun getUserDetail (userId : String) : Promise(Never, Void) {
    sequence {
      status =
        Http.get("#{@ENDPOINT}/auth/user/#{userId}/detail")
        |> Api.send(UserDetail.decodes)

      case (status) {
        Api.Status::Ok(user) => next { rsrcUser = user }
        => next { rsrcUser = UserDetail.empty() }
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
