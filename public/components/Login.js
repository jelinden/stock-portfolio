import React from 'react'
import axios from 'axios'

class Login extends React.Component {

  constructor() {
    super();
    this.state = { user: null };
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

  render() {
    return (
        <div className="pure-g">
            <div className="login">
            <h1>Login</h1>
            <form id="login" method="post" action="/login" className="pure-form pure-form-stacked">
                <label htmlFor="emailLogin">Email</label>
                <input id="emailLogin" type="text" name="email" placeholder="my@email.com"/>
                <label htmlFor="passwordLogin">Password</label>
                <input id="passwordLogin" type="password" name="password" placeholder=""/>
                <button type="submit" className="pure-button pure-button-primary">Login</button>
            </form>
            </div>
        </div>
    )
  }
}

export default Login
