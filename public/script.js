



function websocketStart()
{
   var ws = new WebSocket("ws://localhost:8080/websocket");

   ws.onopen = function()
   {
      // Web Socket is connected, send data using send()
      var myTextArea = document.getElementById("textarea1");
      myTextArea.value = myTextArea.value + "\n" + "Connection opened";
   };

   ws.onmessage = function (evt)
   {
      var myTextArea = document.getElementById("textarea1");
      myTextArea.value = myTextArea.value + "\n" + evt.data
      if(evt.data == "pong") {
        setTimeout(function(){ws.send("ping");}, 2000);
      }
   };

   ws.onclose = function()
   {
      var myTextArea = document.getElementById("textarea1");
      myTextArea.value = myTextArea.value + "\n" + "Connection closed";
   };

}
