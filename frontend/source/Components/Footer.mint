component Footer {
  style stickyFooter {
    background-color: #DDDDDD;
    display: flex;
    flex: 1;
    min-height: 1vh;
    justify-content: center;
  }

  fun render : Html {
    <footer::stickyFooter class="footer">
    <div class="content has-text-centered">
      "powerd by "

      <Link href="https://www.mint-lang.com/">
        "mint-lang"
      </Link>
    </div>
    </footer>
  }
}
