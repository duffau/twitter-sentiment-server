function makeplot() {
  Plotly.d3.csv("./data/sentiment.csv", function(data){ processData(data) } );

};

function processData(allRows) {

  console.log(allRows);
  var x = [], y = [];

  for (var i=0; i<allRows.length; i++) {
    row = allRows[i];
    x.push( row['timestamp'] );
    y.push( row['sentiment'] );
  }
  console.log( 'X',x, 'Y',y );
  makePlotly( x, y );
}

function makePlotly( x, y, standard_deviation ){
  var plotDiv = document.getElementById("plot");
  var traces = [{
    x: x,
    y: y
  }];

  Plotly.newPlot('sentiment-plot', traces,
    {title: '#bitcoin tweets sentiment'});
};

makeplot();