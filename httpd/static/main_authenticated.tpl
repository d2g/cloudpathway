{{define "content"}}
  <!-- Default Page! -->
{{end}}

{{define "navigation"}}
  <nav class="menu">
    <nav>
      <ul>
        <li {{if eq .Action "index"}}class="selected"{{end}}>
          <a href="/">Home</a>
        </li>
        <li {{if eq .Action "internet"}}class="selected"{{end}}>
          <a href="/internet/">Internet</a>
        </li>
      </ul>
    </nav>
    {{if .User.IsAdmin}}
    <nav class="admin">
      <ul>
        <li {{if eq .Action "deviceSettings"}}class="selected"{{end}}>
          <a href="/devices/">Devices</a>
        </li>
        <li {{if eq .Action "userSettings"}}class="selected"{{end}}>
          <a href="/users/">Users</a>
        </li>
      </ul>
    </nav>
    {{end}}
  </nav>
{{end}}