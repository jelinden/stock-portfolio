import React from "react";

class Verify extends React.Component {
    constructor() {
        super();
    }

    render() {
        return (
            <div className="pure-g">
                <div className="login">
                    <div className="register">
                        <h1>Verifying your email address</h1>
                        <p>Please verify your email by clicking the link in the email that was sent to you.</p>
                    </div>
                </div>
            </div>
        );
    }
}

export default Verify;
