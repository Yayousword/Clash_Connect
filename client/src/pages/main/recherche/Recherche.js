import React, { Component } from "react";
import { getToken } from "../../../general/token";
import axios from "axios";
import DetailProfil from "../profil/DetailProfil";
class Recherche extends Component {
  constructor(props) {
    super(props);
    this.recherche = "";
    this.state = {
      refraiche: null,
    };
  }
  rechercher() {
    axios
      .get(
        "http://localhost:8080/profile?token=" +
          getToken() +
          "&userID=" +
          this.recherche.slice(1)
      )
      .then((res) => {
        if (res.data == null) {
          console.log("le profile est null");
        }

        this.props.setPage(
          <DetailProfil
            setPage={this.props.setPage}
            setBody={this.props.setBody}
            profil={JSON.stringify(res.data)}
          />
        );
        const chatContainer = document.getElementById("mainPage");
        chatContainer.scrollTop = 0;
      })
      .catch((err) => alert(err));
  }
  getPageComponent() {
    return this.state.refraiche;
  }
  render() {
    return (
      <div style={{ width: "50%" }}>
        <div style={{ width: "100%" }}>
          <div
            style={{
              display: "flex",
              width: "100%",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <textarea
              onChange={(event) => (this.recherche = event.target.value)}
              placeholder="Rechercher: #User2265"
              maxLength="150"
              style={{
                borderRadius: "10px",
                width: "60%",
                textAlign: "center",
                justifyContent: "center",
              }}
            />
            <div
              className="buttons"
              style={{ width: "20%", textAlign: "center" }}
              onClick={() => this.rechercher()}
            >
              Rechercher
            </div>
          </div>
          {this.getPageComponent()}
        </div>
      </div>
    );
  }
}

export default Recherche;
