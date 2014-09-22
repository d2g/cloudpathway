{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/access.css">
{{end}}

{{define "javascript[additional]"}}
		<script type="text/javascript" src="/js/jquery-2.1.1.min.js"></script>
		<script type="text/javascript">        
    		$(document).ready(function(){				
				{{$assignedCollections := .UserFilterCollections.Collections}}
				var assignedCollections = $.parseJSON('{"assignedCollections":[{{range $index, $assignedCollection := $assignedCollections}}{{if $index}},{{end}}"{{$assignedCollection}}"{{end}}]}').assignedCollections;
				
				
				// Update the displayed collections.
				var UpdateCollections = function() {
					$(assignedCollections).each(function(i, collection) {
						// Add a hidden field to the form for saving this collection.
						$("form").append(
							"<input type='hidden' name='collections[]' value='" + collection + "' />"
						);

						var collectionID = collection.replace( /(:|\.|\+|\[|\])/g, "\\$1" )
						$("#" + collectionID).removeClass("label-default")
						$("#" + collectionID).addClass("label-primary")
					});
				};
				
				// Add the collection.
				$("#addCollection").click(function() {					
					var collection = $('#collectionName').val();
					
					console.log(collection)
					
					// Add the collection to the collections list.
					collections.push(collection);
					
					// Add a hidden field to the form for saving this collection.
					$("#userAccessForm").append(
						"<input type='hidden' name='collections[]' value='" + collection + "' />"
					);
					
					// Update the site listing and reset the input box to blank.
					$('#siteAddress').val("");
					UpdateSites();
				});
				
				// Select or unselect blocked lists on click.
				$(document).on('click', '#filterCollectionsList .gridListItem', function() {
					var collection = $(this).attr('id');

					// Alter the colour of the selected item.
					if ($(this).hasClass("label-default")) {
						// Change to selected, add hidden field etc.
						// Add the collection to the collections list.
						assignedCollections.push(collection);
					
						// Add a hidden field to the form for saving this collection.
						$("#userAccessForm").append(
							"<input type='hidden' name='collections[]' value='" + collection + "' />"
						);
						
						// Change colour.
						$(this).removeClass("label-default")
						$(this).addClass("label-primary")
						
					} else {
						// Change from selected, remove hidden field etc.
						// Remove the collection from the list of collections.
						var index = assignedCollections.indexOf(collection);
						if (index > -1) {
							assignedCollections.splice(index, 1);
						}
					
						// Remove the form field relating to this collection.
						$('input[value="' + collection + '"]').remove();
						
						// Change colour.
						$(this).removeClass("label-primary")
						$(this).addClass("label-default")
					}
            	});

				UpdateCollections();
    		});
    	</script>
{{end}}

{{define "content"}}
<div class="content">
  <div class="section">
    <header>
      <h2>
      {{$filterCollectionsUser := .FilterCollectionsUser}}
      Edit access for {{.FilterCollectionsUser.GetDisplayName}}      
      </h2>
    </header>
    <form class="sectionbody" method="post" action="/useraccess/save">
			<input type="hidden" name="idName" value="{{.FilterCollectionsUser.Username}}" />
			<label for="user">User:</label>{{.FilterCollectionsUser.GetDisplayName}}

			<label for="collectionFilter">Filter:</label>
			<input type="text" name="collectionFilter" />
			
			
			<label>Blocked Lists:</label>
						
			<br style="clear:both;"/>
			
			<div id="filterCollectionsList">
        		{{$allCollections := .AllCollections}}
				{{range $collection := $allCollections}}							
          			<div class="record" id="{{$collection.EscapedName}}" title="{{$collection.Name}}">
						
						<p>{{$collection.Name}}</p>
					</div>
        		{{end}}
      		</div>
									
			
			<a href="/useraccess/" class="action">Cancel</a>
			<button class="action" type="submit">Save</button>
			
      <br style="clear:both;" />
    </form>	
  </div>
</div>    
{{end}}