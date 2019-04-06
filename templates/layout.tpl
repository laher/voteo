<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="x-ua-compatible" content="ie=edge,chrome=1" />

    <title>Voteo</title>
    <meta name="description" content="Video voting app" />

    <link rel="stylesheet" type="text/css" href="./static/style.css" />
    <link rel="home" href="https://voteo.laher.net.nz/" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1"
    />
    <link rel="shortcut icon" type="image/png" href="/static/assets/favicon.png" />
    <link rel="canonical" href="https://github.com/laher/voteo" />

    <!-- Open graph -->
    <meta property="og:title" content="Voteo" />
    <meta property="og:description" content="Video voting app" />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://voteo.laher.net.nz/" />
    <meta property="og:image" content="./static/assets/card.png" />
    <meta property="og:image:secure_url" content="./static/assets/card.png" />
    <meta property="og:image:type" content="image/png" />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta property="og:image:alt" content="Voteo" />

    <!-- Twitter Card -->
    <meta name="twitter:card" value="summary" />
    <meta name="twitter:url" content="https://voteo.laher.net.nz/" />
    <meta name="twitter:title" content="Voteo" />

    <!-- Android web app -->
    <link rel="manifest" href="/static/manifest.webmanifest" />
    <meta name="mobile-web-app-capable" content="yes" />
    <meta name="theme-color" content="#323232" />

    <!-- IOS web app -->
    <meta name="apple-mobile-web-app-capable" content="yes" />
    <meta name="apple-mobile-web-app-status-bar-style" content="#323232" />
    <meta name="apple-mobile-web-app-title" content="Voteo" />
    <link
      rel="apple-touch-icon"
      sizes="180x180"
      href="./static/assets/icon-512x512.png"
    />

    <!-- Windows web app -->
    <meta name="msapplication-TileImage" content="./static/assets/icon-512x512.png" />
    <meta name="msapplication-TileColor" content="#323232" />

    <!-- this style seems appropriate for the main 'app' -->
    <script type="module" src="./static/app.js"></script>
    <!-- second style seems appropriate for a 'library' -->
    <script type="module">
      import { showSignInModal } from './static/auth-okta.js';
      document.getElementById('sign-in').addEventListener('click', event => {
        event.preventDefault();
        showSignInModal();
      });
    </script>
    <script
      src="https://ok1static.oktacdn.com/assets/js/sdk/okta-auth-js/2.0.1/okta-auth-js.min.js"
      type="text/javascript"
    ></script>
    <script
      src="https://ok1static.oktacdn.com/assets/js/sdk/okta-signin-widget/2.6.0/js/okta-sign-in.min.js"
      type="text/javascript"
    ></script>
    <link
      href="https://ok1static.oktacdn.com/assets/js/sdk/okta-signin-widget/2.6.0/css/okta-sign-in.min.css"
      type="text/css"
      rel="stylesheet"
    />
  </head>
  <body>
    <div>
      <div id="nav">
        <ul id="logged-in" style="{{ if .PersonID }}{{ else }}display:none{{ end }}">
          <li>
            <abbr title="Voteo is a video voting app for teams">Voteo</abbr>
          </li>
          <li id="name">{{ .PersonID }}</li>
          <li class="button">
            <a href="#" id="sign-out" class="">Sign out</a>
          </li>
        </ul>
        <div id="logged-out" class="loggedout" style="{{ if .PersonID }}display:none{{ end }}" >
          <ul>
            <li>
              <abbr title="Voteo is a video voting app for teams">Voteo</abbr>
            </li>
            <li>guest</li>
            <li></li>
          </ul>
          <ul class="right">
            <li class="button">
              <a href="/register" target="_blank" id="register" class=""
                >Register</a
              >
            </li>
            <li class="button">
              <a href="#" id="sign-in" class="">Sign in</a>
            </li>
          </ul>
        </div>
      </div>
      <div id="widget-container"></div>
      <div id="app-container" class="container">

        {{ template "index.tpl" . }}
      </div>
    </div>
  </body>
</html>
