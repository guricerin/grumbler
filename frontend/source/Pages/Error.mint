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
      401 => "#{code} unauthorized"
      403 => "#{code} forbidden"
      404 => "#{code} not found"
      500 => "#{code} internal server error"
      => "#{code} unexausted!"
    }
  } where {
    code =
      Number.toString(statusCode)
  }
}
