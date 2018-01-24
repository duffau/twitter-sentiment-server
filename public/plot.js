var time = new Date();
var sentiment;
var filteredSentiment;

var data = [{
  x: [time], 
  y: [sentiment],
  mode: 'lines',
  line: {color: '#80CAF6'}
}]

Plotly.plot('sentiment-graph', data);  

  
