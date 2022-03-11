component Components.GrumbleList {
  connect Application exposing { userStatus }
  property grumbles : Grumbles = Grumbles.empty()

  fun grumbleListItem (grumble : Grumble) : Html {
    case (userStatus) {
      /* unreachable! */
      UserStatus::Guest => Html.empty()

      UserStatus::SignIn(user) =>
        <div>
          <Components.GrumbleBox
            signinUser={user}
            grumble={grumble}/>
        </div>
    }
  }

  fun render : Html {
    <div>
      <{ Array.map(grumbleListItem, grumbles.grumbles) }>
    </div>
  }
}
