import React, { useContext } from "react";
import { BACKEND_ENDPOINT } from "./produktAPI";
import { Hovedknapp } from "nav-frontend-knapper";
import { Child, Next } from "@navikt/ds-icons";
import { UserContext } from "./userContext";
import NaisPrideLogo from "./naisLogo";
import { Link } from "react-router-dom";
import { Systemtittel } from "nav-frontend-typografi";

const BrukerBoks: React.FC = () => {
  const user = useContext(UserContext);

  if (!user) {
    return (
      <a className="innloggingsknapp" href={`${BACKEND_ENDPOINT}/login`}>
        <Hovedknapp className="innloggingsknapp">logg inn</Hovedknapp>
      </a>
    );
  } else {
    return (
      <div className={"brukerboks"}>
        <Child />
        {user.email}
      </div>
    );
  }
};

export const PageHeader: React.FC<{ crumbs: string | null }> = ({ crumbs }) => {
  return (
    <header>
      <NaisPrideLogo />
      <Link to="/">
        <Systemtittel>Dataprodukter</Systemtittel>
      </Link>
      {crumbs ? <Next className="pil" /> : null}
      {crumbs ? <Systemtittel>{crumbs}</Systemtittel> : null}
      <BrukerBoks />
    </header>
  );
};

export default PageHeader;
