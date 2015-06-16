{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/users.css">
{{end}}


{{define "content"}}
<div class="content">
  <div class="section">
    <header>
      {{if .EditUser}}
      <h2>Editing {{.UserToEdit.Username}}</h2>
      {{end}}
      {{if .NewUser}}
      <h2>New User</h2>
      {{end}}      
    </header>
    <form class="sectionbody" method="post" action="/users/save" >
      {{if .EditUser}}
        <input type="hidden" name="idUsername" value="{{.UserToEdit.Username}}">
      {{end}}
      <label for="username">Username:</label><input class="form-control" type="text" name="username" {{if .EditUser}}value="{{.UserToEdit.Username}}"{{end}} /><br />
      <label for="name">Name:</label><input class="form-control" type="text" name="name" {{if .EditUser}}value="{{.UserToEdit.DisplayName}}"{{end}} /><br />
      <label for="dob">Date of Birth:</label><input class="form-control" type="text" name="dob" id="dob" {{if .EditUser}}value="{{.UserToEdit.ShortDOB}}"{{end}} placeholder="YYYY-MM-DD" /><br />
      <label for="isAdmin">Administrator:</label><input type="checkbox" id="isAdmin" name="isAdmin" value="true" {{if .EditUser}}{{if .UserToEdit.IsAdmin}}checked{{end}}{{end}} />
    
      {{if .EditUser}}<a href="/users/{{.UserToEdit.Username}}/delete" class="action">Delete</a>{{end}}
      <a class="action" href="/users/">Cancel</a>
      <button class="action" type="submit">Save</button>
      <br style="clear:both;" />    
    </form>
  </div>
</div>										
{{end}}