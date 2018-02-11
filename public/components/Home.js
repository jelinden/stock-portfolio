import { withRouter, Link } from 'react-router-dom'
import React from 'react'
import axios from 'axios'

class Home extends React.Component {

  constructor() {
    super();
    this.state = { 
      portfolios: null,
      failed: false
    };
  }

  componentDidMount() {
    var _this = this;

    axios.get("/api/portfolios", {timeout: 2000})
        .then(function(result) {
          _this.setState({
            portfolios: result.data.portfolios
          });
        })
        .catch(function(error) {
          _this.setState({
            failed: true
          });
        });
  }

  render() {
      return (
        <div>
          <div id="addportfolio">
            <h1>Add new portfolio</h1>
            <form method="POST" action="/api/portfolio/create" className="pure-form pure-form-stacked">
              <input type="text" name="name" placeholder="Portfolio name" />
              <button type="submit" className="pure-button pure-button-primary">Create</button>
            </form>
          </div>
          <div id="portfolios">
            {this.state.failed?<a href="/login" className="link">Lost connection, go to Login</a>:''}
            {this.state.portfolios ? <h1>List of portfolios</h1>: ''}
            {this.state.portfolios ? this.state.portfolios.map((item, index) => (
              <Link className="link" to={'/portfolio/'+item.portfolioid} key={index}>{item.name}</Link>
            )):''}
          </div>
        </div>
      )
  }
}
export default Home;