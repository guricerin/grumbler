enum SearchResultKind {
  Users(Users)
}

store Stores.Search {
  state searchWord : String = ""
  state apiStatus : Api.Status(SearchResultKind) = Api.Status::Initial
  state rsltUsers : Users = Users.empty()

  fun setSearchWord (v : String) : Promise(Never, Void) {
    next { searchWord = v }
  }

  fun setApiStatus (v : Api.Status(SearchResultKind)) : Promise(Never, Void) {
    next { apiStatus = v }
  }

  fun resetSearchResult : Promise(Never, Void) {
    next { rsltUsers = Users.empty() }
  }

  fun setRsltUsers (users : Users) : Promise(Never, Void) {
    next { rsltUsers = users }
  }
}
