import React, { useEffect, useReducer, useState } from "react";
import ProduktTabell from "./produktTabell";
import { DataProduktListe, hentProdukter } from "./produktAPI";
import ProduktFilter from "./produktFilter";
import NavFrontendSpinner from "nav-frontend-spinner";

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
          return p.data_product?.owner === action.filter;
        }),
      };
  }
  return prevState;
};

export const ProduktListe = (): JSX.Element => {
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
        setError(`${e}`);
      });
  };

  useEffect(loadProducts, []);

  if (error) {
    setTimeout(loadProducts, 1500);
    console.log(error);
    // log:           {error}
    return (
      <div>
        <h1>Kunne ikke hente produkter</h1>
        <h2>
          <NavFrontendSpinner /> Prøver på nytt...
        </h2>
      </div>
    );
  }

  return (
    <div>
      <ProduktFilter state={state} dispatch={dispatch} />
      <ProduktTabell state={state} dispatch={dispatch} />
    </div>
  );
};
