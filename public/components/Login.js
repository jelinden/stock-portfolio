import React from "react";
import axios from "axios";

class Login extends React.Component {
    constructor() {
        super();
        this.state = {
            user: null
        };
        this.getUrlParameter = this.getUrlParameter.bind(this);
    }

    componentDidMount() {
        var _this = this;
        this.serverRequest = axios
            .get("/api/user", {
                timeout: 2000
            })
            .then(function(result) {
                _this.setState({
                    user: result.data
                }, err => console.log(err));
            })
            .catch(err => {
                console.log("error", err);
            });
    }

    getUrlParameter(name) {
        console.log("getUrlParameter", name);
        name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
        let regex = new RegExp("[\\?&]" + name + "=([^&#]*)");
        let results = regex.exec(window.location.search);
        return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
    }

    render() {
        console.log("here");
        var verified = this.getUrlParameter("verified") === "true" ? "Verification succeeded, please login" : "";
        var login = this.getUrlParameter("login") === "failed" ? "Login failed" : "";
        return (
            <div className="pure-g">
                <div className="loginPage">
                    <div className="login">
                        <h1> Login </h1>
                        <form id="login" method="post" action="/login" className="pure-form pure-form-stacked">
                            <label htmlFor="emailLogin"> Email </label> <input id="emailLogin" type="text" name="email" placeholder="my@email.com" />
                            <label htmlFor="passwordLogin"> Password </label>
                            <input id="passwordLogin" type="password" name="password" placeholder="Over 8 characters" />
                            <button type="submit" className="pure-button pure-button-primary">
                                Login
                            </button>
                        </form>
                        <div>
                            New user? Please
                            <a className="a-link" href="/signup"> signup</a>.
                        </div>
                        <div className="alert"> {login} </div>
                        <div className="success"> {verified} </div>
                    </div>
                    <div className="exampleImg">
                        <img src="/img/portfolio-example.png" alt="Portfolio example view" />
                    </div>
                </div>
            </div>
        );
    }
}

export default Login;
