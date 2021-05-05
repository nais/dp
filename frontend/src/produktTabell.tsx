import React, { useEffect, useReducer } from "react";
import "nav-frontend-tabell-style";
import NavFrontendSpinner from "nav-frontend-spinner";
import { Select } from "nav-frontend-skjema";
import { DataProdukt, DataProduktList, hentProdukter } from "./produktAPI";

type ProduktTabellState = {
  loading: boolean;
  error: string | null;
  products: DataProduktList;
};

const initialState: ProduktTabellState = {
  loading: true,
  products: [],
  error: null,
};

interface ProduktTabellAction {
  type: string;
  results: DataProduktList;
}

const ProduktTabellReducer = (
  prevState: ProduktTabellState,
  action: ProduktTabellAction
): ProduktTabellState => {
  switch (action.type) {
    case "FETCH_DONE":
      return {
        ...prevState,
        products: action.results,
        loading: false,
        error: null,
      };
  }
  return prevState;
};

interface ProduktProps {
  produkt: DataProdukt;
}

const Produkt = ({ produkt }: ProduktProps) => {
  if (!produkt.data_product) return null;
  return (
    <tr>
      <td>{produkt.data_product ? produkt.data_product.name : "tom"}</td>
    </tr>
  );
};

export const ProduktTabell = () => {
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
      <table className={"tabell"}>
        <thead>
          <tr>
            <th>Navn</th>
          </tr>
        </thead>
        <tbody>
          {state.loading ? (
            <NavFrontendSpinner />
          ) : (
            state.products.map((x) => <Produkt key={x.id} produkt={x} />)
          )}
        </tbody>
      </table>
    </div>
  );
};

export default ProduktTabell;
