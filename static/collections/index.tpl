{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/collections.css">
{{end}}


{{define "content"}}
<div class="content">
	{{if .SaveComplete}}
  <div class="section alert alert-success">
    <a class="close" href="#">&times;</a>
    Save Complete. <br /> Blocked site list saved successfully.
  </div>
	{{end}}
									
  {{if .SaveError}}
  <div class="alert alert-danger">
    <a class="close" href="#">&times;</a>
    Save Failed. <br /> Blocked site list could not be saved.
  </div>
  {{end}}
  
  <div class="section">
    <header>
		<a class="action" href="/collections/create" style="padding-left:1%;padding-right:1%;">New Category</a>
    	<h2>Blocked Sites</h2>
    </header>
    <div class="sectionbody">
		{{if .AllCollections}}
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Number of Blocked Sites</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {{$collections := .AllCollections}}
          {{range $collection := $collections}}
          <tr>
            <td><a href="/collections/{{$collection.Name}}/edit">{{$collection.Name}}</a></td>
            <td><a href="/collections/{{$collection.Name}}/edit">{{$collection.GetNumberOfDomains}}</a></td>
            <td><a class="action" href="/collections/{{$collection.Name}}/edit">Edit</a></td>
          </tr>
          {{end}}
        </tbody>
      </table>
      {{end}}
    </div>
  </div>
</div>
{{end}}