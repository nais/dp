import React from "react";
import ProduktTabell from "./produktTabell";
import { Sidetittel } from "nav-frontend-typografi";

export const ProduktListe = (): JSX.Element => {
  return (
    <div>
      <Sidetittel>Dataprodukter dashboard</Sidetittel>
      <ProduktTabell />
    </div>
  );
};
