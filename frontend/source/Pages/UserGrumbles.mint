component Pages.UserGrumbles {
  connect Stores.PageUser exposing { grumblesStatus }

  fun showStatus (status : Api.Status(Grumbles)) : Html {
    case (status) {
      Api.Status::Initial => Html.empty()
      Api.Status::Error(err) => <Errors errors={es}/>
      Api.Status::Ok(gs) => <GrumbleList grumbles={gs}/>
    }
  } where {
    es =
      Api.errorsOf("error", grumblesStatus)
  }

  fun render : Html {
    <div>
      <{ showStatus(grumblesStatus) }>
    </div>
  }
}
