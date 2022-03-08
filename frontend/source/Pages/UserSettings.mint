component Pages.UserSettings {
  connect Application exposing { userStatus }
  state apiStatus : Api.Status(UserSettingsRes) = Api.Status::Initial
  state name : String = ""
  state profile : String = ""

  fun setName (v : String) : Promise(Never, Void) {
    next { name = v }
  }

  fun setProfile (v : String) : Promise(Never, Void) {
    next { profile = v }
  }

  fun setApiStatus (v : Api.Status(UserSettingsRes)) : Promise(Never, Void) {
    next { apiStatus = v }
  }

  get error : Html {
    case (apiStatus) {
      Api.Status::Error => <Errors errors={es}/>
      Api.Status::Initial => Html.empty()
      Api.Status::Ok(res) => <Success message="変更を保存しました。"/>
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  fun submit : Promise(Never, Void) {
    sequence {
      userSettings =
        {
          name = name,
          profile = profile
        }

      status =
        Http.post("#{@ENDPOINT}/auth/settings")
        |> Http.jsonBody(encode userSettings)
        |> Api.send(UserSettingsRes.decodes)

      case (status) {
        Api.Status::Ok(res) => Application.updateUser(name, profile)
        => next { }
      }

      setApiStatus(status)
    }
  }

  fun handleInput (
    onChange : Function(String, Promise(Never, Void)),
    event : Html.Event
  ) : a {
    onChange(Dom.getValue(event.target))
  }

  /*
  this is called when the component is mounted
  https://www.mint-lang.com/guide/reference/components/lifecycle-functions
  */
  fun componentDidMount : Promise(Never, Void) {
    sequence {
      case (userStatus) {
        /* unreachable! */
        UserStatus::Guest => next { }

        UserStatus::SignIn(user) =>
          sequence {
            setName(user.name)
            setProfile(user.profile)
          }
      }
    }
  }

  style content {
    flex-direction: column;
  }

  style button {
    margin-top: 20px;
  }

  fun render : Html {
    <div::content class="column">
      <div class="box form-box">
        <{ error }>

        <form>
          <div class="field">
            <label class="label">
              "ユーザ名"
            </label>

            <input
              class="input"
              type="text"
              maxlength="32"
              placeholder="ユーザ名"
              value={name}
              onChange={handleInput(setName)}
              required="required"/>

            <div>
              <small>"＊1文字以上32文字以下の範囲で設定してください。"</small>
              <br/>
              <small>"＊使用可能な文字は、UTF-8文字です（半角英数字だけでなく日本語も使用可能です）。"</small>
            </div>
          </div>

          <div class="field">
            <label class="label">
              "プロフィール"
            </label>

            <textarea
              class="textarea"
              placeholder="プロフィール"
              maxlength="200"
              value={profile}
              onChange={handleInput(setProfile)}/>

            <div>
              <small>"＊1～200文字の範囲で入力可能です。"</small>
            </div>
          </div>
        </form>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={submit}>

          <{ "保存" }>

        </button>
      </div>
    </div>
  }
}
