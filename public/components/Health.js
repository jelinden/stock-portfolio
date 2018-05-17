import React from "react";
import { withRouter, Link } from "react-router-dom";
import axios from "axios";

import { Line } from "react-chartjs";

function healthChartData(dataSet, title) {
    var labels = [];
    dataSet.forEach(function(elem, i) {
        dataSet[i] = Math.floor(elem * 100) / 100;
        labels.push("");
    });
    return {
        labels: labels,
        datasets: [
            {
                label: title,
                bezierCurve: false,
                fillColor: "rgba(0,0,0,0)",
                strokeColor: "#F7464A",
                pointDot: false,
                datasetFill: false,
                pointColor: "rgba(0,0,0,0)",
                pointStrokeColor: "rgba(0,0,0,0)",
                pointHighlightFill: "rgba(0,0,0,0)",
                pointHighlightStroke: "rgba(0,0,0,0)",

                //fillColor: "#F7464A",
                //strokeColor: "#f85e62",
                //highlightFill: "#FF5A5E",
                //highlightStroke: "#f85e62",
                data: dataSet
            }
        ]
    };
}

class Health extends React.Component {
    constructor() {
        super();
        this.state = {
            memchartData: null,
            memAllocData: null,
            cpuTotals: null,
            diskUsage: null,
            failed: false
        };
        this.getHealth = this.getHealth.bind(this);
        this.options = this.options.bind(this);
    }

    componentDidMount() {
        var _this = this;
        const MINUTE = 60000;
        _this.getHealth();
        setTimeout(function() {
            _this.getHealth();
        }, MINUTE);
        _this.interval = setInterval(_this.getHealth, MINUTE);
    }

    getHealth() {
        var _this = this;
        const TIMEOUT = 2000;
        axios
            .get("/api/health", {
                timeout: TIMEOUT
            })
            .then(function(result) {
                _this.setState({
                    memchartData: healthChartData(result.data.MemUsedPercent, "Memory usage"),
                    memAllocData: healthChartData(result.data.ProgramMemUsage, "Program memory usage"),
                    cpuTotals: healthChartData(result.data.CPUTotal, "CPU totals"),
                    diskUsage: healthChartData(result.data.DiskUsage, "Disk Usage")
                });
            })
            .catch(function(error) {
                _this.setState({
                    failed: true
                });
            });
    }

    options(max, steps) {
        var options = {
            scaleOverride: true,
            scaleSteps: steps,
            scaleStepWidth: Math.ceil(max / steps),
            scaleStartValue: 0
        };
        return options;
    }

    render() {
        if (this.state.memchartData === null) {
            return null;
        }
        return (
            <div>
                <div id="health">
                    <h1>System memory usage %</h1>
                    <Line data={this.state.memchartData} options={this.options(100, 10)} width="350" height="180" />
                </div>
                <div id="health">
                    <h1>CPU totals %</h1>
                    <Line data={this.state.cpuTotals} options={this.options(100, 10)} width="350" height="180" />
                </div>
                <div id="health">
                    <h1>Disk usage %</h1>
                    <Line data={this.state.diskUsage} options={this.options(100, 10)} width="350" height="180" />
                </div>
                <div id="health">
                    <h1>Program memory allocation MiB</h1>
                    <Line
                        data={this.state.memAllocData}
                        options={this.options(
                            Math.max.apply(
                                Math,
                                this.state.memAllocData.datasets[0].data.map(function(m) {
                                    return Math.floor(m) + 1;
                                })
                            ),
                            6
                        )}
                        width="350"
                        height="180"
                    />
                </div>
            </div>
        );
    }
}
export default Health;
