component Components.GrumbleForm {
  property setGrumbleContent : Function(String, Promise(Never, Void))

  fun handleInput (
    onChange : Function(String, Promise(Never, Void)),
    event : Html.Event
  ) : a {
    onChange(Dom.getValue(event.target))
  }

  style warning {
    color: red;
  }

  fun render : Html {
    <form>
      <div class="field">
        <textarea
          class="textarea"
          maxlength="300"
          onChange={handleInput(setGrumbleContent)}/>

        <small>"＊1～300文字の範囲で入力可能です。"</small>
        <br/>
        <small::warning>"＊なにをぼやいても自由ですが、誹謗中傷や犯罪予告などを書き込まないようにしましょう。大いなる自由には大いなる責任が伴います。"</small>
      </div>
    </form>
  }
}
