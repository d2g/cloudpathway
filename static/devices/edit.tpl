{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/device.css">
{{end}}

{{define "content"}}
<div class="content">
  <div class="section">
      <header>
        <h2>Editing {{if .Device.Nickname}}{{.Device.Nickname}}{{else}}{{.Device.Hostname}}{{end}}</h2>
      </header>
      <form class="sectionbody" method="POST" action="/devices/save">
        <input type="hidden" name="macAddress" value="{{.Device.MACAddress}}">
        
        <label>Hostname:</label><input type="text" disabled="true" readonly="true" value="{{.Device.Hostname}}" />
        <label>MAC Address:</label><input type="text" disabled="true" readonly="true" value="{{.Device.MACAddress}}" />
        
        <label for="nickname">Nickname:</label><input type="text" name="nickname" value="{{.Device.Nickname}}"/>
				<label for="defaultUser">Assigned To:</label>
				<select name="defaultUser">
          <option value=""> </option>
					{{$device := .Device}}
					{{range .AllUsers}}
					<option value="{{.Username}}" {{if $device.DefaultUser}}{{if eq .Username $device.DefaultUser.Username}}selected{{end}}{{end}}>{{.GetDisplayName}}</option>
					{{end}}
        </select>

        <a href="/devices/{{.Device.MACAddress}}/delete" class="action">Delete</a>
				<a class="action" href="/devices/">Cancel</a>
				<button class="action" type="submit">Save</button>
        <br style="clear:both;" />
      </form>
  </div>
</div>
{{end}}