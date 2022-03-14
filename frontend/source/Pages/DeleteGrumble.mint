component Pages.DeleteGrumble {
  connect Application exposing { userStatus }
  connect Stores.GrumbleDetail exposing { apiStatus }
  state deleteStatus : Api.Status(DeleteGrumbleRes) = Api.Status::Initial

  fun render : Html {
    case (apiStatus) {
      Api.Status::Initial => Html.empty()
      Api.Status::Error(err) => <Errors errors={es}/>
      Api.Status::Ok(grumbleDetail) => core(grumbleDetail)
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  style button {
    margin: 5px;
  }

  fun core (grumbleDetail : GrumbleDetail) : Html {
    case (userStatus) {
      /* unreachable! */
      UserStatus::Guest => Html.empty()

      UserStatus::SignIn(user) =>
        <div>
          <Components.GrumbleBox
            grumble={grumbleDetail.target}
            signinUser={user}/>

          <hr/>

          <div class="box form-box">
            <p>"このぼやきを削除しますか？（この操作は取り消せません）"</p>
            <br/>

            <button::button
              class="button is-primary"
              type="submit"
              onClick={submit(grumbleDetail.target)}>

              <{ "削除" }>

            </button>

            <button::button
              class="button"
              type="submit"
              onClick={cancel()}>

              <{ "キャンセル" }>

            </button>
          </div>
        </div>
    }
  }

  fun cancel (event : Html.Event) : Promise(Never, Void) {
    `history.back()`
  }

  fun submit (target : Grumble, event : Html.Event) : Promise(Never, Void) {
    sequence {
      req =
        { grumblePk = target.pk }

      status =
        Http.post("#{@ENDPOINT}/auth/delete-grumble")
        |> Http.jsonBody(encode req)
        |> Api.send(DeleteGrumbleRes.decodes)

      case (status) {
        Api.Status::Ok(res) => `history.back()`
        => setApiStatus(status)
      }
    }
  }

  fun setApiStatus (v : Api.Status(DeleteGrumbleRes)) : Promise(Never, Void) {
    next { deleteStatus = v }
  }
}
