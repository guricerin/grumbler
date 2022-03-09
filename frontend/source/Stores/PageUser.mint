enum UserDetailShowKind {
  Grumbles
  Follows
  Followers
  Bookmarks
}

store Stores.PageUser {
  state userDetail : UserDetail = UserDetail.empty()
  state showKind : UserDetailShowKind = UserDetailShowKind::Grumbles
  state grumblesStatus : Api.Status(Grumbles) = Api.Status::Initial

  fun setShowKind (sk : UserDetailShowKind) : Promise(Never, Void) {
    next { showKind = sk }
  }

  fun getUserDetail (userId : String) : Promise(Never, Void) {
    sequence {
      status =
        Http.get("#{@ENDPOINT}/auth/user/#{userId}/detail")
        |> Api.send(UserDetail.decodes)

      case (status) {
        Api.Status::Ok(ud) => next { userDetail = ud }
        => next { userDetail = UserDetail.empty() }
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
