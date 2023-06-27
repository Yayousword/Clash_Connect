import React, { Component } from "react";

import {
  getToken,
  getLoginFromToken,
  testToken,
} from "../../../general/token.js";
import Recherche from "../recherche/Recherche.js";

class Header extends Component {
  constructor(props) {
    super(props);
    this.setPage = this.props.setPage;
    this.token = getToken();

    this.state = {
      login: "",
    };
  }
  componentWillMount() {
    if (testToken(this.token)) {
      this.setState({ login: getLoginFromToken() });
    }
  }
  render() {
    return (
      <header>
        <div id="logo">
          <div id="header_main">
            <div id="logo">
              <img
                src="https://img.icons8.com/color/512/clash-royal-red.png"
                alt="logo"
                width="128"
                height="128"
              />
            </div>
            <span id="title">Clash Connect</span>
          </div>
        </div>
        <Recherche  setPage={this.setPage} setBody={this.props.setBody}/>
        
      </header>
    );
  }
}

export default Header;
