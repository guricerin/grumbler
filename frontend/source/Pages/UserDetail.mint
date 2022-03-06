component Pages.UserDetail {
  connect Application exposing { userStatus }
  connect Stores.PageUser exposing { userDetail, showKind }
  state isFollow : Bool = false
  state followApiStatus : Api.Status(FollowRes) = Api.Status::Initial

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

      /* unreachable! */
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
    </div>
  }

  fun showTabs (ud : UserDetail) : Html {
    <div class="tabs is-centered">
      <ul>
        <li>
          <a href="/user/#{ud.user.id}">
            <span>"ぼやき"</span>
          </a>
        </li>

        <li>
          <a href="/user/#{ud.user.id}/follows">
            <span>"#{arraySize(userDetail.follows)} フォロー"</span>
          </a>
        </li>

        <li>
          <a href="/user/#{ud.user.id}/followers">
            <span>"#{arraySize(userDetail.followers)} フォロワー"</span>
          </a>
        </li>
      </ul>
    </div>
  }

  fun showSub : Html {
    case (showKind) {
      UserDetailShowKind::Grumbles => <GrumbleList grumbles={Grumbles(userDetail.grumbles)}/>
      UserDetailShowKind::Follows => <UserList users={Users(userDetail.follows)}/>
      UserDetailShowKind::Followers => <UserList users={Users(userDetail.followers)}/>
    }
  }

  fun render : Html {
    <div>
      <{ showUserDetail(userDetail) }>
      <br/>
      <{ showTabs(userDetail) }>
      <{ showSub() }>
    </div>
  }
}
