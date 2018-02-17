import React from 'react'
import axios from 'axios'

class Login extends React.Component {

  constructor() {
    super();
    this.state = { user: null };
    this.getUrlParameter = this.getUrlParameter.bind(this);
  }

  componentDidMount() {
    var _this = this;
    this.serverRequest = 
      axios
        .get("/api/user", {timeout: 2000})
        .then(function(result) {    
          _this.setState({
            user: result.data
          });
        });
  }

  getUrlParameter(name) {
    name = name.replace(/[\[]/, '\\[').replace(/[\]]/, '\\]');
    let regex = new RegExp('[\\?&]' + name + '=([^&#]*)');
    let results = regex.exec(window.location.search);
    return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '));
  }

  render() {
    return (
        <div className="pure-g">
            <div className="login">
              <h1>Login</h1>
              <form id="login" method="post" action="/login" className="pure-form pure-form-stacked">
                  <label htmlFor="emailLogin">Email</label>
                  <input id="emailLogin" type="text" name="email" placeholder="my@email.com"/>
                  <label htmlFor="passwordLogin">Password</label>
                  <input id="passwordLogin" type="password" name="password" placeholder="Over 8 characters"/>
                  <button type="submit" className="pure-button pure-button-primary">Login</button>
              </form>
              <div className="alert">{this.getUrlParameter('login') === 'failed'?'Login failed':''}</div>
              <div className="success">{this.getUrlParameter('verified') === 'true'?'Verification succeeded, please login':''}</div>
            </div>
        </div>
    )
  }
}

export default Login
