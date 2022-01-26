import { Redirect, withRouter, Link } from "react-router-dom";
import React from "react";
import axios from "axios";
import DatePicker from "react-datepicker";
import moment from "moment";

import News from "./News";
import "react-datepicker/dist/react-datepicker.css";

class Portfolio extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            stocks: null,
            portfolioName: null,
            startDate: moment(),
            total: 0,
            currentTotal: 0,
            changeTotal: 0,
            itemTotal: 0,
            currentItemTotal: 0,
            closeTotal: 0,
            gain: "",
            failed: false,
            symbols: null,
            dividendData: null,
            stockMap: null,
            news: null
        };
        this.handleChange = this.handleChange.bind(this);
        this.convertTimestamp = this.convertTimestamp.bind(this);
        this.portfolio = this.portfolio.bind(this);
        this.getUrlParameter = this.getUrlParameter.bind(this);
        this.numberFormat = this.numberFormat.bind(this);
        this.dividends = this.dividends.bind(this);
    }

    handleChange(date) {
        this.setState({ startDate: date });
    }

    portfolio() {
        var _this = this;
        axios
            .get("/api/portfolio/get/" + _this.props.match.params.id, {
                timeout: 3000
            })
            .then(function(result) {
                let total = 0;
                let itemTotal = 0;
                let currentItemTotal = 0;
                let currentTotal = 0;
                let gainTotal = 0;
                let closeTotal = 0;
                let changeTotal = 0;
                let symbols = "";
                let stockMap = [];
                if (result.data.stocks) {
                    result.data.stocks.forEach(item => {
                        total += item.price;
                        if (item.latestPrice) {
                            currentTotal += item.latestPrice * item.amount;
                            itemTotal += item.change;
                            currentItemTotal += item.latestPrice;
                            changeTotal += item.change * item.amount;
                        }
                        closeTotal += item.close * item.amount;
                        if (symbols === "") {
                            symbols += item.symbol;
                        } else {
                            symbols += "," + item.symbol;
                        }
                        stockMap[item.symbol] = item.amount;
                    });
                }
                gainTotal = currentTotal - total;
                _this.setState({
                    stocks: result.data.stocks ? result.data.stocks : [],
                    stockMap: stockMap,
                    portfolioName: result.data.portfolioName,
                    total: total,
                    currentTotal: currentTotal,
                    changeTotal: changeTotal,
                    itemTotal: itemTotal,
                    currentItemTotal: currentItemTotal,
                    closeTotal: closeTotal,
                    gain: gainTotal,
                    failed: false,
                    symbols: symbols,
                    news: result.data.news
                });
            })
            .catch(function(error) {
                console.log(error);
                _this.setState({ failed: true });
            });
    }

    dividends(symbols) {
        var _this = this;
        if (symbols) {
            axios.get("/api/dividends?symbols=" + symbols, { timeout: 6000 }).then(function(result) {
                if (result.data) {
                    _this.setState({ dividendData: result.data });
                }
            });
        }
    }

    numberFormat(n) {
        if (n !== undefined) {
            var parts = n.toString().split(".");
            return parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, " ") + (parts[1] ? "." + parts[1] : "");
        }
        return "0.00";
    }

    componentDidMount() {
        var _this = this;
        _this.portfolio();
        setTimeout(function() {
            _this.dividends(_this.state.symbols);
        }, 1200);
        _this.interval = setInterval(_this.portfolio, 10000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    convertTimestamp(timestamp, full) {
        //console.log(timestamp, new Date(timestamp));
        if (timestamp !== undefined && timestamp !== "") {
            var d = new Date(timestamp), // Convert the passed timestamp to milliseconds
                yyyy = d.getFullYear(),
                mm = ("0" + (d.getMonth() + 1)).slice(-2), // Months are zero based. Add leading 0.
                dd = ("0" + d.getDate()).slice(-2), // Add leading 0.
                hh = d.getHours(),
                h = hh,
                min = ("0" + d.getMinutes()).slice(-2), // Add leading 0.
                time;

            if (full) {
                time = dd + "." + mm + "." + yyyy;
            } else {
                time = dd + "." + mm + ". " + h + ":" + min;
            }
            return time;
        }
        return "";
    }

    getUrlParameter(name) {
        name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
        let regex = new RegExp("[\\?&]" + name + "=([^&#]*)");
        let results = regex.exec(window.location.search);
        return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
    }

    render() {
        if (!this.state.stocks) {
            return null;
        }
        return (
            <div>
                <div className="headerLinks pure-g">
                    <Link className="pure-menu-link pure-menu-heading" to={"/transactions/" + this.props.match.params.id}>
                        Transactions
                    </Link>
                </div>
                <div className="content pure-g">
                    <div className="alert info">{this.state.failed ? "Connection lost" : ""}</div>
                    <div id="portfolio" className="pure-u-1">
                        <h1>{this.state.portfolioName ? this.state.portfolioName : ""}</h1>
                        <div className="scrolling-wrapper-flexbox">
                            <table className="pure-table">
                                <thead>
                                    <tr>
                                        <th>Name</th>
                                        <th>Symbol</th>
                                        <th>Last price</th>
                                        <th>Time</th>
                                        <th>Change</th>
                                        <th>Shares</th>
                                        <th>Cost basis</th>
                                        <th>Market Value</th>
                                        <th>Gain</th>
                                        <th>Gain %</th>
                                        <th>Day's gain</th>
                                        <th>Close price</th>
                                        <th>Close time</th>
                                        <th>P/e ratio</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {this.state.stocks
                                        ? this.state.stocks.map((item, index) => (
                                              <tr>
                                                  <td>{item.companyName}</td>
                                                  <td>{item.symbol}</td>
                                                  <td className="right">{item.latestPrice}</td>
                                                  <td className="right">{this.convertTimestamp(item.latestUpdate)}</td>
                                                  <td
                                                      className={
                                                          item.changePercent && item.changePercent > 0
                                                              ? "right green"
                                                              : item.changePercent < 0
                                                              ? "right red"
                                                              : "right"
                                                      }>
                                                      {item.changePercent
                                                          ? this.numberFormat(item.change.toFixed(2)) +
                                                            " (" +
                                                            this.numberFormat((item.changePercent * 100).toFixed(2)) +
                                                            "%)"
                                                          : "0.00"}
                                                  </td>
                                                  <td className="right">{item.amount}</td>
                                                  <td className="right">{this.numberFormat(item.price.toFixed(2))}</td>
                                                  <td className="right">
                                                      {item.latestPrice ? this.numberFormat((item.latestPrice * item.amount).toFixed(2)) : ""}
                                                  </td>
                                                  <td
                                                      className={
                                                          item.latestPrice && item.latestPrice * item.amount - item.price > 0
                                                              ? "right green"
                                                              : item.latestPrice * item.amount - item.price < 0
                                                              ? "right red"
                                                              : "right"
                                                      }>
                                                      {item.latestPrice ? (item.latestPrice * item.amount - item.price).toFixed(2) : ""}
                                                  </td>
                                                  <td
                                                      className={
                                                          item.latestPrice && item.latestPrice - item.price / item.amount > 0
                                                              ? "right green"
                                                              : item.latestPrice - item.price / item.amount < 0
                                                              ? "right red"
                                                              : "right"
                                                      }>
                                                      {item.latestPrice
                                                          ? (((item.latestPrice - item.price / item.amount) / (item.price / item.amount)) * 100).toFixed(2) +
                                                            "%"
                                                          : ""}
                                                  </td>
                                                  <td className={item.change && item.change > 0 ? "right green" : item.change < 0 ? "right red" : "right"}>
                                                      {item.change ? this.numberFormat((item.change * item.amount).toFixed(2)) : "0.00"}
                                                  </td>
                                                  <td className="right">{item.close ? this.numberFormat(item.close.toFixed(2)) : ""}</td>
                                                  <td className="right">{this.convertTimestamp(item.closeTime)}</td>
                                                  <td className="right">{item.peRatio ? item.peRatio.toFixed(2) : ""}</td>
                                              </tr>
                                          ))
                                        : ""}
                                </tbody>
                                <tfoot>
                                    <tr>
                                        <th />
                                        <th />
                                        <th />
                                        <th />
                                        <th
                                            className={
                                                this.state.itemTotal && this.state.itemTotal > 0
                                                    ? "right green"
                                                    : this.state.itemTotal < 0
                                                    ? "right red"
                                                    : "right"
                                            }>
                                            {this.state.itemTotal !== "undefined"
                                                ? this.numberFormat(((this.state.itemTotal / this.state.currentItemTotal) * 100).toFixed(2)) + "%"
                                                : ""}
                                        </th>
                                        <th />
                                        <th className="right">{this.state.stocks !== "undefined" ? this.numberFormat(this.state.total.toFixed(2)) : ""}</th>
                                        <th className="right">
                                            {this.state.stocks !== "undefined" ? this.numberFormat(this.state.currentTotal.toFixed(2)) : ""}
                                        </th>
                                        <th className={this.state.gain && this.state.gain > 0 ? "right green" : this.state.gain < 0 ? "right red" : "right"}>
                                            {this.state.gain ? this.numberFormat(this.state.gain.toFixed(2)) : ""}
                                        </th>
                                        <th className={this.state.gain && this.state.gain > 0 ? "right green" : this.state.gain < 0 ? "right red" : "right"}>
                                            {this.numberFormat(((this.state.gain / this.state.total) * 100).toFixed(2)) + "%"}
                                        </th>
                                        <th
                                            className={
                                                this.state.changeTotal && this.state.changeTotal > 0
                                                    ? "right green"
                                                    : this.state.changeTotal < 0
                                                    ? "right red"
                                                    : "right"
                                            }>
                                            {this.state.changeTotal ? this.numberFormat(this.state.changeTotal.toFixed(2)) : ""}
                                            {this.state.changeTotal
                                                ? " (" + this.numberFormat(((this.state.changeTotal / this.state.closeTotal) * 100).toFixed(2)) + "%)"
                                                : ""}
                                        </th>
                                        <th />
                                        <th />
                                        <th />
                                    </tr>
                                </tfoot>
                            </table>
                        </div>
                        <div className="footer">
                            Data provided for free by IEX. <a href="https://iextrading.com/api-exhibit-a">iextrading.com/api-exhibit-a</a>
                        </div>
                    </div>

                    <div id="addstock">
                        <form method="POST" action="/api/portfolio/add" className="pure-form pure-form-stacked">
                            <h2>Add stock</h2>
                            <input type="hidden" name="portfolioid" value={this.props.match.params.id} />
                            <label for="symbol">Stock symbol</label>
                            <input id="symbol" type="text" name="symbol" placeholder="Symbol eg. INTC" />
                            <div className="alert">{this.getUrlParameter("symbolMsg") ? "Not a correct Symbol" : ""}</div>
                            <div>
                                <label for="date">Date</label>
                                <DatePicker id="date" name="date" selected={this.state.startDate} onChange={this.handleChange} />
                                <div className="alert">{this.getUrlParameter("dateMsg") ? "Date was not a date" : ""}</div>
                            </div>
                            <input type="hidden" name="date" value={this.state.date} />
                            <label for="price">Price</label>
                            <input id="price" type="text" name="price" placeholder="Price" />
                            <div className="alert">{this.getUrlParameter("priceMsg") ? "Price was not a number" : ""}</div>
                            <label for="amount">Amount</label>
                            <input id="amount" type="text" name="amount" placeholder="How many bought?" />
                            <div className="alert">{this.getUrlParameter("amountMsg") ? "Amount was not a number" : ""}</div>
                            <label for="commission">Commission</label>
                            <input id="commission" type="text" name="commission" />
                            <div className="alert">{this.getUrlParameter("commissionMsg") ? "Commission was not a number" : ""}</div>
                            <button type="submit" className="pure-button pure-button-primary">
                                Add
                            </button>
                        </form>
                    </div>
                    {/* <News news={this.state.news} /> */}

                    <div id="dividends">
                        <h2>Latest dividend</h2>
                        <div>
                            <table class="pure-table">
                                <thead>
                                    <tr>
                                        <th>Symbol</th>
                                        <th>Ex date</th>
                                        <th>Payment date</th>
                                        <th>Amount</th>
                                        <th>Total</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {this.state.dividendData
                                        ? this.state.dividendData.map((item, index) => (
                                              <tr>
                                                  <td>{item.symbol}</td>
                                                  <td className="right">{this.convertTimestamp(item.exDate, true)}</td>
                                                  <td className="right">{this.convertTimestamp(item.paymentDate, true)}</td>
                                                  <td className="right">{item.amount.toFixed(4)}</td>
                                                  <td className="right">{(this.state.stockMap[item.symbol] * item.amount).toFixed(2)}</td>
                                              </tr>
                                          ))
                                        : ""}
                                </tbody>
                                <tfoot>
                                    <tr>
                                        <th />
                                        <th />
                                        <th />
                                        <th />
                                        <th />
                                    </tr>
                                </tfoot>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
export default Portfolio;
