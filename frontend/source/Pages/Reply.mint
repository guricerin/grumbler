component Pages.Reply {
  connect Application exposing { userStatus }
  connect Stores.GrumbleDetail exposing { apiStatus }
  state replyContent : String = ""

  fun setGrumbleContent (v : String) : Promise(Never, Void) {
    if (String.size(v) <= 300) {
      next { replyContent = v }
    } else {
      Promise.never()
    }
  }

  fun render : Html {
    case (apiStatus) {
      Api.Status::Initial => Html.empty()
      Api.Status::Error(err) => <Errors errors={es}/>

      Api.Status::Ok(grumbleDetail) => core(grumbleDetail.root)
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  fun core (grumble : Grumble) : Html {
    case (userStatus) {
      /* unreachable! */
      UserStatus::Guest => Html.empty()

      UserStatus::SignIn(user) =>
        <div>
          <Components.GrumbleBox
            signinUser={user}
            grumble={grumble}/>

          <hr/>
          <{ replyForm(grumble) }>
        </div>
    }
  }

  fun submit (replyDst : Grumble, event : Html.Event) : Promise(Never, Void) {
    sequence {
      replyReq =
        {
          content = replyContent,
          dstGrumblePk = replyDst.pk
        }

      status =
        Http.post("#{@ENDPOINT}/auth/reply")
        |> Http.jsonBody(encode replyReq)
        |> Api.send(ReplyRes.decodes)

      case (status) {
        Api.Status::Initial => next { }
        Api.Status::Ok(res) => `history.back()`
        Api.Status::Error(err) => Window.navigate("/")
      }
    }
  }

  get disabled : Bool {
    len <= 0 || 300 < len
  } where {
    len =
      String.size(replyContent)
  }

  style button {
    margin-top: 20px;
  }

  fun replyForm (replyDst : Grumble) : Html {
    <div>
      <div class="box form-box">
        <Components.GrumbleForm setGrumbleContent={setGrumbleContent}/>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={submit(replyDst)}
          disabled={disabled}>

          <{ "リプライ" }>

        </button>
      </div>
    </div>
  }
}
