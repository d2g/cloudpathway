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
        <li {{if eq .Action ""}}class="selected"{{end}}>
          <a href="#">Reports</a>
        </li>
        <li {{if eq .Action ""}}class="selected"{{end}}>
          <a href="#">About</a>
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
		<!--
        <li {{if eq .Action "collectionSettings"}}class="selected"{{end}}>
          <a href="/collections/">Access Lists</a>
        </li>
        <li {{if eq .Action "userAccessSettings"}}class="selected"{{end}}>
          <a href="/useraccess/">User Access</a>
        </li>
		-->
      </ul>
    </nav>
    {{end}}
  </nav>
{{end}}