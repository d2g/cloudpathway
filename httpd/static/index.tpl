{{define "css[additional]"}}
  <link rel="stylesheet" href="css/index.css">
{{end}}

{{define "content"}}
<div class="content">
  <div class="section">
    <header>
      <h2>Search</h2>
    </header>
    <!-- 
      TODO: 
        - Search Engine Should be customisable or atleast the country 
        - Client Should Represent The version.
    -->
    <form class="sectionbody" method="GET" action="https://www.google.co.uk/search">
      <input type="text" name="q" autofocus="autofocus"/>
      <input type="hidden" name="safe" value="active"/>
      <input type="hidden" name="client" value="cp[alpha]"/>
      <button type="submit">Search</button>
    </form>
  </div>
  <div class="section">
    <header>
      <h2>Devices You're Using</h2>
    </header>
    <div class="sectionbody">
    {{if .CurrentDevices}}
			{{range .CurrentDevices}}    
      <div class="record">
        <a class="action" href="devices/{{.MACAddress}}/removeuser">Logout</a>
        <p>{{.GetDisplayName}}</p>
      </div>
      {{end}}
    {{end}}
    </div>
  </div>
  <div class="section">
    <header>
      <h2>My Devices</h2>
    </header>
    <div class="sectionbody">
    {{if .DefaultDevices}}
			{{range .DefaultDevices}}
      <div class="record">
        <a class="action" href="devices/{{.MACAddress}}/edit">Edit</a>
        <p>{{.GetDisplayName}}</p>
      </div>
      {{end}}
    {{end}}
    </div>
  </div>
</div>
{{end}}