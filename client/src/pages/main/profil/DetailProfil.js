import React, { Component } from "react";

import { getToken } from "../../../general/token.js";
import ProfilReduit from "./ProfilReduit.js";

function affiche_victoire(nbTeam, nbOpp) {
  if (nbTeam > nbOpp) return "rgba(40, 116, 11, 0.658)";
  if (nbTeam < nbOpp) return "rgba(179, 19, 19, 0.658)";
}
function affiche_result(nbTeam, nbOpp) {
  if (nbTeam > nbOpp) return "Victoire";
  if (nbTeam < nbOpp) return "Defaite";
}

class DetailProfil extends Component {
  constructor(props) {
    super(props);
    this.setPage = this.props.setPage;
    this.token = getToken();
    this.state = {
      container: null,
      button1Name: "Publications",
      button2Name: "Détail Profil",
      user: this.props.user,
      messages: ["test1", "test2", "test3"],
    };
  }

  componentWillReceiveProps(props) {
    this.props = props;
  }

  affiche_image() {
    console.log(JSON.parse(this.props.profil));
    let parsedData = JSON.parse(this.props.profil);
    let images = [];

    for (let image in parsedData.playerInfo.cards) {
      let iconUrl = parsedData.playerInfo.cards[image].iconUrls.medium;
      images.push(
        <img src={iconUrl} alt={`Card ${image}`} className="imageDetail" />
      );
    }
    return images;
  }

  affiche_combats() {
    let parsedData = JSON.parse(this.props.profil);
    let combats = [];

    for (let combat in parsedData.battlesResponse) {
      let team = parsedData.battlesResponse[combat].team[0];
      let opponent = parsedData.battlesResponse[combat].opponent[0];
      combats.push([team, opponent]);
    }

    let affichageimagesTeam = [];
    let affichageimagesOpp = [];

    let affichagedivCombat = [];

    for (let combat in combats) {
      let team = combats[combat][0];
      let opponent = combats[combat][1];
      //this.affiche_detail_combat(team,opponent)
      for (let image in team.cards) {
        let iconUrl = team.cards[image].iconUrls.medium;
        affichageimagesTeam.push(
          <img
            src={iconUrl}
            alt={`Card ${image}`}
            className="imageDetailCombat"
          />
        );
      }
      for (let image in opponent.cards) {
        let iconUrl = opponent.cards[image].iconUrls.medium;
        affichageimagesOpp.push(
          <img
            src={iconUrl}
            alt={`Card ${image}`}
            className="imageDetailCombat"
          />
        );
      }
      affichagedivCombat.push(
        <article
          className="combat_cell"
          style={{
            backgroundColor: affiche_victoire(team.crowns, opponent.crowns),
          }}
        >
          <div
            style={{ textAlign: "center", fontSize: "100%", color: "yellow" }}
          >
            {affiche_result(team.crowns, opponent.crowns)}
          </div>
          <div style={{ textAlign: "left", fontSize: "100%" }}>
            {}
            <div style={{ marginBottom: "5px" }}>
              {team.name} : 
              <ProfilReduit setPage={this.props.setPage} setBody={this.props.setBody} userId={team.tag.slice(1)} />
            </div>
            {affichageimagesTeam}
          </div>{" "}
          <div style={{ textAlign: "center", color: "yellow" }}>VS</div>
          <div style={{ textAlign: "right", fontSize: "100%" }}>
            {affichageimagesOpp}
            <div style={{ marginBottom: "5px" }}>
              {opponent.name} :           
              <ProfilReduit setPage={this.props.setPage} setBody={this.props.setBody} userId={opponent.tag.slice(1)} />

            </div>
          </div>
        </article>
      );
      affichageimagesTeam = [];
      affichageimagesOpp = [];
    }
    return affichagedivCombat;

    //return images
  }

  render() {
    if (this.props.profil !== "") {
      let parsedData = JSON.parse(this.props.profil);
      return (
        <section id="messages">
          <article className="message">
            <div className="profil-stats">
              <div className="stat">
                <h3>ID : {parsedData.playerInfo.name}</h3>
                <h3>Nb de Victoires : {parsedData.playerInfo.wins}</h3>
              </div>
              <div className="stat">
                <h3>TAG : {parsedData.playerInfo.tag}</h3>
                <h3>Nb de Defaites : {parsedData.playerInfo.losses}</h3>
              </div>
            </div>
            <div style={{ textAlign: "center", color: "red" }}>
              <h1>
                WinRate :{" "}
                {Math.round(
                  (parsedData.playerInfo.wins /
                    (parsedData.playerInfo.wins +
                      parsedData.playerInfo.losses)) *
                    100
                )}{" "}
                %{" "}
              </h1>
            </div>
            <div className="profil-stats">
              <div className="stat">
                <h3>Trophés : {parsedData.playerInfo.trophies}</h3>
                <h3>Nb de dons : {parsedData.playerInfo.donations}</h3>
                <h3>Niveau : {parsedData.playerInfo.expLevel}</h3>
              </div>
              <div className="stat">
                <h3>Trophés max : {parsedData.playerInfo.bestTrophies}</h3>
                <h3>
                  Nb de dons recus : {parsedData.playerInfo.donationsReceived}
                </h3>
                <h3>Clan : {parsedData.playerInfo.clan.name}</h3>
              </div>
            </div>
          </article>
          <article className="message" style={{ textAlign: "center" }}>
            <h1>Cartes du joueur</h1>
            <h3>{this.affiche_image().map((img) => img)}</h3>
          </article>
          <article className="message" style={{ textAlign: "center" }}>
            <h1>Historique des combats</h1>
          </article>

          <h3>{this.affiche_combats()}</h3>
        </section>
      );
    }
  }
}

export default DetailProfil;
