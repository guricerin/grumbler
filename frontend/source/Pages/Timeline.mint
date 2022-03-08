component Pages.Timeline {
  connect Stores.Timeline exposing { apiStatus }

  fun showResult (res : TimelineResultKind) : Html {
    case (res) {
      TimelineResultKind::Initial => Html.empty()

      TimelineResultKind::Grumbles(grumbles) =>
        if (Array.size(grumbles.grumbles) < 1) {
          <div>
            <p>"最初のぼやきを投稿してみましょう。"</p>
          </div>
        } else {
          <GrumbleList grumbles={grumbles}/>
        }
    }
  }

  fun showStatus (status : Api.Status(TimelineResultKind)) : Html {
    case (status) {
      Api.Status::Initial => Html.empty()
      Api.Status::Error(err) => <Errors errors={es}/>
      Api.Status::Ok(res) => showResult(res)
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  fun render : Html {
    <div>
      <{ showStatus(apiStatus) }>
    </div>
  }
}
