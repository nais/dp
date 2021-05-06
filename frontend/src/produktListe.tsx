import React, { useEffect, useReducer } from "react";
import ProduktTabell from "./produktTabell";
import { DataProduktListe, hentProdukter } from "./produktAPI";
import ProduktFilter from "./produktFilter";

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

  useEffect(() => {
    hentProdukter().then((produkter) => {
      dispatch({
        type: "FETCH_DONE",
        results: produkter,
      });
    });
  }, []);
  return (
    <div>
      <ProduktFilter state={state} dispatch={dispatch} />
      <ProduktTabell state={state} dispatch={dispatch} />
    </div>
  );
};
