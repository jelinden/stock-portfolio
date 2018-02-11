import { Redirect, withRouter } from 'react-router-dom'
import React from 'react'
import axios from 'axios'
import DatePicker from 'react-datepicker';
import moment from 'moment';

import 'react-datepicker/dist/react-datepicker.css';

class Portfolio extends React.Component {

  constructor(props) {
    super(props);
    this.state = { 
      stocks: null,
      startDate: moment(),
      total: 0,
      currentTotal: 0
    };
    this.handleChange = this.handleChange.bind(this);
    this.convertTimestamp = this.convertTimestamp.bind(this);
    this.portfolio = this.portfolio.bind(this);
    this.removeStock = this.removeStock.bind(this);
    this.numberFormat = this.numberFormat.bind(this);
  }

  handleChange(date) {
    this.setState({startDate: date});
  }

  portfolio() {
    var _this = this;
    axios.get("/api/portfolio/get/" + _this.props.match.params.id, {timeout: 2000})
        .then(function(result) {
          let total = 0;
          let currentTotal = 0;
          result.data.stocks.forEach(item => {
            total += item.price;
            if (item.latestPrice) {
              currentTotal += item.latestPrice*item.amount;
            }
          });
          _this.setState({
            stocks: result.data.stocks,
            total: total,
            currentTotal: currentTotal
          });
        });
  }

  removeStock(symbol) {
    var _this = this;
    return function (e) {
      console.log('removing ' + symbol);
      axios.get("/api/portfolio/remove/"+ _this.props.match.params.id + "/" + symbol, {timeout: 2000})
      .then(function(result) {
        console.log('removed ' + symbol);
      });
    };
  }

  numberFormat(n) {
    var parts = n.toString().split(".");
    return parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, " ") + (parts[1] ? "." + parts[1] : "");
  }

  componentDidMount() {
    var _this = this;
    _this.portfolio();
    _this.interval = setInterval(_this.portfolio, 10000);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  convertTimestamp(timestamp) {
    if (timestamp !== undefined && timestamp !== '') {
      var d = new Date(timestamp),	// Convert the passed timestamp to milliseconds
        yyyy = d.getFullYear(),
        mm = ('0' + (d.getMonth() + 1)).slice(-2),	// Months are zero based. Add leading 0.
        dd = ('0' + d.getDate()).slice(-2),			// Add leading 0.
        hh = d.getHours(),
        h = hh,
        min = ('0' + d.getMinutes()).slice(-2),		// Add leading 0.
        time;

      time = dd + '.' + mm + '.' + yyyy + ' ' + h + ':' + min;
      return time;
    }
    return '';
  }

  render() {
      return (
        <div class="content">
          <div id="portfolio">
            <table class="pure-table">
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
                  <th>Gain %</th>
                  <th>Close price</th>
                  <th>Close time</th>
                  <th>P/e ratio</th>
                  <th>Remove</th>
              </tr>
            </thead>
            <tbody>
                {this.state.stocks ? this.state.stocks.map((item, index) => (
                    <tr>
                      <td>{item.companyName}</td>
                      <td>{item.symbol}</td>
                      <td class="right">{item.latestPrice}</td>
                      <td class="right">{this.convertTimestamp(item.latestUpdate)}</td>
                      <td class="right">{item.latestPrice?this.numberFormat((((item.latestPrice-item.close)/item.close)*100).toFixed(2)):''}%</td>
                      <td class="right">{item.amount}</td>
                      <td class="right">{this.numberFormat((item.price).toFixed(2))}</td>
                      <td class="right">{(item.latestPrice?this.numberFormat((item.latestPrice*item.amount).toFixed(2)):'')}</td>
                      <td class="right">{item.latestPrice?(((item.latestPrice-(item.price/item.amount))/(item.price/item.amount))*100).toFixed(2):''}%</td>
                      <td class="right">{this.numberFormat(item.close)}</td>
                      <td class="right">{this.convertTimestamp(item.closeTime)}</td>
                      <td class="right">{item.peRatio?(item.peRatio).toFixed(2):''}</td>
                      <td class="right"><a href="#" onClick={this.removeStock(item.symbol)}>(<span className="red delete"></span>)</a></td>
                    </tr>
                )):''}
              </tbody>
              <tfoot>
                <tr>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th class="right">{(this.state.stocks !== 'undefined'?this.state.total:'').toFixed(2)}</th>
                  <th class="right">{(this.state.stocks !== 'undefined'?this.state.currentTotal:'').toFixed(2)}</th>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th></th>
                </tr>
              </tfoot>
            </table>
          </div>
          <div id="addstock">
            <form method="POST" action="/api/portfolio/add" className="pure-form pure-form-stacked">
                <input type="hidden" name="portfolioid" value={this.props.match.params.id}/>
                <label for="symbol">Stock symbol</label>
                <input id="symbol" type="text" name="symbol" placeholder="Symbol eg. INTC" />
                <div>
                <label for="date">Date</label>
                <DatePicker id="date"
                  name="date"
                  selected={this.state.startDate}
                  onChange={this.handleChange}
                />
                </div>
                <input type="hidden" name="date" value={this.state.date}/>
                <label for="price">Price</label>
                <input id="price" type="text" name="price" placeholder="Price" />
                <label for="amount">Amount</label>
                <input id="amount" type="text" name="amount" placeholder="How many bought?" />
                <label for="commission">Commission</label>
                <input id="commission" type="text" name="commission" />
                <button type="submit" className="pure-button pure-button-primary">Add</button>
            </form>
          </div>
          <div class="footer">Data provided for free by IEX. <a href="https://iextrading.com/api-exhibit-a">iextrading.com/api-exhibit-a</a></div>
        </div>
      )
  }
}
export default Portfolio;