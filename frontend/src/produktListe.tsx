import React, { useEffect, useReducer } from "react";
import ProduktTabell from "./produktTabell";
import { Sidetittel } from "nav-frontend-typografi";
import { DataProduktListe, hentProdukter } from "./produktAPI";
import ProduktFilter from "./produktFilter";

export type ProduktListeState = {
  loading: boolean;
  error: string | null;
  products: DataProduktListe;
  filter: string | null;
};

const initialState: ProduktListeState = {
  loading: true,
  products: [],
  error: null,
  filter: null,
};

interface ProduktListeAction {
  type: string;
  results: DataProduktListe;
}

const ProduktTabellReducer = (
  prevState: ProduktListeState,
  action: ProduktListeAction
): ProduktListeState => {
  switch (action.type) {
    case "FETCH_DONE":
      return {
        ...prevState,
        products: action.results,
        loading: false,
        error: null,
      };
    case "FILTER_CHANGE":
      return {
        ...prevState,
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
      <Sidetittel>Dataprodukter dashboard</Sidetittel>
      <ProduktFilter state={state} dispatch={dispatch} />
      <ProduktTabell state={state} dispatch={dispatch} />
    </div>
  );
};
