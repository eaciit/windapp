// Highcharts plugin for displaying value information in the legend
// Author: Torstein Hønsi
// License: MIT license
// Last revision: 2013-07-29
(function (H) {
    H.Series.prototype.point = {}; // The active point
    H.Chart.prototype.callbacks.push(function (chart) {
        $(chart.container).bind('mousemove', function () {
            var legendOptions = chart.legend.options,
                hoverPoints = chart.hoverPoints;
            
            if (!hoverPoints && chart.hoverPoint) {
                hoverPoints = [chart.hoverPoint];
            }
            if (hoverPoints) {

                H.each(chart.legend.allItems,function(i, res){
                    H.each(chart.series[res].points, function(val){
                        if(val != undefined){
                            if(val.category == hoverPoints[0].category){
                               val.series.point = val;
                            }
                        }
                    });
                });

                H.each(chart.legend.allItems, function (item) {
                    item.legendItem.attr({
                        text: legendOptions.labelFormat ? 
                            H.format(legendOptions.labelFormat, item) :
                            legendOptions.labelFormatter.call(item)
                    });
                });


                chart.legend.render();
            }
        });
    });
    // Hide the tooltip but allow the crosshair
    H.Tooltip.prototype.defaultFormatter = function () { return false; };
}(Highcharts));