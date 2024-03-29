import React from "react";
import axios from "axios";

class Dividend extends React.Component {
    constructor(props) {
        super(props);
        this.convertTimestamp = this.convertTimestamp.bind(this);
        this.dividends = this.dividends.bind(this);
        this.state = {
            dividendData: null
        }
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

    dividends(symbols) {
        var _this = this;
        if (symbols) {
            axios.get("/api/dividends?symbols=" + symbols, { timeout: 6000 }).then(function(result) {
                if (result.data) {
                    
                    var month = -1;
                    var monthlyTotal = 0;
                    var data = result.data;
                    data.forEach((element, index) => {
                        const paymentMonth = (new Date(element.paymentDate)).getMonth();

                        if (element.currency === 'CAD' && element.currencyRate != null && element.currencyRate !== 0) {
                            element.amount = element.currencyRate*element.amount;
                        }
                        
                        if (month === -1) { month = paymentMonth; }

                        if (month !== paymentMonth) { // month changed
                            month = paymentMonth;
                            data[index-1].total = monthlyTotal;
                            monthlyTotal = _this.props.stockMap[element.symbol] * element.amount;
                        } else if (month === paymentMonth) { // same month as last one
                            monthlyTotal += _this.props.stockMap[element.symbol] * element.amount;
                        }
                        if (data.length === index+1) { // the last line
                            monthlyTotal += _this.props.stockMap[element.symbol] * element.amount;
                            data[index].total = monthlyTotal;    
                        }
                    });
                    _this.setState({ dividendData: data });
                }
            });
        }
    }



    componentDidMount() {
        var _this = this;
        setTimeout(function() {
            _this.dividends(_this.props.symbols);
        }, 1200);
    }

    render() {
        if (!this.state.dividendData) {
            return null;
        }
        return (
        <div id="dividends">
        <h2>Latest dividend</h2>
        <div>
            <table className="pure-table">
                <thead>
                    <tr>
                        <th>Symbol</th>
                        <th>Ex date</th>
                        <th>Payment date</th>
                        <th>Amount</th>
                        <th>Total</th>
                        <th>Monthly total</th>
                    </tr> 
                </thead>
                <tbody>
                    {this.state.dividendData
                        ? this.state.dividendData.map((item, index) => (
                              <tr key={"tr"+index}>
                                  <td>{item.symbol}</td>
                                  <td className="right">{this.convertTimestamp(item.exDate, true)}</td>
                                  <td className="right">{this.convertTimestamp(item.paymentDate, true)}</td>
                                  <td className="right">{item.amount.toFixed(4)}</td>
                                  <td className="right">{(this.props.stockMap[item.symbol] * item.amount).toFixed(2)}</td>
                                  <td className="right">{item.total?item.total.toFixed(2):''}</td>
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
                        <th />
                    </tr>
                </tfoot>
            </table>
        </div>
    </div>
    );
    }
}

export default Dividend;