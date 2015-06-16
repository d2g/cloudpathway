{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/device.css">
{{end}}

{{define "content"}}
<div class="content">
  {{if .SaveComplete}}
  <div class="section alert alert-success">
    <a class="close" href="#">&times;</a>
    Save Complete. <br /> Device changes saved successfully.
  </div>
  {{end}}
  
  <div class="section">
    <header>
      <h2>Devices</h2>
    </header>
    <div class="sectionbody">
    {{if .AllDevices}}
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Hardware Address</th>
              <th>Belongs To</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
	      {{range $device := .AllDevices}}
            <tr>
              <td><a href="/devices/{{$device.MACAddress}}/edit">{{$device.GetDisplayName}}</a></td>
              <td><a href="/devices/{{$device.MACAddress}}/edit">{{$device.MACAddress}}</a></td>
              <td><a href="/devices/{{$device.MACAddress}}/edit">{{if $device.DefaultUser}}{{$device.DefaultUser.GetDisplayName}}{{else}}&nbsp;{{end}}</a></td>
              <td><a class="action" href="/devices/{{$device.MACAddress}}/edit">Edit</a></td>
            </tr>
	      {{end}}
          </tbody>
        </table>
    {{end}}
    </div>
  </div>
  
</div>
{{end}}