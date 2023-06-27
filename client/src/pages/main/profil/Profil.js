import React, { Component } from "react";
import DetailProfil from "./DetailProfil";
import LoginPage from "../../login/Login";
import axios from "axios";



import { getToken, testToken, getLoginFromToken,setToken } from "../../../general/token.js";

class Profil extends Component {
  constructor(props) {
    super(props);
    this.setPage = this.props.setPage;
    this.token = getToken();
    this.state = {
      container: null,
      button1Name: "Publications",
      button2Name: "Détail Profil",
      afficherDetail: false,
      afficherPublications: false,
      userConnect: null,
      user: this.props.user,
      messages: [],
      profil: "",
    };
    this.refresh = this.refresh.bind(this);
  }

  componentWillReceiveProps(props) {
    this.props = props;
  }

  componentDidMount() {
    if (testToken(this.token)) {
      this.setState({ userConnect: getLoginFromToken() });
    }
    this.setContainer();
  }

  refresh() {
    this.setContainer();
  }

  disconnect() {
    setToken({token : "" , login: "" });
    this.props.setBody(<LoginPage setBody={this.props.setBody} />);
  }
  getnom() {}

  setContainer() {


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
        console.log(res.data);
        this.setState({
          profil: JSON.stringify(res.data),
        });
      })
      .catch((err) => alert(err));
  }

  render() {
    // eslint-disable-next-line eqeqeq
    if (this.state.profil != "") {
      let parsedData = JSON.parse(this.state.profil);
      return (
        <div className="millieu">
          <section id="info_user">
            <span className="photoProfil">
              <img
                id="pdp"
                src="https://img.icons8.com/color/512/clash-royal-red.png"
                alt="Profil"
                style={{ maxWidth: "60px", maxHeight: "80px" }}
              />
            </span>
            <div id="info_profil">
              <div className="info_ligne">
                <div id="loginProfil" className="info">
                  <h3>{parsedData.playerInfo.name}</h3>
                  <p className="breaker">{}</p>
                </div>
              </div>
              <div id="stats_mini_profil">
                <table>
                  <tr>
                    <th>Trophés</th>
                    <th>{parsedData.playerInfo.trophies}</th>
                  </tr>
                  <tr>
                    <td class="td">Niveau</td>
                    <td>{parsedData.playerInfo.expLevel}</td>
                  </tr>
                  <tr>
                    <td class="td">Victoires</td>
                    <td>{parsedData.playerInfo.wins}</td>
                  </tr>
                  <tr>
                    <td class="td">Defaites</td>
                    <td>{parsedData.playerInfo.losses}</td>
                  </tr>
                </table>
              </div>
              <div id="button_profil">
                <div
                  className="buttons"
                  onClick={() => {
                    this.setPage(
                      <DetailProfil
                        setPage={this.props.setPage}
                        setBody={this.props.setBody}
                        profil={this.state.profil}
                      />
                    );
                  }}
                >
                  {this.state.button2Name}
                </div>

                <div className="buttons" onClick={() => this.disconnect()}>
                  Déconnection
                </div>
              </div>
            </div>
          </section>
          {this.state.container}
        </div>
      );
    }
  }
}

export default Profil;
