import React, { useContext, useEffect, useReducer, useState } from "react";
import ProduktTabell from "./produktTabell";
import { DataProduktListe, hentProdukter } from "./produktAPI";
import ProduktFilter from "./produktFilter";
import { Add } from "@navikt/ds-icons";
import { Link } from "react-router-dom";
import { UserContext } from "./userContext";
import { useLocation } from "react-router-dom";
import { useHistory } from 'react-router-dom'

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
function useQuery() {
  return new URLSearchParams(useLocation().search);
}
export const Hovedside = (): JSX.Element => {
  const user = useContext(UserContext);
  const query = useQuery();
  const history = useHistory();
  const queryParameters = (query.get('teams') || null)?.split(',');

  const [error, setError] = useState<string | null>();
  const [filters, setFilters] = useState<string[]>(queryParameters? queryParameters : []);
  const [produkter, setProdukter] = useState<DataProduktListe>();

  useEffect(() => {
    if (!queryParameters) {
      const localStorageFilters = window.localStorage.getItem("filters")
      if (localStorageFilters?.length) setFilters(JSON.parse(localStorageFilters))
    }
  }, [])

  useEffect(() => {
    window.localStorage.setItem("filters", JSON.stringify(filters))

    history.push({
      search: filters.length ? '?teams=' + filters.join(',') : ''
    })
  }, [filters])

  useEffect(() => {
    hentProdukter()
        .then((produkter) => {
          setProdukter(produkter);
          setError(null);
        })
        .catch((e) => {
          console.log(e)
          setError(e.toString());
        });
  }, []);

  if (error) {
    setTimeout(() => window.location.reload(false), 1500);
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
        <ProduktFilter produkter={produkter} filters={filters} setFilters={setFilters} />
        {user ? <ProduktNyKnapp /> : <></>}
      </div>
      <ProduktTabell produkter={produkter?.filter(p => (!filters.length) || filters.includes(p.data_product.team))} />
    </div>
  );
};
