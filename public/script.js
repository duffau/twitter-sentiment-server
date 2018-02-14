var ws
var plotpoints = 0;
var maxplotpoints = 5000;


function websocketStart()
{
	ws = new WebSocket("ws://localhost:8080/websocket");
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
    	tweets.push(msg["Text"] + "\n Sentiment = " + msg["Sentiment"].toFixed(2))
 
	   	if (tweets.length > 10) {
    		tweets.shift()
    	}

   		var textArea = document.getElementById("textarea1");
    	textArea.value = tweets.join("\n\n------- New Tweet -------\n\n") 
      	textArea.scrollTop = textArea.scrollHeight;
   		
   		var time = new Date(msg["Timestamp"]*1000);
		var update = {
			x:  [[time], [time]],
			y: [[msg.Sentiment], [msg.SentimentFilter]]
		}
  		plotpoints += 1

		var olderTime = time.setMinutes(time.getMinutes() - 3);
		var futureTime = time.setMinutes(time.getMinutes() + 3);
  
		var minuteView = {
	      xaxis: {
	        type: 'date',
	        range: [olderTime,futureTime]
	      }
	    };
		
		if(plotpoints > maxplotpoints){
		    data[0].x.shift();
		    data[0].y.shift();
		    data[1].x.shift();
		    data[1].y.shift();
		} 
		
		Plotly.relayout('sentiment-graph', minuteView);
		Plotly.extendTraces('sentiment-graph', update, [0, 1])

   };

   ws.onclose = function()
   {
      var textArea = document.getElementById("textarea1");
      textArea.value = textArea.value + "\n" + "Connection closed";
   };

}


function websocketStop()
{
	ws.close()
}
