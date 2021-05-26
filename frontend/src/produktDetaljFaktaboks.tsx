import React from "react";
import { DataProduktResponse, DataProduktTilgangListe } from "./produktAPI";
import moment from "moment";
import { Normaltekst, Systemtittel } from "nav-frontend-typografi";
import { DatalagerInfo } from "./produktDatalager";
import { ProduktTilganger } from "./produktDetaljTilganger";

export const ProduktInfoFaktaboks: React.FC<{
  produkt: DataProduktResponse;
  tilganger: DataProduktTilgangListe;
}> = ({ tilganger, produkt }) => {
  moment.locale("nb");

  return (
    <div className={"faktaboks"}>
      <Systemtittel className={"produktnavn"}>
        {produkt.data_product?.name}
      </Systemtittel>

      <Normaltekst>
        Produkteier: {produkt.data_product?.team || "uvisst"}
      </Normaltekst>
      <Normaltekst>
        Opprettet {moment(produkt.created).format("LLL")}
        {produkt.created !== produkt.updated
          ? ` (Oppdatert: ${moment(produkt.updated).fromNow()})`
          : ""}
      </Normaltekst>
      <Normaltekst className="beskrivelse">
        {produkt.data_product?.description || "Ingen beskrivelse"}
      </Normaltekst>
      {produkt.data_product?.datastore &&
        produkt.data_product?.datastore.map((ds) => <DatalagerInfo ds={ds} />)}

      <ProduktTilganger produkt={produkt} tilganger={tilganger} />
    </div>
  );
};
