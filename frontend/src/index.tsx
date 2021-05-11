import React, { useState } from "react";
import ReactDOM from "react-dom";
import "./index.less";
import reportWebVitals from "./reportWebVitals";
import { Hovedside } from "./hovedside";
import { BrowserRouter as Router, Switch, Link, Route } from "react-router-dom";
import ProduktDetalj from "./produktDetalj";
import ProduktNytt from "./produktNytt";
import { Sidetittel, Systemtittel } from "nav-frontend-typografi";
import NaisPrideLogo from "./naisLogo";
import { Next } from "@navikt/ds-icons";

const next = (
  <svg
    width="1em"
    height="1em"
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    {" "}
    <path
      fill-rule="evenodd"
      clip-rule="evenodd"
      d="M17 12L8.429 3 7 4.5l7.143 7.5L7 19.5 8.429 21 17 12z"
      fill="currentColor"
    ></path>{" "}
  </svg>
);

const App = (): JSX.Element => {
  const [crumb, setCrumb] = useState<string | null>(null);
  return (
    <div className={"dashboard-main"}>
      <div className="app">
        <Router>
          <header>
            <Link to="/">
              <NaisPrideLogo />
              <Systemtittel>Dataprodukter</Systemtittel>
            </Link>
            {crumb != null ? (
              <Systemtittel>
                {next} {crumb}
              </Systemtittel>
            ) : null}
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
