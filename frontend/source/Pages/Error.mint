component Pages.Error {
  property statusCode : Number = 200

  style content {
    align-items: center;
    justify-content: center;
    display: flex;
    height: 50vh;
  }

  fun render : Html {
    <div::content>
      <p>
        <{ message }>
      </p>
    </div>
  }

  get message : String {
    case (statusCode) {
      200 => "#{code} unreachable!"
      403 => "#{code} Formidden."
      404 => "#{code} Page not found."
      500 => "#{code} Internal server error."
      => "#{code} unexausted!"
    }
  } where {
    code =
      Number.toString(statusCode)
  }
}
