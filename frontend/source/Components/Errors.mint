component Errors {
  property errors : Array(String) = []

  style base {
    margin-bottom: 20px;
    background: #f7b6b6;
    text-align: center;
    border-radius: 2px;
    font-weight: bold;
    font-size: 14px;
    color: #902e2e;
    padding: 20px;
  }

  fun renderError (error : String) : Html {
    <li>
      <{ error }>
    </li>
  }

  fun render : Html {
    if (Array.isEmpty(errors)) {
      Html.empty()
    } else {
      <article class="message is-danger">
        <div class="message-header">
          <p>"Error"</p>
        </div>
        <div class="message-body">
          <ul>
        <{ Array.map(renderError, errors) }>
          </ul>
        </div>
      </article>
    }
  }
}
