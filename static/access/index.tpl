{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/access.css">
{{end}}


{{define "content"}}
<div class="content">
  {{if .SaveComplete}}
  <div class="section alert alert-success">
    <a class="close" href="#">&times;</a>
    Save Complete. <br /> Useraccess saved successfully.
  </div>
  {{end}}
  {{if .SaveError}}
  <div class="section alert alert-failed">
    <a class="close" href="#">&times;</a>
    Save Failed. <br /> Useraccess could not be saved.
  </div>
  {{end}}

  <div class="section">
    <header>
      <h2>User Access Settings</h2>
    </header>
    <div class="sectionbody">
    {{if .AllUsers}}
      <table>
        <thead>
          <tr>
            <th>Username</th>
            <th>Number of Assigned Lists</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {{$users := .AllUsers}}
          {{range $user := $users}}
          <tr>
            <td><a href="/useraccess/{{$user.Username}}/edit">{{$user.Username}}</a></td>
            <td><a href="/useraccess/{{$user.Username}}/edit">{{$user.NumberOfFilterCollections}}</a></td>
            <td><a class="action" href="/useraccess/{{$user.Username}}/edit">Edit</a></td>
          </tr>
          {{end}}
        </tbody>
      </table>
      {{else}}
			<!-- No users to display -->
      There are currently no users to display.
      {{end}}
    </div>
  </div>
{{end}}