window.createChart = function (pointsRawData) {
    aapl = pointsRawData.split(";").map(function (point) {
        var xy = point.split(":").map(Number);
        return { date: new Date(xy[0]), value: xy[1] };
    });

    // Declare the chart dimensions and margins.
    const width = 928;
    const height = 500;
    const marginTop = 20;
    const marginRight = 30;
    const marginBottom = 30;
    const marginLeft = 40;

    // Declare the x (horizontal position) scale.
    const x = d3.scaleUtc(
        d3.extent(aapl, function (d) {
            return d.date;
        }),
        [marginLeft, width - marginRight]
    );

    // Declare the y (vertical position) scale.
    const y = d3.scaleLinear(
        [
            d3.min(aapl, function (d) {
                return d.value;
            }),
            d3.max(aapl, function (d) {
                return d.value;
            }),
        ],
        [height - marginBottom, marginTop]
    );

    // Declare the line generator.
    const line = d3
        .line()
        .x(function (d) {
            return x(d.date);
        })
        .y(function (d) {
            return y(d.value);
        });

    // Create the SVG container.
    const svg = d3
        .create("svg")
        .attr("width", width)
        .attr("height", height)
        .attr("viewBox", [0, 0, width, height])
        .attr("style", "max-width: 100%; height: auto; height: intrinsic;");

    // Add the x-axis.
    svg.append("g")
        .attr("transform", "translate(0," + height - marginBottom + ")")
        .call(
            d3
                .axisBottom(x)
                .ticks(width / 80)
                .tickSizeOuter(0)
        );

    // Add the y-axis, remove the domain line, add grid lines and a label.
    svg.append("g")
        .attr("transform", "translate(" + marginLeft + ",0)")
        .call(d3.axisLeft(y).ticks(height / 40))
        .call(function (g) {
            return g.select(".domain").remove();
        })
        .call(function (g) {
            return g
                .selectAll(".tick line")
                .clone()
                .attr("x2", width - marginLeft - marginRight)
                .attr("stroke-opacity", 0.1);
        })
        .call(function (g) {
            return g
                .append("text")
                .attr("x", -marginLeft)
                .attr("y", 10)
                .attr("fill", "currentColor")
                .attr("text-anchor", "start")
                .text("Watchcount");
        });

    // Append a path for the line.
    svg.append("path")
        .attr("fill", "none")
        .attr("stroke", "steelblue")
        .attr("stroke-width", 1.5)
        .attr("d", line(aapl));

    // Append the SVG element.
    document.getElementById("container").appendChild(svg.node());
};
