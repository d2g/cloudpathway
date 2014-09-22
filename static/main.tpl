{{define "metadata"}}
  <title>Cloud Pathway | Home</title>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
{{end}}

{{define "css[additional]"}}{{end}}
{{define "css"}}
  <link href="http://fonts.googleapis.com/css?family=Ubuntu:400,700|Arvo:700" rel="stylesheet" type="text/css">
  <link rel="stylesheet" href="/css/standardize.css">
  <link rel="stylesheet" href="/css/styles.css">
        {{template "css[additional]" .}}        
{{end}}

{{define "javascript[additional]"}}{{end}}

{{define "javascript"}}
  {{template "javascript[additional]" .}}
{{end}}

{{define "head"}}
  {{template "metadata" .}}
  
  {{template "css" .}}

  {{template "javascript" .}}

  <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->
  <!--[if lt IE 9]>
     <script type="text/javascript" src="/js/html5shiv.js"></script>
     <script type="text/javascript" src="/js/respond.min.js"></script>
  <![endif]-->
{{end}}

{{define "page"}}{{end}}

{{define "content"}}{{end}}

{{define "navigation"}}{{end}}

<!DOCTYPE html>
<html>
    <head>
      {{template "head" .}}
    </head>
    <body>
      <header class="siteheader">
        <img class="logo" src="/images/logo.png">
        <h1>Cloud Pathway</h1>
      </header>

      {{template "navigation" .}}
      
        <div class="container">
          {{template "content" .}}  
        </div>
    </body>
</html>
