import React, { Component } from "react";
import Main from "../main/Main.js";
import Register from "./Register.js";
import "../../css/login.css";
import axios from "axios";
import {getToken, setToken} from "../../general/token.js";

class Login extends Component {
  constructor(props) {
    super(props);
    this.login = "";
    this.password = "";
    this.tokenPerso = "";
    this.state = { messageErreur: "" };
  }

  connecte() {
    // eslint-disable-next-line eqeqeq
    if (this.login == "" || this.password == "") {
      this.setState({ messageErreur: "veuillez remplir tous les champs" });
      return;
    }
  console.log(this.login.slice(1))
   axios
      .post("http://localhost:8080/auth/signIn", { Identity: this.login.slice(1), Password: this.password, TokenAPI: this.tokenPerso})
      .then((res) => {
        setToken({token : res.data.token, login : this.login.slice(1)});
        console.log(getToken());
        console.log(res.data);
        this.props.setBody(<Main setBody={this.props.setBody}/>);
      })
      .catch((err) => {
        this.setState({ messageErreur: err.message });
      });
  }

  gotoSignIn() {
    this.props.setBody(
      <Register serveur={this.props.serveur} setBody={this.props.setBody} />
    );
  }


  getErreur() {
    return (
      <div
        className="breaker"
        style={{ color: "white", backgroundColor: "red" }}
      >
        <span>{this.state.messageErreur}</span>
      </div>
    );
  }
  render() {
    return (
      
      <div id="login">
        <header>
          <div id="header_login">
            <div id="logo">
              <img
                src="https://img.icons8.com/color/512/clash-royal-red.png" 
                alt="logo" 
                width="128" 
                height="128"
              />
            </div>
            <div id="title">
              <p>Clash Connect</p>
            </div>
          </div>
        </header>
        <section id="block_con">
          <h2 id="titre_con">Connexion</h2>

          {this.getErreur()}

          <div className="text">
            <input
              type="text"
              name="Login"
              placeholder="ID clash Royale: #RE25G"
              maxLength="30"
              onChange={(event) => (this.login = event.target.value)}
            />
          </div>
          <div className="text">
            <input
              type="password"
              name="password"
              placeholder="mot de passe"
              maxLength="30"
              onChange={(event) => (this.password = event.target.value)}
            />
          </div>
          <div className="text">
            <input
              type="TokenCR"
              name="TokenCR"
              placeholder="token API ClashRoyal: https://developer.clashroyale.com/#/"
              onChange={(event) => (this.tokenPerso = event.target.value)}
            />
          </div>
          <div className="button">
            <input
              type="button"
              onClick={() => this.connecte()}
              name="Connexion"
              value="Connexion"
            />
          </div>
          <div style={{ display: "flex", justifyContent: "center" }}>
            <div className="lien">
              <span
                href=""
                onClick={(event) => {
                  this.gotoSignIn();
                  event.stopPropagation();
                }}
              >
                Inscription
              </span>
            </div>
          </div>
        </section>
      </div>
    );
  }
}

export default Login;
