component Pages.SignUp {
  state userId : String = ""
  state userName : String = ""
  state password : String = ""
  state apiStatus : Api.Status(User) = Api.Status::Initial

  fun setUserId (v : String) : Promise(Never, Void) {
    next { userId = v }
  }

  fun setUserName (v : String) : Promise(Never, Void) {
    next { userName = v }
  }

  fun setPassword (v : String) : Promise(Never, Void) {
    next { password = v }
  }

  fun setApiStatus (v : Api.Status(User)) : Promise(Never, Void) {
    next { apiStatus = v }
  }

  style content {
    flex-direction: column;
  }

  fun handleInput (
    onChange : Function(String, Promise(Never, Void)),
    event : Html.Event
  ) : a {
    onChange(Dom.getValue(event.target))
  }

  fun submit : Promise(Never, Void) {
    sequence {
      signupUser =
        {
          id = userId,
          name = userName,
          password = password
        }

      reqBody =
        encode signupUser

      status =
        Http.post("#{@ENDPOINT}/signup")
        |> Http.jsonBody(reqBody)
        |> Api.send(User.decodes)

      case (status) {
        Api.Status::Ok(user) => Application.signin(user)
        => setApiStatus(status)
      }
    }
  }

  get error : Html {
    case (apiStatus) {
      Api.Status::Error => <Errors errors={es}/>
      => Html.empty()
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  style button {
    margin-top: 20px;
  }

  fun render : Html {
    <div::content class="column">
      <div class="box form-box">
        <form>
          <{ error }>

          <div class="field">
            <label class="label">
              "ユーザ名"
            </label>

            <input
              class="input"
              type="text"
              placeholder="ユーザ名"
              value={userName}
              onChange={handleInput(setUserName)}/>

            <div>
              <small>"＊1文字以上127文字以下の範囲で設定してください。"</small>
              <br/>
              <small>"＊使用可能な文字は、UTF-8文字です（半角英数字だけでなく日本語も使用可能です）。"</small>
            </div>
          </div>

          <div class="field">
            <label class="label">
              "ユーザID"
            </label>

            <input
              class="input"
              type="text"
              placeholder="ユーザID"
              value={userId}
              onChange={handleInput(setUserId)}/>

            <div>
              <small>"＊他ユーザアカウントと重複するIDは設定できません。"</small>
              <br/>
              <small>"＊1文字以上127文字以下の範囲で設定してください。"</small>
              <br/>
              <small>"＊使用可能な文字は、半角英数字とアンダーバー（_）です。"</small>
            </div>
          </div>

          <div class="field">
            <label class="label">
              "パスワード"
            </label>

            <input
              class="input"
              type="password"
              placeholder="パスワード"
              value={password}
              onChange={handleInput(setPassword)}/>

            <div>
              <small>"＊8文字以上127文字以下の半角英数字で設定してください。"</small>
            </div>
          </div>
        </form>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={submit}>

          <{ "新規登録" }>

        </button>
      </div>
    </div>
  }
}
