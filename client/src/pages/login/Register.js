import axios from 'axios';
import React, { Component } from 'react';
import '../../css/login.css'
import Login from './Login.js'
import {setToken} from "../../general/token.js";
class Register extends Component {
    constructor(props) {
        super(props);
        this.state = { messageErreur: "" }
        this.identity = ""
        this.password = ""
        this.confPassword = ""
        this.email = ""
    }
    gotoConnection() {
        this.props.setBody(<Login setBody={this.props.setBody}  />)
    }


    getErreur() {
        return <div className="breaker" style={{ color: "white", backgroundColor: "red", maxLength: "30px", fontWeight: "bold" }}><span>{this.state.messageErreur}</span></div>
    }

    signUp() {

        const regEmail = new RegExp("^[a-zA-Z0-9._:$!%-]+@[a-zA-Z0-9.-]+.[a-zA-Z]$")

        // eslint-disable-next-line eqeqeq
        if (this.password == "" || this.confPassword == "" || this.identity == ""|| this.email == "") {
            this.setState({ messageErreur: "un des champs est vide, veuillez remplir tout les champs" })
            return
        }

        if (!regEmail.test(this.email)) {
            this.setState({ messageErreur: "mail invalide" })
            return
        }
         // eslint-disable-next-line eqeqeq
        if (this.password != this.confPassword) {
            this.setState({ messageErreur: "le mot de passe de confirmation est different du mot de passe" })
            return
        }

        let identity = this.identity.slice(1)
        let password = this.password
        let email = this.email


        axios.post("http://localhost:8080/auth/signUp",
        {
            identity,
			password,
			email,

        }
        ).then((res) => {
            setToken(res.data)
            this.props.setBody(<Login setBody={this.props.setBody} />)
        }
        ).catch((err) => {
            this.setState({ messageErreur: err.message })
        })

    }


    render() {
        return <div id="login">

            <header>
                <div id="header_login">
                    <div id="logo" >
                    <img
                        src="https://img.icons8.com/color/512/clash-royal-red.png" 
                        alt="logo" 
                        width="128" 
                        height="128"
                    />
                    </div >
                    <div id="title"><p>Clash Connect</p></div>
                </div >
            </header>
            <section id="block_con">
                <h2 id="titre_con">Inscription</h2>
                {this.getErreur()}
                <div className="text">
                <input type="text" name="ID" placeholder="ID clash Royale: #RE25G" maxLength="30"
                        alt="l'ID commence par # et est suivi de letteres et chiffres" onChange={event => this.identity = event.target.value} />
                </div>

                <div className="text">
                    <input type="password" name="Mot de Passe" placeholder="Votre mot de passe" maxLength="30" onChange={event => this.password = event.target.value} />
                </div>
                <div className="text">
                    <input type="password" name="Confirmez le mot de passe" placeholder="Confirmez votre mot de passe"
                        maxLength="30" onChange={event => this.confPassword = event.target.value} />
                </div>
                <div className="text">
                    <input type="text" name="Email" placeholder="email : email@domain.tln/" maxLength="30" onChange={event => this.email = event.target.value} />
                </div>
                <br />
                <div className="h-captcha" data-sitekey="5ecff875-84d2-42d0-939c-32fe8d536fb0"></div>
                <div className="button">
                    <input type="submit" name="Enregistre" value="Envoyer" onClick={() => this.signUp()} />
                </div>
                <div style={{ display: 'flex', justifyContent: 'center' }}>
                    <div className="lien">
                        <span onClick={(event) => { this.gotoConnection(); event.stopPropagation() }}>Connexion</span>
                    </div>
                </div>
            </section >

        </div >
    }
}

export default Register