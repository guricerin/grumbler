component Pages.PostGrumble {
  state grumbleContent : String = ""
  state apiStatus : Api.Status(GrumbleRes) = Api.Status::Initial

  fun setGrumbleContent (v : String) : Promise(Never, Void) {
    if (String.size(v) <= 300) {
      next { grumbleContent = v }
    } else {
      Promise.never()
    }
  }

  fun setApiStatus (v : Api.Status(GrumbleRes)) : Promise(Never, Void) {
    next { apiStatus = v }
  }

  fun submit : Promise(Never, Void) {
    sequence {
      req =
        { content = grumbleContent }

      reqBody =
        encode req

      status =
        Http.post("#{@ENDPOINT}/auth/grumble")
        |> Http.jsonBody(reqBody)
        |> Api.send(GrumbleRes.decodes)

      case (status) {
        Api.Status::Ok(res) =>
          case (Application.userStatus) {
            UserStatus::SignIn(u) => Window.navigate("/user/#{u.id}/timeline")

            // unreachable!
            => Window.navigate("")
          }

        => setApiStatus(status)
      }
    }
  }

  fun handleInput (
    onChange : Function(String, Promise(Never, Void)),
    event : Html.Event
  ) : a {
    onChange(Dom.getValue(event.target))
  }

  get disabled : Bool {
    len <= 0 || 300 < len
  } where {
    len =
      String.size(grumbleContent)
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

  style warning {
    color: red;
  }

  style button {
    margin-top: 20px;
  }

  fun render : Html {
    <div>
      <div class="box form-box">
        <{ error }>

        <form>
          <div class="field">
            <textarea
              class="textarea"
              maxlength="300"
              onChange={handleInput(setGrumbleContent)}/>

            <small>"＊1～300文字の範囲で入力可能です。"</small>
            <br/>
            <small::warning>"＊なにをぼやいても自由ですが、誹謗中傷や犯罪予告などを書き込まないようにしましょう。大いなる自由には大いなる責任が伴います。"</small>
          </div>
        </form>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={submit}
          disabled={disabled}>

          <{ "ぼやく" }>

        </button>
      </div>
    </div>
  }
}
