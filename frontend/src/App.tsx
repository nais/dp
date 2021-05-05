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

const App = (): JSX.Element => {
  return (
    <div className="app">
      <Router>
        <ProduktListe />
      </Router>
    </div>
  );
};

export default App;
