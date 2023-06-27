import React from "react";
import { getToken } from "../../../general/token";
import axios from "axios";
import DetailProfil from "./DetailProfil";
const ProfilReduit = (props) => {
  return (
    <div >
      <span className="profil_reduit"
        onClick={() => {
          axios
            .get(
              "http://localhost:8080/profile?token=" +
                getToken() +
                "&userID=" +
                props.userId
            )
            .then((res) => {
              props.setPage(
                <DetailProfil
                  setPage={props.setPage}
                  setBody={props.setBody}
                  profil={JSON.stringify(res.data)}
                />
              );
            })
            .catch((err) => alert(err));
            const chatContainer = document.getElementById("mainPage");
            chatContainer.scrollTop = 0;
        }}
      >
        #{props.userId}
      </span>
    </div>
  );
};

export default ProfilReduit;
