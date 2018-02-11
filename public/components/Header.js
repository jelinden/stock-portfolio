import React from 'react';
import { Redirect, Link } from 'react-router-dom';
import axios from 'axios';

class Header extends React.Component {

  constructor(props) {
    super(props);
    this.state = { user: null, signup: "", login: "", loggedin: false };
  }

  componentDidMount() {
    var _this = this;
    this.serverRequest = 
      axios
        .get("/api/user", {timeout: 2000})
        .then(function(result) {
          if (result.data.error) {
            _this.setState({
              signup: <li><Link className="headerLink" to='/signup'>Signup</Link></li>,
              login: <li><Link className="headerLink" to='/login'>Login</Link></li>
            });
          }
          _this.setState({
            user: result.data,
            loggedin: result.data.error?false:true
          });
        })
        .catch(function (error) {
          console.log(error);
          _this.setState({
            user: undefined,
            loggedin: false
          });
        });
  }

  render() {
    if (!this.state.user) {
      return null;
    }
    if (!this.state.loggedin && location.pathname.indexOf('login') === -1 && location.pathname.indexOf('signup') === -1) {
      return (
        <Redirect to={'/login'}/>
      )
    } else {
      return (
        <header>
          <nav>
            <ul>
              <li><Link className="headerLink" to='/'>Home</Link></li>
              {this.state.login}
              {this.state.signup}
              {this.state.loggedin?<a className="headerLink" href="/logout">Logout</a>:''}
            </ul>
            <div class="floatright">{this.state.loggedin?'Welcome, ' + this.state.user.username:''}</div>
          </nav>
        </header>
      )
    }
  }
}
export default Header;
