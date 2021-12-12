component Link {
  property children : Array(Html) = []
  property href : String

  style link {
    &:hover {
      text-decoration: underline;
    }
  }

  fun render : Html {
    <a::link
      href="#{href}"
      target="_blank">

      <{ children }>

    </a>
  }
}
