import { Redirect, withRouter, Link, useParams } from "react-router-dom";
import React, {useState, useEffect} from "react";
import axios from "axios";

function Transactions() {
    const {portfolioId} = useParams();
    const [portfolioName, setPortfolioName] = useState(null);
    const [transactions, setTransactions] = useState(null);
    const [failed, setFailed] = useState(null);
    

    const portfolio = () => {
        axios
            .get("/api/transactions/" + portfolioId, {
                timeout: 3000
            })
            .then(function(result) {
                setPortfolioName(result.data.portfolioName);
                setTransactions(result.data.stocks ? result.data.stocks : []);
            })
            .catch(function(error) {
                console.log(error);
                setFailed(true);
            });
    }

    const removeStock = (symbol, transactionId) => {
        if (!window.confirm("Are you sure you wish to delete this item?")) {
            return;
        }
        console.log("removing " + symbol);
        axios
            .get("/api/portfolio/remove/" + portfolioId + "/" + symbol + "/" + transactionId, {
                timeout: 3000
            })
            .then(function(result) {
                if (!result.data.error) {
                    console.log("removed " + symbol);
                    portfolio();
                } else {
                    console.log("remove was unsuccessful", result.data.error);
                }
            });
    }

    useEffect(() => {
        portfolio();
    });

    return (
        <div>
            <div className="headerLinks pure-g">
                <Link className="pure-menu-link pure-menu-heading" to={"/portfolio/" + portfolioId}>
                    Portfolio
                </Link>
            </div>
            <div className="content pure-g">
                <div className="alert info"> {failed ? "Connection lost" : ""}</div>
                <div id="portfolio" className="pure-u-1">
                    <h1> {portfolioName ? portfolioName : ""}</h1>
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
                                {transactions
                                    ? transactions.map((item, index) => (
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
                                                            removeStock(item.symbol, item.transactionId);
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
        </div>
    );
}
export default Transactions;
