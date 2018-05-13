import React from "react";
import { withRouter, Link } from "react-router-dom";
import axios from "axios";

import { Bar } from "react-chartjs";

function healthChartData(dataSet) {
    var labels = [];
    dataSet.forEach(function(elem, i) {
        labels.push("");
    });
    return {
        labels: labels,
        datasets: [
            {
                label: "Memory usage",
                fillColor: "#F7464A",
                strokeColor: "#f85e62",
                highlightFill: "#FF5A5E",
                highlightStroke: "#f85e62",
                data: dataSet
            }
        ]
    };
}

class Health extends React.Component {
    constructor() {
        super();
        this.state = {
            health: null,
            failed: false
        };
        this.getHealth = this.getHealth.bind(this);
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
                    health: healthChartData(result.data.MemUsedPercent)
                });
            })
            .catch(function(error) {
                _this.setState({
                    failed: true
                });
            });
    }

    render() {
        if (this.state.health === null) {
            return null;
        }
        const steps = 10;
        const max = 100;
        var options = {
            scaleOverride: true,
            scaleSteps: steps,
            scaleStepWidth: Math.ceil(max / steps),
            scaleStartValue: 0
        };
        return (
            <div id="health">
                <h1>Memory usage</h1>
                <Bar data={this.state.health} options={options} width="350" height="180" />
            </div>
        );
    }
}
export default Health;
