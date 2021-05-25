import React, { useContext, useEffect, useReducer, useState } from "react";
import ProduktTabell from "./produktTabell";
import { DataProduktListe, hentProdukter } from "./produktAPI";
import ProduktFilter from "./produktFilter";
import { Add } from "@navikt/ds-icons";
import { Link } from "react-router-dom";
import { UserContext } from "./userContext";

import NavFrontendSpinner from "nav-frontend-spinner";
import * as z from "zod";
import { Knapp } from "nav-frontend-knapper";
import "./hovedside.less";

export type ProduktListeState = {
  loading: boolean;
  error: string | null;
  products: DataProduktListe;
  filtered_products: DataProduktListe;
  filter: string | null;
};

const initialState: ProduktListeState = {
  loading: true,
  products: [],
  filtered_products: [],
  error: null,
  filter: "",
};

type ProduktListeFetch = {
  type: "FETCH_DONE";
  results: DataProduktListe;
};

type ProduktListeFilter = {
  type: "FILTER_CHANGE";
  filter: string;
};

const ProduktTabellReducer = (
  prevState: ProduktListeState,
  action: ProduktListeFetch | ProduktListeFilter
): ProduktListeState => {
  switch (action.type) {
    case "FETCH_DONE":
      return {
        ...prevState,
        products: action.results,
        filtered_products: action.results,
        loading: false,
        error: null,
      };
    case "FILTER_CHANGE":
      return {
        ...prevState,
        filtered_products: prevState.products.filter((p) => {
          if (action.filter === "") return true;
          return p.data_product?.team === action.filter;
        }),
      };
  }
  return prevState;
};

const ProduktNyKnapp = (): JSX.Element => (
  <div className={"nytt-produkt"}>
    <Link to="/produkt/nytt">
      <Knapp>
        <Add />
        Nytt produkt
      </Knapp>
    </Link>
  </div>
);

export const Hovedside = (): JSX.Element => {
  const user = useContext(UserContext);
  const [state, dispatch] = useReducer(ProduktTabellReducer, initialState);
  const [error, setError] = useState<string | null>();

  const loadProducts = () => {
    hentProdukter()
      .then((produkter) => {
        setError(null);

        dispatch({
          type: "FETCH_DONE",
          results: produkter,
        });
      })
      .catch((e) => {
        if (e instanceof z.ZodError) {
          console.log(JSON.stringify(e.errors, null, 2));
        }
        setError(`${e}`);
      });
  };

  useEffect(loadProducts, []);

  if (error) {
    setTimeout(loadProducts, 1500);
    console.log(error);
    // log:           {error}
    return (
      <div className={"feilBoks"}>
        <div>
          <h1>Kunne ikke hente produkter</h1>
        </div>
        <div>
          <h2>
            <NavFrontendSpinner /> Prøver på nytt...
          </h2>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="filter-and-button">
        <ProduktFilter state={state} dispatch={dispatch} />
        {user ? <ProduktNyKnapp /> : <></>}
      </div>
      <ProduktTabell state={state} dispatch={dispatch} />
    </div>
  );
};
