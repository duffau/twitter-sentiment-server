
function websocketStart()
{
	var ws = new WebSocket("ws://localhost:8080/websocket");

	ws.onopen = function()
	{
      	// Web Socket is connected, send data using send()
      	var textArea = document.getElementById("textarea1");
      	textArea.value = textArea.value + "\n" + "Connection opened";
	};

	ws.onmessage = function(evt)
	{
   		var textArea = document.getElementById("textarea1");
    	textArea.value = textArea.value + "\n" + evt.data
      	var msg = JSON.parse(evt.data);
      	console.log(msg)
      	textArea.scrollTop = textArea.scrollHeight;
   		
   		var time = new Date(msg["Timestamp"]*1000);

		var update = {
			x:  [[time]],
			y: [[msg["Sentiment"]]]
		}
  
		var olderTime = time.setMinutes(time.getMinutes() - 1);
		var futureTime = time.setMinutes(time.getMinutes() + 1);
  
		var minuteView = {
	      xaxis: {
	        type: 'date',
	        range: [olderTime,futureTime]
	      }
	    };
  
	  Plotly.relayout('sentiment-graph', minuteView);
	  Plotly.extendTraces('sentiment-graph', update, [0])

   };

   ws.onclose = function()
   {
      var textArea = document.getElementById("textarea1");
      textArea.value = textArea.value + "\n" + "Connection closed";
   };

}


