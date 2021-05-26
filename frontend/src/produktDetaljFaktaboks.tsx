import React from "react";
import { DataProduktResponse, DataProduktTilgangListe } from "./produktAPI";
import moment from "moment";
import { Normaltekst, Systemtittel } from "nav-frontend-typografi";
import { DatalagerInfo } from "./produktDatalager";
import { ProduktTilganger } from "./produktDetaljTilganger";

export const ProduktFaktaboks: React.FC<{
  produkt: DataProduktResponse;
  tilganger: DataProduktTilgangListe;
}> = ({ tilganger, produkt }) => {
  moment.locale("nb");

  return (
    <div className={"infoBoks"}>
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
    </div>
  );
};
