'use strict';

viewModel.TDA = new Object();
var pg = viewModel.TDA;
var heightSub = 200;

vm.currentMenu('3D Analytic');
vm.currentTitle('3D Analytic');
vm.breadcrumb([{ title: '3D Analytic', href: viewModel.appName + 'page/threedanalytic' }]);

pg.createChart = function() {
    $('#chart').html('');

    Plotly.d3.csv('https://raw.githubusercontent.com/plotly/datasets/master/api_docs/mt_bruno_elevation.csv', function(err, rows){
    function unpack(rows, key) {
    return rows.map(function(row) { return row[key]; });
    }
    
    var z_data=[ ]
    for(var i=0;i<24;i++)
    {
        z_data.push(unpack(rows,i));
    }

    var data = [{
            z: z_data,
            type: 'surface'
            }];
    
    var layout = {
        title: '',
        showlegend: false,
        cliponaxis: false,
        autosize: true,
        width: $('#chart').parent().width(),
        height: 480,
        margin: {
            l: 0,
            r: 0,
            b: 20,
            t: 20,
        },
        scene: {
            xaxis: {title: 'Temperature (&deg;C)'},
            yaxis: {title: 'Wind Speed (m/s)'},
            zaxis: {title: 'Active Power (kW)'}
        }
    };
    Plotly.newPlot('chart', data, layout, {displaylogo: false});
    });
}

pg.createChart1 = function() {
    Plotly.d3.csv('https://raw.githubusercontent.com/plotly/datasets/master/3d-line1.csv', function(err, rows){
    var unpack = function(rows, key) {
    return rows.map(function(row) 
    { return row[key]; }); }
            
    var x = unpack(rows , 'x');
    var y = unpack(rows , 'y');
    var z = unpack(rows , 'z'); 
    var c = unpack(rows , 'color');
    Plotly.plot('chart1', [{
    type: 'scatter3d',
    mode: 'lines',
    x: x,
    y: y,
    z: z,
    opacity: 1,
    line: {
        width: 6,
        color: c,
        reversescale: false
    }
    }], {
    height: 640
    });
    });
}

pg.createChart2 = function() {
    Plotly.d3.json('https://raw.githubusercontent.com/plotly/datasets/master/3d-ribbon.json', function(figure){

    var trace1 = {
        x:figure.data[0].x, y:figure.data[0].y, z:figure.data[0].z,
        name: '',
        colorscale: figure.data[0].colorscale,
        showscale: false
    }
    var trace2 = {
        x:figure.data[1].x, y:figure.data[1].y, z:figure.data[1].z,
        name: '',
        colorscale: figure.data[1].colorscale,
        type: 'surface',
        showscale: false
    }
    var trace3 = {
        x:figure.data[2].x, y:figure.data[2].y, z:figure.data[2].z,
        colorscale: figure.data[2].colorscale,
        type: 'surface',
        showscale: false
    }
    var trace4 = {
        x:figure.data[3].x, y:figure.data[3].y, z:figure.data[3].z,
        colorscale: figure.data[3].colorscale,
        type: 'surface',
        showscale: false
    }
    var trace5 = {
        x:figure.data[4].x, y:figure.data[4].y, z:figure.data[4].z,
        colorscale: figure.data[4].colorscale,
        type: 'surface',
        showscale: false
    }
    var trace6 = {
        x:figure.data[5].x, y:figure.data[5].y, z:figure.data[5].z,
        colorscale: figure.data[5].colorscale,
        type: 'surface',
        showscale: false
    }
    var trace7 = {
        x:figure.data[6].x, y:figure.data[6].y, z:figure.data[6].z,
        name: '',
        colorscale: figure.data[6].colorscale,
        type: 'surface',
        showscale: false
    }
    
    var data = [trace1, trace2, trace3, trace4, trace5, trace6, trace7];

    var layout = {
    title: 'Ribbon Plot',
    showlegend: false,
    autosize: true,
    width: $('#chart2').parent().width(),
    height: 480,
    scene: {
        xaxis: {title: 'Sample #'},
        yaxis: {title: 'Wavelength'},
        zaxis: {title: 'OD'}
    }
    };
    Plotly.newPlot('chart2', data, layout);
    });
}

