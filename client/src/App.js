import './css/App.css';
import React, { Component } from "react";
import Main from "./pages/main/Main.js";
import Login from "./pages/login/Login.js"
import axios from "axios";

class App extends Component {
  constructor(props) {
    super(props);
    this.setBody = this.setBody.bind(this);
    // const host = location.hostname;
    //axios.defaults.baseURL = 'https://spacie.fr/';
    //axios.defaults.port = '443';
  
    this.state = {
      page: null,
    };
    //setToken({token : "test",login : "test"});
  }
  

  componentWillMount() {
    this.setBody(<Login setBody={this.setBody} />);
  }
  setBody(cl) {
    this.setState({ page: cl });
  }
  getBody() {
    return this.state.page;
  }
  render() {
    return this.getBody();
  }
  /*testing(){
    return(
    
        <div id="login" style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          height: "100%"
        }}>

            <header>
                <div id="header_login">
                    <div id="logo" >
                        <img src="https://i.pinimg.com/736x/26/23/6d/26236d6ea9b6b60de0946959e94b16ff.jpg" alt="logo" width="128" height="128" />
                    </div >
                    <div id="title"><p>CLASH CONNECT</p></div>
                </div >
            </header>
            <section id="block_con">
                <h2 id="titre_con">Inscription</h2>
                <div className="text">
                    <input type="text" name="ID" placeholder="ID clash Royale: #RE25G" maxLength="30"
                        alt="l'ID est commence par # et est suivi de letteres et chiffres" onChange={event => this.login = event.target.value} />
                </div>
                <div className="text">
                    <input type={this.state.dateRef} placeholder="Votre Token"  />
                </div>

                <div className="text">
                    <input type="email" name="Email" placeholder="email : email@domain.tln/" maxLength="30" onChange={event => this.email = event.target.value} />
                </div>
                <div className="text">
                    <input type="password" name="Mot de Passe" placeholder="Votre mot de passe" maxLength="30" onChange={event => this.motDePasse = event.target.value} />
                </div>
                <div className="text">
                    <input type="password" name="Confirmez le mot de passe" placeholder="Confirmez votre mot de passe"
                        maxLength="30" onChange={event => this.confPassword = event.target.value} />
                </div>
                <br />
                <div className="h-captcha" data-sitekey="5ecff875-84d2-42d0-939c-32fe8d536fb0"></div>
                <div className="button">
                    <input type="submit" name="Enregistre" value="Envoyer" />
                </div>
                <div style={{ display: 'flex', justifyContent: 'center' }}>
                    <div className="lien">
                        <span>Accueil</span>
                    </div>
                    <div className="lien">
                        <span>Connexion</span>
                    </div>
                </div>
            </section >

        </div >
        )

  }*/
}

export default App;
