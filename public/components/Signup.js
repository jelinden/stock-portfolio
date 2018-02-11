import React from 'react'
import axios from 'axios'

class Signup extends React.Component {

  constructor() {
    super();
  }

  render() {
      return (
        <div class="pure-g">
          <div class="login">
            <div class="register">
              <h1>Signup</h1>
              <form id="signup" method="post" action="/signup" className="pure-form pure-form-stacked">
                <label for="username">Username</label>
                <input id="username" type="text" name="username" placeholder="username"/>
                <label for="email">Email</label>
                <input id="email" type="text" name="email" placeholder="my@email.com"/>
                <label for="password">Password</label>
                <input id="password" type="password" name="password" placeholder="Over 8 characters"/>
                <button type="submit" className="pure-button pure-button-primary">Signup</button>
              </form>
            </div>
          </div>
        </div>
      )
  }
}

export default Signup
