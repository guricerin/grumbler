component Pages.UserDetail {
  connect Application exposing { userStatus }
  connect Stores.PageUser exposing { userDetail, showKind }
  state followApiStatus : Api.Status(FollowRes) = Api.Status::Initial

  fun doFollow (signinUser : User, apiUrl : String, event : Html.Event) : Promise(Never, Void) {
    sequence {
      follow =
        {
          srcUserId = signinUser.id,
          dstUserId = userDetail.user.id
        }

      status =
        Http.post("#{apiUrl}")
        |> Http.jsonBody(encode follow)
        |> Api.send(FollowRes.decodes)

      case (status) {
        Api.Status::Ok(res) =>
          sequence {
            next { followApiStatus = status }
            Stores.PageUser.getUserDetail(userDetail.user.id)
          }

        Api.Status::Initial => next { followApiStatus = status }
        Api.Status::Error(e) => next { followApiStatus = status }
      }
    }
  }

  fun followButton : Html {
    case (userStatus) {
      UserStatus::SignIn(signinUser) =>
        if (signinUser.id != userDetail.user.id) {
          if (userDetail.isFollow) {
            <a
              class="button is-outlined is-info"
              onClick={doFollow(signinUser, "#{@ENDPOINT}/auth/unfollow")}>

              "フォロー解除"

            </a>
          } else {
            <a
              class="button is-outlined is-info"
              onClick={doFollow(signinUser, "#{@ENDPOINT}/auth/follow")}>

              "フォロー"

            </a>
          }
        } else {
          Html.empty()
        }

      /* unreachable! */
      => Html.empty()
    }
  }

  fun moveToSettingsPage (event : Html.Event) : Promise(Never, Void) {
    Window.navigate("/settings")
  }

  fun settingsButton : Html {
    case (userStatus) {
      UserStatus::SignIn(signinUser) =>
        if (signinUser.id == userDetail.user.id) {
          <a
            class="button is-outlined is-info"
            onClick={moveToSettingsPage}>

            "設定"

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

  style text {
    white-space: pre-wrap;
  }

  fun showUserDetail : Html {
    <div>
      <strong>"#{userDetail.user.name}"</strong>
      <small>"@#{userDetail.user.id}"</small>
      <p::text>"#{userDetail.user.profile}"</p>
      <{ followButton() }>
      <{ settingsButton() }>
    </div>
  }

  fun showTabs : Html {
    <div class="tabs is-centered">
      <ul>
        <li>
          <a href="/user/#{userDetail.user.id}">
            <span>"ぼやき"</span>
          </a>
        </li>

        <li>
          <a href="/user/#{userDetail.user.id}/follows">
            <span>"#{arraySize(userDetail.follows)} フォロー"</span>
          </a>
        </li>

        <li>
          <a href="/user/#{userDetail.user.id}/followers">
            <span>"#{arraySize(userDetail.followers)} フォロワー"</span>
          </a>
        </li>

        <li>
          <a href="/user/#{userDetail.user.id}/bookmarks">
            <span>"#{arraySize(userDetail.bookmarks)} ブックマーク"</span>
          </a>
        </li>
      </ul>
    </div>
  }

  fun showSub : Html {
    case (showKind) {
      UserDetailShowKind::Grumbles => <Components.GrumbleList grumbles={Grumbles(userDetail.grumbles)}/>
      UserDetailShowKind::Follows => <Components.UserList users={Users(userDetail.follows)}/>
      UserDetailShowKind::Followers => <Components.UserList users={Users(userDetail.followers)}/>
      UserDetailShowKind::Bookmarks => <Components.GrumbleList grumbles={Grumbles(userDetail.bookmarks)}/>
    }
  }

  fun render : Html {
    <div>
      <{ showUserDetail() }>
      <br/>
      <{ showTabs() }>
      <{ showSub() }>
    </div>
  }
}
