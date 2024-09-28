function getChartName(metric) {
    switch (metric) {
        case "watchcount":
            return "Watch Count over Time";
        case "rating":
            return "Rating over Time";
        case "likecount":
            return "Like Count over Time";
    }
};

function getTickConfig(timeRange) {
    var hours = timeRange / (1000 * 60 * 60);

    if (hours <= 24) {
        return {
            interval: d3.timeHour.every(2),
            format: d3.timeFormat("%H:%M")
        };
    } else if (hours <= 24 * 7) {
        return {
            interval: d3.timeDay.every(1),
            format: d3.timeFormat("%d %b")
        };
    } else {
        return {
            interval: d3.timeDay.every(2),
            format: d3.timeFormat("%d %b")
        };
    }
};

function getDateRange(dateRange) {
    date = new Date();
    switch (dateRange) {
        case "d":
            return [date.setDate(date.getDate() - 1), new Date()];
        case "w":
            return [date.setDate(date.getDate() - 7), new Date()];
        case "m":
            return [date.setDate(date.getDate() - 30), new Date()];
    }
}

function getYRange(ymin, ymax, metric, height) {
    switch (metric) {
        case "watchcount":
            return d3.scaleLinear([ymin, ymax], [height, 0]);
        case "rating":
            return d3.scaleLinear([0, 5], [height, 0]);
        case "likecount":
            return d3.scaleLinear([ymin, ymax], [height, 0]);
    }
}

window.createChart = function (pointsRawData, metricName, dateRangeStr) {
    // Parse the raw data into an array of objects.
    aapl = pointsRawData.split(";").map(function (point) {
        var xy = point.split(":").map(Number);
        return { date: new Date(xy[0]), value: xy[1] };
    });

    // Declare the chart dimensions and margins.
    var margin = { top: 100, right: 10, bottom: 50, left: 100 },
        width = 1500 - margin.left - margin.right,
        height = 800 - margin.top - margin.bottom;

    // Append the SVG element to the container.
    const svg = d3
        .select("#container")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .attr("viewBox", [0, 0, width + margin.left + margin.right, height + margin.top + margin.bottom])
        .attr("style", "width: 100%; max-width: 100%; height:auto; max-height: 100%;")
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")")

    // Define area and line gradient.
    var lg = svg
        .append("defs")
        .append("linearGradient")
        .attr("id", "chartGradient")
        .attr("x1", "0%")
        .attr("x2", "0%")
        .attr("y1", "0%")
        .attr("y2", "100%"); //vertical linear gradient

    lg.append("stop")
        .attr("offset", "0%")
        .style("stop-color", "#FF8000")
        .style("stop-opacity", 1);

    lg.append("stop")
        .attr("offset", "100%")
        .style("stop-color", "#40BCF4")
        .style("stop-opacity", 1);

    dateRange = getDateRange(dateRangeStr)
    var x = d3.scaleUtc(dateRange, [0, width]);
    
    var tickConfig = getTickConfig(dateRange[1] - dateRange[0]);

    svg.append("g")
        .attr("class", "axisWhite")
        .attr("transform", "translate(0," + (height + 10) + ")")
        .style("font-size", "12px")
        .call(
            d3.axisBottom(x)
                .ticks(tickConfig.interval)
                .tickFormat(tickConfig.format)
                .tickSizeOuter(0)
        )
        .call(function(g) {
            return g.select(".domain").remove();
        });

    ymin = d3.min(aapl, function (d) {
        return d.value;
    });
    ymax = d3.max(aapl, function (d) {
        return d.value;
    });

    // Declare the y (vertical position) scale.
    const y = getYRange(ymin, ymax, metricName, height);

    // Add the y-axis, remove the domain line, add grid lines.
    svg.append("g")
        .attr("class", "axisWhite")
        .attr("transform", "translate(-5,0)")
        .style("font-size", "12px")
        .call(d3.axisLeft(y).ticks(15))
        .call(function (g) {
            return g.select(".domain").remove();
        })
        // Add grid lines.
        .call(function (g) {
            return g
                .selectAll(".tick line")
                .clone()
                .attr("x2", width)
                .attr("stroke-opacity", 0.1)
        });

    // Declare the area generator.
    const area = d3
        .area()
        .x(function (d) {
            return x(d.date);
        })
        .y0(height)
        .y1(function (d) {
            return y(d.value);
        });

    // Append a path for the area.
    svg.append("path")
        .attr("fill", "url(#chartGradient)")
        .attr("fill-opacity", 0.2)
        .attr("stroke", "none")
        .attr("d", area(aapl));

    // Declare the line generator.
    const line = d3
        .line()
        .x(function (d) {
            return x(d.date);
        })
        .y(function (d) {
            return y(d.value);
        });

    // Append a path for the line.
    svg.append("path")
        .attr("fill", "none")
        .attr("stroke", "url(#chartGradient)")
        .attr("stroke-width", "0.2vw")
        .attr("d", line(aapl));

    svg.append('text')
        .attr('class', 'chart-title')
        .attr('x', width / 2)
        .attr('y', - margin.top / 2)
        .attr('text-anchor', 'middle')
        .text(getChartName(metricName));
};
