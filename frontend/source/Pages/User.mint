component Pages.User {
  connect Application exposing { userStatus }
  connect Stores.User exposing { rsrcUser }

  style profileItem {
    margin-left: 5px;
  }

  fun rsrcUserProfile (user : User) : Html {
    <div>
      <strong>"#{user.name}"</strong>
      <small>"@#{user.id}"</small>
      <p>"#{user.profile}"</p>
      <hr/>
      <a::profileItem>"ぼやき"</a>
      <a::profileItem>"フォロー"</a>
      <a::profileItem>"フォロワー"</a>
    </div>
  }

  fun render : Html {
    <div>
      <{ rsrcUserProfile(rsrcUser) }>
    </div>
  }
}
