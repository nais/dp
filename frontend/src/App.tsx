import React from "react";
import "./App.less";
import { ProduktListe } from "./produktListe";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useParams,
} from "react-router-dom";
import ProduktDetalj from "./produktDetalj";

const App = (): JSX.Element => {
  return (
    <div className="app">
      <Router>
        <Switch>
          <Route path="/produkt/:produktID" children={<ProduktDetalj />} />
          <Route path="/" children={<ProduktListe />} />
        </Switch>
      </Router>
    </div>
  );
};

export default App;
