var args = require("system").args,
    page = require("webpage").create(),
    url = args[1],
    pointData = args[2],
    metric = args[3],
    dateRange = args[4];

page.onConsoleMessage = function (msg) {
    console.log(msg);
};

page.open(url, function (status) {
    if (status === "success") {
        page.injectJs("d3.v5.min.js")
        page.injectJs("createGraph.js")

        page.evaluate(function (data, metricName, dateRange) {
            createChart(data, metricName, dateRange);
            console.log(document.getElementById("container").innerHTML);
        }, pointData, metric, dateRange);

        phantom.exit();
    }
});
