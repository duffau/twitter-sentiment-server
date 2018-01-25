var time = new Date();
var sentiment;
var filteredSentiment;

var data = [{
  x: [time], 
  y: [sentiment],
  mode: 'lines',
  line: {color: '#85929e'},
  name: 'Live sentiment'
},
{
  x: [time], 
  y: [filteredSentiment],
  mode: 'lines',
  line: {color: '#ff9900'},
  name: 'Mov. avg.'
}
]

var layout = {
  title: 'Live Twitter crypto sentiment',
  //legend: {"orientation": "h", "x":100, "y":1},
  showlegend: false
 }

Plotly.plot('sentiment-graph', data, layout);  

  
