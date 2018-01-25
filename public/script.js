
function websocketStart()
{
	var ws = new WebSocket("ws://localhost:8080/websocket");
	var tweets = new Array();

	ws.onopen = function()
	{
      	// Web Socket is connected, send data using send()
      	var textArea = document.getElementById("textarea1");
      	textArea.value = textArea.value + "\n" + "Connection opened";
	};

	ws.onmessage = function(evt)
	{
      	var msg = JSON.parse(evt.data);
    	tweets.push(msg["Text"] + "\n------- New Tweet -------\n")
 
	   	if (tweets.length > 5) {
    		tweets.shift()
    	}

   		var textArea = document.getElementById("textarea1");
    	textArea.value = tweets.join("\n") 
      	textArea.scrollTop = textArea.scrollHeight;
   		
   		var time = new Date(msg["Timestamp"]*1000);
		var update = {
			x:  [[time], [time]],
			y: [[msg.Sentiment], [msg.SentimentFilter]]
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
	  Plotly.extendTraces('sentiment-graph', update, [0, 1])

   };

   ws.onclose = function()
   {
      var textArea = document.getElementById("textarea1");
      textArea.value = textArea.value + "\n" + "Connection closed";
   };

}


