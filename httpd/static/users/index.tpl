{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/users.css">
{{end}}


{{define "content"}}
<div class="content">
  {{if .SaveComplete}}
  <div class="section alert alert-success">
    <a class="close" href="#">&times;</a>
    Save Complete. <br /> User saved successfully.
  </div>
  {{end}}
  {{if .SaveError}}
  <div class="section alert alert-failed">
    <a class="close" href="#">&times;</a>
    Save Failed. <br /> User could not be saved.
  </div>
  {{end}}
  
  <div class="section">
    <header>
		<a class="action" href="/users/edit" style="padding-left:1%;padding-right:1%;">New User</a>
    	<h2>Users</h2>
    </header>
    <div class="sectionbody">
    {{if .AllUsers}}
      <table>
        <thead>
          <tr>
            <th>Username</th>
            <th>Name</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {{$users := .AllUsers}}
          {{range $user := $users}}
          <tr>
            <td><a href="/users/{{$user.Username}}/edit">{{$user.Username}}</a></td>
            <td><a href="/users/{{$user.Username}}/edit">{{if $user.DisplayName}}{{$user.DisplayName}}{{else}}&nbsp;{{end}}</a></td>
            <td><a class="action" href="/users/{{$user.Username}}/edit">Edit</a></td>
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
</div>
{{end}}