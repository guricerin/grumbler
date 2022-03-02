enum TimelineResultKind {
  Initial
  Grumbles(Grumbles)
}

store Stores.Timeline {
  state apiStatus : Api.Status(TimelineResultKind) = Api.Status::Initial

  fun getTimeline (userId : String) : Promise(Never, Void) {
    sequence {
      status =
        Http.get("#{@ENDPOINT}/auth/user/#{userId}/timeline")
        |> Api.send(Grumbles.decodes)

      case (status) {
        Api.Status::Ok(grumbles) =>
          sequence {
            result =
              TimelineResultKind::Grumbles(grumbles)

            next { apiStatus = Api.Status::Ok(result) }
          }

        Api.Status::Initial => next { apiStatus = Api.Status::Initial }
        Api.Status::Error(err) => next { apiStatus = Api.Status::Error(err) }
      }
    }
  }
}
