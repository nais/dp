import React from "react";
import "nav-frontend-tabell-style";
import NavFrontendSpinner from "nav-frontend-spinner";
import { DataProdukt } from "./produktAPI";
import { ProduktListeState } from "./produktListe";
import { Link } from "react-router-dom";
import "./produktTabell.less";
interface ProduktProps {
  produkt: DataProdukt;
}

const Produkt = ({ produkt }: ProduktProps) => {
  if (!produkt.data_product) return null;
  return (
    <tr>
      <td>{produkt.data_product.owner}</td>
      <td>
        <Link to={`/produkt/${produkt.id}`}>{produkt.data_product.name}</Link>
      </td>
      <td>{produkt.data_product.description}</td>
    </tr>
  );
};

interface ProduktTabellProps {
  state: ProduktListeState;
  dispatch: any; // chickening out again
}

export const ProduktTabell = ({ state, dispatch }: ProduktTabellProps) => {
  return (
    <div>
      <table className={"tabell"}>
        <thead>
          <tr>
            <th>Produkteier</th>

            <th>Navn</th>
            <th>Beskrivelse</th>
          </tr>
        </thead>
        <tbody>
          {state.loading ||
            state.filtered_products.map((x) => (
              <Produkt key={x.id} produkt={x} />
            ))}
        </tbody>
      </table>
      {!state.loading && !state.products.length ? (
        <p style={{ textAlign: "center", fontStyle: "italic", margin: "2%" }}>
          Ingen dataprodukter i katalogen
        </p>
      ) : null}
      {state.loading ? (
        <p style={{ textAlign: "center", margin: "2%" }}>
          <NavFrontendSpinner />
        </p>
      ) : null}
    </div>
  );
};

export default ProduktTabell;
