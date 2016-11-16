'use strict';

vm.currentMenu('Organization');
vm.currentTitle("Organization");
vm.breadcrumb([{ title: 'Godrej', href: '#' }, { title: 'Organization', href: '/organization' }]);

viewModel.Or = new Object();
var or = viewModel.Or;

or.optionsdiagram = {
    options1: false,
    options2: false,
    options3: false,
    options4: false,
    options5: false,
    options6: false
};
or.newDiagram = function (idyo, id, container, src, options) {
    var result = void 0;

    result = ggOrgChart.render({
        data_id: idyo,
        container: id,
        max_text_width: 20,
        use_zoom_print: true,
        container_supra: container,
        initial_zoom: 0.75,
        box_color: "#C5C5C5",
        box_color_hover: "#FBFBFB",
        box_border_color: "#E6E6E6",
        title_color: "#A0A7A4",
        subtitle_color: "#8A8C8E"
    }, "/res/diagram/" + src);
    if (result === false) {
        alert("INFO: render() #4 failed (bad 'options' or 'data' definition)");return;
    } else {
        or.optionsdiagram[options] = result;
    }
};

or.zoomIn = function (options) {
    ggOrgChart.zoom_in(options);
};

or.zoomOut = function (options) {
    ggOrgChart.zoom_out(options);
};

$(function () {
    or.newDiagram(1, "diagram1", "or_container_1", "marketing.json", "options1");
    or.newDiagram(2, "diagram2", "or_container_2", "finance.json", "options2");
    or.newDiagram(3, "diagram3", "or_container_3", "pso.json", "options3");
    or.newDiagram(4, "diagram4", "or_container_4", "lbo.json", "options4");
    or.newDiagram(5, "diagram5", "or_container_5", "sales1.json", "options5");
    or.newDiagram(6, "diagram6", "or_container_6", "sales2.json", "options6");
});