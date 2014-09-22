{{define "javascript[additional]"}}
<script type="text/javascript" src="/js/jquery-2.1.1.min.js"></script>
<script type="text/javascript" >
	$( document ).ready(function() {
		
		$(".record, .record .expand, .device, .device .expand").click(function(e){

				expand = $(this);
				
				if (expand.is("div")) {
					expand = expand.children("a").eq(0);
				}
			
				expand.parent().next("div").toggle();
				
				if (expand.parent().next("div").css("display") === 'none'){
					expand.html("+");
				} else {
					expand.html("-");
				}
				
				return false;
			}
		);
	});
</script>
{{end}}



{{define "css[additional]"}}
  <link rel="stylesheet" href="/css/internet.css">
{{end}}

{{define "content"}}
<div class="content">
  <div class="section">
    <header>
      <h2>Current Internet Connections</h2>
    </header>
    <div class="sectionbody">
		{{if .Connections}}	
			{{range $username := .Connections.Usernames}}
			{{$usersdevices := ($.Connections.DeviceIDsForUsername $username)}}
			<section>
				{{if gt (len $usersdevices) 1}}
					<div class="record"><a class="expand" href="#">+ </a> {{$username}}</div>
					<div style="display:none;">
					{{range $deviceid := $usersdevices}}
					{{$device := index $.Devices ($deviceid.String)}}
						<div class="device">
							<a class="expand" href="#">+ </a> {{$device.DisplayName}}
						</div>
						<div class="connections" style="display:none;">
						{{$connections := $.Connections.ConnectionsForUsernameAndDeviceID $username ($deviceid.String)}}
						{{range $connection := $connections}}
							<div>
								<a class="action" href="{{$connection.SourceIP.String}}/{{$connection.SourcePort}}/{{$connection.DestinationIP.String}}/{{$connection.DestinationPort}}//">Info</a>
								{{$hostname := index $.Hostnames ($connection.DestinationIP.String)}}
								<a href="{{$connection.SourceIP.String}}/{{$connection.SourcePort}}/{{$connection.DestinationIP.String}}/{{$connection.DestinationPort}}/" style="display:block;">
								{{if $hostname}}
									{{$hostname}}
								{{else}}
									{{$connection.DestinationIP.String}}
								{{end}}
								</a>
							</div>
						{{end}}
						</div>					
					{{end}}
					</div>
				{{else}}
				{{$deviceid := index $usersdevices 0}}
				{{$device := index $.Devices ($deviceid.String)}}
				<div class="record"><a class="expand" href="#">+ </a> {{$username}} on {{$device.DisplayName}}</div>
				<div style="display:none;">
					<div class="connections">
						{{$connections := $.Connections.ConnectionsForUsernameAndDeviceID $username ($deviceid.String)}}
						{{range $connection := $connections}}
							<div>
								<a class="action" href="{{$connection.SourceIP.String}}/{{$connection.SourcePort}}/{{$connection.DestinationIP.String}}/{{$connection.DestinationPort}}/">Info</a>
								<a href="{{$connection.SourceIP.String}}/{{$connection.SourcePort}}/{{$connection.DestinationIP.String}}/{{$connection.DestinationPort}}/" style="display:block;">
								{{$hostname := index $.Hostnames ($connection.DestinationIP.String)}}
								{{if $hostname}}
									{{$hostname}}
								{{else}}
									{{$connection.DestinationIP.String}}
								{{end}}
								</a>
							</div>
						{{end}}
					</div>					
				</div>
				{{end}}				
			</section>
			{{end}}
			{{end}}
    </div>
  </div>
</div>
{{end}}