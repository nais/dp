import React from "react";
import "./App.less";
import { Hovedside } from "./hovedside";
import { BrowserRouter as Router, Switch, Link, Route } from "react-router-dom";
import ProduktDetalj from "./produktDetalj";
import ProduktNytt from "./produktNytt";
import { Sidetittel, Systemtittel } from "nav-frontend-typografi";

const App = (): JSX.Element => {
  return (
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
  );
};

export default App;
