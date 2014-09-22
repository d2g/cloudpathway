{{define "css[additional]"}}
	<link rel="stylesheet" href="/css/collections.css">
{{end}}

{{define "javascript[additional]"}}
	<script type="text/javascript" src="/js/jquery-2.1.1.min.js"></script>
	
	<script type="text/javascript">
   		$(document).ready(function(){	
		
			
   		});
   	</script>
{{end}}

{{define "content"}}
<div class="content">
  	<div class="section">
	    <header>
	      <h2>
	        Editing <span class="boldFont">{{.Collection.Name}}</span>
	      </h2>
	    </header>
	    <form class="sectionbody" method="post" action="/collections/save">
	    	<div class="section alert alert-success" id="duplicateAddressWarning" style="display: none;">
	    		<a class="close" href="#">&times;</a>
	          	Duplicate address entered. <br /> This address is already blocked.
	      	</div>
	      
			<input type="hidden" name="id" value="{{.Collection.Name}}"/>
							
			<label for="siteAddress">Site Address:</label>
			<input type="text" name="siteAddress" id="siteAddress" />
				
			<button id="addBlockedSite" style="clear:left;" class="action" type="button">Add</button>
				
	      	<br style="clear:both;"/>
	      	<br />
		
	      	<label>Blocked Sites:</label>
	      	<br style="clear:both;"/>
		
		  	<div id="blocked_sites">
			{{$domains := .Collection.Domains}}
	        {{range $index, $domain := $domains}}
				{{if lt $index 300}}
				<div class="record">
					<a class="close" href="/collections/{{$.Collection.Name}}/url/remove/{{$domain}}">&times;</a>
					<p>{{$domain}}</p>
				</div>		        
				{{end}}
	        {{end}}
			</div>
	      
	      	<br style="clear:both;" />
	    </form>	
  	</div>
</div>    
{{end}}