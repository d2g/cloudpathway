{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/internet.css">
{{end}}

{{define "content"}}
<div class="content">
  <div class="section">
    <header>
      <h2>Internet Connection</h2>
    </header>
    <div class="sectionbody">
	{{if .Connection.SourceIP }}
		{{with .Connection}}
	    	{{.SourceIP}}
		{{end}}
	{{else}}
		Connection Has Closed..	
	{{end}}
    </div>
  </div>
</div>
{{end}}