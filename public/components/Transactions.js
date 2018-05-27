import { Redirect, withRouter, Link } from "react-router-dom";
import React from "react";
import axios from "axios";

class Transactions extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            portfolioName: null,
            transactions: null,
            failed: null
        };
    }

    portfolio() {
        var _this = this;
        axios
            .get("/api/transactions/" + _this.props.match.params.id, {
                timeout: 3000
            })
            .then(function(result) {
                _this.setState({
                    portfolioName: result.data.portfolioName,
                    transactions: result.data.stocks ? result.data.stocks : []
                });
            })
            .catch(function(error) {
                console.log(error);
                _this.setState({
                    failed: true
                });
            });
    }

    removeStock(symbol) {
        var _this = this;
        if (!window.confirm("Are you sure you wish to delete this item?")) {
            return;
        }
        console.log("removing " + symbol);
        axios
            .get("/api/portfolio/remove/" + _this.props.match.params.id + "/" + symbol, {
                timeout: 3000
            })
            .then(function(result) {
                if (!result.data.error) {
                    console.log("removed " + symbol);
                    _this.portfolio();
                } else {
                    console.log("remove was unsuccessful", result.data.error);
                }
            });
    }

    componentDidMount() {
        this.portfolio();
    }

    render() {
        if (!this.state.transactions) {
            return null;
        }
        return (
            <div className="content pure-g">
                <div className="headerLinks">
                    <Link className="pure-menu-link pure-menu-heading" to={"/portfolio/" + this.props.match.params.id}>
                        Portfolio
                    </Link>
                </div>
                <div className="alert info"> {this.state.failed ? "Connection lost" : ""}</div>
                <div id="portfolio" className="pure-u-1">
                    <h1> {this.state.portfolioName ? this.state.portfolioName : ""}</h1>
                    <div className="scrolling-wrapper-flexbox">
                        <table className="pure-table">
                            <thead>
                                <tr>
                                    <th>Symbol</th>
                                    <th>Company Name</th>
                                    <th>Buying Price</th>
                                    <th>Amount</th>
                                    <th>Commission</th>
                                    <th>Date</th>
                                    <th>Remove</th>
                                </tr>
                            </thead>
                            <tbody>
                                {this.state.transactions
                                    ? this.state.transactions.map((item, index) => (
                                          <tr>
                                              <td className="right">{item.symbol}</td>
                                              <td className="right">{item.companyName}</td>
                                              <td className="right">{item.price}</td>
                                              <td className="right">{item.amount}</td>
                                              <td className="right">{item.commission}</td>
                                              <td className="right">{item.date}</td>
                                              <td className="right">
                                                  <a
                                                      href="#"
                                                      onClick={() => {
                                                          this.removeStock(item.symbol);
                                                      }}>
                                                      (<span className="red delete" />)
                                                  </a>
                                              </td>
                                          </tr>
                                      ))
                                    : ""}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
}
export default Transactions;
