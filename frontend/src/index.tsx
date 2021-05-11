import React from "react";
import ReactDOM from "react-dom";
import "./index.less";
import reportWebVitals from "./reportWebVitals";
import { Hovedside } from "./hovedside";
import { BrowserRouter as Router, Switch, Link, Route } from "react-router-dom";
import ProduktDetalj from "./produktDetalj";
import ProduktNytt from "./produktNytt";
import { Sidetittel, Systemtittel } from "nav-frontend-typografi";

ReactDOM.render(
  <React.StrictMode>
    <div className={"dashboard-main"}>
      <div className="app">
        <Router>
          <header>
            <Link to="/">
              <object data="/nav-logo.svg" style={{ height: "5.625rem" }} />
              <Systemtittel>Dataprodukter</Systemtittel>
            </Link>
          </header>

          <main>
            <Switch>
              <Route path="/produkt/nytt" children={<ProduktNytt />} />
              <Route path="/produkt/:produktID" children={<ProduktDetalj />} />
              <Route exact path="/" children={<Hovedside />} />
              <Route path="*">
                <Systemtittel>404 - ikke funnet</Systemtittel>
              </Route>
            </Switch>
          </main>
        </Router>
      </div>
    </div>
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
