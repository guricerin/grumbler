store Stores.GrumbleDetail {
  state apiStatus : Api.Status(GrumbleDetail) = Api.Status::Initial

  fun getGrumbleDetail (grumblePk : String) : Promise(Never, Void) {
    sequence {
      status =
        Http.get("#{@ENDPOINT}/auth/grumble/#{grumblePk}")
        |> Api.send(GrumbleDetail.decodes)

      next { apiStatus = status }
    }
  }
}
