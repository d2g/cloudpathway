{{define "metadata"}}
	<title>Cloud Pathway | Authentication Required</title>
{{end}}

{{define "content"}}
	<div class="content">
		You need to 
		<a href="http://cloudpathway.d2g.org.uk/">login</a>
		gain access to the internet.
	</div>
{{end}}


{{define "css"}}
  	<link href="http://fonts.googleapis.com/css?family=Ubuntu:400,700|Arvo:700" rel="stylesheet" type="text/css">
  	<link rel="stylesheet" href="http://cloudpathway.d2g.org.uk/css/standardize.css">
  	<link rel="stylesheet" href="http://cloudpathway.d2g.org.uk/css/styles.css">
  	<link rel="stylesheet" href="http://cloudpathway.d2g.org.uk/css/authentication.css">	
{{end}}

{{define "head"}}
  {{template "metadata" .}}

  {{template "css"}}
{{end}}

{{define "page"}}{{end}}

{{define "navigation"}}{{end}}

<!DOCTYPE html>
<html>
    <head>
      {{template "head" .}}
    </head>
    <body>
      <header class="siteheader">
        <img class="logo" src="http://cloudpathway.d2g.org.uk/images/logo.png">
        <h1>Cloud Pathway</h1>
      </header>

      {{template "navigation" .}}
      
        <div class="container">
          {{template "content" .}}  
        </div>
    </body>
</html>