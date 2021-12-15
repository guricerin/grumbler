component Pages.SignOut {
  connect Application exposing { userStatus }

  fun cancel (user : User, event : Html.Event) : Promise(Never, Void) {
    sequence {
      Window.navigate("/user/#{user.id}/timeline")
    }
  }

  fun doSignOut (user : User, event : Html.Event) : Promise(Never, Void) {
    Window.navigate("/user/#{user.id}/timeline")
  }

  style content {
    flex-direction: column;
  }

  style button {
    margin: 5px;
  }

  fun core (user : User) : Html {
    <div::content class="column">
      <div class="box form-box">
        <p>"サインアウトしますか？"</p>
        <br/>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={doSignOut(user)}>

          <{ "はい" }>

        </button>

        <button::button
          class="button"
          type="submit"
          onClick={cancel(user)}>

          <{ "いいえ" }>

        </button>
      </div>
    </div>
  }

  fun render : Html {
    case (userStatus) {
      /* unreachable */
      UserStatus::Guest => Html.empty()
      UserStatus::SignIn(user) => core(user)
    }
  }
}
