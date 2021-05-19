import React, { useState, useEffect } from "react";
import ReactDOM from "react-dom";
import "./index.less";
import reportWebVitals from "./reportWebVitals";
import { hentBrukerInfo, BrukerInfo } from "./produktAPI";
import { Hovedside } from "./hovedside";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import ProduktDetalj from "./produktDetalj";
import ProduktNytt from "./produktNytt";
import { Systemtittel } from "nav-frontend-typografi";

import { UserContext } from "./userContext";
import PageHeader from "./pageHeader";

const App = (): JSX.Element => {
  const [crumb, setCrumb] = useState<string | null>(null);
  const [user, setUser] = useState<BrukerInfo | null>(null);

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
            <PageHeader crumbs={crumb} />
            <main>
              <Switch>
                <Route
                  path="/produkt/nytt"
                  children={() => {
                    setCrumb("Nytt produkt");
                    return <ProduktNytt />;
                  }}
                />
                <Route
                  path="/produkt/:produktID"
                  children={<ProduktDetalj setCrumb={setCrumb} />}
                />
                <Route
                  exact
                  path="/"
                  children={() => {
                    setCrumb(null);
                    return <Hovedside />;
                  }}
                />
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
