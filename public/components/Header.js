import React from "react";
import { Redirect, Link } from "react-router-dom";
import axios from "axios";

class Header extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            user: null,
            signup: "",
            login: "",
            loggedin: false,
            loaded: false
        };
    }

    componentDidMount() {
        var _this = this;
        this.serverRequest = axios
            .get("/api/user", {
                timeout: 2000
            })
            .then(function(result) {
                if (result.data.error) {
                    _this.setState({
                        signup: (
                            <li className="pure-menu-item">
                                <Link className="pure-menu-link" to="/signup">
                                    Signup
                                </Link>
                            </li>
                        ),
                        login: (
                            <li className="pure-menu-item">
                                <Link className="pure-menu-link" to="/login">
                                    Login
                                </Link>
                            </li>
                        )
                    });
                }
                _this.setState({
                    user: result.data,
                    loggedin: result.data.error ? false : true,
                    loaded: true
                });
            })
            .catch(function(error) {
                if (window.location.pathname !== "/login" && window.location.pathname !== "/signup") {
                    console.log("redirecting now");
                    location.href = "/login";
                }
                _this.setState({
                    loggedin: false,
                    loaded: true
                });
            });
    }

    render() {
        if (!this.state.loaded) {
            return null;
        } else if (
            !this.state.loggedin &&
            location.pathname.indexOf("login") === -1 &&
            location.pathname.indexOf("signup") === -1 &&
            location.pathname.indexOf("verify") === -1
        ) {
            return <Redirect to={"/login"} />;
        } else {
            return (
                <header className="header">
                    <div className="pure-menu pure-menu-horizontal pure-menu-scrollable">
                        <Link className="pure-menu-link pure-menu-heading" to="/">
                            Home
                        </Link>
                        <ul className="pure-menu-list">
                            {this.state.login} {this.state.signup}
                            {this.state.loggedin ? (
                                <li className="pure-menu-item">
                                    <a className="pure-menu-link" href="/logout">
                                        Logout
                                    </a>
                                </li>
                            ) : (
                                ""
                            )}
                            <li className="pure-menu-item">
                                <div className="floatright"> {this.state.loggedin ? "Welcome, " + this.state.user.username : "Your personal stock portfolio"} </div>
                            </li>
                        </ul>
                    </div>
                </header>
            );
        }
    }
}
export default Header;
