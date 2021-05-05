import React, { useEffect, useReducer } from "react";
import * as z from "zod";
import "nav-frontend-tabell-style";

//import ReactDOM from "react-dom";
//import "./index.less";
/*
    "id": "1JgfSOUm0ztjhiUljZHX",
    "name": "container resource usage 2",
    "description": "beskrivelse",
    "uri": "https://uri.com",
    "updated": "2021-05-04T13:28:47.474283Z",
    "created": "2021-05-04T09:54:26.129149Z"
 */
const DataProduktSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    description: z.string(),
    uri: z.string(),
    updated: z.string(),
    created: z.string(),
  })
  .nonstrict()
  .partial();
export const DataProduktListSchema = DataProduktSchema.array();

type DataProdukt = z.infer<typeof DataProduktSchema>;
type DataProduktList = z.infer<typeof DataProduktListSchema>;

export const getProducts = async (): Promise<DataProduktList> => {
  let apiURL = "http://localhost:8080/dataproducts";
  const res = await fetch(apiURL);
  const json = await res.json();
  return DataProduktListSchema.parse(json);
};

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

const Produkt = ({ produkt }: ProduktProps) => (
  <tr>
    <td>{produkt.name}</td>
  </tr>
);

export const ProduktTabell = () => {
  const [state, dispatch] = useReducer(ProduktTabellReducer, initialState);

  useEffect(() => {
    getProducts().then((products) => {
      dispatch({
        type: "FETCH_DONE",
        results: products,
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
          {state.products.map((x) => (
            <Produkt key={x.id} produkt={x} />
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ProduktTabell;
