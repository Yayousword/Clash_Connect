import React from 'react';
import ProfilReduit from '../profil/ProfilReduit';

const Message = (props) => {
  return (
    <li className="chat-message">
      <ProfilReduit setPage={props.setPage} setBody={props.setBody} userId={props.userId} />
      <span className="message-text"> : {props.message}</span>
    </li>
  );
};

export default Message;