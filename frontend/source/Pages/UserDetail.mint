component Pages.UserDetail {
  connect Application exposing { userStatus }
  connect Stores.PageUser exposing { rsrcUser }
  state isFollow : Bool = false
  state followApiStatus : Api.Status(FollowRes) = Api.Status::Initial

  style profileItem {
    margin-left: 5px;
  }

  fun doFollow (ud : UserDetail, event : Html.Event) : Promise(Never, Void) {
    sequence {
      follow =
        {
          srcUserId = ud.user.id,
          dstUserId = rsrcUser.user.id
        }

      status =
        Http.post("#{@ENDPOINT}/auth/follow")
        |> Http.jsonBody(encode follow)
        |> Api.send(FollowRes.decodes)

      case (status) {
        Api.Status::Ok(res) => next { isFollow = true }
        => next { followApiStatus = status }
      }
    }
  }

  fun followButton : Html {
    case (userStatus) {
      UserStatus::SignIn(signinUser) =>
        if (signinUser.id != rsrcUser.user.id) {
          <a
            class="button is-outlined is-info"
            onClick={doFollow(rsrcUser)}>

            "フォロー"

          </a>
        } else {
          Html.empty()
        }

      => Html.empty()
    }
  }

  fun rsrcUserProfile (ud : UserDetail) : Html {
    <div>
      <strong>"#{ud.user.name}"</strong>
      <small>"@#{ud.user.id}"</small>
      <p>"#{ud.user.profile}"</p>
      <{ followButton() }>
      <hr/>

      <a::profileItem href="/user/#{ud.user.id}/grumbles">
        "ぼやき"
      </a>

      <a::profileItem>"フォロー"</a>
      <a::profileItem>"フォロワー"</a>
    </div>
  }

  fun render : Html {
    <div>
      <{ rsrcUserProfile(rsrcUser) }>
    </div>
  }
}
