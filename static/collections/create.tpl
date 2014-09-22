{{define "css[additional]"}}
	<link rel="stylesheet" href="/css/collections.css">
{{end}}

{{define "javascript[additional]"}}

{{end}}

{{define "content"}}
<div class="content">
  	<div class="section">
	    <header>
	      <h2>
	        New Blocked List
	      </h2>
	    </header>
	    <form class="sectionbody" method="post" action="/collections/save">
			<label for="name">Name:</label>
			<input type="text" name="name" />
				
	      	<br />
														
	    	<a href="/collections/" class="action">Cancel</a>
	      	<button class="action" type="submit">Create</button>

	      	<br style="clear:both;"/>
	
	    </form>	
  	</div>
</div>    
{{end}}