pg.createChart3 = function() {
    var z1 = [
        [8.83,8.89,8.81,8.87,8.9,8.87],
        [8.89,8.94,8.85,8.94,8.96,8.92],
        [8.84,8.9,8.82,8.92,8.93,8.91],
        [8.79,8.85,8.79,8.9,8.94,8.92],
        [8.79,8.88,8.81,8.9,8.95,8.92],
        [8.8,8.82,8.78,8.91,8.94,8.92],
        [8.75,8.78,8.77,8.91,8.95,8.92],
        [8.8,8.8,8.77,8.91,8.95,8.94],
        [8.74,8.81,8.76,8.93,8.98,8.99],
        [8.89,8.99,8.92,9.1,9.13,9.11],
        [8.97,8.97,8.91,9.09,9.11,9.11],
        [9.04,9.08,9.05,9.25,9.28,9.27],
        [9,9.01,9,9.2,9.23,9.2],
        [8.99,8.99,8.98,9.18,9.2,9.19],
        [8.93,8.97,8.97,9.18,9.2,9.18]
    ];
    
    var z2 = [];
    for (var i=0;i<z1.length;i++ ) { 
      var z2_row = [];
        for(var j=0;j<z1[i].length;j++) { 
          z2_row.push(z1[i][j]+1);
        }
        z2.push(z2_row);
    }
    
    var z3 = []
    for (var i=0;i<z1.length;i++ ) { 
      var z3_row = [];
        for(var j=0;j<z1[i].length;j++) { 
          z3_row.push(z1[i][j]-1);
        }
        z3.push(z3_row);
    }
    var data_z1 = {z: z1, type: 'surface'};
    var data_z2 = {z: z2, showscale: false, opacity:0.9, type: 'surface'};
    var data_z3 = {z: z3, showscale: false, opacity:0.9, type: 'surface'};
    
    Plotly.newPlot('chart3', [data_z1, data_z2, data_z3]);
}

pg.customVis = function(x, y) {
    return (Math.sin(x/50) * Math.cos(y/50) * 50 + 50);
};

pg.createChartVis = function() {
    var data = null;
    var graph = null;

    // Create and populate a data table.
    data = new vis.DataSet();
    // create some nice looking data with sin/cos
    var counter = 0;
    var steps = 50;  // number of datapoints will be steps*steps
    var axisMax = 314;
    var axisStep = axisMax / steps;
    for (var x = 0; x < axisMax; x+=axisStep) {
      for (var y = 0; y < axisMax; y+=axisStep) {
        var value = pg.customVis(x,y);
        data.add({id:counter++,x:x,y:y,z:value,style:value});
      }
    }

    // specify options
    var options = {
      width:  '100%',
      height: '480px',
      style: 'surface',
      showPerspective: true,
      showGrid: true,
      tooltip: true,
      showShadow: false,
      keepAspectRatio: true,
      verticalRatio: 0.5
    };

    // Instantiate our graph object.
    var container = document.getElementById('chartVis');
    graph = new vis.Graph3d(container, data, options);
}

pg.createChartVis1 = function() {
    var data1 = null;
    var graph1 = null;

    var style = 'bar';
      var showPerspective = true;
      var xBarWidth = parseFloat(0) || undefined;
      var yBarWidth = parseFloat(0) || undefined;
      var withValue = ['bar-color', 'bar-size', 'dot-size', 'dot-color'].indexOf(style) != -1;

      // Create and populate a data table.
      data1 = [];

      // create some nice looking data with sin/cos
      var steps = 5;  // number of datapoints will be steps*steps
      var axisMax = 10;
      var axisStep = axisMax / steps;
      for (var x = 0; x <= axisMax; x+=axisStep) {
        for (var y = 0; y <= axisMax; y+=axisStep) {
          var z = pg.customVis(x,y);
          if (withValue) {
            var value = (y - x);
            data1.push({x:x, y:y, z: z, style:value});
          }
          else {
            data1.push({x:x, y:y, z: z});
          }
        }
      }

      // specify options
      var options = {
        width:  '600px',
        height: '600px',
        style: style,
        xBarWidth: xBarWidth,
        yBarWidth: yBarWidth,
        showPerspective: showPerspective,
        showGrid: true,
        showShadow: false,
        tooltip: true,
        keepAspectRatio: true,
        verticalRatio: 0.5
      };
  
      var camera = graph1 ? graph1.getCameraPosition() : null;

      // create our graph
      var container = document.getElementById('chartVis1');
      graph1 = new vis.Graph3d(container, data1, options);

      if (camera) graph1.setCameraPosition(camera); // restore camera position
}

