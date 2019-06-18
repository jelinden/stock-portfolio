import React from "react";

class Signup extends React.Component {
    constructor() {
        super();
        this.getUrlParameter = this.getUrlParameter.bind(this);
    }

    getUrlParameter(name) {
        name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
        let regex = new RegExp("[\\?&]" + name + "=([^&#]*)");
        let results = regex.exec(window.location.search);
        return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
    }

    render() {
        return (
            <div class="pure-g">
                <div className="loginPage">
                    <div class="login">
                        <div class="register">
                            <h1>Signup</h1>
                            <form id="signup" method="post" action="/signup" className="pure-form pure-form-stacked">
                                <label for="username">Username</label>
                                <input id="username" type="text" name="username" placeholder="username" />
                                <label for="email">Email</label>
                                <input id="email" type="text" name="email" placeholder="my@email.com" />
                                <label for="password">Password</label>
                                <input id="password" type="password" name="password" placeholder="Over 8 characters" />
                                <button type="submit" className="pure-button pure-button-primary">
                                    Signup
                                </button>
                            </form>
                            <div className="alert">
                                {this.getUrlParameter("emailused") === "true" ? "Email address already used" : ""}
                                {this.getUrlParameter("validation") === "credentials" ? "Password should be over 8 characters long" : ""}
                                {this.getUrlParameter("validation") === "email" ? "Check the email" : ""}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

export default Signup;
