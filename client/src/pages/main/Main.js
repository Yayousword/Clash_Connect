import React, { Component } from "react";
import "../../css/index.css";
import Header from "./header/Header.js";
import Profil from "./profil/Profil.js";
import Chat from "./chat/ChatBox.js";
import DetailProfil from "./profil/DetailProfil";
import axios from "axios";
import { getLoginFromToken, getToken } from "../../general/token";

class Main extends Component {
  constructor(props) {
    super(props);
    this.setPage = this.setPage.bind(this);
    this.state = {
      page: null,

    };
  }

  componentDidMount() {
    axios
      .get(
        "http://localhost:8080/profile?token=" +
          getToken() +
          "&userID=" +
          getLoginFromToken()
      )
      .then((res) => {
        if (res.data == null) {
          console.log("le profile est null");
        }


        this.setPage(
          <DetailProfil
            setPage={this.setPage}
            setBody={this.props.setBody}
            profil={JSON.stringify(res.data)}
          />
        );
      })
      .catch((err) => alert(err));
    
  }
  setPage(cl) {
    this.setState({ page: cl });
  }

  getPage() {
    return this.state.page;
  }

  render() {
    return (
      <div id="mainPage" style={{ scrollBehavior: "smooth", overflowY: "scroll" }}>
        <Header
          setPage={this.setPage}
          setBody={this.props.setBody}
        />
        <div id="corps">
          <Profil
            setPage={this.setPage}
            setBody={this.props.setBody}
          />
          <Chat
            setPage={this.setPage}
            setBody={this.props.setBody}
          />
          {this.getPage()}
        </div>
      </div>
    );
  }
}

export default Main;