pg.createChartVis2 = function() {
    var data = null;
    var graph = null;

    data = new vis.DataSet();

    // create some shortcuts to math functions
    var sqrt = Math.sqrt;
    var pow = Math.pow;
    var random = Math.random;

    // create the animation data
    var imax = 100;
    for (var i = 0; i < imax; i++) {
      var x = pow(random(), 2);
      var y = pow(random(), 2);
      var z = pow(random(), 2);
      var style = (i%2==0) ? sqrt(pow(x, 2) + pow(y, 2) + pow(z, 2)) : "#00ffff";

      data.add({x:x,y:y,z:z,style:style});
    }
    console.log(data)
    // specify options
    var options = {
      width:  '600px',
      height: '600px',
      style: 'dot-color',
      showPerspective: true,
      showGrid: true,
      keepAspectRatio: true,
      verticalRatio: 1.0,
      legendLabel: 'distance',
    //   onclick: onclick,
    tooltip: true,
      cameraPosition: {
        horizontal: -0.35,
        vertical: 0.22,
        distance: 1.8
      }
    };

    // create our graph
    var container = document.getElementById('chartVis2');
    graph = new vis.Graph3d(container, data, options);
}

pg.createChartVis3 = function() {
    var data = null;
    var graph = null;

    // create the data table.
    data = new vis.DataSet();

    var custom = function (x, y, t) {
        return Math.sin(x/50 + t/10) * Math.cos(y/50 + t/10) * 50 + 50;
    }

    var steps = 25;
      var axisMax = 314;
      var tMax = 31;
      var axisStep = axisMax / steps;
      for (var t = 0; t < tMax; t++) {
        for (var x = 0; x < axisMax; x+=axisStep) {
          for (var y = 0; y < axisMax; y+=axisStep) {
            var value = custom(x, y, t);
            data.add([
              {x:x,y:y,z:value,filter:t,style:value}
            ]);
          }
        }
      }

      // specify options
      var options = {
        width:  '600px',
        height: '600px',
        style: 'surface',
        showPerspective: true,
        showGrid: true,
        showShadow: false,
        // showAnimationControls: false,
        keepAspectRatio: true,
        verticalRatio: 0.5,
        animationInterval: 100, // milliseconds
        animationPreload: true,
        animationAutoStart: true,
      };

      // create our graph
      var container = document.getElementById('chartVis3');
      graph = new vis.Graph3d(container, data, options);
}

pg.getData = function() {
    // pg.createChart();
    // pg.createChart1();
    // pg.createChart2();
    // pg.createChart3();
    // pg.createChartVis();
    // pg.createChartVis1();
    pg.createChartVis2();
    // pg.createChartVis3();
}

pg.initLoad = function() {
    window.setTimeout(function(){
        fa.LoadData();
        di.getAvailDate();
        app.loading(false);

        pg.refresh();
    }, 200);
}

pg.refresh = function() {
    fa.checkTurbine();

    pg.getData();
}

$(function(){
    $('#projectList').kendoDropDownList({
        change: function () {  
            di.getAvailDate();
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.project = project;
            fa.populateTurbine(project);
        }
    });
    $('#btnRefresh').on('click', function () {
        pg.refresh();
    });

    pg.initLoad();
})
