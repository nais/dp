import React from "react";
import "./App.less";
import { ProduktListe } from "./produktListe";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import ProduktDetalj from "./produktDetalj";
import { Sidetittel } from "nav-frontend-typografi";

const App = (): JSX.Element => {
  return (
    <div className="app">
      <Router>
        <Sidetittel>Dataprodukter</Sidetittel>
        <Switch>
          <Route path="/produkt/:produktID" children={<ProduktDetalj />} />
          <Route path="/" children={<ProduktListe />} />
        </Switch>
      </Router>
    </div>
  );
};

export default App;
