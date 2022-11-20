import { Redirect, withRouter, Link, useParams } from "react-router-dom";
import React, {useState, useEffect} from "react";
import axios from "axios";
import DatePicker from "react-datepicker";
import moment from "moment";

import Dividend from "./Dividend";
import "react-datepicker/dist/react-datepicker.css";

function Portfolio() {

    const {portfolioId} = useParams();
    const [stocks, setStocks] = useState(null);
    const [portfolioName, setPortfolioName] = useState(null);
    const [total, setTotal] = useState(null);
    const [currentTotal, setCurrentTotal] = useState(null);
    const [changeTotal, setChangeTotal] = useState(null);
    const [itemTotal, setItemTotal] = useState(null);
    const [currentItemTotal, setCurrentItemTotal] = useState(null);
    const [closeTotal, setCloseTotal] = useState(null);
    const [gain, setGain] = useState(null);
    const [failed, setFailed] = useState(null);
    const [symbols, setSymbols] = useState(null);
    const [stockMap, setStockMap] = useState(null);
    const [startDate, setStartDate] = useState(new Date());

    const portfolio = () => {
        axios
            .get("/api/portfolio/get/" + portfolioId, {
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

                setStocks(result.data.stocks ? result.data.stocks : []);
                setStockMap(stockMap);
                setPortfolioName(result.data.portfolioName);
                setTotal(total);
                setCurrentTotal(currentTotal);
                setChangeTotal(changeTotal);
                setItemTotal(itemTotal);
                setCurrentItemTotal(currentItemTotal);
                setCloseTotal(closeTotal);
                setGain(gainTotal);
                setFailed(false);
                setSymbols(symbols);
            })
            .catch(function(error) {
                console.log(error);
                setFailed(true);
            });
    }

    const numberFormat = (n) => {
        if (n !== undefined) {
            var parts = n.toString().split(".");
            return parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, " ") + (parts[1] ? "." + parts[1] : "");
        }
        return "0.00";
    }

    useEffect(() => {
        portfolio();
        const interval = setInterval(portfolio, 20000);
        return function cleanup() {
            clearInterval(interval);
        };
    });

    const convertTimestamp = (timestamp, full) => {
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

    const getUrlParameter =(name) => {
        name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
        let regex = new RegExp("[\\?&]" + name + "=([^&#]*)");
        let results = regex.exec(window.location.search);
        return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
    }

    return (
        <div>
            <div className="headerLinks pure-g">
                <Link className="pure-menu-link pure-menu-heading" to={"/transactions/" + portfolioId}>
                    Transactions
                </Link>
            </div>
            <div className="content pure-g">
                <div className="alert info">{failed ? "Connection lost" : ""}</div>
                <div id="portfolio" className="pure-u-1">
                    <h1>{portfolioName ? portfolioName : ""}</h1>
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
                                    <th>P/e ratio</th>
                                </tr>
                            </thead>
                            <tbody>
                                {stocks
                                    ? stocks.map((item, index) => (
                                            <tr key={"tr" + index}>
                                                <td>{item.companyName}</td>
                                                <td>{item.symbol}</td>
                                                <td className="right">{item.latestPrice}</td>
                                                <td className="right">{convertTimestamp(item.latestUpdate)}</td>
                                                <td
                                                    className={
                                                        item.changePercent && item.changePercent > 0
                                                            ? "right green"
                                                            : item.changePercent < 0
                                                            ? "right red"
                                                            : "right"
                                                    }>
                                                    {item.changePercent
                                                        ? numberFormat(item.change.toFixed(2)) +
                                                        " (" +
                                                        numberFormat((item.changePercent * 100).toFixed(2)) +
                                                        "%)"
                                                        : "0.00"}
                                                </td>
                                                <td className="right">{item.amount}</td>
                                                <td className="right">{numberFormat(item.price.toFixed(2))}</td>
                                                <td className="right">
                                                    {item.latestPrice ? numberFormat((item.latestPrice * item.amount).toFixed(2)) : ""}
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
                                                    {item.change ? numberFormat((item.change * item.amount).toFixed(2)) : "0.00"}
                                                </td>
                                                <td className="right">{item.close ? numberFormat(item.close.toFixed(2)) : ""}</td>
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
                                            itemTotal && itemTotal > 0
                                                ? "right green"
                                                : itemTotal < 0
                                                ? "right red"
                                                : "right"
                                        }>
                                        {(itemTotal && itemTotal > 0)
                                            ? numberFormat(((itemTotal / currentItemTotal) * 100).toFixed(2)) + "%"
                                            : ""}
                                    </th>
                                    <th />
                                    <th className="right">{total ? numberFormat(total.toFixed(2)) : ""}</th>
                                    <th className="right">
                                        {currentTotal ? numberFormat(currentTotal.toFixed(2)) : ""}
                                    </th>
                                    <th className={gain && gain > 0 ? "right green" : gain < 0 ? "right red" : "right"}>
                                        {gain ? numberFormat(gain.toFixed(2)) : ""}
                                    </th>
                                    <th className={gain && gain > 0 ? "right green" : gain < 0 ? "right red" : "right"}>
                                        {numberFormat(((gain / total) * 100).toFixed(2)) + "%"}
                                    </th>
                                    <th
                                        className={
                                            changeTotal && changeTotal > 0
                                                ? "right green"
                                                : changeTotal < 0
                                                ? "right red"
                                                : "right"
                                        }>
                                        {changeTotal ? numberFormat(changeTotal.toFixed(2)) : ""}
                                        {changeTotal
                                            ? " (" + numberFormat(((changeTotal / closeTotal) * 100).toFixed(2)) + "%)"
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
                        <input type="hidden" name="portfolioid" value={portfolioId} />
                        <label htmlFor="symbol">Stock symbol</label>
                        <input id="symbol" type="text" name="symbol" placeholder="Symbol eg. INTC" />
                        <div className="alert">{getUrlParameter("symbolMsg") ? "Not a correct Symbol" : ""}</div>
                        <div>
                            <label htmlFor="date">Date</label>
                            <DatePicker id="date" name="date" selected={startDate} onChange={(date) => setStartDate(date)} />
                            <div className="alert">{getUrlParameter("dateMsg") ? "Date was not a date" : ""}</div>
                        </div>
                        <label htmlFor="price">Price</label>
                        <input id="price" type="text" name="price" placeholder="Price" />
                        <div className="alert">{getUrlParameter("priceMsg") ? "Price was not a number" : ""}</div>
                        <label htmlFor="amount">Amount</label>
                        <input id="amount" type="text" name="amount" placeholder="How many bought?" />
                        <div className="alert">{getUrlParameter("amountMsg") ? "Amount was not a number" : ""}</div>
                        <label htmlFor="commission">Commission</label>
                        <input id="commission" type="text" name="commission" />
                        <div className="alert">{getUrlParameter("commissionMsg") ? "Commission was not a number" : ""}</div>
                        <button type="submit" className="pure-button pure-button-primary">
                            Add
                        </button>
                    </form>
                </div>

                <Dividend 
                   stockMap={stockMap} 
                   symbols={symbols} />
                
            </div>
        </div>
    );
}
export default Portfolio;
