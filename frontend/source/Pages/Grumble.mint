component Pages.Grumble {
  state grumbleContent : String = ""
  state apiStatus : Api.Status(GrumbleRes) = Api.Status::Initial

  fun setGrumbleContent (v : String) : Promise(Never, Void) {
    next { grumbleContent = v }
  }

  fun setApiStatus (v : Api.Status(GrumbleRes)) : Promise(Never, Void) {
    next { apiStatus = v }
  }

  fun render : Html {
    <div>
      <div class="box form-box">
        <form>
          <p>"grumble"</p>
        </form>
      </div>
    </div>
  }
}
