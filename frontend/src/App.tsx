import React from "react";
import "./App.less";
import { Sidetittel } from "nav-frontend-typografi";
import ProduktTabell from "./produktTabell";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useParams,
} from "react-router-dom";

const App = (): JSX.Element => {
  return (
    <div className="app">
      <Router>
        <Sidetittel>Dataprodukter dashboard</Sidetittel>

        <ProduktTabell />
      </Router>
    </div>
  );
};

export default App;
