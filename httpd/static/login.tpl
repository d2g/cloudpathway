{{define "css[additional]"}}
  <link rel="stylesheet" href="css/login.css">
{{end}}

{{define "content"}}
  <form method="POST" action="/login">
		{{if .LoginFailed}}
    <div class="error">
      <h4>Login Failed</h4>
    </div>
		{{end}}
    <label class="username">Username:</label>
    <input class="username" name="username" type="text">
    <label class="password">Password:</label>
    <input class="password" name="password" type="password">
    <button type="submit">Login</button>
  </form>
{{end}}