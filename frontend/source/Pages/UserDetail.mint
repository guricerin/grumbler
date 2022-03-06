component Pages.UserDetail {
  connect Application exposing { userStatus }
  connect Stores.PageUser exposing { userDetail, showKind }
  state isFollow : Bool = false
  state followApiStatus : Api.Status(FollowRes) = Api.Status::Initial

  style profileItem {
    margin-left: 5px;
  }

  fun doFollow (signinuser : User, event : Html.Event) : Promise(Never, Void) {
    sequence {
      follow =
        {
          srcUserId = signinuser.id,
          dstUserId = userDetail.user.id
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
        if (signinUser.id != userDetail.user.id) {
          <a
            class="button is-outlined is-info"
            onClick={doFollow(signinUser)}>

            "フォロー"

          </a>
        } else {
          Html.empty()
        }

      => Html.empty()
    }
  }

  fun arraySize (ls : Array(a)) : String {
    ls
    |> Array.size
    |> Number.toString
  }

  fun showUserDetail (ud : UserDetail) : Html {
    <div>
      <strong>"#{ud.user.name}"</strong>
      <small>"@#{ud.user.id}"</small>
      <p>"#{ud.user.profile}"</p>
      <{ followButton() }>
      <hr/>

      <nav class="level is-mobile">
        <div class="level-item ">
          <a>"ぼやき"</a>
        </div>

        <div class="level-item">
          <a>"#{arraySize(userDetail.follows)} フォロー"</a>
        </div>

        <div class="level-item">
          <a>"#{arraySize(userDetail.followers)} フォロワー"</a>
        </div>
      </nav>

      <a::profileItem href="/user/#{ud.user.id}/grumbles">
        "ぼやき"
      </a>

      <a::profileItem>"フォロー"</a>
      <a::profileItem>"フォロワー"</a>
    </div>
  }

  fun showSub : Html {
    case (showKind) {
      UserDetailShowKind::Grumbles => <GrumbleList grumbles={gs}/>
      => Html.empty()
    }
  } where {
    gs =
      Grumbles(userDetail.grumbles)
  }

  fun render : Html {
    <div>
      <{ showUserDetail(userDetail) }>
      <hr/>
      <{ showSub() }>
    </div>
  }
}
