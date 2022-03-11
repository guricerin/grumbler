component Components.Success {
  property message : String = ""

  fun render : Html {
    if (message == "") {
      Html.empty()
    } else {
      <article class="message is-info">
        <div class="message-header">
          <p>"Success"</p>
        </div>

        <div class="message-body">
          <p>
            <{ message }>
          </p>
        </div>
      </article>
    }
  }
}
