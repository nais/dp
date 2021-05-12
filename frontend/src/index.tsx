import React, { useState, useEffect, useContext } from "react";
import ReactDOM from "react-dom";
import "./index.less";
import reportWebVitals from "./reportWebVitals";
import { hentBrukerInfo, BrukerInfo } from "./produktAPI";
import { Hovedside } from "./hovedside";
import { BrowserRouter as Router, Switch, Link, Route } from "react-router-dom";
import ProduktDetalj from "./produktDetalj";
import ProduktNytt from "./produktNytt";
import { Systemtittel } from "nav-frontend-typografi";
import NaisPrideLogo from "./naisLogo";
import { Next } from "@navikt/ds-icons";
import { Hovedknapp } from "nav-frontend-knapper";
import { Child } from "@navikt/ds-icons";
import { UserContext } from "./userContext";

export const Bruker: React.FC<{ user: BrukerInfo }> = ({ user }) => {
  return (
    <div className={"brukerboks"}>
      <Child />
      {user.email}
    </div>
  );
};

const App = (): JSX.Element => {
  const [crumb, setCrumb] = useState<string | null>(null);
  const [user, setUser] = useState<BrukerInfo | null>(null);
  const userContext = useContext(UserContext);

  useEffect(() => {
    hentBrukerInfo()
      .then((bruker) => {
        setUser(bruker);
      })
      .catch(() => {
        setUser(null);
      });
  }, []);

  return (
    <div className={"dashboard-main"}>
      <div className="app">
        <UserContext.Provider value={user}>
          <Router>
            <header>
              <NaisPrideLogo />
              <Link to="/">
                <Systemtittel>Dataprodukter</Systemtittel>
              </Link>
              {crumb ? <Next className="pil" /> : null}
              {crumb ? <Systemtittel>{crumb}</Systemtittel> : null}
              {user ? (
                <Bruker user={user} />
              ) : (
                <Hovedknapp className="innloggingsknapp">logg inn</Hovedknapp>
              )}
            </header>

            <main>
              <Switch>
                <Route path="/produkt/nytt" children={<ProduktNytt />} />
                <Route
                  path="/produkt/:produktID"
                  children={<ProduktDetalj setCrumb={setCrumb} />}
                />
                <Route exact path="/" children={<Hovedside />} />
                <Route path="*">
                  <Systemtittel>404 - ikke funnet</Systemtittel>
                </Route>
              </Switch>
            </main>
          </Router>
        </UserContext.Provider>
      </div>
    </div>
  );
};

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
